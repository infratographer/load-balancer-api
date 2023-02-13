package api

import (
	"github.com/labstack/echo/v4"
	"go.infratographer.com/load-balancer-api/internal/models"
)

// originsList returns a list of origins
func (r *Router) originsList(c echo.Context) error {
	ctx := c.Request().Context()

	mods, err := r.originsParamsBinding(c)
	if err != nil {
		r.logger.Errorw("error parsing query params", "error", err)
		return v1BadRequestResponse(c, err)
	}

	os, err := models.Origins(mods...).All(ctx, r.db)
	if err != nil {
		r.logger.Errorw("error getting origins", "error", err)
		return v1InternalServerErrorResponse(c, err)
	}

	return v1OriginsResponse(c, os)
}

// originsGet returns an origin by id
func (r *Router) originsGet(c echo.Context) error {
	ctx := c.Request().Context()

	mods, err := r.originsParamsBinding(c)
	if err != nil {
		r.logger.Errorw("error parsing query params", "error", err)
		return v1BadRequestResponse(c, err)
	}

	os, err := models.Origins(mods...).All(ctx, r.db)
	if err != nil {
		r.logger.Errorw("error getting origins", "error", err)
		return v1InternalServerErrorResponse(c, err)
	}

	switch len(os) {
	case 0:
		return v1NotFoundResponse(c)
	default:
		return v1OriginsResponse(c, os)
	}
}
