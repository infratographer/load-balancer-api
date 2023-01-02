package api

import (
	"github.com/labstack/echo/v4"
	"go.infratographer.com/loadbalancerapi/internal/models"
)

// poolsGet returns a list of pools
func (r *Router) poolsGet(c echo.Context) error {
	ctx := c.Request().Context()

	mods, err := r.poolsParamsBinding(c)
	if err != nil {
		r.logger.Errorw("error parsing query params", "error", err)
		return v1BadRequestResponse(c, err)
	}

	ps, err := models.Pools(mods...).All(ctx, r.db)
	if err != nil {
		r.logger.Errorw("error getting pools", "error", err)
		return v1InternalServerErrorResponse(c, err)
	}

	switch len(ps) {
	case 0:
		return v1NotFoundResponse(c)
	default:
		return v1PoolsResponse(c, ps)
	}
}
