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

Go REST API using **Gin** + **MongoDB**, organized by layer packages. `fastapify/` wraps Gin with automatic OpenAPI generation, request binding, and URI param protection.

**Dependency flow (strictly one-way):**
```
routes/ → handlers/ → database/ → schemas/
    ↘                    ↗
     websocket/ → database/
         ↘
       schemas/
```

Register all routes in `main.go` via `routes.RegisterRoutes(api)`.

## Layer Packages

### `schemas/` — All type definitions
Domain files: `<domain>.schemas.go` (e.g. `user.schemas.go`, `movie.schemas.go`, `chat.schemas.go`)

- **DB model**: `<Domain>Schema` (e.g. `UserSchema`, `MovieSchema`) — carries `bson` + `json` tags
- **Request payloads**: `<Action><Domain>PayloadSchema` (e.g. `RegisterPayloadSchema`, `UpdateUserPayloadSchema`)
  - Create: fields use `binding:"required,..."`
  - Update: fields use `binding:"omitempty,..."` and `bson:"field,omitempty"` (enables partial `$set`)
  - GET query: use `form` tag instead of `json`
- **Embedded value objects**: `<Name>Schema` (e.g. `GenreSchema`, `RankingSchema`) — nested with `binding:"required,dive"` or `binding:"omitempty,dive"`
- **WebSocket structs**: `WS<Name>Schema` (e.g. `WSClientSchema`, `WSHubSchema`, `WSBroadcastMessageSchema`, `WSMessageSchema`, `WSResponseSchema`) — used for WebSocket connection management and message types

### `database/` — MongoDB init + all domain DB functions
Domain files: `<domain>.database.go` (e.g. `user.database.go`)

- Collection vars initialized in `init()`: `var xxxCollection *mongo.Collection; func init() { xxxCollection = OpenCollection("name") }`
- All functions take `context.Context` as first param, return `(*schemas.DomainSchema, *utils.ApiError)`
- Use semantic error constructors: `utils.NotFound("not found")`, `utils.Conflict("already exists")`, `utils.InternalError(err.Error())`
- Updates use `FindOneAndUpdate` with `options.After`; partial updates work via omitempty bson tags
- `database.go` contains MongoDB client init (via `init()`) and `OpenCollection()` helper

### `handlers/` — REST handler functions
Domain files: `<domain>.handler.go` (e.g. `user.handler.go`)

- Signature: `func XxxHandler(c *gin.Context) any` — never call `c.JSON()` directly
- Return `utils.NewApiResponse(statusCode, data, message)` for success
- Return database errors as-is (they are already `*utils.ApiError`): `return err`
- For handler-level errors use shorthand constructors: `utils.InternalError("msg")`, `utils.Unauthorized("msg")`, etc.
- Get validated body via `fastapify.Req[schemas.PayloadSchema](c)`, URI params via `fastapify.Params[schemas.ParamsSchema](c)`
- Pass `c.Request.Context()` to all database calls
- Qualify types: `schemas.XxxSchema`, `database.XxxFunc()`

### `websocket/` — WebSocket hub, client, and page
- `hub.go` — `NewWSHub()`, `DefaultWSHub`, `RunWSHub()`, `JoinRoom()`, `LeaveRoom()` (standalone functions, not methods)
- `client.go` — `ReadPump()`, `WritePump()`, `handleAction()`, `sendError()`
- `page.go` — `ChatPageHTML` constant for the chat UI

### `routes/` — Route registration
Domain files: `<domain>.routes.go` (e.g. `user.routes.go`)

- Unexported per-domain functions: `registerUserRoutes`, `registerMovieRoutes`, `registerChatRoutes`
- Single exported aggregator: `RegisterRoutes(api *fastapify.Wrapper)` in `routes.go`
- Use `api.Group("/prefix")` to avoid repeating the base path
- Chain `.Params(schemas.ParamsSchema{})`, `.Body(schemas.PayloadSchema{})`, `.Response(schemas.DomainSchema{})` for OpenAPI docs
- URI params use `{param_name}` syntax
- Middleware passed as extra args: `group.GET("/path", handlers.XxxHandler, auth.AuthMiddleware())`

## Supporting Packages

- `globals/` — Env var loading via godotenv into `globals.Vars` singleton
- `utils/` — `ApiResponse[T]`, `ApiError`, `HandleError()`, `InvokeUID(prefix, length)`, error code constants, `HashPassword()`, `VerifyPassword()`, `GenerateJWT()`, `ValidateToken()`, `GetAccessTokenFromHeader()`
- `fastapify/` — `HandlerFunc = func(c *gin.Context) any`, auto-validation, `Group`, `Params[T]`, `Req[T]`, OpenAPI gen, Scalar docs at `/docs`
- `middleware/` — `AuthMiddleware()` gin middleware for JWT validation

## Error Handling

Standardized `*utils.ApiError` flows through all layers:

| Layer | Returns | Error constructors |
|-------|---------|-------------------|
| **Database** | `(*schemas.Schema, *utils.ApiError)` | `utils.NotFound()`, `utils.Conflict()`, `utils.InternalError()` |
| **Handler** | `return err` (passthrough from DB) | `utils.Unauthorized()`, `utils.InternalError()` for handler-only errors |
| **Middleware** | `utils.HandleError(utils.Unauthorized(...))` → writes JSON + aborts | |

Shorthand constructors: `NotFound`, `BadRequest`, `Unauthorized`, `Forbidden`, `Conflict`, `InternalError` — each sets correct HTTP status + error code.

## Response Format

All endpoints return: `{statusCode, data, message, success, code}`. Fastapify auto-serializes: `*ApiError` routes through `HandleError()`, anything else writes as `200 OK`.
