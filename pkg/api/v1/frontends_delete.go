package api

import (
	"github.com/gin-gonic/gin"
	"go.infratographer.com/load-balancer-api/internal/models"
	"go.infratographer.com/load-balancer-api/internal/pubsub"
	"go.uber.org/zap"
)

// frontendDelete deletes a frontend
func (r *Router) frontendDelete(c *gin.Context) {
	ctx := c.Request.Context()

	mods, err := r.frontendParamsBinding(c)
	if err != nil {
		r.logger.Error("failed to bind frontend params", zap.Error(err))
		v1BadRequestResponse(c, err)

		return
	}

	frontends, err := models.Frontends(mods...).All(ctx, r.db)
	if err != nil {
		r.logger.Error("failed to get frontend", zap.Error(err))
		v1InternalServerErrorResponse(c, err)

		return
	}

	if len(frontends) == 0 {
		v1NotFoundResponse(c)

		return
	} else if len(frontends) != 1 {
		v1BadRequestResponse(c, ErrAmbiguous)

		return
	}

	frontend := frontends[0]

	loadBalancer, err := models.LoadBalancers(
		models.LoadBalancerWhere.LoadBalancerID.EQ(frontend.LoadBalancerID),
	).One(ctx, r.db)
	if err != nil {
		r.logger.Error("error looking up load balancer for frontend", zap.Error(err))
		v1BadRequestResponse(c, err)

		return
	}

	if _, err := frontend.Delete(ctx, r.db, false); err != nil {
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
		r.logger.Error("failed to create load balancer frontend message", zap.Error(err))
	}

	if err := r.pubsub.PublishDelete(ctx, "load-balancer-frontend", "global", msg); err != nil {
		// TODO: add status to reconcile and requeue this
		r.logger.Error("failed to publish load balancer frontend message", zap.Error(err))
	}

	v1DeletedResponse(c)
}
