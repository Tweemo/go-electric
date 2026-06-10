# Production readiness plan — go-electric API

This document is a step-by-step implementation guide for taking the Gin API in `main.go` to a production-ready state. Work through phases in order; later phases assume earlier ones are done.

**Current baseline**

- Single route: `POST /costs` (multipart CSV upload)
- Binds to `localhost` only via `router.Run`
- Writes every upload to `data/data.csv` (race-prone)
- `utils` uses `log.Fatal` on CSV/config errors (crashes the process)
- CORS with dev fallbacks when env is unset
- No health checks, graceful shutdown, tests, or container/CI artifacts

---

## How to use this plan

| Column | Meaning |
|--------|---------|
| **Files** | What to create or edit |
| **Steps** | Exact implementation sequence |
| **Verify** | How to confirm it works |
| **Done when** | Acceptance criteria |

Recommended order: **Phase 1 → 2 → 3 → 4 → 5**. Phase 1 alone is enough for a cautious private deploy; complete all phases before a public internet-facing launch.

---

## Phase 1 — Critical safety and stability

### 1.1 Bind address and environment model

**Problem:** `router.Run("localhost:" + port)` only accepts connections on the loopback interface. Containers and load balancers cannot reach the process.

**Files:** `main.go`, `.env.example`, deployment docs (README or platform config)

**Steps:**

1. Add environment variables:

   ```env
   PORT=3000
   HOST=0.0.0.0
   ENV=development
   GIN_MODE=debug
   ```

   Production example:

   ```env
   ENV=production
   GIN_MODE=release
   HOST=0.0.0.0
   PORT=3000
   ```

2. In `main()`, resolve listen address:

   ```go
   host := os.Getenv("HOST")
   if host == "" {
       host = "localhost" // safe local default
   }
   addr := net.JoinHostPort(host, port)
   ```

3. Only call `godotenv.Load()` when `ENV != "production"` so production relies on platform-injected env, not a `.env` file on disk:

   ```go
   if os.Getenv("ENV") != "production" {
       _ = godotenv.Load()
   }
   ```

4. Set Gin mode from env before creating the router:

   ```go
   if mode := os.Getenv("GIN_MODE"); mode != "" {
       gin.SetMode(mode)
   }
   ```

**Verify:** `curl http://127.0.0.1:3000/health` works when `HOST=0.0.0.0`. Inside Docker, `docker run -p 3000:3000 ...` reaches the API from the host.

**Done when:** Listen host and Gin mode are fully env-driven; production never depends on `.env` on the filesystem.

---

### 1.2 Replace `http.Server` + graceful shutdown (replace `router.Run`)

**Problem:** `router.Run` blocks forever and kills in-flight requests on SIGTERM (bad for deploys).

**Files:** `main.go` (or new `server/server.go`)

**Steps:**

1. Build the Gin engine as today (`router := gin.New()` or `gin.Default()`).

2. Create an explicit server:

   ```go
   srv := &http.Server{
       Addr:              addr,
       Handler:           router,
       ReadHeaderTimeout: 10 * time.Second,
       ReadTimeout:       30 * time.Second,
       WriteTimeout:      60 * time.Second,
       IdleTimeout:       120 * time.Second,
   }
   ```

3. Start in a goroutine:

   ```go
   go func() {
       if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
           slog.Error("server failed", "err", err)
           os.Exit(1)
       }
   }()
   ```

4. Wait for `SIGINT` / `SIGTERM`, then shutdown:

   ```go
   quit := make(chan os.Signal, 1)
   signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
   <-quit

   ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
   defer cancel()
   if err := srv.Shutdown(ctx); err != nil {
       slog.Error("shutdown failed", "err", err)
   }
   ```

**Verify:** Start server, send a slow request, send SIGTERM — process exits within ~15s without corrupting responses mid-flight.

**Done when:** No `router.Run` in codebase; deploy platform can roll out without hard kills.

---

### 1.3 Remove `log.Fatal` from the request path

**Problem:** `utils/usage_data.go` `readCsvFile` and `utils/config.go` call `log.Fatal`, which terminates the entire process on one bad upload or missing rates file.

**Files:** `utils/usage_data.go`, `utils/config.go`, `main.go` (`Costs` handler), any caller of `GetRate`

**Steps:**

1. Change signatures to return errors:

   ```go
   func readCsvFile(filePath string) ([][]string, error)
   func GetUsageData(filepath string) ([][]string, error)
   func CalculateDayPower(usageData [][]string) ([]DayPower, error) // optional: validate parse errors
   ```

