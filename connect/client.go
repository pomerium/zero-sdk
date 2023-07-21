package connect

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"time"

	"google.golang.org/grpc"
	grpc_backoff "google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/pomerium/zero-sdk/token"
)

type client struct {
	endpoint    *url.URL
	tokenCache  *token.Cache
	minTokenTTL time.Duration
	requireTLS  bool
}

func NewAuthorizedConnectClient(
	ctx context.Context,
	endpoint string,
	cache *token.Cache,
) (ConnectClient, error) {
	url, err := url.Parse(endpoint)
	if err != nil {
		return nil, fmt.Errorf("error parsing endpoint url: %w", err)
	}

	cc := &client{
		tokenCache: cache,
		endpoint:   url,
		// streaming connection would reset based on token duration,
		// so we need it be close to max duration 1hr
		minTokenTTL: time.Minute * 55,
		requireTLS:  url.Scheme == "https",
	}

	grpcConn, err := cc.getGRPCConn(ctx, endpoint)
	if err != nil {
		return nil, err
	}

	return NewConnectClient(grpcConn), nil
}

func (c *client) getGRPCConn(ctx context.Context, endpoint string) (*grpc.ClientConn, error) {
	opts := []grpc.DialOption{
		grpc.WithPerRPCCredentials(c),
		grpc.WithConnectParams(grpc.ConnectParams{
			Backoff:           grpc_backoff.DefaultConfig,
			MinConnectTimeout: 1 * time.Second,
		}),
	}

	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, fmt.Errorf("error parsing endpoint url: %w", err)
	}

	host, port, err := net.SplitHostPort(u.Host)
	if err != nil {
		return nil, fmt.Errorf("error splitting host and port: %w", err)
	}
	if c.endpoint.Scheme == "http" {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if port == "" {
			port = "80"
		}
	} else if c.endpoint.Scheme == "https" {
		if port == "" {
			port = "443"
		}
	} else {
		return nil, fmt.Errorf("unsupported url scheme: %s", c.endpoint.Scheme)
	}

	if port == "" {
		return nil, fmt.Errorf("port should be specified")
	}

	// endpoint should be a URI https://github.com/grpc/grpc/blob/master/doc/naming.md
	conn, err := grpc.DialContext(ctx, fmt.Sprintf("dns:%s:%s", host, port), opts...)
	if err != nil {
		return nil, fmt.Errorf("error dialing grpc server: %w", err)
	}
	return conn, nil
}

func (c *client) GetRequestMetadata(ctx context.Context, _ ...string) (map[string]string, error) {
	token, err := c.tokenCache.GetToken(ctx, c.minTokenTTL)
	if err != nil {
		return nil, err
	}
	return map[string]string{
		"authorization": fmt.Sprintf("Bearer %s", token),
	}, nil
}

// RequireTransportSecurity indicates whether the credentials requires
// transport security.
func (c *client) RequireTransportSecurity() bool {
	return c.requireTLS
}
