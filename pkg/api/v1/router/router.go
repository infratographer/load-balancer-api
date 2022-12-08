// Package router has a router for dnscontroller
package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"go.hollow.sh/toolbox/ginjwt"
	"go.infratographer.sh/loadbalancerapi/pkg/api/v1/locations"
	"go.uber.org/zap"
)

const (
	// Version is the API version
	Version = "v1"

	// V1URI is the path prefix for all v1 endpoints
	V1URI = "/api/v1"

	// TenantPath is the path prefix for all tenant endpoints
	TenantPath = "/tenant"

	// TenantParam is the path param for all tenant endpoints
	TenantParam = "/:tenant"

	// TenantPrefix is the path prefix for all tenant endpoints
	TenantPrefix = TenantPath + TenantParam
)

// Router provides a router for the v1 API
type Router struct {
	authMW *ginjwt.Middleware
	db     *sqlx.DB
	logger *zap.SugaredLogger
}

// New builds a Router
func New(amw *ginjwt.Middleware, db *sqlx.DB, l *zap.SugaredLogger) *Router {
	locations.SetLogger(l)

	return &Router{authMW: amw, db: db, logger: l}
}

// Routes will add the routes for this API version to a router group
func (r *Router) Routes(rg *gin.RouterGroup) {
	r.addLoadBalancerRoutes(rg)
	r.addLocationRoutes(rg)
}
