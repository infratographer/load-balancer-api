package api

import (
	"database/sql"
	"errors"

	"github.com/labstack/echo/v4"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"go.infratographer.com/load-balancer-api/internal/models"
	"go.uber.org/zap"
)

// poolsGet returns a list of pools
func (r *Router) poolsList(c echo.Context) error {
	ctx := c.Request().Context()

	mods, err := r.poolsParamsBinding(c)
	if err != nil {
		r.logger.Error("error parsing query params", zap.Error(err))
		return v1BadRequestResponse(c, err)
	}

	mods = append(mods,
		qm.Load("Origins"),
	)

	ps, err := models.Pools(mods...).All(ctx, r.db)
	if err != nil {
		r.logger.Error("error getting pools", zap.Error(err))
		return v1InternalServerErrorResponse(c, err)
	}

	return v1PoolsResponse(c, ps)
}

// poolsGet returns a pool by ID
func (r *Router) poolsGet(c echo.Context) error {
	ctx := c.Request().Context()

	poolID, err := r.parseUUID(c, "pool_id")
	if err != nil {
		return v1BadRequestResponse(c, err)
	}

	mods := []qm.QueryMod{
		models.PoolWhere.PoolID.EQ(poolID),
		qm.Load("Origins"),
	}

	p, err := models.Pools(mods...).One(ctx, r.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return v1NotFoundResponse(c)
		}

		return v1InternalServerErrorResponse(c, err)
	}

	return v1PoolResponse(c, p)
}
