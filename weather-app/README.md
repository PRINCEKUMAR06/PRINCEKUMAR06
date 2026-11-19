# City Weather — Go API + React client

A small full-stack weather application. The React client accepts a city name and calls a Go HTTP API. The API retrieves current conditions from OpenWeatherMap and returns a compact JSON response for the UI.

## Architecture

```text
Browser (React)
    |  GET /api/weather?city=Hyderabad
    v
Go HTTP API (:8080)
    |  GET /data/2.5/weather?q=...&units=metric
    v
OpenWeatherMap
```

The backend keeps the provider API key on the server. The frontend receives only the app API base URL; never put the OpenWeatherMap key in React environment variables or browser code.

## Project layout

- `backend/` — Go HTTP server, provider integration, Dockerfile and handler tests
- `frontend/` — React UI, Axios API client, CSS and Dockerfile
- `docker-compose.yml` — local two-container development stack
- `.github/workflows/backend-tests.yml` — Go test workflow

## API

### `GET /healthz`

Returns `{"status":"ok"}` when the API process is responding.

### `GET /api/weather?city=Hyderabad`

Example success response:

```json
{
  "city": "Hyderabad",
  "temperature_c": 31.5,
  "condition": "clear sky",
  "humidity_percent": 42,
  "wind_speed_mps": 3.2
}
```

The numbers above are an illustrative response schema, not a live weather reading.

## Run locally

Requirements: Go 1.23+; Node.js 20+ and npm; an OpenWeatherMap API key.

### 1. Start the API

```bash
cd backend
export OPENWEATHER_API_KEY="your_api_key"
go test ./...
go run .
```

The API listens on `http://localhost:8080` by default. Optional environment variables: `PORT`, `OPENWEATHER_API_URL` (useful for local testing), and `CORS_ALLOWED_ORIGIN`.

### 2. Start the React client

In a second terminal:

```bash
cd frontend
cp .env.example .env
npm install
npm start
```

The client uses `REACT_APP_API_BASE_URL`, defaulting to `http://localhost:8080`.

## Run with Docker Compose

From this directory, create `.env` from `.env.example`, add your provider key, and run:

```bash
docker compose up --build
```

Open `http://localhost:3000`. The API health endpoint is at `http://localhost:8080/healthz`.

## Tests and validation

The Go handler tests cover required input, missing provider configuration, response mapping and provider failure handling:

```bash
cd backend
go test ./... -race -cover
```

The repository includes a GitHub Actions workflow for the backend test suite. No external API key is needed by the tests; the provider is simulated with `httptest`.

## Deployment notes

Deploy the API and frontend as separate services. Configure `OPENWEATHER_API_KEY` and `CORS_ALLOWED_ORIGIN` on the API service. Set `REACT_APP_API_BASE_URL` to the publicly reachable API base URL at frontend build time. HTTPS should be enabled for both services. Do not commit `.env` files or provider credentials.

This repository contains the source and deployment configuration. A live hosting deployment is not claimed by this commit.

## Data source and attribution

Weather data is retrieved from [OpenWeatherMap](https://openweathermap.org/api). The project is an educational/demo application and is not intended for emergency or safety-critical decisions.
