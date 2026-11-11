package server

import (
	"sync"
	"time"
)

type cacheItem struct {
	body      []byte
	expiresAt time.Time
}

type cache struct {
	items      map[string]cacheItem
	mu         sync.Mutex
	sizeBytes  int
	maxEntries int
	maxBytes   int
	now        func() time.Time
}

func newCache() *cache {
	return &cache{
		items:      make(map[string]cacheItem),
		maxEntries: 256,
		maxBytes:   16 << 20,
		now:        time.Now,
	}
}

func (c *cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, found := c.items[key]
	if !found {
		return nil, false
	}
	if !c.now().Before(item.expiresAt) {
		c.remove(key)
		return nil, false
	}
	return item.body, true
}

func (c *cache) Set(key string, body []byte, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.remove(key)
	cost := len(key) + len(body)
	if ttl <= 0 || c.maxEntries <= 0 || cost > c.maxBytes {
		return
	}
	now := c.now()
	for key, item := range c.items {
		if !now.Before(item.expiresAt) {
			c.remove(key)
		}
	}
	// Each cache has a fixed TTL, so the earliest expiry is its oldest entry.
	for len(c.items) >= c.maxEntries || c.sizeBytes+cost > c.maxBytes {
		var oldestKey string
		var oldest time.Time
		for key, item := range c.items {
			if oldest.IsZero() || item.expiresAt.Before(oldest) {
				oldestKey, oldest = key, item.expiresAt
			}
		}
		c.remove(oldestKey)
	}
	c.items[key] = cacheItem{
		body:      body,
		expiresAt: now.Add(ttl),
	}
	c.sizeBytes += cost
}

// remove is called with c.mu held.
func (c *cache) remove(key string) {
	if item, found := c.items[key]; found {
		c.sizeBytes -= len(key) + len(item.body)
		delete(c.items, key)
	}
}
