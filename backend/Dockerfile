# --- build stage ---
FROM golang:1.25-alpine AS build
WORKDIR /src

# Cache deps first.
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/go-electric .

# --- run stage ---
FROM alpine:3.21
WORKDIR /app
COPY --from=build /out/go-electric /app/go-electric
# Runtime config the server reads relative to its working dir (rates.json).
# Uploads are parsed in-memory and never written to disk (data.csv is excluded
# via .dockerignore and is not needed at runtime).
COPY --from=build /src/data ./data

# The server defaults to localhost:8080; in a container bind to all interfaces.
ENV HOST=0.0.0.0 \
    PORT=8080 \
    ENV=production \
    GIN_MODE=release
EXPOSE 8080

ENTRYPOINT ["/app/go-electric"]
