package api

import (
	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"github.com/labstack/echo/v4"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"go.infratographer.com/load-balancer-api/internal/models"
	"go.infratographer.com/load-balancer-api/internal/pubsub"
)

// frontendCreate creates a new frontend
func (r *Router) frontendCreate(c echo.Context) error {
	ctx := c.Request().Context()

	payload := struct {
		DisplayName string `json:"display_name"`
		Port        int64  `json:"port"`
	}{}
	if err := c.Bind(&payload); err != nil {
		r.logger.Errorw("failed to bind frontend create input", "error", err)
		return v1BadRequestResponse(c, err)
	}

	loadBalancerID, err := r.parseUUID(c, "load_balancer_id")
	if err != nil {
		r.logger.Errorw("bad request", "error", err)
		return v1BadRequestResponse(c, err)
	}

	loadBalancer, err := models.LoadBalancers(
		models.LoadBalancerWhere.LoadBalancerID.EQ(loadBalancerID),
	).One(ctx, r.db)
	if err != nil {
		r.logger.Errorw("error looking up load balancer", "error", err)
		return v1BadRequestResponse(c, err)
	}

	frontend := models.Frontend{
		DisplayName:    payload.DisplayName,
		Port:           payload.Port,
		LoadBalancerID: loadBalancer.LoadBalancerID,
		Slug:           slug.Make(payload.DisplayName),
		CurrentState:   "pending",
	}

	if err := validateFrontend(&frontend); err != nil {
		r.logger.Errorw("failed to validate frontend", "error", err)
		return v1BadRequestResponse(c, err)
	}

	if err := frontend.Insert(ctx, r.db, boil.Infer()); err != nil {
		r.logger.Errorw("failed to insert frontend", "error", err)
		return v1InternalServerErrorResponse(c, err)
	}

	msg, err := pubsub.NewFrontendMessage(
		someTestJWTURN,
		pubsub.NewTenantURN(loadBalancer.TenantID),
		pubsub.NewFrontendURN(frontend.FrontendID),
		pubsub.NewLoadBalancerURN(loadBalancer.LoadBalancerID),
	)
	if err != nil {
		// TODO: add status to reconcile and requeue this
		r.logger.Errorw("failed to create load balancer message", "error", err)
	}

	if err := r.pubsub.PublishCreate(ctx, "load-balancer-frontend", "global", msg); err != nil {
		// TODO: add status to reconcile and requeue this
		r.logger.Errorw("failed to publish load balancer frontend message", "error", err)
	}

	return v1FrontendCreatedResponse(c, frontend.FrontendID)
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

	if frontend.DisplayName == "" {
		// TODO: generate a display name
		return ErrDisplayNameMissing
	}

	return nil
}
