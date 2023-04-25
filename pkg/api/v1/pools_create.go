package api

import (
	"github.com/gosimple/slug"
	"github.com/labstack/echo/v4"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"go.infratographer.com/load-balancer-api/internal/models"
	"go.infratographer.com/load-balancer-api/internal/pubsub"
	"go.uber.org/zap"
)

// poolCreate creates a new pool
func (r *Router) poolCreate(c echo.Context) error {
	ctx := c.Request().Context()
	payload := struct {
		Name     string `json:"name"`
		Protocol string `json:"protocol"`
		Origins  []struct {
			Disabled bool   `json:"disabled"`
			Name     string `json:"name"`
			Target   string `json:"target"`
			Port     int64  `json:"port"`
		} `json:"origins"`
	}{}

	if err := c.Bind(&payload); err != nil {
		r.logger.Error("error binding payload", zap.Error(err))
		return v1BadRequestResponse(c, err)
	}

	tenantID, err := r.parseUUID(c, "tenant_id")
	if err != nil {
		r.logger.Error("error parsing tenant id", zap.Error(err))
		return v1BadRequestResponse(c, err)
	}

	pool := &models.Pool{
		Name:     payload.Name,
		Protocol: payload.Protocol,
		TenantID: tenantID,
		Slug:     slug.Make(payload.Name),
	}

	if err := validatePool(pool); err != nil {
		r.logger.Error("error validating pool", zap.Error(err))

		return v1BadRequestResponse(c, err)
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		r.logger.Error("failed to begin transaction", zap.Error(err))
		return v1InternalServerErrorResponse(c, err)
	}

	if err := pool.Insert(ctx, tx, boil.Infer()); err != nil {
		r.logger.Error("failed to create pool, rolling back transaction", zap.Error(err))

		if err := tx.Rollback(); err != nil {
			r.logger.Error("error rolling back transaction", zap.Error(err))
			return v1InternalServerErrorResponse(c, err)
		}

		return v1InternalServerErrorResponse(c, err)
	}

	additionalURNs := []string{}

	for _, o := range payload.Origins {
		originID, err := r.createOrigin(ctx, tx, pool.PoolID, o.Name, o.Target, o.Port, o.Disabled)
		if err != nil {
			r.logger.Error("failed to create pool origin, rolling back transaction", zap.Error(err))

			if err := tx.Rollback(); err != nil {
				r.logger.Error("error rolling back transaction", zap.Error(err))
				return v1InternalServerErrorResponse(c, err)
			}

			return v1BadRequestResponse(c, err)
		}

		additionalURNs = append(additionalURNs, pubsub.NewPortURN(originID))
	}

	if err := tx.Commit(); err != nil {
		r.logger.Error("failed to commit transaction", zap.Error(err))
		return v1InternalServerErrorResponse(c, err)
	}

	msg, err := pubsub.NewMessage(
		pubsub.NewTenantURN(tenantID),
		pubsub.WithActorURN(someTestJWTURN),
		pubsub.WithSubjectURN(
			pubsub.NewPoolURN(pool.PoolID),
		),
		pubsub.WithAdditionalSubjectURNs(
			additionalURNs...,
		),
		pubsub.WithSubjectFields(map[string]string{"tenant_id": tenantID}),
	)
	if err != nil {
		// TODO: add status to reconcile and requeue this
		r.logger.Error("error creating pool message", zap.Error(err))
	}

	if err := r.pubsub.PublishCreate(ctx, "load-balancer-pool", "global", msg); err != nil {
		// TODO: add status to reconcile and requeue this
		r.logger.Error("error publishing pool event", zap.Error(err))
	}

	return v1PoolCreatedResponse(c, pool.PoolID)
}

// validatePool validates a pool
func validatePool(p *models.Pool) error {
	if p.Name == "" {
		return ErrNameMissing
	}

	if p.Protocol == "" {
		p.Protocol = "tcp"
	}

	if p.Protocol != "tcp" {
		return ErrPoolProtocolInvalid
	}

	return nil
}
