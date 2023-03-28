package api

import (
	"github.com/labstack/echo/v4"
	"go.infratographer.com/load-balancer-api/internal/models"
	"go.uber.org/zap"
)

// originsDelete deletes an origin
func (r *Router) originsDelete(c echo.Context) error {
	ctx := c.Request().Context()

	mods, err := r.originsParamsBinding(c)
	if err != nil {
		r.logger.Error("error parsing query params", zap.Error(err))
		return v1BadRequestResponse(c, err)
	}

	os, err := models.Origins(mods...).All(ctx, r.db)
	if err != nil {
		r.logger.Error("error getting origins", zap.Error(err))
		return v1InternalServerErrorResponse(c, err)
	}

	switch len(os) {
	case 0:
		return v1NotFoundResponse(c)
	case 1:
		if _, err := os[0].Delete(ctx, r.db, false); err != nil {
			r.logger.Error("error deleting origin", zap.Error(err))
			return v1InternalServerErrorResponse(c, err)
		}

		return v1DeletedResponse(c)
	default:
		return v1BadRequestResponse(c, ErrAmbiguous)
	}
}
