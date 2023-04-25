package api

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"go.infratographer.com/load-balancer-api/internal/models"
	"go.infratographer.com/load-balancer-api/internal/pubsub"
	"go.uber.org/zap"
)

// poolDelete deletes a pool
func (r *Router) poolDelete(c echo.Context) error {
	ctx := c.Request().Context()

	mods, err := r.poolsParamsBinding(c, "Origins")
	if err != nil {
		r.logger.Error("error parsing query params", zap.Error(err))
		return v1BadRequestResponse(c, err)
	}

	mods = append(
		mods,
		qm.Load("Assignments"),
	)

	pools, err := models.Pools(mods...).All(ctx, r.db)
	if err != nil {
		r.logger.Error("error getting pool", zap.Error(err))
		return v1InternalServerErrorResponse(c, err)
	}

	if len(pools) == 0 {
		return v1NotFoundResponse(c)
	} else if len(pools) != 1 {
		return v1BadRequestResponse(c, ErrAmbiguous)
	}

	pool := pools[0]

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		r.logger.Error("error starting transaction", zap.Error(err))
		return v1InternalServerErrorResponse(c, err)
	}

	assignments, err := r.cleanupPoolAssignments(ctx, tx, pool)
	if err != nil {
		r.logger.Error("error cleaning up assignments for pool, rolling back", zap.Error(err))

		if err := tx.Rollback(); err != nil {
			r.logger.Error("error rolling back transaction", zap.Error(err))
			return v1InternalServerErrorResponse(c, err)
		}

		return v1InternalServerErrorResponse(c, err)
	}

	origins, err := r.cleanUpOrigins(ctx, tx, pool)
	if err != nil {
		r.logger.Error("error cleaning up origins for pool, rolling back", zap.Error(err))

		if err := tx.Rollback(); err != nil {
			r.logger.Error("error rolling back transaction", zap.Error(err))
			return v1InternalServerErrorResponse(c, err)
		}

		return v1InternalServerErrorResponse(c, err)
	}

	r.logger.Debug("deleting pool", zap.String("pool.id", pool.PoolID))

	// delete all the pool members
	if _, err := pool.Delete(ctx, tx, false); err != nil {
		r.logger.Error("error deleting pool, rolling back", zap.Error(err))

		if err := tx.Rollback(); err != nil {
			r.logger.Error("error rolling back transaction", zap.Error(err))
			return v1InternalServerErrorResponse(c, err)
		}

		return v1InternalServerErrorResponse(c, err)
	}

	if err := tx.Commit(); err != nil {
		r.logger.Error("failed to commit transaction", zap.Error(err))
		return v1InternalServerErrorResponse(c, err)
	}

	msg, err := pubsub.NewMessage(
		pubsub.NewTenantURN(pool.TenantID),
		pubsub.WithActorURN(someTestJWTURN),
		pubsub.WithSubjectURN(
			pubsub.NewPoolURN(pool.PoolID),
		),
		pubsub.WithAdditionalSubjectURNs(
			append(assignments, origins...)...,
		),
		pubsub.WithSubjectFields(map[string]string{"tenant_id": pool.TenantID}),
	)
	if err != nil {
		// TODO: add status to reconcile and requeue this
		r.logger.Error("error creating pool message", zap.Error(err))
	}

	if err := r.pubsub.PublishDelete(ctx, "load-balancer-pool", "global", msg); err != nil {
		// TODO: add status to reconcile and requeue this
		r.logger.Error("error publishing pool event", zap.Error(err))
	}

	return v1DeletedResponse(c)
}

func (r *Router) cleanupPoolAssignments(ctx context.Context, exec boil.ContextExecutor, pool *models.Pool) ([]string, error) {
	assignmentUrns := []string{}

	// delete assignments
	for _, assignment := range pool.R.Assignments {
		r.logger.Debug("deleting assignment for pool",
			zap.String("pool.id", pool.PoolID),
			zap.String("assignment.id", assignment.AssignmentID),
		)

		if _, err := assignment.Delete(ctx, exec, false); err != nil {
			return nil, err
		}

		assignmentUrns = append(assignmentUrns, pubsub.NewAssignmentURN(assignment.AssignmentID))
	}

	return assignmentUrns, nil
}

func (r *Router) cleanUpOrigins(ctx context.Context, exec boil.ContextExecutor, pool *models.Pool) ([]string, error) {
	originUrns := []string{}

	// delete origins
	for _, origin := range pool.R.Origins {
		r.logger.Debug("deleting origin for pool",
			zap.String("pool.id", pool.PoolID),
			zap.String("origin.id", origin.OriginID),
		)

		if _, err := origin.Delete(ctx, exec, false); err != nil {
			return nil, err
		}

		originUrns = append(originUrns, pubsub.NewOriginURN(origin.OriginID))
	}

	return originUrns, nil
}
