package token

import (
	"context"
	"fmt"
	"sync"
	"time"
)

const (
	maxLockWait = 30 * time.Second
)

// Cache is a thread-safe cache of a authorization token
// that may be used across http and grpc clients
type Cache struct {
	sync.RWMutex
	refreshToken string
	token        *Token
	fetcher      Fetcher
	TimeNow      func() time.Time
}

type Fetcher func(ctx context.Context, refreshToken string) (*Token, error)

type Token struct {
	Bearer  string
	Expires time.Time
}

func (t *Token) ExpiresAfter(tm time.Time) bool {
	return t != nil && t.Expires.After(tm)
}

func NewCache(fetcher Fetcher, refreshToken string) *Cache {
	return &Cache{
		fetcher:      fetcher,
		refreshToken: refreshToken,
	}
}

func (c *Cache) timeNow() time.Time {
	if c.TimeNow != nil {
		return c.TimeNow()
	}
	return time.Now()
}

// GetToken returns the current token if its at least `minTTL` from expiration, or fetches a new one.
func (c *Cache) GetToken(ctx context.Context, minTTL time.Duration) (string, error) {
	minExpiration := c.timeNow().Add(minTTL)

	c.RLock()
	token := c.token
	c.RUnlock()

	if token.ExpiresAfter(minExpiration) {
		return token.Bearer, nil
	}

	return c.forceRefreshToken(ctx, minExpiration)
}

func (c *Cache) forceRefreshToken(ctx context.Context, minExpiration time.Time) (string, error) {
	c.Lock()
	defer c.Unlock()

	ctx, cancel := context.WithTimeout(ctx, maxLockWait)
	defer cancel()

	if c.token.ExpiresAfter(minExpiration) {
		return c.token.Bearer, nil
	}

	token, err := c.fetcher(ctx, c.refreshToken)
	if err != nil {
		return "", err
	}
	c.token = token

	if token.Expires.Before(minExpiration) {
		return "", fmt.Errorf("new token cannot satisfy TTL: %v", minExpiration.Sub(token.Expires))
	}

	return token.Bearer, nil
}
