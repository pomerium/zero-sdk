package token

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"
)

const (
	maxLockWait = 30 * time.Second
)

// Cache is a thread-safe cache of a authorization token
// that may be used across http and grpc clients
type Cache struct {
	TimeNow func() time.Time

	refreshToken string
	fetcher      Fetcher

	lock  chan struct{}
	token atomic.Value
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
		lock:         make(chan struct{}, 1),
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

	token, ok := c.token.Load().(*Token)
	if ok && token.ExpiresAfter(minExpiration) {
		return token.Bearer, nil
	}

	return c.forceRefreshToken(ctx, minExpiration)
}

func (c *Cache) forceRefreshToken(ctx context.Context, minExpiration time.Time) (string, error) {
	select {
	case c.lock <- struct{}{}:
	case <-ctx.Done():
		return "", ctx.Err()
	}
	defer func() {
		<-c.lock
	}()

	ctx, cancel := context.WithTimeout(ctx, maxLockWait)
	defer cancel()

	token, ok := c.token.Load().(*Token)
	if ok && token.ExpiresAfter(minExpiration) {
		return token.Bearer, nil
	}

	token, err := c.fetcher(ctx, c.refreshToken)
	if err != nil {
		return "", err
	}
	c.token.Store(token)

	if token.Expires.Before(minExpiration) {
		return "", fmt.Errorf("new token cannot satisfy TTL: %v", minExpiration.Sub(token.Expires))
	}

	return token.Bearer, nil
}
