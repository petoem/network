package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"

	"github.com/jwhited/corebgp"
)

func NewAnnouncer(config Config) error {
	setupDefaultLogger()
	setCoreBgpInternalDefaultLogger()

	srv, err := corebgp.NewServer(config.RouterID)
	if err != nil {
		return fmt.Errorf("error constructing server: %w", err)
	}

	for _, peer := range config.Peers {
		p := &plugin{
			config: peer,
			logger: slog.Default().With("subsystem", "announcer"),
		}
		err = srv.AddPeer(peer.PeerConfig, p, corebgp.WithLocalAddress(peer.LocalAddress))
		if err != nil {
			return fmt.Errorf("error adding peer: %w", err)
		}
	}

	lconfig := &net.ListenConfig{}
	listener, err := lconfig.Listen(context.Background(), "tcp", ":179")
	if err != nil {
		return fmt.Errorf("error constructing listener: %w", err)
	}
	err = srv.Serve([]net.Listener{listener})
	if err != nil {
		return fmt.Errorf("serve: %w", err)
	}

	return nil
}

func setupDefaultLogger() {
	jsonLogger := slog.NewJSONHandler(os.Stdout, nil)
	logger := slog.New(jsonLogger)
	slog.SetDefault(logger)
}

func setCoreBgpInternalDefaultLogger() {
	coreBgpLogger := slog.Default().With("subsystem", "corebgp")
	corebgp.SetLogger(func(a ...any) {
		coreBgpLogger.Info(fmt.Sprint(a...))
	})
}
