package connect

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	grpc_backoff "google.golang.org/grpc/backoff"

	"github.com/pomerium/zero-sdk/token"
)

type client struct {
	config      *Config
	tokenCache  *token.Cache
	minTokenTTL time.Duration
}

func NewAuthorizedConnectClient(
	ctx context.Context,
	endpoint string,
	cache *token.Cache,
) (ConnectClient, error) {
	cfg, err := NewConfig(endpoint)
	if err != nil {
		return nil, err
	}

	cc := &client{
		tokenCache: cache,
		config:     cfg,
		// streaming connection would reset based on token duration,
		// so we need it be close to max duration 1hr
		minTokenTTL: time.Minute * 55,
	}

	grpcConn, err := cc.getGRPCConn(ctx)
	if err != nil {
		return nil, err
	}

	return NewConnectClient(grpcConn), nil
}

func (c *client) getGRPCConn(ctx context.Context) (*grpc.ClientConn, error) {
	conn, err := grpc.DialContext(ctx,
		c.config.GetConnectionURI(),
		append(c.config.opts,
			grpc.WithPerRPCCredentials(c),
			grpc.WithConnectParams(grpc.ConnectParams{
				Backoff:           grpc_backoff.DefaultConfig,
				MinConnectTimeout: 1 * time.Second,
			}),
		)...)
	if err != nil {
		return nil, fmt.Errorf("error dialing grpc server: %w", err)
	}
	return conn, nil
}

// GetRequestMetadata implements credentials.PerRPCCredentials
func (c *client) GetRequestMetadata(ctx context.Context, _ ...string) (map[string]string, error) {
	token, err := c.tokenCache.GetToken(ctx, c.minTokenTTL)
	if err != nil {
		return nil, err
	}
	return map[string]string{
		"authorization": fmt.Sprintf("Bearer %s", token),
	}, nil
}

// RequireTransportSecurity implements credentials.PerRPCCredentials
func (c *client) RequireTransportSecurity() bool {
	return c.config.RequireTLS()
}
