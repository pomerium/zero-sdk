package cluster

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

const (
	defaultMinTokenTTL = time.Minute * 5
)

type client struct {
	tokenProvider TokenProviderFn
	httpClient    *http.Client
	minTokenTTL   time.Duration
}

// TokenProviderFn is a function that returns a token that is expected to be valid for at least minTTL
type TokenProviderFn func(ctx context.Context, minTTL time.Duration) (string, error)

func NewAuthorizedClient(
	endpoint string,
	tokenProvider TokenProviderFn,
	httpClient *http.Client,
) (ClientWithResponsesInterface, error) {
	c := &client{
		minTokenTTL: defaultMinTokenTTL,
		httpClient:  httpClient,
	}

	c.tokenProvider = tokenProvider

	return NewClientWithResponses(endpoint, WithHTTPClient(c))
}

func (c *client) Do(req *http.Request) (*http.Response, error) {
	ctx := req.Context()
	token, err := c.tokenProvider(ctx, c.minTokenTTL)
	if err != nil {
		return nil, fmt.Errorf("error getting token: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	return c.httpClient.Do(req)
}
