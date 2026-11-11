package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

const testForecast = `{
	"timezone":"America/Los_Angeles",
	"currently":{"time":1794351600,"icon":"rain","summary":"Rain","temperature":50,"apparentTemperature":48,"windSpeed":5,"humidity":0.8,"dewPoint":44,"uvIndex":1,"visibility":10,"pressure":1010},
	"minutely":{"summary":"Rain","data":[]},
	"hourly":{"summary":"Rain","data":[{"time":1794351600,"temperature":50,"precipIntensity":0.1,"precipType":"rain","cloudCover":1}]},
	"daily":{"summary":"Rain","data":[{"time":1794297600,"icon":"rain","summary":"Rain","temperatureMin":40,"temperatureMax":55,"temperatureLow":40,"temperatureHigh":55,"temperatureMinTime":1794297600,"temperatureMaxTime":1794351600,"sunriseTime":1794322800,"sunsetTime":1794358800,"precipIntensity":0.1}]}
}`

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func testResponse(status int, body string) *http.Response {
	return &http.Response{StatusCode: status, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(body))}
}

type closingBody struct {
	io.ReadCloser
	done chan struct{}
}

func (b closingBody) Close() error {
	err := b.ReadCloser.Close()
	close(b.done)
	return err
}

// Tests never run in parallel: they replace package-level clients and caches. Every
// provider request is intercepted by a transport, including unexpected requests.
func setupForecastTest(t *testing.T, transport roundTripFunc, geocodeTransport ...roundTripFunc) <-chan struct{} {
	t.Helper()
	oldConfig, oldForecastClient, oldGeocodeClient := appConfig, forecastHTTPClient, geocodeHTTPClient
	oldForecastCache, oldTimeMachineCache, oldGeocodeCache := forecastCache, timeMachineCache, reverseGeocodeCache
	appConfig = config{pirateWeatherKey: "dummy-test-key"}
	forecastCache, timeMachineCache, reverseGeocodeCache = newCache(), newCache(), newCache()
	forecastHTTPClient = &http.Client{Transport: transport}
	geocodeDone := make(chan struct{})
	geocodeHTTPClient = &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		resp := testResponse(http.StatusOK, `{"lat":"47","lon":"-122","display_name":"Test location"}`)
		if len(geocodeTransport) > 0 {
			var err error
			resp, err = geocodeTransport[0](r)
			if err != nil {
				close(geocodeDone)
				return nil, err
			}
		}
		resp.Body = closingBody{ReadCloser: resp.Body, done: geocodeDone}
		return resp, nil
	})}
	t.Cleanup(func() {
		// A failed forecast may return before the parallel geocode completes.
		select {
		case <-geocodeDone:
		case <-time.After(2 * time.Second):
			t.Error("reverse geocode did not complete")
		}
		appConfig, forecastHTTPClient, geocodeHTTPClient = oldConfig, oldForecastClient, oldGeocodeClient
		forecastCache, timeMachineCache, reverseGeocodeCache = oldForecastCache, oldTimeMachineCache, oldGeocodeCache
	})
	return geocodeDone
}

func TestForecastNetworkErrorDoesNotExposeAPIKey(t *testing.T) {
	setupForecastTest(t, func(r *http.Request) (*http.Response, error) {
		return nil, errors.New("synthetic connection failure")
	})
	w := httptest.NewRecorder()
	ForecastHandler(w, httptest.NewRequest(http.MethodGet, "/forecast?lat=47&lng=-122", nil))
	if w.Code < 500 {
		t.Fatalf("expected upstream failure, got %d", w.Code)
	}
	if strings.Contains(w.Body.String(), "dummy-test-key") || strings.Contains(w.Body.String(), "pirateweather.net") {
		t.Fatalf("public error exposes the provider URL: %s", w.Body.String())
	}
}

func TestTimeMachineNetworkErrorDoesNotLogAPIKey(t *testing.T) {
	setupForecastTest(t, func(r *http.Request) (*http.Response, error) {
		if r.URL.Host == "api.pirateweather.net" {
			return testResponse(http.StatusOK, testForecast), nil
		}
		return nil, errors.New("synthetic connection failure")
	})
	var logs bytes.Buffer
	oldOutput := log.Writer()
	log.SetOutput(&logs)
	t.Cleanup(func() { log.SetOutput(oldOutput) })
	w := httptest.NewRecorder()
	ForecastHandler(w, httptest.NewRequest(http.MethodGet, "/forecast?lat=47&lng=-122", nil))
	if w.Code != http.StatusOK {
		t.Fatalf("best-effort backfill should not block forecast: %d", w.Code)
	}
	if strings.Contains(logs.String(), "dummy-test-key") {
		t.Fatal("time machine error logged the API key")
	}
}

