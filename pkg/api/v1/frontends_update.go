package api

import (
	"database/sql"
	"errors"

	"github.com/labstack/echo/v4"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"go.infratographer.com/load-balancer-api/internal/models"
	"go.infratographer.com/load-balancer-api/internal/pubsub"
	"go.uber.org/zap"
)

// frontendUpdate updates a frontend
func (r *Router) frontendUpdate(c echo.Context) error {
	ctx := c.Request().Context()

	mods, err := r.frontendParamsBinding(c)
	if err != nil {
		r.logger.Error("failed to bind frontend params", zap.Error(err))
		return v1BadRequestResponse(c, err)
	}

	frontend, err := models.Frontends(mods...).One(ctx, r.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return v1NotFoundResponse(c)
		}

		r.logger.Error("failed to get frontend", zap.Error(err))

		return v1InternalServerErrorResponse(c, err)
	}

	loadBalancer, err := models.LoadBalancers(
		models.LoadBalancerWhere.LoadBalancerID.EQ(frontend.LoadBalancerID),
	).One(ctx, r.db)
	if err != nil {
		r.logger.Error("error looking up load balancer for frontend", zap.Error(err))
		return v1BadRequestResponse(c, err)
	}

	payload := struct {
		DisplayName string `json:"display_name"`
		Port        int64  `json:"port"`
	}{}
	if err := c.Bind(&payload); err != nil {
		r.logger.Error("failed to bind frontend update input", zap.Error(err))
		return v1BadRequestResponse(c, err)
	}

	frontend.DisplayName = payload.DisplayName
	frontend.Port = payload.Port
	// TODO do we need to update a CurrentState here?

	if err := validateFrontend(frontend); err != nil {
		return v1BadRequestResponse(c, err)
	}

	if _, err := frontend.Update(ctx, r.db, boil.Infer()); err != nil {
		r.logger.Error("failed to update frontend", zap.Error(err))
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
		r.logger.Error("failed to create load balancer frontend message", zap.Error(err))
	}

	if err := r.pubsub.PublishUpdate(ctx, "load-balancer-frontend", "global", msg); err != nil {
		// TODO: add status to reconcile and requeue this
		r.logger.Error("failed to publish load balancer frontend message", zap.Error(err))
	}

	return v1UpdateFrontendResponse(c, frontend.FrontendID)
}
