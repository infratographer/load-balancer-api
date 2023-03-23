// Package api provides the API for the load balancers
package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
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
	logger *zap.SugaredLogger
}

// NewRouter creates a new router for the API
func NewRouter(db *sqlx.DB, l *zap.SugaredLogger, ps *pubsub.Client) *Router {
	return &Router{
		db:     db,
		pubsub: ps,
		logger: l.Named("api"),
	}
}

// func notYet(c echo.Context) error {
// 	return c.JSON(http.StatusOK, map[string]string{"status": "endpoint not implemented yet"})
// }

// func errorHandler(err error, c echo.Context) {
// 	c.Echo().DefaultHTTPErrorHandler(err, c)
// }

// func defaultRequestType(c *gin.Context) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		contentType := c.Request.Header.Get("Content-Type")
// 		if contentType == "" {
// 			c.Request.Header.Set("Content-Type", "application/json")
// 		}

// 		c.Next()
// 	}
// }

// Routes will add the routes for this API version to a router group
func (r *Router) Routes(rg *gin.RouterGroup) {
	// Health endpoints
	rg.GET("/healthz", r.livenessCheck)
	rg.GET("/readyz", r.readinessCheck)
	v1 := rg.Group(apiVersion)
	{
		r.addAssignRoutes(v1)
		r.addFrontendRoutes(v1)
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
func (r *Router) livenessCheck(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]string{
		"status": "UP",
	})
}

// readinessCheck ensures that the server is up and that we are able to process
// requests. currently this only checks the database connection.
func (r *Router) readinessCheck(c *gin.Context) {
	ctx := c.Request.Context()

	if err := r.db.PingContext(ctx); err != nil {
		r.logger.Errorf("readiness check db ping failed", "err", err)

		c.JSON(http.StatusServiceUnavailable, map[string]string{
			"status": "DOWN",
		})

		return
	}

	c.JSON(http.StatusOK, map[string]string{
		"status": "UP",
	})
}
