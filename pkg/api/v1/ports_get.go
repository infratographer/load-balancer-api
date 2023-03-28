package api

import (
	"github.com/labstack/echo/v4"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"go.infratographer.com/load-balancer-api/internal/models"
	"go.uber.org/zap"
)

// portList returns a list of ports
func (r *Router) portList(c echo.Context) error {
	ctx := c.Request().Context()

	mods, err := r.portParamsBinding(c)
	if err != nil {
		r.logger.Error("failed to bind port params", zap.Error(err))
		return v1BadRequestResponse(c, err)
	}

	mods = append(
		mods,
		qm.Load("Assignments"),
	)

	ports, err := models.Ports(mods...).All(ctx, r.db)
	if err != nil {
		r.logger.Error("failed to get ports", zap.Error(err))
		return v1InternalServerErrorResponse(c, err)
	}

	return v1Ports(c, ports)
}

// portGet returns a list of ports for a given load balancer
func (r *Router) portGet(c echo.Context) error {
	ctx := c.Request().Context()

	mods, err := r.portParamsBinding(c)
	if err != nil {
		r.logger.Error("failed to bind port params", zap.Error(err))
		return v1BadRequestResponse(c, err)
	}

	mods = append(
		mods,
		qm.Load("Assignments"),
	)

	ports, err := models.Ports(mods...).All(ctx, r.db)
	if err != nil {
		r.logger.Error("failed to get ports", zap.Error(err))
		return v1InternalServerErrorResponse(c, err)
	}

	switch len(ports) {
	case 0:
		return v1NotFoundResponse(c)
	default:
		return v1Ports(c, ports)
	}
}
