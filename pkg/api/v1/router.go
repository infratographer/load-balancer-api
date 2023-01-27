// Package api provides the API for the load balancers
package api

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

const apiVersion = "v1"

// Router provides a router for the API
type Router struct {
	db     *sqlx.DB
	events nats.JetStreamContext
	logger *zap.SugaredLogger
}

// NewRouter creates a new router for the API
func NewRouter(db *sqlx.DB, l *zap.SugaredLogger, js nats.JetStreamContext) *Router {
	return &Router{
		db:     db,
		events: js,
		logger: l.Named("api"),
	}
}

// func notYet(c echo.Context) error {
// 	return c.JSON(http.StatusOK, map[string]string{"status": "endpoint not implemented yet"})
// }

func errorHandler(err error, c echo.Context) {
	c.Echo().DefaultHTTPErrorHandler(err, c)
}

func defaultRequestType(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		contentType := c.Request().Header.Get("Content-Type")
		if contentType == "" {
			c.Request().Header.Set("Content-Type", "application/json")
		}

		return next(c)
	}
}

// Routes will add the routes for this API version to a router group
func (r *Router) Routes(e *echo.Echo) {
	// authenticate a request, not included the v1 group since this has custom
	// authentication as it's accepting external auth
	e.HideBanner = true

	e.HTTPErrorHandler = errorHandler

	e.Use(defaultRequestType)

	v1 := e.Group(apiVersion)
	{
		r.addAssignRoutes(v1)
		r.addFrontendRoutes(v1)
		r.addLoadBalancerRoutes(v1)
		r.addOriginRoutes(v1)
		r.addPoolsRoutes(v1)
	}

	r.events.AddStream(&nats.StreamConfig{
		Name: "lbapi",
		Subjects: []string{
			"com.infratographer.events.>",
		},
	})
}
