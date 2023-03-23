package api

import (
	"github.com/gin-gonic/gin"
	"go.infratographer.com/load-balancer-api/internal/models"
)

// originsDelete deletes an origin
func (r *Router) originsDelete(c *gin.Context) {
	ctx := c.Request.Context()

	mods, err := r.originsParamsBinding(c)
	if err != nil {
		r.logger.Errorw("error parsing query params", "error", err)
		v1BadRequestResponse(c, err)

		return
	}

	os, err := models.Origins(mods...).All(ctx, r.db)
	if err != nil {
		r.logger.Errorw("error getting origins", "error", err)
		v1InternalServerErrorResponse(c, err)

		return
	}

	switch len(os) {
	case 0:
		v1NotFoundResponse(c)
	case 1:
		if _, err := os[0].Delete(ctx, r.db, false); err != nil {
			r.logger.Errorw("error deleting origin", "error", err)
			v1InternalServerErrorResponse(c, err)

			return
		}

		v1DeletedResponse(c)
	default:
		v1BadRequestResponse(c, ErrAmbiguous)
	}
}
