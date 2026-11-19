package main

import (
  "net/http"
  "net/http/httptest"
  "strings"
  "testing"
)

func TestWeatherHandlerRequiresCity(t *testing.T) {
  app := &server{apiKey:"test-key", apiURL:"http://example.test/weather"}
  req := httptest.NewRequest(http.MethodGet, "/api/weather", nil)
  res := httptest.NewRecorder()
  app.weatherHandler(res, req)
  if res.Code != http.StatusBadRequest { t.Fatalf("expected 400, got %d", res.Code) }
}

func TestWeatherHandlerRequiresAPIKey(t *testing.T) {
  app := &server{apiURL:"http://example.test/weather"}
  req := httptest.NewRequest(http.MethodGet, "/api/weather?city=Hyderabad", nil)
  res := httptest.NewRecorder()
  app.weatherHandler(res, req)
  if res.Code != http.StatusServiceUnavailable { t.Fatalf("expected 503, got %d", res.Code) }
}

func TestWeatherHandlerMapsProviderResponse(t *testing.T) {
  provider := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    if r.URL.Query().Get("q") != "Hyderabad" { t.Fatalf("unexpected city query") }
    if r.URL.Query().Get("units") != "metric" { t.Fatalf("expected metric units") }
    if r.URL.Query().Get("appid") != "test-key" { t.Fatalf("expected provider API key") }
    w.Header().Set("Content-Type", "application/json")
    _, _ = w.Write([]byte(`{"name":"Hyderabad","main":{"temp":31.5,"humidity":42},"weather":[{"description":"clear sky"}],"wind":{"speed":3.2}}`))
  }))
  defer provider.Close()
  app := &server{apiKey:"test-key", apiURL:provider.URL, client:provider.Client()}
  req := httptest.NewRequest(http.MethodGet, "/api/weather?city=Hyderabad", nil)
  res := httptest.NewRecorder()
  app.weatherHandler(res, req)
  if res.Code != http.StatusOK { t.Fatalf("expected 200, got %d: %s", res.Code, res.Body.String()) }
  for _, want := range []string{`"city":"Hyderabad"`, `"temperature_c":31.5`, `"condition":"clear sky"`, `"humidity_percent":42`, `"wind_speed_mps":3.2`} {
    if !strings.Contains(res.Body.String(), want) { t.Errorf("missing %s in %s", want, res.Body.String()) }
  }
}

func TestWeatherHandlerHandlesProviderFailure(t *testing.T) {
  provider := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
    http.Error(w, "upstream unavailable", http.StatusServiceUnavailable)
  }))
  defer provider.Close()
  app := &server{apiKey:"test-key", apiURL:provider.URL, client:provider.Client()}
  req := httptest.NewRequest(http.MethodGet, "/api/weather?city=Hyderabad", nil)
  res := httptest.NewRecorder()
  app.weatherHandler(res, req)
  if res.Code != http.StatusBadGateway { t.Fatalf("expected 502, got %d", res.Code) }
}
