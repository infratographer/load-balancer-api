package api

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// RouterOption defines the router option type
type RouterOption func(r *Router)

// WithMiddleware includes the provided middleware with the api
func WithMiddleware(mdw ...echo.MiddlewareFunc) RouterOption {
	return func(r *Router) {
		r.middleware = append(r.middleware, mdw...)
	}
}

// WithLogger sets the logger for the service
func WithLogger(logger *zap.Logger) RouterOption {
	return func(r *Router) {
		r.logger = logger
	}
}
