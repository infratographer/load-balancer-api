package api

import (
	"github.com/gosimple/slug"
	"github.com/labstack/echo/v4"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"go.infratographer.com/load-balancer-api/internal/models"
	"go.infratographer.com/load-balancer-api/internal/pubsub"
)

// loadBalancerCreate creates a new load balancer for a tenant
func (r *Router) loadBalancerCreate(c echo.Context) error {
	ctx := c.Request().Context()

	payload := struct {
		DisplayName      string `json:"display_name"`
		LoadBalancerSize string `json:"load_balancer_size"`
		LoadBalancerType string `json:"load_balancer_type"`
		IPAddr           string `json:"ip_addr"`
		LocationID       string `json:"location_id"`
	}{}

	if err := c.Bind(&payload); err != nil {
		r.logger.Errorw("failed to bind load balancer input", "error", err)
		return v1BadRequestResponse(c, err)
	}

	// Ensure the tenant ID is a set from the path,this prevents
	// a tenant from creating a load balancer for another tenant
	tenantID, err := r.parseUUID(c, "tenant_id")
	if err != nil {
		r.logger.Errorw("bad request", "error", err)
		return v1BadRequestResponse(c, err)
	}

	lb := &models.LoadBalancer{
		TenantID:         tenantID,
		DisplayName:      payload.DisplayName,
		LoadBalancerSize: payload.LoadBalancerSize,
		LoadBalancerType: payload.LoadBalancerType,
		IPAddr:           payload.IPAddr,
		LocationID:       payload.LocationID,
		Slug:             slug.Make(payload.DisplayName),
		CurrentState:     "provisioning",
	}

	if err := validateLoadBalancer(lb); err != nil {
		r.logger.Errorw("failed to validate load balancer", "error", err)
		return v1BadRequestResponse(c, err)
	}

	err = lb.Insert(ctx, r.db, boil.Infer())
	if err != nil {
		r.logger.Errorw("failed to create load balancer, rolling back transaction", "error", err)
		return v1InternalServerErrorResponse(c, err)
	}

	lbMods := []qm.QueryMod{
		models.LoadBalancerWhere.TenantID.EQ(tenantID),
		models.LoadBalancerWhere.IPAddr.EQ(payload.IPAddr),
	}

	lbModel, err := models.LoadBalancers(lbMods...).One(ctx, r.db)
	if err != nil {
		r.logger.Errorw("failed to retrieve load balancer", "error", err)
		return v1InternalServerErrorResponse(c, err)
	}

	msg, err := pubsub.NewLoadBalancerMessage(
		someTestJWTURN,
		pubsub.NewTenantURN(tenantID),
		pubsub.NewLoadBalancerURN(lbModel.LoadBalancerID),
	)
	if err != nil {
		// TODO: add status to reconcile and requeue this
		r.logger.Errorw("failed to create load balancer message", "error", err)
	}

	if err := r.pubsub.PublishCreate(ctx, "load-balancer", "global", msg); err != nil {
		// TODO: add status to reconcile and requeue this
		r.logger.Errorw("failed to publish load balancer message", "error", err)
	}

	return v1LoadBalancerCreatedResponse(c, lb.LoadBalancerID)
}
