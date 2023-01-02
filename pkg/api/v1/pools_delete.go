package api

import (
	"context"

	"github.com/labstack/echo/v4"
	"go.infratographer.com/loadbalancerapi/internal/models"
)

// poolDelete deletes a pool
func (r *Router) poolDelete(c echo.Context) error {
	ctx := c.Request().Context()

	mods, err := r.poolsParamsBinding(c)
	if err != nil {
		r.logger.Errorw("error parsing query params", "error", err)
		return v1BadRequestResponse(c, err)
	}

	pool, err := models.Pools(mods...).All(ctx, r.db)
	if err != nil {
		r.logger.Errorw("error getting pool", "error", err)
		return v1InternalServerErrorResponse(c, err)
	}

	switch len(pool) {
	case 0:
		return v1NotFoundResponse(c)
	default:
		tx, err := r.db.BeginTx(ctx, nil)
		if err != nil {
			r.logger.Errorw("error starting transaction", "error", err)
			return v1InternalServerErrorResponse(c, err)
		}

		if err := r.cleanUpPool(ctx, pool[0]); err != nil {
			if err := tx.Rollback(); err != nil {
				r.logger.Errorw("error rolling back transaction", "error", err)
				return v1InternalServerErrorResponse(c, err)
			}

			r.logger.Errorw("error cleaning up pool", "error", err)

			return v1InternalServerErrorResponse(c, err)
		}

		if err := tx.Commit(); err != nil {
			r.logger.Errorw("failed to commit transaction", "error", err)
			return v1InternalServerErrorResponse(c, err)
		}

		return v1DeletedResponse(c)
	}
}

func (r *Router) cleanUpPool(ctx context.Context, pool *models.Pool) error {
	// delete all the pool members
	if _, err := pool.Delete(ctx, r.db, false); err != nil {
		return err
	}

	// delete origins
	// TODO

	return nil
}