2. Replace `log.Fatal(...)` with `return nil, fmt.Errorf("...: %w", err)`.

3. In `Costs`, handle errors:

   ```go
   data, err := utils.GetUsageData(path)
   if err != nil {
       slog.Error("usage data", "err", err, "request_id", requestID)
       c.JSON(http.StatusBadRequest, gin.H{"error": "invalid or unreadable CSV"})
       return
   }
   ```

4. For `GetRate` / config: load `data/rates.json` once at startup into memory, or return `error` and let the handler respond 500. Never `Fatal` during a request.

5. Run `go test ./...` and fix compile errors from signature changes.

**Verify:** Upload a non-CSV or empty file — API returns 400, process stays up. Rename `data/rates.json` temporarily — `/ready` fails (see 2.1) but process does not exit on unrelated requests if rates are cached at startup.

**Done when:** `grep -R log.Fatal utils/` shows no Fatal on paths reachable from HTTP handlers.

---

### 1.4 Per-request temp files (no shared `data/data.csv`)

**Problem:** Concurrent uploads overwrite the same path; disk growth; predictable filename.

**Files:** `main.go` (`Costs`)

**Steps:**

1. After opening the upload, create a temp file:

   ```go
   tmp, err := os.CreateTemp("", "usage-*.csv")
   if err != nil { /* 500 */ }
   tmpPath := tmp.Name()
   defer os.Remove(tmpPath)
   defer tmp.Close()
   ```

2. `io.Copy(tmp, src)` instead of writing to `data/data.csv`.

3. Pass `tmpPath` to `utils.GetUsageData(tmpPath)`.

4. Remove `os.MkdirAll("data", ...)` from the handler unless `data/` is only for static `rates.json`.

**Verify:** Run two parallel `curl -F file=@...` uploads — both succeed, no corruption, no permanent `data/data.csv` growth from uploads.

**Done when:** Handler never writes user content to a fixed shared path.

---

### 1.5 Upload validation

**Problem:** Any file can be uploaded; size is only loosely bounded by multipart memory.

**Files:** `main.go`, optional new `internal/validate/upload.go`, `.env.example`

**Steps:**

1. Add `MAX_UPLOAD_BYTES` (default `10485760`, 10 MiB).

2. After `c.FormFile("file")`:

   ```go
   if file.Size > maxUploadBytes {
       c.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": "file too large"})
       return
   }
   ```

3. Validate extension (case-insensitive): `.csv` via `filepath.Ext(file.Filename)`.

4. Optional content sniff (recommended): read first 512 bytes, use `github.com/gabriel-vasile/mimetype` (already indirect in `go.mod`):

   ```go
   go get github.com/gabriel-vasile/mimetype
   ```

   Reset file offset with `io.NewSectionReader` or copy to temp then detect.

5. After `GetUsageData`, require minimum row count (e.g. `len(data) > 0`) or expected column shape in `filterColumns`; return 400 if empty.

**Verify:** Upload 20MB file → 413. Upload `.exe` renamed to `.csv` → 400 if sniffing enabled. Valid `temp/novaenergy.csv` → 200.

**Done when:** Oversized and obviously invalid uploads never reach `CalculateDayPower`.

---

### 1.6 Safe API error responses

**Problem:** `gin.H{"error": err.Error()}` leaks internal details (paths, OS messages).

**Files:** `main.go`, optional `internal/api/errors.go`

**Steps:**

1. Define stable client messages:

   | Condition | Status | Body |
   |-----------|--------|------|
   | Missing `file` field | 400 | `{"error":"file is required"}` |
   | Invalid CSV / parse | 400 | `{"error":"invalid or unreadable CSV"}` |
   | File too large | 413 | `{"error":"file too large"}` |
   | Internal failure | 500 | `{"error":"internal server error"}` |

2. Log full errors with `slog` (see 2.2), including request ID.

3. Never return `err.Error()` to clients in production.

**Verify:** Force a read error — response body is generic; logs contain the real `err`.

**Done when:** All `c.JSON` error payloads use fixed strings; logs hold diagnostics.

---

### 1.7 CORS required in production

**Problem:** Empty `CORS_ALLOWED_ORIGINS` falls back to localhost dev origins.

**Files:** `main.go` (`corsMiddleware`), `.env.example`

