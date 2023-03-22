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

	if err := pool.Insert(ctx, r.db, boil.Infer()); err != nil {
		r.logger.Error("error inserting pool", zap.Error(err))

		return v1InternalServerErrorResponse(c, err)
	}

	msg, err := pubsub.NewPoolMessage(
		someTestJWTURN,
		pubsub.NewTenantURN(tenantID),
		pubsub.NewPoolURN(pool.PoolID),
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
