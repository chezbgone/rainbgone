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

var geocodeHTTPClient = &http.Client{
	Timeout: 5 * time.Second,
}

type Location struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type Bounds struct {
	Northeast Location `json:"northeast"`
	Southwest Location `json:"southwest"`
}

type Geometry struct {
	Location     Location `json:"location"`
	LocationType string   `json:"location_type"`
	Viewport     Bounds   `json:"viewport"`
}

type GeocodeResult struct {
	FormattedAddress string   `json:"formatted_address"`
	Geometry         Geometry `json:"geometry"`
}

type GeocodeResponse struct {
	Status  string          `json:"status"`
	Results []GeocodeResult `json:"results"`
}

// NominatimSearchResult represents a result from Nominatim search endpoint
type NominatimSearchResult struct {
	Lat         string `json:"lat"`
	Lon         string `json:"lon"`
	DisplayName string `json:"display_name"`
}

// NominatimReverseResult represents a result from Nominatim reverse endpoint
type NominatimReverseResult struct {
	Lat         string `json:"lat"`
	Lon         string `json:"lon"`
	DisplayName string `json:"display_name"`
}

func geocode(address string) (*GeocodeResponse, error) {
	baseURL := "https://nominatim.openstreetmap.org/search"
	params := url.Values{}
	params.Add("q", address)
	params.Add("format", "json")
	params.Add("limit", "1")

	urlStr := fmt.Sprintf("%s?%s", baseURL, params.Encode())
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("User-Agent", "rainbgone/1.0")

	resp, err := geocodeHTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error geocoding address: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading geocoding response: %w", err)
	}

	var nominatimResults []NominatimSearchResult
	err = json.Unmarshal(body, &nominatimResults)
	if err != nil {
		return nil, fmt.Errorf("error parsing geocoding response: %w", err)
	}

	if len(nominatimResults) == 0 {
		return nil, fmt.Errorf("no results found for address: %s", address)
	}

	nomResult := nominatimResults[0]
	lat, err := strconv.ParseFloat(nomResult.Lat, 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing lat: %w", err)
	}
	lng, err := strconv.ParseFloat(nomResult.Lon, 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing lon: %w", err)
	}

	geocodeResult := GeocodeResult{
		FormattedAddress: nomResult.DisplayName,
		Geometry: Geometry{
			Location: Location{
				Lat: lat,
				Lng: lng,
			},
		},
	}

	return &GeocodeResponse{
		Status:  "OK",
		Results: []GeocodeResult{geocodeResult},
	}, nil
}

func reverseGeocode(lat, lng float64) (*GeocodeResult, error) {
	baseURL := "https://nominatim.openstreetmap.org/reverse"
	params := url.Values{}
	params.Add("lat", fmt.Sprintf("%f", lat))
	params.Add("lon", fmt.Sprintf("%f", lng))
	params.Add("format", "json")

	urlStr := fmt.Sprintf("%s?%s", baseURL, params.Encode())
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("User-Agent", "rainbgone/1.0")

	resp, err := geocodeHTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error reverse geocoding: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading reverse geocoding response: %w", err)
	}

	var nomResult NominatimReverseResult
	err = json.Unmarshal(body, &nomResult)
	if err != nil {
		return nil, fmt.Errorf("error parsing reverse geocoding response: %w", err)
	}

	latParsed, err := strconv.ParseFloat(nomResult.Lat, 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing lat: %w", err)
	}
	lngParsed, err := strconv.ParseFloat(nomResult.Lon, 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing lon: %w", err)
	}

	return &GeocodeResult{
		FormattedAddress: nomResult.DisplayName,
		Geometry: Geometry{
			Location: Location{
				Lat: latParsed,
				Lng: lngParsed,
			},
		},
	}, nil
}

func geocode_one(address string) (*GeocodeResult, error) {
	response, err := geocode(address)
	if err != nil {
		return nil, err
	}
	if len(response.Results) == 0 {
		return nil, fmt.Errorf("no results found for address: %s", address)
	}
	return &response.Results[0], nil
}

func GeocodeHandler(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	address := params.Get("address")
	geocodeResp, err := geocode_one(address)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "%v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(geocodeResp)
}
