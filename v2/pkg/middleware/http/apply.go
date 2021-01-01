package http

import (
	"net/http"
)

// Middleware is http middleware
type Middleware func(http.Handler) http.Handler

// Apply applies a chain of middleware in order
func Apply(handler http.Handler, middlewares ...Middleware) http.Handler {
	if len(middlewares) < 1 {
		return handler
	}
	wrapped := handler
	for i := len(middlewares) - 1; i >= 0; i-- {
		wrapped = middlewares[i](wrapped)
	}
	return wrapped
}
