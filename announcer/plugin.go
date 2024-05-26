package main

import (
	"fmt"
	"log/slog"
	"net/netip"

	"github.com/jwhited/corebgp"
	"github.com/osrg/gobgp/v3/pkg/packet/bgp"
)

var _ corebgp.Plugin = (*plugin)(nil)

type plugin struct { // TODO: rename to peerPlugin ?
	// config is our configuration
	config peerConfig
	// capabilities we got from the remote peer
	addPathIPv4, addPathIPv6 bool
	// Writer for sending update messages to the remote peer
	writer corebgp.UpdateMessageWriter
	// logger for BGP messages
	logger *slog.Logger
}

func (p *plugin) GetCapabilities(c corebgp.PeerConfig) []corebgp.Capability {
	// capabilities send to the remote peer
	capabilities := make([]corebgp.Capability, 0)
	if p.config.EnableIPv4 {
		capabilities = append(capabilities, corebgp.NewMPExtensionsCapability(corebgp.AFI_IPV4, corebgp.SAFI_UNICAST))
	}
	if p.config.EnableIPv6 {
		capabilities = append(capabilities, corebgp.NewMPExtensionsCapability(corebgp.AFI_IPV6, corebgp.SAFI_UNICAST))
	}
	if p.config.EnableAddPath {
		tuples := make([]corebgp.AddPathTuple, 0)
		if p.config.EnableIPv4 {
			tuples = append(tuples, corebgp.AddPathTuple{
				AFI:  corebgp.AFI_IPV4,
				SAFI: corebgp.SAFI_UNICAST,
				Tx:   true,
				Rx:   true,
			})
		}
		if p.config.EnableIPv6 {
			tuples = append(tuples, corebgp.AddPathTuple{
				AFI:  corebgp.AFI_IPV6,
				SAFI: corebgp.SAFI_UNICAST,
				Tx:   true,
				Rx:   true,
			})
		}
		capabilities = append(capabilities, corebgp.NewAddPathCapability(tuples))
	}
	return capabilities
}

func (p *plugin) OnOpenMessage(peer corebgp.PeerConfig, routerID netip.Addr, capabilities []corebgp.Capability) *corebgp.Notification {
	// open message received from remote peer
	if p.config.EnableAddPath {
		p.addPathIPv4 = false
		p.addPathIPv6 = false
		for _, c := range capabilities {
			if c.Code != corebgp.CAP_ADD_PATH {
				continue
			}
			tuples, err := corebgp.DecodeAddPathTuples(c.Value)
			if err != nil {
				return err.(*corebgp.Notification)
			}
			for _, tuple := range tuples {
				if tuple.SAFI != corebgp.SAFI_UNICAST || !tuple.Tx {
					continue
				}
				if tuple.AFI == corebgp.AFI_IPV4 {
					p.addPathIPv4 = true
				} else if tuple.AFI == corebgp.AFI_IPV6 {
					p.addPathIPv6 = true
				}
			}
		}
	}
	return nil
}

func (p *plugin) OnEstablished(peer corebgp.PeerConfig, writer corebgp.UpdateMessageWriter) corebgp.UpdateMessageHandler {
	// established connection to peer
	p.writer = writer
	go p.sendRoutes()
	return p.handleUpdate
}

func (p *plugin) OnClose(peer corebgp.PeerConfig) {
	// peer closed
	p.writer = nil
}

func (p *plugin) handleUpdate(peer corebgp.PeerConfig, u []byte) *corebgp.Notification {
	// got update message
	msg := bgp.BGPUpdate{}
	if err := msg.DecodeFromBytes(u); err != nil {
		return convertGoBgpToCoreBgpError(err)
	}
	// TODO: send message data to other program for further processing
	p.logger.Info(fmt.Sprintf("%+v", msg))
	return nil
}

func (p *plugin) sendRoutes() {
	// TODO: advertise ipv4 routes
	if p.config.EnableIPv6 && len(p.config.Networks.V6) != 0 {
		nlriV6 := make([]bgp.AddrPrefixInterface, 0, len(p.config.Networks.V6))
		for _, pv6 := range p.config.Networks.V6 {
			nlriV6 = append(nlriV6, bgp.NewIPv6AddrPrefix(uint8(pv6.Bits()), pv6.Addr().String()))
		}

		// TODO: for ipv4 next hop we will need to set bgp.BGP_ATTR_TYPE_NEXT_HOP with func bgp.NewPathAttributeNextHop()
		// and put the ipv4 nlri directly into the bgp update message

		// for ipv6 use bgp.BGP_ATTR_TYPE_MP_REACH_NLRI
		// to advertise ipv4 with ipv6 next hop see https://datatracker.ietf.org/doc/html/rfc5549
		pathAttr := []bgp.PathAttributeInterface{
			bgp.NewPathAttributeOrigin(bgp.BGP_ORIGIN_ATTR_TYPE_IGP),
			bgp.NewPathAttributeAs4Path([]*bgp.As4PathParam{bgp.NewAs4PathParam(bgp.BGP_ASPATH_ATTR_TYPE_SEQ, []uint32{p.config.LocalAS})}),
			bgp.NewPathAttributeMpReachNLRI(p.config.LocalAddress.String(), nlriV6),
		}

		msg := bgp.BGPUpdate{
			WithdrawnRoutesLen:    0,
			WithdrawnRoutes:       nil,
			TotalPathAttributeLen: 0,
			PathAttributes:        pathAttr,
			NLRI:                  nil,
		}

		// send ipv6 routes
		b, err := msg.Serialize()
		if err != nil {
			return
		}
		err = p.writer.WriteUpdate(b)
		if err != nil {
			return
		}

		// send End-of-Rib for ipv6
		b, err = bgp.NewEndOfRib(bgp.RF_IPv6_UC).Body.Serialize()
		if err != nil {
			return
		}
		err = p.writer.WriteUpdate(b)
		if err != nil {
			return
		}
	}
}

func convertGoBgpToCoreBgpError(err error) *corebgp.Notification {
	if msgError, ok := err.(*bgp.MessageError); ok {
		return &corebgp.Notification{
			Code:    msgError.TypeCode,
			Subcode: msgError.SubTypeCode,
			Data:    msgError.Data,
		}
	}
	return nil
}