**Steps:**

1. In `corsMiddleware`, after parsing origins:

   ```go
   if os.Getenv("ENV") == "production" && (len(origins) != 1 || origins[0] == "") {
       slog.Error("CORS_ALLOWED_ORIGINS must be set in production")
       os.Exit(1)
   }
   ```

2. Document comma-separated list: `https://app.example.com,https://www.example.com`.

3. Keep dev fallback only when `ENV != "production"`.

**Verify:** Start with `ENV=production` and no CORS env — process refuses to start. With origins set — browser preflight from allowed origin succeeds.

**Done when:** Production cannot start with implicit CORS allowlist.

---

## Phase 2 — Operability

### 2.1 Health and readiness endpoints

**Files:** `main.go` or `handlers/health.go`

**Steps:**

1. Add routes (no auth, no heavy work):

   ```go
   router.GET("/health", func(c *gin.Context) {
       c.JSON(http.StatusOK, gin.H{"status": "ok"})
   })
   ```

2. Readiness checks dependencies:

   ```go
   router.GET("/ready", func(c *gin.Context) {
       if _, err := os.Stat("data/rates.json"); err != nil {
           c.JSON(http.StatusServiceUnavailable, gin.H{"status": "not ready"})
           return
       }
       c.JSON(http.StatusOK, gin.H{"status": "ready"})
   })
   ```

3. Point load balancer health checks at `/health` (liveness) and optionally `/ready` (traffic only when ready).

4. Extend `/ready` later if you cache rates at startup — verify cache loaded.

**Verify:** `curl -i localhost:3000/health` → 200. Remove `rates.json` → `/ready` → 503.

**Done when:** Orchestrator can probe both endpoints.

---

### 2.2 Structured logging

**Files:** `main.go`, middleware file `middleware/logging.go`

**Steps:**

1. Use `log/slog` (stdlib) with JSON handler in production:

   ```go
   var handler slog.Handler
   if os.Getenv("ENV") == "production" {
       handler = slog.NewJSONHandler(os.Stdout, nil)
   } else {
       handler = slog.NewTextHandler(os.Stdout, nil)
   }
   slog.SetDefault(slog.New(handler))
   ```

2. Replace println / default Gin logging with a custom middleware that logs: `method`, `path`, `status`, `latency`, `client_ip`, `request_id`.

3. In `Costs`, log at `Info` for successful calculation (duration, row count); `Error` for failures.

**Verify:** One request produces one structured log line; errors include `request_id`.

**Done when:** Logs are machine-parseable in production.

---

### 2.3 Request ID middleware

**Files:** `middleware/request_id.go`, `main.go`

**Steps:**

1. Middleware:

   ```go
   func RequestID() gin.HandlerFunc {
       return func(c *gin.Context) {
           id := c.GetHeader("X-Request-ID")
           if id == "" {
               id = uuid.NewString() // google/uuid or crypto/rand hex
           }
           c.Set("request_id", id)
           c.Header("X-Request-ID", id)
           c.Next()
       }
   }
   ```

2. Register first: `router.Use(RequestID(), Logging(), corsMiddleware())`.

3. Pass `c.GetString("request_id")` into all `slog` calls in handlers.

**Verify:** Response includes `X-Request-ID`; same value appears in logs.

**Done when:** Every request is correlatable in logs.

---

### 2.4 Production Gin stack

**Files:** `main.go`

**Steps:**

1. Prefer explicit setup:

   ```go
   router := gin.New()
   router.Use(gin.Recovery())
   router.Use(RequestID())
   router.Use(Logging())
   router.Use(corsMiddleware())
   ```

2. `gin.Recovery()` catches panics in handlers (returns 500, logs stack).

3. Keep `MaxMultipartMemory` aligned with `MAX_UPLOAD_BYTES`.

4. Run `go mod tidy` and ensure direct requires:

   ```bash
   go get github.com/gin-gonic/gin@v1.12.0
   go get github.com/gin-contrib/cors@v1.7.7
   ```

**Verify:** `go.mod` lists `gin` and `gin-contrib/cors` under `require`, not only `// indirect`.

**Done when:** Dev uses readable logs; production uses release mode + recovery + custom middleware.

---

### 2.5 Expand `.env.example` and README

**Files:** `.env.example`, `README.md`

**Steps:**

