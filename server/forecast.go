package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

var forecastHTTPClient = &http.Client{
	Timeout: 10 * time.Second,
}

var forecastCache = newCache()

const forecastCacheTTL = 1 * time.Minute

func getForecast(lat, lon float64) ([]byte, error) {
	cacheKey := fmt.Sprintf("%f,%f", lat, lon)
	if cached, found := forecastCache.Get(cacheKey); found {
		return cached, nil
	}

	baseURL := "https://api.pirateweather.net/forecast"
	apiKey := appConfig.pirateWeatherKey
	params := url.Values{}
	params.Add("extend", "hourly")
	url := fmt.Sprintf("%s/%s/%f,%f?%s", baseURL, apiKey, lat, lon, params.Encode())

	resp, err := forecastHTTPClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching weather data: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading weather data: %w", err)
	}

	forecastCache.Set(cacheKey, body, forecastCacheTTL)
	return body, nil
}

type forecastResult struct {
	body []byte
	err  error
}

type reverseGeocodeResult struct {
	result *GeocodeResult
	err    error
}

func ForecastHandler(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()

	lat, err := strconv.ParseFloat(params.Get("lat"), 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid latitude: %v", err), http.StatusBadRequest)
		return
	}
	lng, err := strconv.ParseFloat(params.Get("lng"), 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid longitude: %v", err), http.StatusBadRequest)
		return
	}

	forecastCh := make(chan forecastResult, 1)
	go func() {
		forecastBytes, err := getForecast(lat, lng)
		forecastCh <- forecastResult{
			body: forecastBytes,
			err:  err,
		}
	}()

	geocodeCh := make(chan reverseGeocodeResult, 1)
	go func() {
		geocodeResp, err := reverseGeocode(lat, lng)
		geocodeCh <- reverseGeocodeResult{
			result: geocodeResp,
			err:    err,
		}
	}()

	forecastResp := <-forecastCh
	if forecastResp.err != nil {
		http.Error(w, fmt.Sprintf("Error fetching weather data: %v", forecastResp.err), http.StatusInternalServerError)
		return
	}

	geocodeResp := <-geocodeCh
	if geocodeResp.err != nil {
		http.Error(w, fmt.Sprintf("%v", geocodeResp.err), http.StatusBadRequest)
		return
	}

	var forecastMap map[string]interface{}
	err = json.Unmarshal(forecastResp.body, &forecastMap)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing weather data: %v", err), http.StatusInternalServerError)
		return
	}

	forecastMap["formatted_address"] = geocodeResp.result.FormattedAddress

	responseBytes, err := json.Marshal(forecastMap)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error generating response: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(responseBytes)
}
