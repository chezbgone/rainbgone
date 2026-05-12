package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

var PIRATE_WEATHER_KEY string

var forecastHTTPClient = &http.Client{
	Timeout: 10 * time.Second,
}

func init() {
	env, err := godotenv.Read()
	if err != nil {
		panic(err)
	}
	PIRATE_WEATHER_KEY = env["PIRATE_WEATHER_KEY"]
}

func getForecast(lat, lon float64) ([]byte, error) {
	baseURL := "https://api.pirateweather.net/forecast"
	apiKey := PIRATE_WEATHER_KEY
	params := url.Values{}
	params.Add("extend", "hourly")
	url := fmt.Sprintf("%s/%s/%f,%f?%s", baseURL, apiKey, lat, lon, params.Encode())

	resp, err := forecastHTTPClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching weather data: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading weather data: %w", err)
	}

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
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid latitude: %v", err)
		return
	}
	lng, err := strconv.ParseFloat(params.Get("lng"), 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid longitude: %v", err)
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
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error fetching weather data: %v", forecastResp.err)
		return
	}

	geocodeResp := <-geocodeCh
	if geocodeResp.err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "%v", geocodeResp.err)
		return
	}

	var forecastMap map[string]interface{}
	err = json.Unmarshal(forecastResp.body, &forecastMap)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error parsing weather data: %v", err)
		return
	}

	forecastMap["formatted_address"] = geocodeResp.result.FormattedAddress

	responseBytes, err := json.Marshal(forecastMap)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error generating response: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseBytes)
}
