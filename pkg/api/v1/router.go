// Package api provides the API for the load balancers
package api

import (
	"context"

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
func (r *Router) Routes(e *echo.Group) {
	// authenticate a request, not included the v1 group since this has custom
	// authentication as it's accepting external auth
	e.Use(defaultRequestType)

	v1 := e.Group(apiVersion)
	{
		r.addAssignRoutes(v1)
		r.addPortRoutes(v1)
		r.addLoadBalancerRoutes(v1)
		r.addOriginRoutes(v1)
		r.addPoolsRoutes(v1)
		// r.addMetadataRoutes(v1)
	}

	_, err := r.pubsub.AddStream()
	if err != nil {
		r.logger.Fatal("failed to add stream", zap.Error(err))
	}
}

// DatabaseCheck implements echox.CheckFunc and ensure the database connection is established.
func (r *Router) DatabaseCheck(ctx context.Context) error {
	if err := r.db.PingContext(ctx); err != nil {
		r.logger.Error("readiness check db ping failed", zap.Error(err))

		return err
	}

	return nil
}
