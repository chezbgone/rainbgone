package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/joho/godotenv"
)

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

var GOOGLE_MAPS_API_KEY string

func init() {
	env, err := godotenv.Read()
	if err != nil {
		panic(err)
	}
	GOOGLE_MAPS_API_KEY = env["GOOGLE_MAPS_API_KEY"]
}

func geocode(address string) (*GeocodeResponse, error) {
	baseURL := "https://maps.googleapis.com/maps/api/geocode/json"
	params := url.Values{}
	params.Add("key", GOOGLE_MAPS_API_KEY)
	params.Add("address", address)

	url := fmt.Sprintf("%s?%s", baseURL, params.Encode())
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error geocoding address: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading geocoding response: %w", err)
	}

	var geocodeResp GeocodeResponse
	err = json.Unmarshal(body, &geocodeResp)
	if err != nil {
		return nil, fmt.Errorf("error parsing geocoding response: %w", err)
	}

	if geocodeResp.Status != "OK" {
		return nil, fmt.Errorf("invalid geocoding response: status=%s, results=%d", geocodeResp.Status, len(geocodeResp.Results))
	}

	return &geocodeResp, nil
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
