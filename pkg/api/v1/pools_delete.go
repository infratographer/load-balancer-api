package api

import (
	"context"
	"database/sql"

	"github.com/gin-gonic/gin"
	"go.infratographer.com/load-balancer-api/internal/models"
	"go.infratographer.com/load-balancer-api/internal/pubsub"
	"go.uber.org/zap"
)

// poolDelete deletes a pool
func (r *Router) poolDelete(c *gin.Context) {
	ctx := c.Request.Context()

	mods, err := r.poolsParamsBinding(c, "Origins")
	if err != nil {
		r.logger.Error("error parsing query params", zap.Error(err))
		v1BadRequestResponse(c, err)

		return
	}

	pools, err := models.Pools(mods...).All(ctx, r.db)
	if err != nil {
		r.logger.Error("error getting pool", zap.Error(err))
		v1InternalServerErrorResponse(c, err)

		return
	}

	if len(pools) == 0 {
		v1NotFoundResponse(c)

		return
	} else if len(pools) != 1 {
		v1BadRequestResponse(c, ErrAmbiguous)

		return
	}

	pool := pools[0]

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		r.logger.Error("error starting transaction", zap.Error(err))
		v1InternalServerErrorResponse(c, err)

		return
	}

	origins, err := r.cleanUpPool(ctx, pool, tx)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			r.logger.Error("error rolling back transaction", zap.Error(err))
			v1InternalServerErrorResponse(c, err)

			return
		}

		r.logger.Error("error cleaning up pool", zap.Error(err))

		v1InternalServerErrorResponse(c, err)

		return
	}

	if err := tx.Commit(); err != nil {
		r.logger.Error("failed to commit transaction", zap.Error(err))
		v1InternalServerErrorResponse(c, err)

		return
	}

	msg, err := pubsub.NewPoolMessage(
		someTestJWTURN,
		pubsub.NewTenantURN(pool.TenantID),
		pubsub.NewPoolURN(pool.PoolID),
		origins...,
	)
	if err != nil {
		// TODO: add status to reconcile and requeue this
		r.logger.Error("error creating pool message", zap.Error(err))
	}

	if err := r.pubsub.PublishDelete(ctx, "load-balancer-pool", "global", msg); err != nil {
		// TODO: add status to reconcile and requeue this
		r.logger.Error("error publishing pool event", zap.Error(err))
	}

	v1DeletedResponse(c)
}

func (r *Router) cleanUpPool(ctx context.Context, pool *models.Pool, tx *sql.Tx) ([]string, error) {
	originUrns := []string{}

	// delete origins
	for _, origin := range pool.R.Origins {
		originUrns = append(originUrns, pubsub.NewOriginURN(origin.OriginID))

		r.logger.Debug("deleting origin for pool",
			zap.String("pool.id", pool.PoolID),
			zap.String("origin.id", origin.OriginID),
		)

		if _, err := origin.Delete(ctx, tx, false); err != nil {
			return nil, err
		}
	}

	r.logger.Debug("deleting pool",
		zap.String("pool.id", pool.PoolID),
	)

	// delete all the pool members
	if _, err := pool.Delete(ctx, tx, false); err != nil {
		return nil, err
	}

	return originUrns, nil
}
