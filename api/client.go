package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/pomerium/zero-sdk/token"
)

const (
	defaultMinTokenTTL = time.Minute * 5
)

type client struct {
	tokenCache  *token.Cache
	httpClient  *http.Client
	minTokenTTL time.Duration
}

func NewAuthorizedClient(
	endpoint string,
	refreshToken string,
	httpClient *http.Client,
) (ClientWithResponsesInterface, error) {
	c := &client{
		httpClient:  httpClient,
		minTokenTTL: defaultMinTokenTTL,
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
