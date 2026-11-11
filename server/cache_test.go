package server

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestCacheRemovesExpiredEntriesOnRead(t *testing.T) {
	c := newCache()
	c.Set("expired", []byte("body"), -time.Second)
	if _, hit := c.Get("expired"); hit {
		t.Fatal("expired entry was returned")
	}
	if len(c.items) != 0 {
		t.Fatal("expired entry is still resident")
	}
}

func TestCacheExpiresWithoutSleeping(t *testing.T) {
	c := newCache()
	now := time.Unix(100, 0)
	c.now = func() time.Time { return now }
	c.Set("a", []byte("body"), time.Minute)
	if body, hit := c.Get("a"); !hit || string(body) != "body" {
		t.Fatal("fresh entry was not returned")
	}
	now = now.Add(time.Minute)
	if _, hit := c.Get("a"); hit || len(c.items) != 0 || c.sizeBytes != 0 {
		t.Fatal("entry was not removed at its expiry")
	}
}

func TestCachePrunesExpiredEntriesOnWrite(t *testing.T) {
	c := newCache()
	now := time.Unix(100, 0)
	c.now = func() time.Time { return now }
	c.Set("old", []byte("body"), time.Minute)
	now = now.Add(2 * time.Minute)
	c.Set("new", []byte("body"), time.Minute)
	if len(c.items) != 1 || c.sizeBytes != len("newbody") {
		t.Fatal("writing a new entry retained expired data")
	}
}

func TestCacheEnforcesEntryAndByteLimits(t *testing.T) {
	for _, test := range []struct {
		name       string
		maxEntries int
		maxBytes   int
	}{
		{"entries", 2, 100},
		{"bytes", 100, 10},
	} {
		t.Run(test.name, func(t *testing.T) {
			c := newCache()
			c.maxEntries, c.maxBytes = test.maxEntries, test.maxBytes
			now := time.Unix(100, 0)
			c.now = func() time.Time { return now }
			for _, key := range []string{"a", "b", "c"} {
				c.Set(key, []byte("body"), time.Minute)
				now = now.Add(time.Second)
			}
			if _, hit := c.Get("a"); hit {
				t.Fatal("oldest entry was not evicted")
			}
			if body, hit := c.Get("c"); !hit || string(body) != "body" {
				t.Fatal("newest entry was not retained")
			}
			if len(c.items) > c.maxEntries || c.sizeBytes > c.maxBytes {
				t.Fatal("cache exceeded its limits")
			}
		})
	}
}

func TestCacheReplacementAndOversizedEntries(t *testing.T) {
	c := newCache()
	c.maxBytes = 10
	c.Set("a", []byte("body"), time.Minute)
	c.Set("a", []byte("x"), time.Minute)
	if c.sizeBytes != 2 || len(c.items) != 1 {
		t.Fatal("replacing an entry double-counted its bytes")
	}
	c.Set("b", make([]byte, 10), time.Minute)
	if _, hit := c.Get("b"); hit || c.sizeBytes != 2 {
		t.Fatal("oversized entry was cached")
	}
}

func TestCacheConcurrentAccess(t *testing.T) {
	c := newCache()
	var workers sync.WaitGroup
	for i := range 8 {
		workers.Add(1)
		go func() {
			defer workers.Done()
			for j := range 100 {
				key := fmt.Sprintf("%d/%d", i, j)
				c.Set(key, []byte("body"), time.Minute)
				c.Get(key)
			}
		}()
	}
	workers.Wait()
	if len(c.items) > c.maxEntries {
		t.Fatal("concurrent writes exceeded the entry limit")
	}
}
