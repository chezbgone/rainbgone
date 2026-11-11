package server

import (
	"bytes"
	"errors"
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
func setupForecastTest(t *testing.T, transport roundTripFunc) {
	t.Helper()
	oldConfig, oldForecastClient, oldGeocodeClient := appConfig, forecastHTTPClient, geocodeHTTPClient
	oldForecastCache, oldTimeMachineCache, oldGeocodeCache := forecastCache, timeMachineCache, reverseGeocodeCache
	appConfig = config{pirateWeatherKey: "dummy-test-key"}
	forecastCache, timeMachineCache, reverseGeocodeCache = newCache(), newCache(), newCache()
	forecastHTTPClient = &http.Client{Transport: transport}
	geocodeDone := make(chan struct{})
	geocodeHTTPClient = &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		resp := testResponse(http.StatusOK, `{"lat":"47","lon":"-122","display_name":"Test location"}`)
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
