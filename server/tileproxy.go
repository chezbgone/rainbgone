package server

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/joho/godotenv"
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

var tileCache = newCache()

var tileHTTPClient = &http.Client{
	Timeout: 10 * time.Second,
}

var MAPTILER_KEY string

func init() {
	env, err := godotenv.Read()
	if err != nil {
		panic(err)
	}
	MAPTILER_KEY = env["MAPTILER_KEY"]
}

func TileHandler(w http.ResponseWriter, r *http.Request) {
	// Expected format: /map/tiles/{type}/{z}/{x}/{y}.png
	path := r.URL.Path
	parts := strings.Split(path, "/")

	tileType := parts[3]
	z := parts[4]
	x := parts[5]
	y := strings.TrimSuffix(parts[6], ".png")

	if tileType != "temp" && tileType != "precipitation" {
		http.Error(w, "Invalid tile type", http.StatusBadRequest)
		return
	}

	cacheKey := fmt.Sprintf("%s/%s/%s/%s", tileType, z, x, y)
	if cached, found := tileCache.Get(cacheKey); found {
		w.Header().Set("Content-Type", "image/png")
		w.Write(cached)
		return
	}

	var layer string
	switch tileType {
	case "temp":
		layer = "temp_new"
	case "precipitation":
		layer = "precipitation_new"
	}

	url := fmt.Sprintf("https://api.maptiler.com/tiles/%s/%s/%s/%s.png?api_key=%s",
		layer, z, x, y, MAPTILER_KEY)

	resp, err := tileHTTPClient.Get(url)
	if err != nil {
		http.Error(w, "Failed to fetch tile", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Failed to fetch tile", resp.StatusCode)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read tile data", http.StatusInternalServerError)
		return
	}

	tileCache.Set(cacheKey, body, 5*time.Minute)

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Cache-Control", "public, max-age=300") // 5 minutes
	w.Write(body)
}
