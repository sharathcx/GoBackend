package fastapify

import (
	"context"
	"net/http"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"GoBackend/utils"
)

// --- Core Types ---

type Wrapper struct {
	Engine *gin.Engine
	Routes []RouteMeta
}

type RouteMeta struct {
	Method       string
	Path         string
	Tag          string
	BodyType     reflect.Type
	ResponseType reflect.Type
}

type RouteBuilder struct {
	wrapper *Wrapper
	index   int
}

func New(r *gin.Engine) *Wrapper {
	return &Wrapper{Engine: r}
}

// --- Bind Helper ---

// Bind validates and binds the request into the given struct.
// Automatically detects the HTTP method:
//   - GET, DELETE → binds from query params
//   - POST, PUT, PATCH → binds from JSON body
//
// Returns true on success, false on failure (error response is auto-sent).
func Bind(c *gin.Context, req any) bool {
	var err error

	switch c.Request.Method {
	case http.MethodGet, http.MethodDelete:
		err = c.ShouldBindQuery(req)
	default:
		err = c.ShouldBindJSON(req)
	}

	if err != nil {
		statusCode, response := utils.HandleError(err)
		c.JSON(statusCode, response)
		return false
	}
	return true
}

// --- Route Registration ---

func (w *Wrapper) handle(method, path string, handler gin.HandlerFunc, middleware ...gin.HandlerFunc) *RouteBuilder {
	ginPath := toGinPath(path)

	w.Routes = append(w.Routes, RouteMeta{
		Method: method,
		Path:   path,
		Tag:    deriveTag(path),
	})

	handlers := make([]gin.HandlerFunc, 0, len(middleware)+1)
	handlers = append(handlers, middleware...)
	handlers = append(handlers, handler)
	w.Engine.Handle(method, ginPath, handlers...)

	return &RouteBuilder{wrapper: w, index: len(w.Routes) - 1}
}

func (w *Wrapper) GET(path string, handler gin.HandlerFunc, middleware ...gin.HandlerFunc) *RouteBuilder {
	return w.handle(http.MethodGet, path, handler, middleware...)
}

func (w *Wrapper) POST(path string, handler gin.HandlerFunc, middleware ...gin.HandlerFunc) *RouteBuilder {
	return w.handle(http.MethodPost, path, handler, middleware...)
}

func (w *Wrapper) PUT(path string, handler gin.HandlerFunc, middleware ...gin.HandlerFunc) *RouteBuilder {
	return w.handle(http.MethodPut, path, handler, middleware...)
}

func (w *Wrapper) PATCH(path string, handler gin.HandlerFunc, middleware ...gin.HandlerFunc) *RouteBuilder {
	return w.handle(http.MethodPatch, path, handler, middleware...)
}

func (w *Wrapper) DELETE(path string, handler gin.HandlerFunc, middleware ...gin.HandlerFunc) *RouteBuilder {
	return w.handle(http.MethodDelete, path, handler, middleware...)
}

// --- RouteBuilder Chainable Methods ---

// Body sets the request body schema type for Swagger documentation.
func (rb *RouteBuilder) Body(schema any) *RouteBuilder {
	rb.wrapper.Routes[rb.index].BodyType = reflect.TypeOf(schema)
	return rb
}

// Response sets the response schema type for Swagger documentation.
func (rb *RouteBuilder) Response(schema any) *RouteBuilder {
	rb.wrapper.Routes[rb.index].ResponseType = reflect.TypeOf(schema)
	return rb
}

// --- Helpers ---

func toGinPath(path string) string {
	ginPath := strings.ReplaceAll(path, "{", ":")
	return strings.ReplaceAll(ginPath, "}", "")
}

func deriveTag(path string) string {
	trimmed := strings.TrimPrefix(path, "/")
	if idx := strings.Index(trimmed, "/"); idx != -1 {
		return strings.Title(trimmed[:idx])
	}
	return strings.Title(trimmed)
}

var paramPattern = regexp.MustCompile(`\{(\w+)\}`)

func extractParamNames(path string) []string {
	matches := paramPattern.FindAllStringSubmatch(path, -1)
	names := make([]string, 0, len(matches))
	for _, m := range matches {
		names = append(names, m[1])
	}
	return names
}

// --- Middleware ---

func TimeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
