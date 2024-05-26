package main

import (
	"net/netip"

	"github.com/jwhited/corebgp"
)

type Config struct {
	RouterID netip.Addr
	Peers    []peerConfig
}

type peerConfig struct {
	corebgp.PeerConfig
	LocalAddress  netip.Addr
	EnableAddPath bool
	EnableIPv4    bool
	EnableIPv6    bool
	Networks      struct {
		V4 []netip.Prefix
		V6 []netip.Prefix
	}
}
