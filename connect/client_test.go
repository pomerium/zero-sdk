package connect_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	cluster_api "github.com/pomerium/zero-sdk/cluster"
	"github.com/pomerium/zero-sdk/connect"
	"github.com/pomerium/zero-sdk/token"
)

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
