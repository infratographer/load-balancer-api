package api

import (
	"github.com/labstack/echo/v4"
	"go.infratographer.com/load-balancer-api/internal/models"
	"go.infratographer.com/load-balancer-api/internal/pubsub"
	"go.uber.org/zap"
)

// frontendDelete deletes a frontend
func (r *Router) frontendDelete(c echo.Context) error {
	ctx := c.Request().Context()

	mods, err := r.frontendParamsBinding(c)
	if err != nil {
		r.logger.Error("failed to bind frontend params", zap.Error(err))
		return v1BadRequestResponse(c, err)
	}

	frontends, err := models.Frontends(mods...).All(ctx, r.db)
	if err != nil {
		r.logger.Error("failed to get frontend", zap.Error(err))
		return v1InternalServerErrorResponse(c, err)
	}

	if len(frontends) == 0 {
		return v1NotFoundResponse(c)
	} else if len(frontends) != 1 {
		return v1BadRequestResponse(c, ErrAmbiguous)
	}

	frontend := frontends[0]

	loadBalancer, err := models.LoadBalancers(
		models.LoadBalancerWhere.LoadBalancerID.EQ(frontend.LoadBalancerID),
	).One(ctx, r.db)
	if err != nil {
		r.logger.Error("error looking up load balancer for frontend", zap.Error(err))
		return v1BadRequestResponse(c, err)
	}

	if _, err := frontend.Delete(ctx, r.db, false); err != nil {
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

	if err := r.pubsub.PublishDelete(ctx, "load-balancer-frontend", "global", msg); err != nil {
		// TODO: add status to reconcile and requeue this
		r.logger.Error("failed to publish load balancer frontend message", zap.Error(err))
	}

	return v1DeletedResponse(c)

}
