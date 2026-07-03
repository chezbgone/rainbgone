package server

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

var forecastHTTPClient = &http.Client{
	Timeout: 10 * time.Second,
}

var forecastCache = newCache()
var timeMachineCache = newCache()

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

// getTimeMachine fetches the Time Machine forecast anchored at `ts`. Unlike the regular
// forecast (whose hourly series starts at the current hour), this returns a clean
// midnight-to-midnight block for the location's local day — used to backfill the
// hours of today that have already elapsed.
//
// The timestamp matters more than the doc suggests: the upstream API is inconsistent
// when asked about "now" — depending on which backend instance answers, it sometimes
// returns the historical local day (what we want) and sometimes a forward-looking
// forecast starting at the current hour (useless for backfill, and silently discarded
// by mergeHourlyFromMidnight's cutoff filter). Passing a timestamp anchored well inside
// the elapsed part of the local day (the caller uses local midnight) got a 6/6 historical
// response in testing, vs. 5/6 for a raw `time.Now()`. See ForecastHandler for how `ts`
// is derived from the location's timezone.
func getTimeMachine(lat, lon float64, ts int64) ([]byte, error) {
	cacheKey := fmt.Sprintf("%f,%f", lat, lon)
	if cached, found := timeMachineCache.Get(cacheKey); found {
		return cached, nil
	}

	baseURL := "https://timemachine.pirateweather.net/forecast"
	apiKey := appConfig.pirateWeatherKey
	params := url.Values{}
	params.Add("exclude", "minutely,alerts")
	requestURL := fmt.Sprintf("%s/%s/%f,%f,%d?%s", baseURL, apiKey, lat, lon, ts, params.Encode())

	resp, err := forecastHTTPClient.Get(requestURL)
	if err != nil {
		return nil, fmt.Errorf("error fetching time machine data: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("time machine returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading time machine data: %w", err)
	}

	timeMachineCache.Set(cacheKey, body, forecastCacheTTL)
	return body, nil
}

// hourlyData digs out the hourly.data array from an unmarshaled forecast/time
// machine response, tolerating any missing/malformed structure by returning nil.
func hourlyData(m map[string]interface{}) []interface{} {
	hourly, ok := m["hourly"].(map[string]interface{})
	if !ok {
		return nil
	}
	data, ok := hourly["data"].([]interface{})
	if !ok {
		return nil
	}
	return data
}

// hourTime extracts the Unix `time` field from an hourly entry.
func hourTime(h interface{}) (float64, bool) {
	obj, ok := h.(map[string]interface{})
	if !ok {
		return 0, false
	}
	t, ok := obj["time"].(float64)
	return t, ok
}

// mergeHourlyFromMidnight builds a single flat hourly series anchored at today's
// local midnight, by prepending today's already-elapsed hours from the Time Machine
// body onto the regular forecast (which starts at the current hour). The seam is a
// plain timestamp threshold — Time Machine hours strictly before the forecast's first
// (current) hour — so only today's past morning comes from the GFS-only Time Machine;
// every hour the standard multi-model forecast covers is kept from it. The result is
// chronological and contiguous. If timeMachineBody is nil/unparseable, the series
// simply starts at the current hour (today's morning absent), matching pre-Time-Machine
// behavior.
func mergeHourlyFromMidnight(forecastMap map[string]interface{}, timeMachineBody []byte) []interface{} {
	fcHours := hourlyData(forecastMap)

	// cutoff = the standard forecast's first (current) hour.
	var cutoff float64
	hasCutoff := false
	if len(fcHours) > 0 {
		if t, ok := hourTime(fcHours[0]); ok {
			cutoff = t
			hasCutoff = true
		}
	}

	result := make([]interface{}, 0, len(fcHours)+24)
	backfilled := 0

	if hasCutoff && timeMachineBody != nil {
		var tm map[string]interface{}
		if err := json.Unmarshal(timeMachineBody, &tm); err != nil {
			log.Printf("time machine body unparseable, today's hourly will start at the current hour: %v", err)
		} else {
			for _, h := range hourlyData(tm) {
				if t, ok := hourTime(h); ok && t < cutoff {
					result = append(result, h)
					backfilled++
				}
			}
			// A non-nil, parseable body that still backfills nothing is the known failure
			// mode where Time Machine answers "now" with a forward-looking forecast instead
			// of the historical local day (see getTimeMachine) — surface it instead of
			// silently falling back, since it used to fail invisibly.
			if backfilled == 0 {
				log.Printf("time machine returned no hours before the forecast cutoff, today's hourly will start at the current hour")
			}
		}
	}

	result = append(result, fcHours...)
	return result
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

	// The Time Machine call runs after (not concurrently with) the forecast fetch because
	// it needs the location's timezone, which only the forecast response carries. It's
	// anchored to local midnight rather than time.Now() — see getTimeMachine for why: the
	// upstream API is unreliable about "now" and reliably returns the historical local day
	// when asked about a timestamp well inside the elapsed part of the day. It's still
	// best-effort: on failure, the midnight-anchored series simply starts at the current
	// hour (today's elapsed morning absent).
	tz, _ := forecastMap["timezone"].(string)
	loc, err := time.LoadLocation(tz)
	if err != nil {
		loc = time.UTC
	}
	now := time.Now().In(loc)
	midnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc).Unix()

	timeMachineBytes, err := getTimeMachine(lat, lng, midnight)
	if err != nil {
		log.Printf("time machine fetch failed, today's hourly will start at the current hour: %v", err)
	}
	forecastMap["hourlyFromMidnight"] = mergeHourlyFromMidnight(forecastMap, timeMachineBytes)

	responseBytes, err := json.Marshal(forecastMap)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error generating response: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(responseBytes)
}
