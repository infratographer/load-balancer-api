package api

import (
	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"github.com/labstack/echo/v4"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"go.infratographer.com/load-balancer-api/internal/models"
	"go.infratographer.com/load-balancer-api/internal/pubsub"
	"go.uber.org/zap"
)

// portCreate creates a new port
func (r *Router) portCreate(c echo.Context) error {
	ctx := c.Request().Context()

	payload := struct {
		Name string `json:"name"`
		Port int64  `json:"port"`
	}{}
	if err := c.Bind(&payload); err != nil {
		r.logger.Error("failed to bind port create input", zap.Error(err))
		return v1BadRequestResponse(c, err)
	}

	loadBalancerID, err := r.parseUUID(c, "load_balancer_id")
	if err != nil {
		r.logger.Error("bad request", zap.Error(err))
		return v1BadRequestResponse(c, err)
	}

	loadBalancer, err := models.LoadBalancers(
		models.LoadBalancerWhere.LoadBalancerID.EQ(loadBalancerID),
	).One(ctx, r.db)
	if err != nil {
		r.logger.Error("error looking up load balancer", zap.Error(err))
		return v1BadRequestResponse(c, err)
	}

	port := models.Port{
		Name:           payload.Name,
		Port:           payload.Port,
		LoadBalancerID: loadBalancer.LoadBalancerID,
		Slug:           slug.Make(payload.Name),
		CurrentState:   "pending",
	}

	if err := validatePort(&port); err != nil {
		r.logger.Error("failed to validate port", zap.Error(err))
		return v1BadRequestResponse(c, err)
	}

	if err := port.Insert(ctx, r.db, boil.Infer()); err != nil {
		r.logger.Error("failed to insert port", zap.Error(err))
		return v1InternalServerErrorResponse(c, err)
	}

	msg, err := pubsub.NewPortMessage(
		someTestJWTURN,
		pubsub.NewTenantURN(loadBalancer.TenantID),
		pubsub.NewPortURN(port.PortID),
		pubsub.NewLoadBalancerURN(loadBalancer.LoadBalancerID),
	)
	if err != nil {
		// TODO: add status to reconcile and requeue this
		r.logger.Error("failed to create load balancer message", zap.Error(err))
	}

	if err := r.pubsub.PublishCreate(ctx, "load-balancer-port", "global", msg); err != nil {
		// TODO: add status to reconcile and requeue this
		r.logger.Error("failed to publish load balancer port message", zap.Error(err))
	}

	return v1PortCreatedResponse(c, port.PortID)
}

// validatePort validates a port
func validatePort(port *models.Port) error {
	if port.Port < 1 || port.Port > 65535 {
		return ErrPortOutOfRange
	}

	if port.LoadBalancerID == "" {
		return ErrLoadBalancerIDMissing
	}

	if _, err := uuid.Parse(port.LoadBalancerID); err != nil {
		return ErrInvalidUUID
	}

	if port.Name == "" {
		// TODO: generate a display name
		return ErrNameMissing
	}

	return nil
}
