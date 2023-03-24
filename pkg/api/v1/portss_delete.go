package api

import (
	"github.com/labstack/echo/v4"
	"go.infratographer.com/load-balancer-api/internal/models"
	"go.infratographer.com/load-balancer-api/internal/pubsub"
	"go.uber.org/zap"
)

// portDelete deletes a port
func (r *Router) portDelete(c echo.Context) error {
	ctx := c.Request().Context()

	mods, err := r.portParamsBinding(c)
	if err != nil {
		r.logger.Error("failed to bind port params", zap.Error(err))
		return v1BadRequestResponse(c, err)
	}

	ports, err := models.Ports(mods...).All(ctx, r.db)
	if err != nil {
		r.logger.Error("failed to get port", zap.Error(err))
		return v1InternalServerErrorResponse(c, err)
	}

	if len(ports) == 0 {
		return v1NotFoundResponse(c)
	} else if len(ports) != 1 {
		return v1BadRequestResponse(c, ErrAmbiguous)
	}

	port := ports[0]

	loadBalancer, err := models.LoadBalancers(
		models.LoadBalancerWhere.LoadBalancerID.EQ(port.LoadBalancerID),
	).One(ctx, r.db)
	if err != nil {
		r.logger.Error("error looking up load balancer for port", zap.Error(err))
		return v1BadRequestResponse(c, err)
	}

	if _, err := port.Delete(ctx, r.db, false); err != nil {
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
		r.logger.Error("failed to create load balancer port message", zap.Error(err))
	}

	if err := r.pubsub.PublishDelete(ctx, "load-balancer-port", "global", msg); err != nil {
		// TODO: add status to reconcile and requeue this
		r.logger.Error("failed to publish load balancer port message", zap.Error(err))
	}

	return v1DeletedResponse(c)
}
