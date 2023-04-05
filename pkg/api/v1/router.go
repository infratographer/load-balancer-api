// Package api provides the API for the load balancers
package api

import (
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"go.infratographer.com/load-balancer-api/internal/pubsub"
	"go.uber.org/zap"
)

const (
	apiVersion     = "v1"
	someTestJWTURN = "urn:infratographer:identity:some-jwt"
)

// Router provides a router for the API
type Router struct {
	db     *sqlx.DB
	pubsub *pubsub.Client
	logger *zap.Logger
}

// NewRouter creates a new router for the API
func NewRouter(db *sqlx.DB, l *zap.SugaredLogger, ps *pubsub.Client) *Router {
	return &Router{
		db:     db,
		pubsub: ps,
		logger: l.Named("api").Desugar(),
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

	// Health endpoints
	e.GET("/healthz", r.livenessCheck)
	e.GET("/readyz", r.readinessCheck)

	v1 := e.Group(apiVersion)
	{
		r.addAssignRoutes(v1)
		r.addPortRoutes(v1)
		r.addLoadBalancerRoutes(v1)
		r.addOriginRoutes(v1)
		r.addPoolsRoutes(v1)
	}

	_, err := r.pubsub.AddStream()
	if err != nil {
		r.logger.Fatal("failed to add stream", zap.Error(err))
	}
}

// livenessCheck ensures that the server is up and responding
func (r *Router) livenessCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status": "UP",
	})
}

// readinessCheck ensures that the server is up and that we are able to process
// requests. currently this only checks the database connection.
func (r *Router) readinessCheck(c echo.Context) error {
	ctx := c.Request().Context()

	if err := r.db.PingContext(ctx); err != nil {
		r.logger.Error("readiness check db ping failed", zap.Error(err))

		return c.JSON(http.StatusServiceUnavailable, map[string]string{
			"status": "DOWN",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"status": "UP",
	})
}

// sliceCompare takes two slices and returns a map of the slice value as the key
// and an int as the value. the int will be 0 if the key exists in both slices,
// greater than one if it exists in s2 but not s1 and less than one if it exists
// in s1 but not s2.
func sliceCompare(s1, s2 []string) map[string]int {
	m := make(map[string]int)
	for _, p := range s2 {
		m[p]++
	}

	for _, p := range s1 {
		m[p]--
	}

	return m
}
