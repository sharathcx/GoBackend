# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & Run

```bash
go run main.go          # Start server (default :8000)
go build -o GoBackend   # Build binary
```

Requires a `.env` file with `MONGO_URI`, `DATABASE_NAME`, `PORT`, and `SECRET`.

No tests exist yet. No Makefile or Dockerfile.

## Architecture

This is a Go REST API using **Gin** with **MongoDB**, organized by domain modules.

### Fastapify (`fastapify/`)

Custom wrapper around Gin that provides:
- **Automatic OpenAPI 3.0.3 spec generation** from Go struct reflection (no schema files needed)
- **Type-safe request binding** via generics: `Req[T](c)` retrieves validated/bound request data
- **Route builder pattern**: `api.POST("/path", handler).Body(Schema{}).Response(ResponseType{})`
- **URI parameter protection**: prevents body fields from overriding URI params
- Serves Swagger UI at `/docs` and spec at `/openapi.json`

### Module Pattern (`modules/<domain>/`)

Each domain module (user, movie) contains four files following a strict naming convention:
- `<domain>.routes.go` — route registration via fastapify
- `<domain>.handlers.go` — HTTP handlers using `Req[T]()` for typed request access
- `<domain>.schemas.go` — request/response structs with `json`, `bson`, and `binding` tags
- `<domain>.database.go` — MongoDB CRUD operations with context support

To add a new module: create a new directory under `modules/`, follow the four-file pattern, and register routes in `main.go`.

### Supporting Packages

- `database/` — MongoDB client init (via `init()`) and `OpenCollection()` helper
- `globals/` — Env var loading via godotenv into `globals.Vars` singleton
- `utils/` — Standardized API responses (`ApiResponse[T]`, `ApiError`) and UID generation (`GenerateUID("PREFIX")`)

### Response Format

All endpoints return a standardized JSON envelope: `{statusCode, data, message, success, code}`. Use `utils.Response()` and `utils.Error()` helpers.
