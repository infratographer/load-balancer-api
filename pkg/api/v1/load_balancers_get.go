package api

import (
	"database/sql"
	"errors"

	"github.com/labstack/echo/v4"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"go.infratographer.com/load-balancer-api/internal/models"
	"go.uber.org/zap"
)

// loadBalancerList returns a list of load balancers for a tenant
func (r *Router) loadBalancerList(c echo.Context) error {
	ctx := c.Request().Context()

	mods, err := r.loadBalancerParamsBinding(c)
	if err != nil {
		r.logger.Error("failed to bind params", zap.Error(err))
		return v1BadRequestResponse(c, err)
	}

	mods = append(mods,
		qm.Load("Ports"),
		qm.Load("Ports.Assignments"),
	)

	lbs, err := models.LoadBalancers(mods...).All(ctx, r.db)
	if err != nil {
		return v1InternalServerErrorResponse(c, err)
	}

	return v1LoadBalancers(c, lbs)
}

// loadBalancerGet returns a load balancer for a tenant by ID
func (r *Router) loadBalancerGet(c echo.Context) error {
	ctx := c.Request().Context()

	loadBalancerID, err := r.parseUUID(c, "load_balancer_id")
	if err != nil {
		return v1BadRequestResponse(c, err)
	}

	mods := []qm.QueryMod{
		models.LoadBalancerWhere.LoadBalancerID.EQ(loadBalancerID),
		qm.Load("Ports"),
		qm.Load("Ports.Assignments"),
	}

	lb, err := models.LoadBalancers(mods...).One(ctx, r.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return v1NotFoundResponse(c)
		}

		return v1InternalServerErrorResponse(c, err)
	}

	return v1LoadBalancer(c, lb)
}
