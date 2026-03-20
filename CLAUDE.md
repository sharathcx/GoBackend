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

Go REST API using **Gin** + **MongoDB**, organized by domain modules. `fastapify/` wraps Gin with automatic OpenAPI generation, request binding, and URI param protection.

## Module Rules (`modules/<domain>/`)

Each module has these files. Register routes in `main.go` via `<domain>.RegisterRoutes(api)`.

### Routes (`<domain>.routes.go`)
- Single exported function: `RegisterRoutes(api *fastapify.Wrapper)`
- Use `api.Group("/prefix")` to avoid repeating the base path
- Chain `.Params(ParamsSchema{})`, `.Body(PayloadSchema{})`, `.Response(DomainSchema{})` for OpenAPI docs
- URI params use `{param_name}` syntax
- Middleware passed as extra args: `group.GET("/path", handler, middleware.AuthMiddleware())`

### Handlers (`<domain>.handlers.go`)
- Signature: `func XxxHandler(c *gin.Context) any` — never call `c.JSON()` directly
- Return `utils.NewApiResponse(statusCode, data, message)` for success
- Return database errors as-is (they are already `*utils.ApiError`): `return err`
- For handler-level errors use shorthand constructors: `utils.InternalError("msg")`, `utils.Unauthorized("msg")`, etc.
- Get validated body via `fastapify.Req[PayloadSchema](c)`, URI params via `fastapify.Params[ParamsSchema](c)`
- Pass `c.Request.Context()` to all database calls

### Schemas (`<domain>.schemas.go`)
- **DB model**: `<Domain>Schema` (e.g. `UserSchema`, `MovieSchema`) — carries `bson` + `json` tags
- **Request payloads**: `<Action><Domain>PayloadSchema` (e.g. `RegisterPayloadSchema`, `UpdateUserPayloadSchema`)
  - Create: fields use `binding:"required,..."`
  - Update: fields use `binding:"omitempty,..."` and `bson:"field,omitempty"` (enables partial `$set`)
  - GET query: use `form` tag instead of `json`
- **Embedded value objects**: `<Name>Schema` (e.g. `GenreSchema`, `RankingSchema`) — nested with `binding:"required,dive"` or `binding:"omitempty,dive"`

### Database (`<domain>.database.go`)
- Package-level collection: `var xxxCollection = database.OpenCollection("name")`
- All functions take `context.Context` as first param, return `(*DomainSchema, *utils.ApiError)`
- Use semantic error constructors: `utils.NotFound("not found")`, `utils.Conflict("already exists")`, `utils.InternalError(err.Error())`
- Updates use `FindOneAndUpdate` with `options.After`; partial updates work via omitempty bson tags

### Utils (optional: `<domain>.utils.go`)
- Domain-specific helpers (e.g. password hashing) — only create when needed

## Supporting Packages

- `database/` — MongoDB client init (via `init()`) and `OpenCollection()` helper
- `globals/` — Env var loading via godotenv into `globals.Vars` singleton
- `utils/` — `ApiResponse[T]`, `ApiError`, `HandleError()`, `InvokeUID(prefix, length)`, error code constants
- `fastapify/` — `HandlerFunc = func(c *gin.Context) any`, auto-validation, `Group`, `Params[T]`, `Req[T]`, OpenAPI gen, Scalar docs at `/docs`
- `middleware/auth/` — `AuthMiddleware()`, `GenerateJWT()`, `ValidateToken()`, `SignedDetailsSchema`

## Error Handling

Standardized `*utils.ApiError` flows through all layers:

| Layer | Returns | Error constructors |
|-------|---------|-------------------|
| **Database** | `(*Schema, *utils.ApiError)` | `utils.NotFound()`, `utils.Conflict()`, `utils.InternalError()` |
| **Handler** | `return err` (passthrough from DB) | `utils.Unauthorized()`, `utils.InternalError()` for handler-only errors |
| **Middleware** | `utils.HandleError(utils.Unauthorized(...))` → writes JSON + aborts | |

Shorthand constructors: `NotFound`, `BadRequest`, `Unauthorized`, `Forbidden`, `Conflict`, `InternalError` — each sets correct HTTP status + error code.

## Response Format

All endpoints return: `{statusCode, data, message, success, code}`. Fastapify auto-serializes: `*ApiError` routes through `HandleError()`, anything else writes as `200 OK`.
