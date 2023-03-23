package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"go.infratographer.com/load-balancer-api/internal/models"
	"go.infratographer.com/load-balancer-api/internal/pubsub"
	"go.uber.org/zap"
)

// poolCreate creates a new pool
func (r *Router) poolCreate(c *gin.Context) {
	ctx := c.Request.Context()
	payload := struct {
		Name     string `json:"name"`
		Protocol string `json:"protocol"`
	}{}

	if err := c.Bind(&payload); err != nil {
		r.logger.Error("error binding payload", zap.Error(err))
		v1BadRequestResponse(c, err)

		return
	}

	tenantID, err := r.parseTenantID(c)
	if err != nil {
		r.logger.Error("error parsing tenant id", zap.Error(err))
		v1BadRequestResponse(c, err)

		return
	}

	pool := &models.Pool{
		Name:     payload.Name,
		Protocol: payload.Protocol,
		TenantID: tenantID,
		Slug:     slug.Make(payload.Name),
	}

	if err := validatePool(pool); err != nil {
		r.logger.Error("error validating pool", zap.Error(err))
		v1BadRequestResponse(c, err)

		return
	}

	if err := pool.Insert(ctx, r.db, boil.Infer()); err != nil {
		r.logger.Error("error inserting pool", zap.Error(err))
		v1InternalServerErrorResponse(c, err)

		return
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

	v1PoolCreatedResponse(c, pool.PoolID)
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
