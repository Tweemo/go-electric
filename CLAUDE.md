# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Repository Layout

This is a monorepo with two services plus root-level deployment config:

- `backend/` - Go API (Gin). Module `github.com/tweemo/go-electric`; run commands from here.
- `frontend/` - Next.js + TypeScript UI.
- `docker-compose.yml` - Brings up both services together.

## Quick Start (Backend)

Run these from `backend/`:

**Build the project:**
```bash
cd backend && go build -o go-electric .
```

**Run the server:**
```bash
cd backend && go run .
```

**Environment variables** (see `backend/.env.example`):
- `PORT` - Server port (default: 8080)
- `HOST` - Server host (default: localhost)
- `ENV` - Set to "production" to skip loading `.env`
- `GIN_MODE` - "debug" or "release"
- `CORS_ALLOWED_ORIGINS` - Comma-separated origins (dev default: localhost:3001)

## Frontend (`frontend/`)

A Next.js (App Router) + TypeScript UI that uploads a usage CSV to `POST /costs` and
displays the estimated costs, grouped by provider with the cheapest plan highlighted.
Built with Tailwind v4 + shadcn/ui (base-ui). NOTE: `frontend/` pins a newer Next.js
than may be in training data — see `frontend/AGENTS.md`.

**Run the UI (dev):**
```bash
cd frontend && npm install && npm run dev   # serves on http://localhost:3001
```

The dev server runs on **port 3001** so the Go CORS default (`localhost:3001`) works
with no backend change. Run the Go API (`cd backend && go run .`) alongside it.

## Docker

```bash
docker compose up --build   # backend on :8080, frontend on :3001
```

Each service has its own `Dockerfile` (`backend/Dockerfile`, `frontend/Dockerfile`);
`docker-compose.yml` wires them. The browser calls the API directly, so
`NEXT_PUBLIC_API_URL` is baked at frontend build time (`http://localhost:8080`) and the
backend's `CORS_ALLOWED_ORIGINS` is set to `http://localhost:3001` in compose. The
frontend image uses Next.js `output: "standalone"` for a slim runtime.

**Frontend env** (`frontend/.env.local`, see `.env.example`):
- `NEXT_PUBLIC_API_URL` - Base URL of the Go API (default: `http://localhost:8080`)

Key files: `lib/api.ts` (calls the API, field name must be `file`), `lib/plans.ts`
(types, friendly labels, tier filtering/sorting), `components/ResultsView.tsx`
(banner + per-provider cards + Standard/Low toggle, derived client-side without refetch).

## Architecture

This is a Go backend service that calculates electricity costs across different power providers based on actual usage data.

### High-Level Flow

1. **API Entry**: HTTP server using Gin framework listens on `/costs` endpoint
2. **File Upload**: Receives CSV file with electricity usage data
3. **Data Processing**:
   - Parse CSV and extract relevant columns (timestamps and usage values)
   - Convert hourly usage data into daily power records
   - Organize data by day/month/weekday for analysis
4. **Cost Calculation**: Multiple pricing models per provider are applied to the normalized usage data
5. **Response**: Returns a nested map with costs for all pricing tiers across all providers

### Directory Structure

All backend paths below are under `backend/`:

- `main.go` - Thin entrypoint: loads config + rates, builds the router, runs the server
- `api/` - HTTP layer (Gin)
  - `server.go` - `Server` struct holding handler dependencies (config, rates)
  - `router.go` - `NewRouter(cfg, rates)` builds the engine + registers routes
  - `middleware.go` - CORS middleware
  - `costs.go` - `POST /costs` handler (parses upload in-memory, no disk write)
  - `health.go` - `GET /health` liveness probe
- `config/` - `Config` struct + `Load()`; the single place env vars are read
- `rates/` - Typed pricing table loaded once from `data/rates.json`; `Get`/`Levy` return errors
- `cost_calculators/` - Pricing logic for different power companies
  - `costs.go` - `AllPrices(records, rates)`; a `specs` table drives every priced plan
  - `contact/` - Contact Energy pricing models (Simple Rates, Good Charge, Good Nights, Good Weekends)
  - `nova/` - Nova Energy pricing models (General Rates)
- `utils/` - Data processing utilities
  - `usage_data.go` - CSV reading/filtering; `ParseUsageData(io.Reader)` + `GetUsageData(path)`
  - `day_power.go` - DayPower struct and aggregation functions (totals by time of use, all days)
  - `rounding.go` - Rounding utilities for cost calculations

### Key Concepts

**DayPower**: Core data structure representing one day's power usage broken into 24 hourly buckets. Contains:
- `date` - Format: "02/01/2006"
- `month` - Format: "01/2006"
- `day` - Day name ("Monday", etc.)
- `usage` - Array of 24 float64 values (one per hour)

**CSV Format**: Expected input has at least 13 columns, using columns:
- Column 9: Start DateTime
- Column 10: End DateTime
- Column 12: Usage value

**Time Periods**: Data can be disaggregated by:
- Weekday vs Weekend
- Hour ranges (e.g., 0-6 peak, 7-22 standard, 23-24 off-peak)
- By month

## Cost Calculation Notes

- Usage totals (`WeekdayUsage`/`WeekendUsage`/`TotalUsage`) and the daily fixed charge cover
  **all days present in the upload** — costs reflect the full uploaded period (e.g. a
  6-month CSV yields a 6-month cost), not a single calendar month.
- `cost_calculators.CalculateSimpleRatesCost` is currently fed the **GoodNights** rate (see
  the `specs` table in `costs.go`) — preserved as-is; likely a bug worth revisiting.

## Adding an Endpoint

1. Add a handler method on `*api.Server` in a new file under `api/` (e.g. `(s *Server) Foo`)
2. Register it in `api/router.go` (one line, e.g. `engine.GET("/foo", s.Foo)`)

## Adding a New Power Provider

1. Create/extend the package under `cost_calculators/{provider}/`
2. Implement a `Calculate{Plan}Cost(records []utils.DayPower, rate rates.Rate, levy float64) (float64, error)`
3. Add the plan/tier entries to the `specs` table in `cost_calculators/costs.go`
4. Add the company/plan rates to `data/rates.json`
5. Use `WeekdayUsage()`, `WeekendUsage()`, `TotalUsage()`, and `DayCount()` from `utils`

## Notes

- Tests: `cd backend && go test ./...` (config, rates, cost_calculators, api covered)
- CORS is configurable; ensure frontend origin is in `CORS_ALLOWED_ORIGINS` in production
- Max upload size is 10 MiB; uploads are processed in-memory (not written to disk)
- The `/costs` endpoint only accepts POST requests with multipart form data (`file` field)
- `GET /health` returns `{"status":"ok"}`