func TestForecastRejectsUpstreamErrorsWithoutCaching(t *testing.T) {
	for _, status := range []int{http.StatusUnauthorized, http.StatusTooManyRequests, http.StatusServiceUnavailable} {
		t.Run(http.StatusText(status), func(t *testing.T) {
			calls := 0
			geocodeDone := setupForecastTest(t, func(r *http.Request) (*http.Response, error) {
				if r.URL.Host != "api.pirateweather.net" {
					return testResponse(http.StatusServiceUnavailable, ""), nil
				}
				calls++
				if calls == 1 {
					return testResponse(status, `{"message":"Forbidden"}`), nil
				}
				return testResponse(http.StatusOK, testForecast), nil
			})
			for i, wantStatus := range []int{http.StatusBadGateway, http.StatusOK} {
				w := httptest.NewRecorder()
				ForecastHandler(w, httptest.NewRequest(http.MethodGet, "/forecast?lat=47&lng=-122", nil))
				<-geocodeDone
				if w.Code != wantStatus {
					t.Errorf("request %d: expected %d, got %d", i, wantStatus, w.Code)
				}
			}
			if calls != 2 {
				t.Errorf("expected retry after failed upstream response, got %d calls", calls)
			}
		})
	}
}

func TestForecastRejectsMalformedSuccessWithoutCaching(t *testing.T) {
	var valid map[string]interface{}
	if err := json.Unmarshal([]byte(testForecast), &valid); err != nil {
		t.Fatal(err)
	}
	bodies := []string{"not JSON", "null", "[]", `{"message":"Forbidden"}`}
	for _, field := range []string{"currently", "minutely", "hourly", "daily", "timezone"} {
		copy := make(map[string]interface{})
		for key, value := range valid {
			copy[key] = value
		}
		delete(copy, field)
		body, _ := json.Marshal(copy)
		bodies = append(bodies, string(body))
	}
	for i, body := range bodies {
		t.Run(fmt.Sprintf("body-%d", i), func(t *testing.T) {
			calls := 0
			geocodeDone := setupForecastTest(t, func(r *http.Request) (*http.Response, error) {
				if r.URL.Host != "api.pirateweather.net" {
					return testResponse(http.StatusServiceUnavailable, ""), nil
				}
				calls++
				if calls == 1 {
					return testResponse(http.StatusOK, body), nil
				}
				return testResponse(http.StatusOK, testForecast), nil
			})
			for _, wantStatus := range []int{http.StatusBadGateway, http.StatusOK} {
				w := httptest.NewRecorder()
				ForecastHandler(w, httptest.NewRequest(http.MethodGet, "/forecast?lat=47&lng=-122", nil))
				<-geocodeDone
				if w.Code != wantStatus {
					t.Errorf("expected %d, got %d", wantStatus, w.Code)
				}
			}
			if calls != 2 {
				t.Fatal("invalid response was cached")
			}
		})
	}
}

func TestValidForecastIsCached(t *testing.T) {
	calls := 0
	setupForecastTest(t, func(r *http.Request) (*http.Response, error) {
		if r.URL.Host != "api.pirateweather.net" {
			return testResponse(http.StatusServiceUnavailable, ""), nil
		}
		calls++
		return testResponse(http.StatusOK, testForecast), nil
	})
	for range 2 {
		w := httptest.NewRecorder()
		ForecastHandler(w, httptest.NewRequest(http.MethodGet, "/forecast?lat=47&lng=-122", nil))
		if w.Code != http.StatusOK {
			t.Fatalf("expected successful forecast, got %d", w.Code)
		}
	}
	if calls != 1 {
		t.Fatalf("expected one provider call, got %d", calls)
	}
}

func TestForecastSurvivesReverseGeocodeFailures(t *testing.T) {
	for _, name := range []string{"no result", "provider error", "network failure", "empty address"} {
		t.Run(name, func(t *testing.T) {
			setupForecastTest(t, func(r *http.Request) (*http.Response, error) {
				if r.URL.Host == "api.pirateweather.net" {
					return testResponse(http.StatusOK, testForecast), nil
				}
				return testResponse(http.StatusServiceUnavailable, ""), nil
			}, func(r *http.Request) (*http.Response, error) {
				switch name {
				case "network failure":
					return nil, errors.New("synthetic geocoding outage")
				case "provider error":
					return testResponse(http.StatusServiceUnavailable, "Unavailable"), nil
				case "empty address":
					return testResponse(http.StatusOK, `{"lat":"47","lon":"-122","display_name":""}`), nil
				default:
					return testResponse(http.StatusOK, `{"error":"Unable to geocode"}`), nil
				}
			})
			w := httptest.NewRecorder()
			ForecastHandler(w, httptest.NewRequest(http.MethodGet, "/forecast?lat=47&lng=-122", nil))
			if w.Code != http.StatusOK {
				t.Fatalf("reverse geocode blocked weather: %d", w.Code)
			}
			var forecast map[string]interface{}
			if err := json.Unmarshal(w.Body.Bytes(), &forecast); err != nil {
				t.Fatal(err)
			}
			if forecast["formatted_address"] != "47.000000, -122.000000" || forecast["currently"] == nil {
				t.Fatal("expected intact forecast with coordinate fallback")
			}
		})
	}
}

func TestForecastUsesReverseGeocodedAddress(t *testing.T) {
	setupForecastTest(t, func(r *http.Request) (*http.Response, error) {
		if r.URL.Host == "api.pirateweather.net" {
			return testResponse(http.StatusOK, testForecast), nil
		}
		return testResponse(http.StatusServiceUnavailable, ""), nil
	})
	w := httptest.NewRecorder()
	ForecastHandler(w, httptest.NewRequest(http.MethodGet, "/forecast?lat=47&lng=-122", nil))
	var forecast map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &forecast); err != nil {
		t.Fatal(err)
	}
	if forecast["formatted_address"] != "Test location" {
		t.Fatal("successful reverse geocoding was not used")
	}
}
