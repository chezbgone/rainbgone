package server

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

var backgroundTileCache = newCache()

var backgroundTileHTTPClient = &http.Client{
	Timeout: 10 * time.Second,
}

func BackgroundTileHandler(w http.ResponseWriter, r *http.Request) {
	// Expected format: /map/background-tiles/{variant}/{z}/{x}/{y}.png
	path := r.URL.Path
	parts := strings.Split(path, "/")
	if len(parts) != 7 {
		http.Error(w, "Invalid background tile path", http.StatusBadRequest)
		return
	}

	variant := parts[3]
	z := parts[4]
	x := parts[5]
	y := strings.TrimSuffix(parts[6], ".png")

	if variant != "temp" && variant != "precipitation" {
		http.Error(w, "Invalid background tile variant", http.StatusBadRequest)
		return
	}

	cacheKey := fmt.Sprintf("%s/%s/%s/%s", variant, z, x, y)
	if cached, found := backgroundTileCache.Get(cacheKey); found {
		w.Header().Set("Content-Type", "image/png")
		w.Write(cached)
		return
	}

	var layer string
	switch variant {
	case "temp":
		layer = "temp_new"
	case "precipitation":
		layer = "precipitation_new"
	}

	url := fmt.Sprintf("https://api.maptiler.com/tiles/%s/%s/%s/%s.png?api_key=%s",
		layer, z, x, y, appConfig.mapTilerKey)

	resp, err := backgroundTileHTTPClient.Get(url)
	if err != nil {
		http.Error(w, "Failed to fetch background tile", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Failed to fetch background tile", resp.StatusCode)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read background tile data", http.StatusInternalServerError)
		return
	}

	backgroundTileCache.Set(cacheKey, body, 5*time.Minute)

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Cache-Control", "public, max-age=300") // 5 minutes
	w.Write(body)
}
