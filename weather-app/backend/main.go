package main

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const defaultWeatherAPIURL = "https://api.openweathermap.org/data/2.5/weather"

type openWeatherResponse struct {
	Name string `json:"name"`
	Main struct {
		Temp     float64 `json:"temp"`
		Humidity int     `json:"humidity"`
	} `json:"main"`
	Weather []struct {
		Description string `json:"description"`
	} `json:"weather"`
	Wind struct {
		Speed float64 `json:"speed"`
	} `json:"wind"`
}

type weatherResponse struct {
	City            string  `json:"city"`
	TemperatureC    float64 `json:"temperature_c"`
	Condition       string  `json:"condition"`
	HumidityPercent int     `json:"humidity_percent"`
	WindSpeedMPS    float64 `json:"wind_speed_mps"`
}

type server struct {
	apiKey    string
	apiURL    string
	client    *http.Client
	transport http.RoundTripper
}

func newServer() *server {
	apiURL := strings.TrimSpace(os.Getenv("OPENWEATHER_API_URL"))
	if apiURL == "" {
		apiURL = defaultWeatherAPIURL
	}

	return &server{
		apiKey: strings.TrimSpace(os.Getenv("OPENWEATHER_API_KEY")),
		apiURL: apiURL,
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (s *server) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})
	mux.HandleFunc("GET /api/weather", s.weatherHandler)

	allowedOrigin := strings.TrimSpace(os.Getenv("CORS_ALLOWED_ORIGIN"))
	if allowedOrigin == "" {
		allowedOrigin = "*"
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		mux.ServeHTTP(w, r)
	})
}

func (s *server) weatherHandler(w http.ResponseWriter, r *http.Request) {
	city := strings.TrimSpace(r.URL.Query().Get("city"))
	if city == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "city query parameter is required"})
		return
	}
	if s.apiKey == "" {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "weather provider is not configured"})
		return
	}

	requestURL, err := url.Parse(s.apiURL)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "weather provider URL is invalid"})
		return
	}
	query := requestURL.Query()
	query.Set("q", city)
	query.Set("appid", s.apiKey)
	query.Set("units", "metric")
	requestURL.RawQuery = query.Encode()

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL.String(), nil)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not create provider request"})
		return
	}
	request.Header.Set("Accept", "application/json")

	client := s.client
	if client == nil {
		client = &http.Client{Timeout: 10 * time.Second}
	}
	if s.transport != nil {
		client = &http.Client{Timeout: 10 * time.Second, Transport: s.transport}
	}

	response, err := client.Do(request)
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]string{"error": "weather provider request failed"})
		return
	}
	defer response.Body.Close()
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		writeJSON(w, http.StatusBadGateway, map[string]string{"error": "weather provider returned an unsuccessful response"})
		return
	}

	var providerData openWeatherResponse
	if err := json.NewDecoder(io.LimitReader(response.Body, 1<<20)).Decode(&providerData); err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]string{"error": "weather provider returned invalid data"})
		return
	}

	condition := "Unknown"
	if len(providerData.Weather) > 0 {
		condition = providerData.Weather[0].Description
	}

	writeJSON(w, http.StatusOK, weatherResponse{
		City:            providerData.Name,
		TemperatureC:    providerData.Main.Temp,
		Condition:       condition,
		HumidityPercent: providerData.Main.Humidity,
		WindSpeedMPS:    providerData.Wind.Speed,
	})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil && !errors.Is(err, context.Canceled) {
		log.Printf("write JSON response: %v", err)
	}
}

func main() {
	port := strings.TrimSpace(os.Getenv("PORT"))
	if port == "" {
		port = "8080"
	}

	app := newServer()
	address := ":" + port
	log.Printf("weather API listening on %s", address)
	log.Fatal(http.ListenAndServe(address, app.routes()))
}
