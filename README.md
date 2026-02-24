# Go API Skeleton

A Go API skeleton based on the `go-backend-clean-architecture` style, with OpenAPI-first development.

## What this gives you

- Clean top-level structure: `bootstrap`, `api`, `domain`, `usecase`, `repository`, `cmd`
- Gin + Zerolog runtime defaults
- OpenAPI contract + Redocly lint/bundle/docs
- `oapi-codegen` strict Gin server bindings
- OpenAPI-driven automatic route registration
- Runtime OpenAPI request validation middleware
- OpenAPI security-driven private/public routes
- Kubernetes probes: `/livez`, `/readyz`
- Utility endpoints: `/v1/ping`, `/v1/hello`
- Prometheus metrics: `/metrics`
- Graceful shutdown with readiness drain
- Optional OTLP tracing

## Project structure

```text
cmd/                 # application entrypoint
bootstrap/           # app wiring, health, observability, router setup
controller/          # feature controllers + controller/setup.go composition
controller/middleware/
                     # request log, panic recovery, http metrics middleware
api/openapi/         # generated OpenAPI server/types + embedded spec
api/spec.go          # embedded bundled OpenAPI spec for runtime validation
bootstrap/route.go   # gin middleware + OpenAPI auto route registration
domain/              # core interfaces/entities + config.go (viper-backed config)
usecase/             # business use-cases
repository/          # infrastructure implementations for domain interfaces
api/openapi.yaml     # root OpenAPI spec
api/paths/           # path fragments
api/components/      # schema fragments
```

## OpenAPI workflow

```bash
npm --prefix api install
npm --prefix api run lint
npm --prefix api run bundle
go generate ./api/generate
```

## Run

```bash
go run ./cmd
# optional explicit config file
go run ./cmd --config ./config.yaml
```

## Test

```bash
go test ./...
```

## Key env vars

- `PORT` (default `8080`)
- `SERVICE_NAME` (default `go-api-skeleton`)
- `PRIVATE_API_TOKEN` (default `dev-private-token`)
- `LOG_LEVEL` (`debug|info|warn|error`, default `info`)
- `ENABLE_OTEL_TRACING` (default `false`)
- `TRACE_SAMPLE_RATIO` (default `0.1`)
- `SHUTDOWN_DRAIN_DELAY` (default `5s`)
- `SHUTDOWN_TIMEOUT` (default `20s`)

## Config file

- The app automatically tries `./config.yaml` when present.
- You can explicitly set a file with `--config`.
- Environment variables override config file values.

## Public/Private endpoints

- Public: `/v1/ping`, `/livez`, `/readyz`
- Private: `/v1/hello` (`Authorization: Bearer <PRIVATE_API_TOKEN>`)
