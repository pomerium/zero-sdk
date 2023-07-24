package connect_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	cluster_api "github.com/pomerium/zero-sdk/cluster"
	"github.com/pomerium/zero-sdk/connect"
	"github.com/pomerium/zero-sdk/token"
)

func TestConfig(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		endpoint      string
		connectionURI string
		requireTLS    bool
		expectError   bool
	}{
		{"", "", false, true},
		{"http://", "", false, true},
		{"https://", "", true, true},
		{"localhost:8721", "", false, true},
		{"http://localhost:8721", "dns:localhost:8721", false, false},
		{"https://localhost:8721", "dns:localhost:8721", true, false},
		{"http://localhost:8721/", "dns:localhost:8721", false, false},
		{"https://localhost:8721/", "dns:localhost:8721", true, false},
		{"http://localhost:8721/path", "dns:localhost:8721", false, true},
		{"https://localhost:8721/path", "dns:localhost:8721", true, true},
		{"http://localhost", "dns:localhost:80", false, false},
		{"https://localhost:443", "dns:localhost:443", true, false},
	} {
		tc := tc
		t.Run(tc.endpoint, func(t *testing.T) {
			t.Parallel()
			cfg, err := connect.NewConfig(tc.endpoint)
			if tc.expectError {
				require.Error(t, err)
				return
			}
			if assert.NoError(t, err) {
				assert.Equal(t, tc.connectionURI, cfg.GetConnectionURI(), "connection uri")
				assert.Equal(t, tc.requireTLS, cfg.RequireTLS(), "require tls")
			}
		})
	}
}

func TestConnectClient(t *testing.T) {
	refreshToken := os.Getenv("CONNECT_CLUSTER_IDENTITY_TOKEN")
	if refreshToken == "" {
		t.Skip("CONNECT_CLUSTER_IDENTITY_TOKEN not set")
	}

	connectServerEndpoint := os.Getenv("CONNECT_SERVER_ENDPOINT")
	if connectServerEndpoint == "" {
		connectServerEndpoint = "http://localhost:8721"
	}

	clusterAPIEndpoint := os.Getenv("CLUSTER_API_ENDPOINT")
	if clusterAPIEndpoint == "" {
		clusterAPIEndpoint = "http://localhost:8720/cluster/v1"
	}

	fetcher, err := cluster_api.NewTokenFetcher(clusterAPIEndpoint)
	require.NoError(t, err, "error creating token fetcher")

	ctx := context.Background()
	deadline, ok := t.Deadline()
	if ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithDeadline(ctx, deadline.Add(-1*time.Second))
		t.Cleanup(cancel)
	}

	tokenCache := token.NewCache(fetcher, refreshToken)

	connectClient, err := connect.NewAuthorizedConnectClient(ctx, connectServerEndpoint, tokenCache)
	require.NoError(t, err, "error creating connect client")

	stream, err := connectClient.Subscribe(ctx, &connect.SubscribeRequest{})
	require.NoError(t, err, "error subscribing")

	for {
		msg, err := stream.Recv()
		require.NoError(t, err, "error receiving message")
		t.Log(msg)
	}
}
