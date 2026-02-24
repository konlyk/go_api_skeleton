# AGENTS State File

Last updated: 2026-02-24

## Purpose
This file tracks architecture intent, runtime flow, and notable changes so future work stays consistent.

## Current Architecture
- Style: `go-backend-clean-architecture` style layout with OpenAPI-first transport.
- Entrypoint: Cobra command in `cmd/main.go` + `cmd/root.go`.
- Bootstrap layer:
  - `bootstrap/observability.go` wires zerolog, Prometheus, OTEL.
  - `bootstrap/health.go` manages liveness/readiness state.
  - `bootstrap/app.go` loads domain config, initializes shared runtime components, and stores controller set.
- API layer:
  - `controller/ping_controller.go` handles ping endpoint logic.
  - `controller/hello_controller.go` handles hello endpoint logic.
  - `controller/health_controller.go` handles health endpoint logic.
  - `controller/setup.go` builds and returns the concrete controller set.
  - controllers implement generated OpenAPI strict operation methods.
  - `controller/middleware/middleware.go` provides request logging, panic recovery, and HTTP metrics.
  - `api/spec.go` embeds bundled OpenAPI contract for runtime request validation.
  - routing and OpenAPI validation setup lives in `bootstrap/route.go`.
  - `api/openapi/openapi.gen.go` contains generated server/types + embedded spec.
- Domain layers:
  - `domain/config.go` loads config via Viper from optional `config.yaml` + environment decode hooks.
  - `domain/` defines entities/interfaces (including shared health state contract).
  - `usecase/` implements business logic.
  - `repository/` contains infrastructure implementations.
- OpenAPI sources:
  - `api/openapi.yaml`
  - `api/paths/*.yaml`
  - `api/components/schemas/*.yaml`

## Runtime Flow
1. `cmd/root.go` executes Cobra root command and accepts optional `--config` path.
2. Command creates signal-aware context and initializes app dependencies via `bootstrap.NewApplicationWithConfig`.
3. `domain.LoadConfig` loads `config.yaml` (or explicit config path), then applies environment overrides into `domain.Config`.
4. `bootstrap/route.go` builds Gin engine with global middleware:
   - panic recovery
   - request logging
   - Prometheus HTTP metrics
   - OTEL Gin middleware
5. Router exposes `/metrics` and applies OpenAPI request validation middleware using embedded bundled spec plus OpenAPI security authentication hook.
6. Router builds an OpenAPI strict server from controllers composed by `controller.Setup(...)` and uses generated `RegisterHandlersWithOptions(...)` to auto-register endpoints from spec.
7. Controllers delegate to usecases (`domain` -> `usecase` -> `repository` path).
8. On shutdown signal:
   - readiness set to `not_ready`
   - readiness drain delay waits
   - HTTP server shuts down with timeout
   - observability resources shut down

## OpenAPI + Codegen Flow
1. Edit source spec in `api/openapi.yaml`, `api/paths/*`, `api/components/schemas/*`.
2. Bundle spec with Redocly into `api/openapi.bundled.yaml`.
3. Regenerate Go bindings with `go generate ./api/generate` when operation signatures or models change.
4. Runtime request validation and private/public enforcement are driven from bundled OpenAPI contract (security requirements + auth hook).
5. Route registration is generated from OpenAPI (`RegisterHandlersWithOptions`), so handlers stay aligned by implementing generated operation methods.

## Operational Endpoints
- `GET /v1/ping`
- `GET /v1/hello`
- `GET /livez`
- `GET /readyz`
- `GET /metrics`

### Access Model
- Public endpoints: `/v1/ping`, `/livez`, `/readyz`
- Private endpoints: `/v1/hello` (Bearer token via `PRIVATE_API_TOKEN`)

## State Snapshot
- HTTP framework: `gin`.
- Logger: `zerolog`.
- Command framework: `cobra`.
- Config loading: Viper-based decode from `config.yaml` + env overrides into `domain.Config`.
- OpenAPI runtime request validation: enabled.
- OpenAPI security auth hook: enabled (Bearer token for secured operations).
- Codegen model: single strict Gin server package (`api/openapi`).
- Route composition model: OpenAPI-generated auto-registration via strict handler + `controller.Setup(...)` controller set.
- Graceful shutdown: readiness drain + timeout.
- Metrics: Prometheus registry + HTTP middleware.
- Tracing: optional OTLP tracing (`ENABLE_OTEL_TRACING`).

## Change Log
- 2026-02-22: Replaced placeholder project with OpenAPI-first skeleton and operational middleware.
- 2026-02-22: Migrated transport/logging stack to `gin` + `zerolog`.
- 2026-02-24: Refactored project to `go-backend-clean-architecture` style top-level layout (`bootstrap`, `api`, `domain`, `usecase`, `repository`, `cmd`).
- 2026-02-24: Removed module/lifecycle vertical-slice scaffolding to reduce complexity.
- 2026-02-24: Moved generated OpenAPI package to `api/openapi` so `api/controller` only contains controller implementation files.
- 2026-02-24: Moved route setup from `api/route` to `bootstrap/route.go`.
- 2026-02-24: Moved health state contract to `domain/health.go` so controller does not import `bootstrap`, avoiding package cycles.
- 2026-02-24: Moved controller to root `controller/` and middleware to `controller/middleware/`.
- 2026-02-24: Added explicit composition in `bootstrap/wire.go` and split controller responsibilities into `ping` and `health` feature controllers.
- 2026-02-24: Added `/v1/hello` endpoint with `HelloUsecase` and `HelloController`.
- 2026-02-24: Removed `controller/api_controller.go`; routing now uses independent modules without controller-level or bootstrap-level API aggregator.
- 2026-02-24: Simplified composition to explicit typed dependencies (`Repositories/Usecases/Controllers`) and removed module-based route indirection.
- 2026-02-24: Switched route wiring to OpenAPI-generated registration (`NewStrictHandler` + `RegisterHandlersWithOptions`) and removed manual `apiGroup.GET(...)` bindings.
- 2026-02-24: Added OpenAPI security-driven public/private routes and validator authentication hook (`PRIVATE_API_TOKEN`).
- 2026-02-24: Added `api/spec.go` to embed bundled OpenAPI contract for runtime validation independent of generated embedded spec refresh.
- 2026-02-24: Removed `BuildDependencies`/`bootstrap/wire.go` and moved controller composition to `controller/setup.go`.
- 2026-02-24: Replaced `bootstrap/env.go` with `domain/config.go` and switched configuration decoding to Viper.
- 2026-02-24: Added Cobra command entrypoint with `--config` and moved startup lifecycle to command execution flow.
- 2026-02-24: Extended `domain/config.go` to optionally load `config.yaml` (or explicit config file) with environment override precedence.

## Update Protocol
When making changes, update this file in the same commit/PR:
- `Last updated` date.
- `State Snapshot` if behavior or stack changed.
- `Change Log` with one line per meaningful change.
- `Runtime Flow` or `OpenAPI + Codegen Flow` if lifecycle steps changed.
