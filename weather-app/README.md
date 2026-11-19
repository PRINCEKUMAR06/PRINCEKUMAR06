# City Weather - Go API + React client

A full-stack weather application. The React client accepts a city name and calls a Go HTTP API. The API retrieves current conditions from OpenWeatherMap and returns a compact JSON response.

## Architecture
```text
Browser (React)
  | GET /api/weather?city=Hyderabad
  v
Go HTTP API (:8080)
  | GET /data/2.5/weather?q=...&units=metric
  v
OpenWeatherMap
```

The provider API key stays on the server. Never place it in React environment variables or browser code.

## Layout
- `backend/` - Go HTTP server, provider integration, Dockerfile and handler tests
- `frontend/` - React UI, Axios API client, CSS and Dockerfile
- `docker-compose.yml` - local two-container development stack
- `../.github/workflows/weather-backend-tests.yml` - root-level Go test workflow

## API
- `GET /healthz` returns `{"status":"ok"}`
- `GET /api/weather?city=Hyderabad` returns city, temperature in Celsius, condition, humidity and wind speed.

Example response (illustrative schema, not a live reading):
```json
{
  "city": "Hyderabad",
  "temperature_c": 31.5,
  "condition": "clear sky",
  "humidity_percent": 42,
  "wind_speed_mps": 3.2
}
```

## Run locally
Requirements: Go 1.23+, Node.js 20+, npm and an OpenWeatherMap API key.

Start the API:
```bash
cd backend
export OPENWEATHER_API_KEY="your_api_key"
go test ./... -race -cover
go run .
```
The API listens on `http://localhost:8080`. Optional variables: `PORT`, `OPENWEATHER_API_URL` (useful for tests), and `CORS_ALLOWED_ORIGIN`.

Start the React client in another terminal:
```bash
cd frontend
cp .env.example .env
npm install
npm start
```
The client uses `REACT_APP_API_BASE_URL`, defaulting to `http://localhost:8080`.

## Docker Compose
From this directory, copy `.env.example` to `.env`, add the provider key, then:
```bash
docker compose up --build
```
Open `http://localhost:3000`; API health is at `http://localhost:8080/healthz`.

## Testing
Go handler tests cover required input, missing provider configuration, response mapping and provider failures. They use `httptest`, so no external API key is required. GitHub Actions runs the Go test suite on relevant pushes and pull requests.

## Deployment
Deploy API and frontend as separate services. Configure `OPENWEATHER_API_KEY` and `CORS_ALLOWED_ORIGIN` on the API service. Set `REACT_APP_API_BASE_URL` to the public API URL at frontend build time. Enable HTTPS and do not commit `.env` or provider credentials.

This commit provides source and deployment configuration; it does not claim a live hosting deployment.

Weather data is provided by [OpenWeatherMap](https://openweathermap.org/api). This educational/demo application is not intended for safety-critical decisions.
