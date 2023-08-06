package mux

import (
	"context"
	"fmt"

	"github.com/pomerium/zero-sdk/connect"
)

// Watch watches for changes to the config until either context is cancelled,
// or an error occurs while muxing
func (svc *Mux) Watch(ctx context.Context, opts ...WatchOption) error {
	cfg := newConfig(opts...)
	return svc.mux.Receive(ctx, func(ctx context.Context, msg message) error {
		return dispatch(ctx, cfg, msg)
	})
}

func dispatch(ctx context.Context, cfg *config, msg message) error {
	switch {
	case msg.stateChange != nil:
		switch *msg.stateChange {
		case connected:
			cfg.onConnected(ctx)
		case disconnected:
			cfg.onDisconnected(ctx)
		default:
			return fmt.Errorf("unknown state change")
		}
	case msg.Message != nil:
		switch msg.Message.Message.(type) {
		case *connect.Message_ConfigUpdated:
			cfg.onBundleUpdated(ctx, "config")
		case *connect.Message_BootstrapConfigUpdated:
			cfg.onBootstrapConfigUpdated(ctx)
		default:
			return fmt.Errorf("unknown message type")
		}
	default:
		return fmt.Errorf("unknown message payload")
	}
	return nil
}

type message struct {
	*stateChange
	*connect.Message
}

type stateChange string

const (
	connected    stateChange = "connected"
	disconnected stateChange = "disconnected"
)

func (svc *Mux) onConnected(ctx context.Context) error {
	s := connected
	return svc.mux.Publish(ctx, message{stateChange: &s})
}

func (svc *Mux) onDisconnected(ctx context.Context) error {
	s := disconnected
	return svc.mux.Publish(ctx, message{stateChange: &s})
}

func (svc *Mux) onMessage(ctx context.Context, msg *connect.Message) error {
	return svc.mux.Publish(ctx, message{Message: msg})
}
