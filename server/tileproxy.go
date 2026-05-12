package server

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

var tileCache = newCache()

var tileHTTPClient = &http.Client{
	Timeout: 10 * time.Second,
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
		layer, z, x, y, appConfig.mapTilerKey)

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