1. Document every variable:

   ```env
   ENV=development
   GIN_MODE=debug
   HOST=localhost
   PORT=3000
   CORS_ALLOWED_ORIGINS=http://localhost:3001,http://127.0.0.1:3001
   MAX_UPLOAD_BYTES=10485760
   # Optional later:
   # API_KEY=
   # RATE_LIMIT_RPM=30
   ```

2. README sections: local run, env table, `curl` upload example, CORS notes (`localhost` vs `127.0.0.1`), production deploy checklist pointing to this file.

**Verify:** New developer can run API from README alone.

**Done when:** Configuration surface is fully documented.

---

## Phase 3 — Security and abuse resistance

### 3.1 Authentication (choose one strategy)

**Problem:** `/costs` is public; anyone can trigger expensive work.

**Options (pick one for v1):**

| Strategy | How | Best for |
|----------|-----|----------|
| **API key** | `Authorization: Bearer <key>` or `X-API-Key`; constant-time compare with `subtle.ConstantTimeCompare` | B2B, internal frontends |
| **JWT** | Validate signature via middleware; issuer from IdP | User accounts |
| **Network** | Private VPC only, no public route | Internal tools |

**Files:** `middleware/auth.go`, `.env.example`, frontend fetch headers

**Steps (API key example):**

1. `API_KEY` in env; if `ENV=production` and empty, fail startup.

2. Middleware:

   ```go
   func APIKeyAuth() gin.HandlerFunc {
       expected := os.Getenv("API_KEY")
       return func(c *gin.Context) {
           got := c.GetHeader("X-API-Key")
           if subtle.ConstantTimeCompare([]byte(got), []byte(expected)) != 1 {
               c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
               return
           }
           c.Next()
       }
   }
   ```

3. Apply only to `POST /costs`, not `/health` or `/ready`.

4. Update CORS `AllowHeaders` to include `X-API-Key` if browser sends it.

**Verify:** Request without key → 401. With key → 200.

**Done when:** Production `/costs` is not anonymously callable (unless intentionally private-network only and documented as such).

---

### 3.2 Rate limiting

**Files:** `middleware/ratelimit.go` or reverse proxy config

**Steps (in-app option):**

1. `go get github.com/ulule/limiter/v3` and memory store, or use Redis for multi-instance.

2. Limit by IP (and by API key if present): e.g. 30 requests/minute per IP on `POST /costs`.

3. Return `429` with `Retry-After` header.

**Alternative:** Configure rate limits at nginx, Cloudflare, or API gateway — often simpler for multi-instance deploys.

**Verify:** Burst > limit → 429; after window → 200 again.

**Done when:** Abuse cannot exhaust CPU/disk with unlimited parallel uploads.

---

### 3.3 Security headers middleware

**Files:** `middleware/security_headers.go`

**Steps:**

1. Set on all responses:

   ```go
   c.Header("X-Content-Type-Options", "nosniff")
   c.Header("X-Frame-Options", "DENY")
   c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
   ```

2. If API is HTTPS-only behind proxy: `Strict-Transport-Security` (only when TLS terminates correctly).

**Verify:** `curl -I` shows headers on `/health` and `/costs`.

**Done when:** Baseline browser security headers present.

---

### 3.4 HTTPS and edge termination

**Not implemented in Go** for most deployments.

**Steps:**

1. Terminate TLS at load balancer, Cloudflare, or nginx.

2. Forward `X-Forwarded-For` / `X-Forwarded-Proto`; use `router.SetTrustedProxies([]string{"10.0.0.0/8"})` if you need client IP behind proxy (Gin docs).

3. Do not expose port 3000 directly to the internet without TLS in front.

**Verify:** Public URL is `https://`; HTTP redirects to HTTPS.

**Done when:** All client traffic is encrypted to the edge.

---

## Phase 4 — Correctness

### 4.1 Fix hardcoded billing month

**Problem:** `utils/day_power.go` uses `monthMap["08/2025"]` in `WeekdayUsage`, `WeekendUsage`, `TotalUsage`.

**Files:** `utils/day_power.go`, cost calculator callers if needed

**Steps:**

1. Decide business rule:
   - **Option A:** Use the latest month present in the upload.
   - **Option B:** Use the month with the most days of data.
   - **Option C:** Accept `?month=08/2025` query param (document in API).

2. Implement helper:

   ```go
   func primaryMonth(records []DayPower) string {
       // e.g. max by record count per month key
   }
   ```

3. Replace `"08/2025"` with `primaryMonth(records)` everywhere.

