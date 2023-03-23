package api

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"go.infratographer.com/load-balancer-api/internal/models"
	"go.infratographer.com/load-balancer-api/internal/pubsub"
	"go.uber.org/zap"
)

// frontendCreate creates a new frontend
func (r *Router) frontendCreate(c *gin.Context) {
	ctx := c.Request.Context()

	payload := struct {
		Name string `json:"name"`
		Port int64  `json:"port"`
	}{}
	if err := c.Bind(&payload); err != nil {
		r.logger.Error("failed to bind frontend create input", zap.Error(err))
		v1BadRequestResponse(c, err)

		return
	}

	loadBalancerID, err := r.parseLoadBalancerID(c)
	if err != nil {
		r.logger.Error("bad request", zap.Error(err))
		v1BadRequestResponse(c, err)

		return
	}

	loadBalancer, err := models.LoadBalancers(
		models.LoadBalancerWhere.LoadBalancerID.EQ(loadBalancerID),
	).One(ctx, r.db)
	if err != nil {
		r.logger.Error("error looking up load balancer", zap.Error(err))
		v1BadRequestResponse(c, err)

		return
	}

	frontend := models.Frontend{
		Name:           payload.Name,
		Port:           payload.Port,
		LoadBalancerID: loadBalancer.LoadBalancerID,
		Slug:           slug.Make(payload.Name),
		CurrentState:   "pending",
	}

	if err := validateFrontend(&frontend); err != nil {
		r.logger.Error("failed to validate frontend", zap.Error(err))
		v1BadRequestResponse(c, err)

		return
	}

	if err := frontend.Insert(ctx, r.db, boil.Infer()); err != nil {
		r.logger.Error("failed to insert frontend", zap.Error(err))
		v1InternalServerErrorResponse(c, err)

		return
	}

	msg, err := pubsub.NewFrontendMessage(
		someTestJWTURN,
		pubsub.NewTenantURN(loadBalancer.TenantID),
		pubsub.NewFrontendURN(frontend.FrontendID),
		pubsub.NewLoadBalancerURN(loadBalancer.LoadBalancerID),
	)
	if err != nil {
		// TODO: add status to reconcile and requeue this
		r.logger.Error("failed to create load balancer message", zap.Error(err))
	}

	if err := r.pubsub.PublishCreate(ctx, "load-balancer-frontend", "global", msg); err != nil {
		// TODO: add status to reconcile and requeue this
		r.logger.Error("failed to publish load balancer frontend message", zap.Error(err))
	}

	v1FrontendCreatedResponse(c, frontend.FrontendID)
}

// validateFrontend validates a frontend
func validateFrontend(frontend *models.Frontend) error {
	if frontend.Port < 1 || frontend.Port > 65535 {
		return ErrPortOutOfRange
	}

	if frontend.LoadBalancerID == "" {
		return ErrLoadBalancerIDMissing
	}

	if _, err := uuid.Parse(frontend.LoadBalancerID); err != nil {
		return ErrInvalidUUID
	}

	if frontend.Name == "" {
		// TODO: generate a display name
		return ErrNameMissing
	}

	return nil
}
