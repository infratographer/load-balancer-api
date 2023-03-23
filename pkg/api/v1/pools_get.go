package api

import (
	"github.com/gin-gonic/gin"
	"go.infratographer.com/load-balancer-api/internal/models"
	"go.uber.org/zap"
)

// poolsGet returns a list of pools
func (r *Router) poolsList(c *gin.Context) {
	ctx := c.Request.Context()

	mods, err := r.poolsParamsBinding(c)
	if err != nil {
		r.logger.Error("error parsing query params", zap.Error(err))
		v1BadRequestResponse(c, err)

		return
	}

	ps, err := models.Pools(mods...).All(ctx, r.db)
	if err != nil {
		r.logger.Error("error getting pools", zap.Error(err))
		v1InternalServerErrorResponse(c, err)

		return
	}

	v1PoolsResponse(c, ps)
}

// poolsGet returns a pool by ID
func (r *Router) poolsGet(c *gin.Context) {
	ctx := c.Request.Context()

	mods, err := r.poolsParamsBinding(c)
	if err != nil {
		r.logger.Error("error parsing query params", zap.Error(err))
		v1BadRequestResponse(c, err)

		return
	}

	ps, err := models.Pools(mods...).All(ctx, r.db)
	if err != nil {
		r.logger.Error("error getting pools", zap.Error(err))
		v1InternalServerErrorResponse(c, err)

		return
	}

	switch len(ps) {
	case 0:
		v1NotFoundResponse(c)
	default:
		v1PoolsResponse(c, ps)
	}
}
