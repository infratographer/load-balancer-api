package api

import (
	"github.com/labstack/echo/v4"
	"go.infratographer.com/loadbalancerapi/internal/models"
)

// frontendDelete deletes a frontend
func (r *Router) frontendDelete(c echo.Context) error {
	ctx := c.Request().Context()

	mods, err := r.frontendParamsBinding(c)
	if err != nil {
		r.logger.Errorw("failed to bind frontend params", "error", err)
		return v1BadRequestResponse(c, err)
	}

	frontends, err := models.Frontends(mods...).All(ctx, r.db)
	if err != nil {
		return v1InternalServerErrorResponse(c, err)
	}

	switch len(frontends) {
	case 0:
		return v1NotFoundResponse(c)
	case 1:
		if _, err := frontends[0].Delete(ctx, r.db, false); err != nil {
			return v1InternalServerErrorResponse(c, err)
		}

		return v1DeletedResponse(c)
	default:
		return v1BadRequestResponse(c, ErrAmbiguous)
	}
}
