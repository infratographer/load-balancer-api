package api

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"github.com/labstack/echo/v4"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"go.infratographer.com/load-balancer-api/internal/models"
	"go.infratographer.com/load-balancer-api/internal/pubsub"
	"go.uber.org/zap"
)

// loadBalancerCreate creates a new load balancer for a tenant
func (r *Router) loadBalancerCreate(c echo.Context) error {
	ctx := c.Request().Context()

	payload := struct {
		Name             string `json:"name"`
		LoadBalancerSize string `json:"load_balancer_size"`
		LoadBalancerType string `json:"load_balancer_type"`
		IPAddressID      string `json:"ip_address_uuid"`
		LocationID       string `json:"location_id"`
	}{}

	if err := c.Bind(&payload); err != nil {
		r.logger.Error("failed to bind load balancer input", zap.Error(err))
		return v1BadRequestResponse(c, err)
	}

	// Ensure the tenant ID is a set from the path,this prevents
	// a tenant from creating a load balancer for another tenant
	tenantID, err := r.parseUUID(c, "tenant_id")
	if err != nil {
		r.logger.Error("bad request", zap.Error(err))
		return v1BadRequestResponse(c, err)
	}

	// TODO get/validate IP address uuid from IPAM - just mock it out for now
	if payload.IPAddressID != "" {
		if _, err := uuid.Parse(payload.IPAddressID); err != nil {
			r.logger.Error("bad ip address uuid in request", zap.Error(err))
			return v1BadRequestResponse(c, err)
		}
	} else {
		u, err := uuid.NewUUID()
		if err != nil {
			return v1BadRequestResponse(c, err)
		}
		payload.IPAddressID = u.String()
	}

	lb := &models.LoadBalancer{
		TenantID:         tenantID,
		Name:             payload.Name,
		LoadBalancerSize: payload.LoadBalancerSize,
		LoadBalancerType: payload.LoadBalancerType,
		IPAddressID:      payload.IPAddressID,
		LocationID:       payload.LocationID,
		Slug:             slug.Make(payload.Name),
		CurrentState:     "provisioning",
	}

	if err := validateLoadBalancer(lb); err != nil {
		r.logger.Error("failed to validate load balancer", zap.Error(err))
		return v1BadRequestResponse(c, err)
	}

	fmt.Printf("inserting lb: %+v", lb)

	err = lb.Insert(ctx, r.db, boil.Infer())
	if err != nil {
		r.logger.Error("failed to create load balancer, rolling back transaction", zap.Error(err))
		return v1InternalServerErrorResponse(c, err)
	}

	msg, err := pubsub.NewLoadBalancerMessage(
		someTestJWTURN,
		pubsub.NewTenantURN(tenantID),
		pubsub.NewLoadBalancerURN(lb.LoadBalancerID),
	)
	if err != nil {
		// TODO: add status to reconcile and requeue this
		r.logger.Error("failed to create load balancer message", zap.Error(err))
	}

	if err := r.pubsub.PublishCreate(ctx, "load-balancer", "global", msg); err != nil {
		// TODO: add status to reconcile and requeue this
		r.logger.Error("failed to publish load balancer message", zap.Error(err))
	}

	return v1LoadBalancerCreatedResponse(c, lb.LoadBalancerID)
}
