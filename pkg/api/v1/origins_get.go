package api

import (
	"github.com/gin-gonic/gin"
	"go.infratographer.com/load-balancer-api/internal/models"
)

// originsList returns a list of origins
func (r *Router) originsList(c *gin.Context) {
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

	v1OriginsResponse(c, os)
}

// originsGet returns an origin by id
func (r *Router) originsGet(c *gin.Context) {
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
	default:
		v1OriginsResponse(c, os)
	}
}
