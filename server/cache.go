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
	items map[string]cacheItem
	mu    sync.RWMutex
}

func newCache() *cache {
	return &cache{
		items: make(map[string]cacheItem),
	}
}

func (c *cache) Get(key string) ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, found := c.items[key]
	if !found || time.Now().After(item.expiresAt) {
		return nil, false
	}
	return item.body, true
}

func (c *cache) Set(key string, body []byte, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = cacheItem{
		body:      body,
		expiresAt: time.Now().Add(ttl),
	}
}
