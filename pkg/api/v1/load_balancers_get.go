package api

import (
	"github.com/labstack/echo/v4"
	"go.infratographer.com/load-balancer-api/internal/models"
)

// loadBalancerList returns a list of load balancers for a tenant
func (r *Router) loadBalancerList(c echo.Context) error {
	ctx := c.Request().Context()

	mods, err := r.loadBalancerParamsBinding(c)
	if err != nil {
		r.logger.Errorw("failed to bind params", "error", err)
		return v1BadRequestResponse(c, err)
	}

	lbs, err := models.LoadBalancers(mods...).All(ctx, r.db)
	if err != nil {
		return v1InternalServerErrorResponse(c, err)
	}

	return v1LoadBalancers(c, lbs)
}

// loadBalancerGet returns a load balancer for a tenant by ID
func (r *Router) loadBalancerGet(c echo.Context) error {
	ctx := c.Request().Context()

	mods, err := r.loadBalancerParamsBinding(c)
	if err != nil {
		r.logger.Errorw("failed to bind params", "error", err)
		return v1BadRequestResponse(c, err)
	}

	lbs, err := models.LoadBalancers(mods...).All(ctx, r.db)
	if err != nil {
		return v1InternalServerErrorResponse(c, err)
	}

	switch len(lbs) {
	case 0:
		return v1NotFoundResponse(c)
	default:
		return v1LoadBalancers(c, lbs)
	}
}
