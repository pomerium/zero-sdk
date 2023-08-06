// Package mux provides the way to listen for updates from the cloud
package mux

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/pomerium/zero-sdk/connect"
	"github.com/pomerium/zero-sdk/fanout"
)

// Run starts the updates service, listening for updates from the cloud
// until the context is canceled
func Start(ctx context.Context, client connect.ConnectClient) *Mux {
	ctx, cancel := context.WithCancelCause(ctx)
	svc := &Mux{
		client: client,
		mux:    fanout.Start[message](ctx),
	}
	go svc.run(ctx, cancel)
	return svc
}

type Mux struct {
	client connect.ConnectClient
	mux    *fanout.FanOut[message]
}

// run does not do any kind of backoff
// due to service having a ttl constraint on the server side,
// thus we want to reconnect as soon as possible
func (svc *Mux) run(ctx context.Context, cancel context.CancelCauseFunc) {
	logger := log.Ctx(ctx).With().Str("service", "connect-mux").Logger().Level(zerolog.DebugLevel)

	for ctx.Err() == nil {
		err := svc.subscribeAndDispatch(ctx)
		if err != nil {
			logger.Err(err).Msg("running")
		}
		if errors.Is(err, nonRetryableError{}) {
			cancel(err)
			return
		}
	}
	cancel(fmt.Errorf("connect-mux run: %w", context.Cause(ctx)))
}

type nonRetryableError struct {
	error
}

func (e nonRetryableError) Is(target error) bool {
	//nolint:errorlint // we want to check for the exact type
	_, ok := target.(nonRetryableError)
	return ok
}

func (svc *Mux) subscribeAndDispatch(ctx context.Context) (err error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	stream, err := svc.client.Subscribe(ctx, &connect.SubscribeRequest{})
	if err != nil {
		return fmt.Errorf("subscribe: %w", err)
	}

	if err = svc.onConnected(ctx); err != nil {
		return fmt.Errorf("on connected: %w", err)
	}
	defer func() {
		err = multierror.Append(
			err,
			nonRetryableError{svc.onDisconnected(ctx)},
		).ErrorOrNil()
	}()

	for {
		msg, err := stream.Recv()
		if err != nil {
			return fmt.Errorf("receive: %w", err)
		}
		err = svc.onMessage(ctx, msg)
		if err != nil {
			return nonRetryableError{fmt.Errorf("on message: %w", err)}
		}
	}
}
