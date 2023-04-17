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
		Name  string   `json:"name"`
		Port  int64    `json:"port"`
		Pools []string `json:"pools"`
	}{}
	if err := c.Bind(&payload); err != nil {
		r.logger.Error("failed to bind payload to port creation input", zap.Error(err))
		return v1BadRequestResponse(c, err)
	}

	loadBalancerID, err := r.parseUUID(c, "load_balancer_id")
	if err != nil {
		r.logger.Error("bad request", zap.Error(err))
		return v1BadRequestResponse(c, err)
	}

	lb, err := models.LoadBalancers(
		models.LoadBalancerWhere.LoadBalancerID.EQ(loadBalancerID),
	).One(ctx, r.db)
	if err != nil {
		r.logger.Error("error looking up load balancer", zap.Error(err))
		return v1BadRequestResponse(c, err)
	}

	port := models.Port{
		Name:           payload.Name,
		Port:           payload.Port,
		LoadBalancerID: lb.LoadBalancerID,
		Slug:           slug.Make(payload.Name),
		CurrentState:   "pending",
	}

	if err := validatePort(&port); err != nil {
		r.logger.Error("failed to validate port", zap.Error(err))
		return v1BadRequestResponse(c, err)
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		r.logger.Error("failed to begin transaction", zap.Error(err))
		return v1InternalServerErrorResponse(c, err)
	}

	if err := port.Insert(ctx, tx, boil.Infer()); err != nil {
		r.logger.Error("failed to create port, rolling back transaction", zap.Error(err))

		if err := tx.Rollback(); err != nil {
			r.logger.Error("error rolling back transaction", zap.Error(err))
			return v1InternalServerErrorResponse(c, err)
		}

		return v1InternalServerErrorResponse(c, err)
	}

	additionalURNs := []string{
		pubsub.NewLoadBalancerURN(lb.LoadBalancerID),
	}

	for _, poolID := range payload.Pools {
		if _, err := uuid.Parse(poolID); err != nil {
			r.logger.Error("invalid uuid in port payload", zap.Error(err))
			return v1BadRequestResponse(c, err)
		}

		assignmentID, err := r.createAssignment(ctx, tx, lb.TenantID, poolID, port.PortID)
		if err != nil {
			r.logger.Error("failed to create port assignment, rolling back transaction", zap.Error(err))

			if err := tx.Rollback(); err != nil {
				r.logger.Error("error rolling back transaction", zap.Error(err))
				return v1InternalServerErrorResponse(c, err)
			}

			return v1BadRequestResponse(c, err)
		}

		additionalURNs = append(additionalURNs, pubsub.NewAssignmentURN(assignmentID))
	}

	if err := tx.Commit(); err != nil {
		r.logger.Error("failed to commit transaction", zap.Error(err))
		return v1InternalServerErrorResponse(c, err)
	}

	r.logger.Info("created new port for load balancer",
		zap.String("port.id", port.PortID),
		zap.String("loadbalancer.id", lb.LoadBalancerID),
	)

	msg, err := pubsub.NewPortMessage(
		someTestJWTURN,
		pubsub.NewTenantURN(lb.TenantID),
		pubsub.NewPortURN(port.PortID),
		additionalURNs...,
	)
	if err != nil {
		// TODO: add status to reconcile and requeue this
		r.logger.Error("failed to create load balancer port message", zap.Error(err))
	}

	if err := r.pubsub.PublishCreate(ctx, "load-balancer-port", "global", msg); err != nil {
		// TODO: add status to reconcile and requeue this
		r.logger.Error("failed to publish load balancer port message", zap.Error(err))
	}

	return v1PortCreatedResponse(c, port.PortID)
}

// validatePort validates a port
func validatePort(p *models.Port) error {
	if p.Port == 0 {
		return ErrMissingPortValue
	}

	if p.Port < 1 || p.Port > 65535 {
		return ErrPortOutOfRange
	}

	if p.LoadBalancerID == "" {
		return ErrLoadBalancerIDMissing
	}

	if _, err := uuid.Parse(p.LoadBalancerID); err != nil {
		return ErrInvalidUUID
	}

	if p.Name == "" {
		// TODO: generate a display name
		return ErrNameMissing
	}

	return nil
}
