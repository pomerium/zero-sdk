package cluster

import (
	"net/url"
	"sync"
	"time"
)

type URLCache struct {
	mx    sync.RWMutex
	cache map[string]DownloadCacheEntry
}

type DownloadCacheEntry struct {
	// URL is the URL to download the bundle from.
	URL url.URL
	// ExpiresAt is the time at which the URL expires.
	ExpiresAt time.Time
	// CaptureHeaders is a list of headers to capture from the response.
	CaptureHeaders []string
}

func NewURLCache() *URLCache {
	return &URLCache{
		cache: make(map[string]DownloadCacheEntry),
	}
}

func (c *URLCache) Get(key string, minTTL time.Duration) (*DownloadCacheEntry, bool) {
	c.mx.RLock()
	defer c.mx.RUnlock()

	entry, ok := c.cache[key]
	if !ok {
		return nil, false
	}

	if time.Until(entry.ExpiresAt) < minTTL {
		return nil, false
	}

	return &entry, true
}

func (c *URLCache) Set(key string, entry DownloadCacheEntry) {
	c.mx.Lock()
	defer c.mx.Unlock()

	c.cache[key] = entry
}
