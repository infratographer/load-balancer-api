package api

import (
	"github.com/gosimple/slug"
	"github.com/labstack/echo/v4"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"go.infratographer.com/load-balancer-api/internal/models"
	"go.infratographer.com/load-balancer-api/internal/pubsub"
	"go.uber.org/zap"
)

// loadBalancerUpdate updates a load balancer's name, size and type
func (r *Router) loadBalancerUpdate(c echo.Context) error {
	ctx := c.Request().Context()

	mods, err := r.loadBalancerParamsBinding(c)
	if err != nil {
		r.logger.Error("failed to bind params", zap.Error(err))
		return v1BadRequestResponse(c, err)
	}

	lb, err := models.LoadBalancers(mods...).One(ctx, r.db)
	if err != nil {
		return v1InternalServerErrorResponse(c, err)
	}

	payload := struct {
		Name             string `json:"name"`
		LoadBalancerSize string `json:"load_balancer_size"`
		LoadBalancerType string `json:"load_balancer_type"`
		IPAddressID      string `json:"ip_address_id"`
		LocationID       string `json:"location_id"`
	}{}

	if err := c.Bind(&payload); err != nil {
		r.logger.Error("failed to bind load balancer input", zap.Error(err))
		return v1BadRequestResponse(c, err)
	}

	lb.Name = payload.Name
	lb.Slug = slug.Make(payload.Name)
	lb.LoadBalancerSize = payload.LoadBalancerSize
	lb.LoadBalancerType = payload.LoadBalancerType
	lb.IPAddressID = payload.IPAddressID
	lb.LocationID = payload.LocationID
	// TODO do we need to update a CurrentState here?

	return r.updateLoadBalancer(c, lb)
}

func (r *Router) loadBalancerPatch(c echo.Context) error {
	ctx := c.Request().Context()

	mods, err := r.loadBalancerParamsBinding(c)
	if err != nil {
		r.logger.Error("failed to bind params", zap.Error(err))
		return v1BadRequestResponse(c, err)
	}

	lb, err := models.LoadBalancers(mods...).One(ctx, r.db)
	if err != nil {
		return v1InternalServerErrorResponse(c, err)
	}

	payload := struct {
		Name             *string `json:"name"`
		LoadBalancerSize *string `json:"load_balancer_size"`
		LoadBalancerType *string `json:"load_balancer_type"`
		IPAddressID      *string `json:"ip_address_id"`
		LocationID       *string `json:"location_id"`
	}{}

	if err := c.Bind(&payload); err != nil {
		r.logger.Error("failed to bind load balancer input", zap.Error(err))
		return v1BadRequestResponse(c, err)
	}

	if payload.Name != nil {
		lb.Name = *payload.Name
		lb.Slug = slug.Make(*payload.Name)
	}

	if payload.LoadBalancerSize != nil {
		lb.LoadBalancerSize = *payload.LoadBalancerSize
	}

	if payload.LoadBalancerType != nil {
		lb.LoadBalancerType = *payload.LoadBalancerType
	}

	if payload.IPAddressID != nil {
		lb.IPAddressID = *payload.IPAddressID
	}

	if payload.LocationID != nil {
		lb.LocationID = *payload.LocationID
	}

	// TODO do we need to update a CurrentState here?

	return r.updateLoadBalancer(c, lb)
}

func (r *Router) updateLoadBalancer(c echo.Context, lb *models.LoadBalancer) error {
	ctx := c.Request().Context()

	if err := validateLoadBalancer(lb); err != nil {
		r.logger.Error("failed to validate load balancer", zap.Error(err))
		return v1BadRequestResponse(c, err)
	}

	if _, err := lb.Update(ctx, r.db, boil.Infer()); err != nil {
		r.logger.Error("failed to update load balancer", zap.Error(err))
		return v1InternalServerErrorResponse(c, err)
	}

	msg, err := pubsub.NewLoadBalancerMessage(
		someTestJWTURN,
		pubsub.NewTenantURN(lb.TenantID),
		pubsub.NewLoadBalancerURN(lb.LoadBalancerID),
	)
	if err != nil {
		// TODO: add status to reconcile and requeue this
		r.logger.Error("failed to create load balancer message", zap.Error(err))
	}

	if err := r.pubsub.PublishUpdate(ctx, "load-balancer", "global", msg); err != nil {
		// TODO: add status to reconcile and requeue this
		r.logger.Error("failed to publish load balancer message", zap.Error(err))
	}

	return v1UpdateLoadBalancerResponse(c, lb.LoadBalancerID)
}