4. Add unit test with CSV spanning two months.

**Verify:** Upload CSV for a different month — costs change appropriately vs hardcoded August.

**Done when:** No literal month strings in usage aggregation.

---

### 4.2 Validate CSV rows in `CalculateDayPower`

**Problem:** `time.Parse` and `ParseFloat` errors are ignored (`_`).

**Files:** `utils/usage_data.go`

**Steps:**

1. On parse failure, return wrapped error with row index:

   ```go
   startTime, err := time.Parse("02/01/2006 15:04:05", record[0])
   if err != nil {
       return nil, fmt.Errorf("row %d: invalid datetime: %w", rowIndex, err)
   }
   ```

2. Map to 400 in `Costs` with generic client message.

3. Optionally validate datetime format regex before parse for clearer errors in logs.

**Verify:** CSV with one malformed date row → 400, descriptive server log.

**Done when:** Bad data never silently produces zero dates or NaN usage.

---

### 4.3 Load `rates.json` once at startup

**Problem:** `GetRate` opens and decodes JSON on every call (slow; Fatal on failure).

**Files:** `utils/config.go`, `main.go`

**Steps:**

1. Add `func LoadRates(path string) error` that reads `data/rates.json` into a package-level struct or `RatesStore`.

2. Call from `main()` before `ListenAndServe`; fail fast if missing in production.

3. Change `GetRate` to read from memory.

4. `/ready` checks that store is non-nil.

**Verify:** `strace` or logs show no per-request `open("data/rates.json")` during `/costs` traffic.

**Done when:** Rates are loaded exactly once per process life.

---

## Phase 5 — Testing, CI, and deployment

### 5.1 HTTP handler tests

**Files:** `main_test.go` or `handlers/costs_test.go`

**Steps:**

1. Use `httptest.NewRecorder` + `gin.CreateTestContext`.

2. Test cases:
   - Valid CSV multipart → 200 + JSON keys `contact`, `nova`
   - No file → 400
   - Empty file → 400
   - File > max size → 413
   - Invalid CSV content → 400
   - With auth middleware: missing key → 401

3. Embed small fixture CSV in `testdata/sample.csv` (few rows, not full 8740-line file).

**Verify:** `go test ./... -count=1` passes locally and in CI.

**Done when:** Core handler regressions are automated.

---

### 5.2 Golden / integration test for pricing

**Files:** `cost_calculators/costs_test.go` or `integration/pricing_test.go`

**Steps:**

1. Run pipeline: `GetUsageData(testdata/sample.csv)` → `CalculateDayPower` → `AllPrices`.

2. Compare to golden JSON snapshot (`testdata/prices.golden.json`).

3. Update golden intentionally when rate formulas change.

**Verify:** Changing a rate in `data/rates.json` fails test until golden updated.

**Done when:** Pricing output is regression-locked.

---

### 5.3 CI pipeline

**Files:** `.github/workflows/ci.yml` (create `.github/workflows/`)

**Steps:**

1. Workflow on push/PR:

   ```yaml
   - uses: actions/checkout@v4
   - uses: actions/setup-go@v5
     with:
       go-version: '1.25.x'
   - run: go vet ./...
   - run: go test ./... -race
   - run: go build -o /dev/null .
   ```

2. Optional: `staticcheck`, `golangci-lint`.

**Verify:** PR shows green CI.

**Done when:** Every merge is vetted and built.

---

### 5.4 Dockerfile

**Files:** `Dockerfile`, `.dockerignore`

**Steps:**

1. Multi-stage build:

   ```dockerfile
   FROM golang:1.25-alpine AS build
   WORKDIR /src
   COPY go.mod go.sum ./
   RUN go mod download
   COPY . .
   RUN CGO_ENABLED=0 go build -o /server .

   FROM alpine:3.20
   RUN adduser -D -g '' appuser
   USER appuser
   WORKDIR /app
   COPY --from=build /server .
   COPY data/rates.json data/rates.json
   ENV ENV=production GIN_MODE=release HOST=0.0.0.0 PORT=3000
   EXPOSE 3000
   CMD ["./server"]
   ```

2. `.dockerignore`: `.env`, `temp/`, `data/data.csv`, `*.csv` uploads, `.git`

3. Build: `docker build -t go-electric .`

4. Run with env: `-e CORS_ALLOWED_ORIGINS=... -e API_KEY=...`

**Verify:** Container starts, `/health` 200, upload works via published port.

