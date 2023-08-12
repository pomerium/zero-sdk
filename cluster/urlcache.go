package cluster

import (
	"net/url"
	"sync"
	"time"
)

type URLCache struct {
	mx    sync.RWMutex
	cache map[string]urlEntry
}

type urlEntry struct {
	URL       url.URL
	ExpiresAt time.Time
}

func NewURLCache() *URLCache {
	return &URLCache{
		cache: make(map[string]urlEntry),
	}
}

func (c *URLCache) Get(key string, minTTL time.Duration) (url.URL, bool) {
	c.mx.RLock()
	defer c.mx.RUnlock()

	entry, ok := c.cache[key]
	if !ok {
		return url.URL{}, false
	}

	if time.Until(entry.ExpiresAt) < minTTL {
		return url.URL{}, false
	}

	return entry.URL, true
}

func (c *URLCache) Set(key string, value url.URL, expires time.Time) {
	c.mx.Lock()
	defer c.mx.Unlock()

	c.cache[key] = urlEntry{
		URL:       value,
		ExpiresAt: expires,
	}
}
