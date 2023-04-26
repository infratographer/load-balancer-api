package api

import (
	"database/sql"
	"errors"

	"github.com/labstack/echo/v4"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"go.infratographer.com/load-balancer-api/internal/models"
	"go.infratographer.com/load-balancer-api/internal/pubsub"
	"go.uber.org/zap"
)

func (r *Router) poolUpdate(c echo.Context) error {
	ctx := c.Request().Context()

	poolID, err := r.parseUUID(c, "pool_id")
	if err != nil {
		return v1BadRequestResponse(c, err)
	}

	payload := struct {
		Name     string `json:"name"`
		Protocol string `json:"protocol"`
	}{}

	if err := c.Bind(&payload); err != nil {
		r.logger.Error("failed to bind pool update input", zap.Error(err))
		return v1BadRequestResponse(c, err)
	}

	mods := []qm.QueryMod{
		models.PoolWhere.PoolID.EQ(poolID),
	}

	pool, err := models.Pools(mods...).One(ctx, r.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return v1NotFoundResponse(c)
		}

		return v1InternalServerErrorResponse(c, err)
	}

	pool.Name = payload.Name
	pool.Protocol = payload.Protocol

	return r.updatePool(c, pool)
}

func (r *Router) poolPatch(c echo.Context) error {
	ctx := c.Request().Context()

	poolID, err := r.parseUUID(c, "pool_id")
	if err != nil {
		return v1BadRequestResponse(c, err)
	}

	payload := struct {
		Name     *string `json:"name"`
		Protocol *string `json:"protocol"`
	}{}

	if err := c.Bind(&payload); err != nil {
		r.logger.Error("failed to bind pool update input", zap.Error(err))
		return v1BadRequestResponse(c, err)
	}

	mods := []qm.QueryMod{
		models.PoolWhere.PoolID.EQ(poolID),
	}

	pool, err := models.Pools(mods...).One(ctx, r.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return v1NotFoundResponse(c)
		}

		return v1InternalServerErrorResponse(c, err)
	}

	if payload.Name != nil {
		pool.Name = *payload.Name
	}

	if payload.Protocol != nil {
		pool.Protocol = *payload.Protocol
	}

	return r.updatePool(c, pool)
}

func (r *Router) updatePool(c echo.Context, pool *models.Pool) error {
	ctx := c.Request().Context()

	if err := validatePool(pool); err != nil {
		return v1BadRequestResponse(c, err)
	}

	if _, err := pool.Update(ctx, r.db, boil.Infer()); err != nil {
		r.logger.Error("failed to update pool", zap.Error(err))
		return v1InternalServerErrorResponse(c, err)
	}

	msg, err := pubsub.NewMessage(
		pubsub.NewTenantURN(pool.TenantID),
		pubsub.WithActorURN(someTestJWTURN),
		pubsub.WithSubjectURN(
			pubsub.NewPoolURN(pool.PoolID),
		),
		pubsub.WithSubjectFields(
			map[string]string{
				"tenant_id":  pool.TenantID,
				"tenant_urn": pubsub.NewTenantURN(pool.TenantID),
			},
		),
	)
	if err != nil {
		// TODO: add status to reconcile and requeue this
		r.logger.Error("failed to create load balancer pool message", zap.Error(err))
	}

	if err := r.pubsub.PublishUpdate(ctx, "load-balancer-pool", "global", msg); err != nil {
		// TODO: add status to reconcile and requeue this
		r.logger.Error("failed to publish load balancer pool message", zap.Error(err))
	}

	return v1UpdatePoolResponse(c, pool.PoolID)
}