**Done when:** Deployable artifact exists without local Go toolchain.

---

### 5.5 Observability (optional but recommended)

**Files:** `middleware/metrics.go`, or platform APM

**Steps:**

1. Prometheus: `go get github.com/prometheus/client_golang`, expose `GET /metrics` on admin port or same server (protect in production).

2. Track: `http_requests_total`, `http_request_duration_seconds`, `costs_upload_bytes`, `costs_processing_errors_total`.

3. Wire alerts on 5xx rate and high latency.

**Verify:** Scrape `/metrics` shows counters incrementing after traffic.

**Done when:** Production issues are visible without reproducing locally.

---

## Suggested file layout (after refactor)

```
go-electric/
├── main.go                 # wiring: env, server, routes
├── middleware/
│   ├── request_id.go
│   ├── logging.go
│   ├── auth.go             # Phase 3
│   ├── ratelimit.go        # Phase 3
│   └── security_headers.go
├── handlers/
│   ├── costs.go
│   └── health.go
├── internal/validate/
│   └── upload.go
├── utils/                  # errors, no Fatal on request path
├── cost_calculators/
├── data/rates.json
├── testdata/
├── Dockerfile
├── .github/workflows/ci.yml
├── .env.example
├── PRODUCTION_PLAN.md      # this file
└── README.md
```

Splitting `main.go` is optional but keeps each phase’s diffs small.

---

## Implementation checklist (copy for tracking)

### Phase 1 — Critical
- [ ] 1.1 Env-based `HOST`, `PORT`, `ENV`, `GIN_MODE`
- [ ] 1.2 `http.Server` + graceful shutdown + timeouts
- [ ] 1.3 Errors instead of `log.Fatal` in utils
- [ ] 1.4 Temp files per upload
- [ ] 1.5 Upload size/type validation
- [ ] 1.6 Generic client errors + server logs
- [ ] 1.7 Require `CORS_ALLOWED_ORIGINS` in production

### Phase 2 — Operability
- [ ] 2.1 `/health` and `/ready`
- [ ] 2.2 Structured `slog` logging
- [ ] 2.3 Request ID middleware
- [ ] 2.4 `gin.New()` + Recovery + tidy deps
- [ ] 2.5 `.env.example` + README

### Phase 3 — Security
- [ ] 3.1 API key or chosen auth
- [ ] 3.2 Rate limiting (app or edge)
- [ ] 3.3 Security headers
- [ ] 3.4 TLS at reverse proxy

### Phase 4 — Correctness
- [ ] 4.1 Dynamic billing month
- [ ] 4.2 Row-level CSV validation
- [ ] 4.3 Load rates at startup

### Phase 5 — Ship
- [ ] 5.1 Handler tests
- [ ] 5.2 Golden pricing test
- [ ] 5.3 GitHub Actions CI
- [ ] 5.4 Dockerfile + .dockerignore
- [ ] 5.5 Metrics (optional)

---

## Production launch gate

Before pointing a public URL at this API, confirm:

1. All Phase 1 items complete  
2. `/health` and `/ready` configured on load balancer  
3. `ENV=production`, `GIN_MODE=release`, `HOST=0.0.0.0`  
4. `CORS_ALLOWED_ORIGINS` set to real frontend origin(s) only  
5. Auth and/or rate limiting on `POST /costs`  
6. TLS in front of the service  
7. CI green on `main`  
8. Hardcoded month removed and validated with real billing CSV  

---

## Reference: minimal `Costs` handler shape (target)

After Phase 1–2, `Costs` should resemble:

```go
func Costs(c *gin.Context) {
    reqID := c.GetString("request_id")

    file, err := c.FormFile("file")
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
        return
    }
    if err := validateUpload(file, maxUploadBytes); err != nil {
        // map to 400/413, log with reqID
        return
    }

    path, cleanup, err := saveUploadToTemp(file)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
        return
    }
    defer cleanup()

    data, err := utils.GetUsageData(path)
    if err != nil {
        slog.Error("csv", "err", err, "request_id", reqID)
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid or unreadable CSV"})
        return
    }

    records, err := utils.CalculateDayPower(data)
    if err != nil {
        slog.Error("parse", "err", err, "request_id", reqID)
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid or unreadable CSV"})
        return
    }

    c.JSON(http.StatusOK, cost_calculators.AllPrices(records))
}
```

This is the intended end state for the upload path described in phases above.
