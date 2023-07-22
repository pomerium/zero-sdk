package cluster

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/pomerium/zero-sdk/token"
)

func NewTokenFetcher(endpoint string, opts ...ClientOption) (token.Fetcher, error) {
	client, err := NewClientWithResponses(endpoint, opts...)
	if err != nil {
		return nil, fmt.Errorf("error creating client: %w", err)
	}

	return func(ctx context.Context, refreshToken string) (*token.Token, error) {
		now := time.Now()
		resp, err := client.ExchangeClusterIdentityTokenWithResponse(ctx, ExchangeTokenRequest{
			RefreshToken: refreshToken,
		})
		if err != nil {
			return nil, fmt.Errorf("error fetching id token: %w", err)
		}

		if resp.JSON400 != nil {
			return nil, fmt.Errorf("error fetching id token: %s", resp.JSON400.Error)
		}

		if resp.JSON200 == nil {
			return nil, fmt.Errorf("unexpected response from GetIdToken: %d: %s", resp.StatusCode(), string(resp.Body))
		}

		expiresSeconds, err := strconv.ParseInt(resp.JSON200.ExpiresInSeconds, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("error parsing expires in: %w", err)
		}

		return &token.Token{
			Bearer:  resp.JSON200.IdToken,
			Expires: now.Add(time.Duration(expiresSeconds) * time.Second),
		}, nil
	}, nil
}
