package api

import (
	"github.com/labstack/echo/v4"
	"go.infratographer.com/load-balancer-api/internal/models"
	"go.uber.org/zap"
)

// frontendList returns a list of frontends
func (r *Router) frontendList(c echo.Context) error {
	ctx := c.Request().Context()

	mods, err := r.frontendParamsBinding(c)
	if err != nil {
		r.logger.Error("failed to bind frontend params", zap.Error(err))
		return v1BadRequestResponse(c, err)
	}

	frontends, err := models.Frontends(mods...).All(ctx, r.db)
	if err != nil {
		r.logger.Error("failed to get frontends", zap.Error(err))
		return v1InternalServerErrorResponse(c, err)
	}

	return v1Frontends(c, frontends)
}

// frontendGet returns a list of frontends for a given load balancer
func (r *Router) frontendGet(c echo.Context) error {
	ctx := c.Request().Context()

	mods, err := r.frontendParamsBinding(c)
	if err != nil {
		r.logger.Error("failed to bind frontend params", zap.Error(err))
		return v1BadRequestResponse(c, err)
	}

	frontends, err := models.Frontends(mods...).All(ctx, r.db)
	if err != nil {
		r.logger.Error("failed to get frontends", zap.Error(err))
		return v1InternalServerErrorResponse(c, err)
	}

	switch len(frontends) {
	case 0:
		return v1NotFoundResponse(c)
	default:
		return v1Frontends(c, frontends)
	}
}
