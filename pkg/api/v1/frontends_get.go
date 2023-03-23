package api

import (
	"github.com/gin-gonic/gin"
	"go.infratographer.com/load-balancer-api/internal/models"
	"go.uber.org/zap"
)

// frontendList returns a list of frontends
func (r *Router) frontendList(c *gin.Context) {
	ctx := c.Request.Context()

	mods, err := r.frontendParamsBinding(c)
	if err != nil {
		r.logger.Error("failed to bind frontend params", zap.Error(err))
		v1BadRequestResponse(c, err)

		return
	}

	frontends, err := models.Frontends(mods...).All(ctx, r.db)
	if err != nil {
		r.logger.Error("failed to get frontends", zap.Error(err))
		v1InternalServerErrorResponse(c, err)

		return
	}

	v1Frontends(c, frontends)
}

// frontendGet returns a list of frontends for a given load balancer
func (r *Router) frontendGet(c *gin.Context) {
	ctx := c.Request.Context()

	mods, err := r.frontendParamsBinding(c)
	if err != nil {
		r.logger.Error("failed to bind frontend params", zap.Error(err))
		v1BadRequestResponse(c, err)

		return
	}

	frontends, err := models.Frontends(mods...).All(ctx, r.db)
	if err != nil {
		r.logger.Error("failed to get frontends", zap.Error(err))
		v1InternalServerErrorResponse(c, err)

		return
	}

	switch len(frontends) {
	case 0:
		v1NotFoundResponse(c)
	default:
		v1Frontends(c, frontends)
	}
}
