package api

import (
	"github.com/gin-gonic/gin"
	"go.infratographer.com/load-balancer-api/internal/models"
	"go.uber.org/zap"
)

// loadBalancerList returns a list of load balancers for a tenant
func (r *Router) loadBalancerList(c *gin.Context) {
	ctx := c.Request.Context()

	mods, err := r.loadBalancerParamsBinding(c)
	if err != nil {
		r.logger.Error("failed to bind params", zap.Error(err))
		v1BadRequestResponse(c, err)

		return
	}

	lbs, err := models.LoadBalancers(mods...).All(ctx, r.db)
	if err != nil {
		v1InternalServerErrorResponse(c, err)

		return
	}

	v1LoadBalancers(c, lbs)
}

// loadBalancerGet returns a load balancer for a tenant by ID
func (r *Router) loadBalancerGet(c *gin.Context) {
	ctx := c.Request.Context()

	mods, err := r.loadBalancerParamsBinding(c)
	if err != nil {
		r.logger.Error("failed to bind params", zap.Error(err))
		v1BadRequestResponse(c, err)

		return
	}

	lbs, err := models.LoadBalancers(mods...).All(ctx, r.db)
	if err != nil {
		v1InternalServerErrorResponse(c, err)

		return
	}

	switch len(lbs) {
	case 0:
		v1NotFoundResponse(c)
	default:
		v1LoadBalancers(c, lbs)
	}
}
