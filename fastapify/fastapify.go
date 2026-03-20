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

// Req retrieves the automatically bound and validated request data from the context.
func Req[T any](c *gin.Context) *T {
	val, exists := c.Get("fastapify_body")
	if !exists {
		return new(T)
	}
	return val.(*T)
}

// Bind validates and binds the request into the given struct.
// Automatically:
//  1. Binds URI parameters (and makes them immutable)
//  2. Detects HTTP method for Body vs Query binding
//  3. Sends 400 error response on failure
//
// Returns true on success, false on failure.
func Bind(c *gin.Context, req any) bool {
	// Step 1: Bind URI params and save their values for protection
	uriValues := make(map[int]reflect.Value)
	reqVal := reflect.ValueOf(req).Elem()
	if reqVal.Kind() == reflect.Struct {
		_ = c.ShouldBindUri(req)

		// Snapshot URI-tagged fields
		reqType := reqVal.Type()
		for i := 0; i < reqType.NumField(); i++ {
			if reqType.Field(i).Tag.Get("uri") != "" {
				uriValues[i] = reflect.ValueOf(reqVal.Field(i).Interface())
			}
		}
	}

	// Step 2: Bind Body or Query
	var err error
	switch c.Request.Method {
	case http.MethodGet, http.MethodDelete:
		err = c.ShouldBindQuery(req)
	default:
		err = c.ShouldBindJSON(req)
	}

	if err != nil && err.Error() != "EOF" {
		statusCode, response := utils.HandleError(err)
		c.JSON(statusCode, response)
		return false
	}

	// Step 3: Restore URI values (Protection against body override)
	for i, val := range uriValues {
		reqVal.Field(i).Set(val)
	}

	return true
}

// --- Route Registration ---

func (w *Wrapper) handle(method, path string, handler gin.HandlerFunc, middleware ...gin.HandlerFunc) *RouteBuilder {
	ginPath := toGinPath(path)

	routeIdx := len(w.Routes)
	w.Routes = append(w.Routes, RouteMeta{
		Method: method,
		Path:   path,
		Tag:    deriveTag(path),
	})
	meta := &w.Routes[routeIdx]

	// Automatic Validation Middleware
	autoValidator := func(c *gin.Context) {
		if meta.BodyType != nil {
			req := reflect.New(meta.BodyType).Interface()
			if !Bind(c, req) {
				c.Abort()
				return
			}
			c.Set("fastapify_body", req)
		}
		c.Next()
	}

	handlers := make([]gin.HandlerFunc, 0, len(middleware)+2)
	handlers = append(handlers, middleware...)
	handlers = append(handlers, autoValidator)
	handlers = append(handlers, handler)
	w.Engine.Handle(method, ginPath, handlers...)

	return &RouteBuilder{wrapper: w, index: routeIdx}
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
