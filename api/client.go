package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/pomerium/zero-sdk/token"
)

type client struct {
	tokenCache  *token.Cache
	httpClient  *http.Client
	minTokenTTL time.Duration
}

type APIClientOption func(*client)

func WithMinTokenTTL(minTokenTTL time.Duration) APIClientOption {
	return func(c *client) {
		c.minTokenTTL = minTokenTTL
	}
}

func WithClient(httpClient *http.Client) APIClientOption {
	return func(c *client) {
		c.httpClient = httpClient
	}
}

func NewAPIClient(endpoint string, refreshToken string, opts ...APIClientOption) (ClientWithResponsesInterface, error) {
	c := new(client)
	opts = append([]APIClientOption{
		WithMinTokenTTL(time.Minute * 5),
		WithClient(http.DefaultClient),
	}, opts...)
	for _, opt := range opts {
		opt(c)
	}

	fetcher, err := NewTokenFetcher(endpoint, WithHTTPClient(c.httpClient))
	if err != nil {
		return nil, fmt.Errorf("error creating token fetcher: %w", err)
	}
	c.tokenCache = token.NewCache(fetcher, refreshToken)

	return NewClientWithResponses(endpoint, WithHTTPClient(c))
}

func (c *client) Do(req *http.Request) (*http.Response, error) {
	ctx := req.Context()
	token, err := c.tokenCache.GetToken(ctx, c.minTokenTTL)
	if err != nil {
		return nil, fmt.Errorf("error getting token: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	return c.httpClient.Do(req)
}
