package api

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"go.infratographer.com/load-balancer-api/internal/models"
	"go.infratographer.com/load-balancer-api/internal/pubsub"
	"go.uber.org/zap"
)

// loadBalancerDelete deletes a load balancer for a tenant
func (r *Router) loadBalancerDelete(c *gin.Context) {
	ctx := c.Request.Context()

	// Look up the load balancer by ID from the path and IP address from the query param
	// this is a unique index in the database, so it will only return one load balancer
	mods, err := r.loadBalancerParamsBinding(c)
	if err != nil {
		r.logger.Error("bad request", zap.Error(err))
		v1BadRequestResponse(c, err)

		return
	}

	lb, err := models.LoadBalancers(mods...).All(ctx, r.db)
	if err != nil {
		r.logger.Error("failed to delete load balancer", zap.Error(err))
		v1InternalServerErrorResponse(c, err)

		return
	}

	switch len(lb) {
	case 0:
		v1NotFoundResponse(c)
	case 1:
		tx, err := r.db.BeginTx(ctx, nil)
		if err != nil {
			r.logger.Error("failed to begin transaction", zap.Error(err))
			v1InternalServerErrorResponse(c, err)

			return
		}

		if err := r.cleanupLoadBalancer(ctx, lb[0]); err != nil {
			v1InternalServerErrorResponse(c, err)

			return
		}

		if err := tx.Commit(); err != nil {
			r.logger.Error("failed to commit transaction", zap.Error(err))
			v1InternalServerErrorResponse(c, err)

			return
		}

		msg, err := pubsub.NewLoadBalancerMessage(
			someTestJWTURN,
			pubsub.NewLoadBalancerURN(lb[0].TenantID),
			pubsub.NewLoadBalancerURN(lb[0].LoadBalancerID),
		)
		if err != nil {
			// TODO: add status to reconcile and requeue this
			r.logger.Error("failed to create load balancer message", zap.Error(err))
		}

		if err := r.pubsub.PublishDelete(ctx, "load-balancer", "global", msg); err != nil {
			// TODO: add status to reconcile and requeue this
			r.logger.Error("failed to publish load balancer message", zap.Error(err))
		}

		v1DeletedResponse(c)
	default:
		v1UnprocessableEntityResponse(c, ErrAmbiguous)
	}
}

// cleanupLoadBalancer deletes all related objects for a load balancer
func (r *Router) cleanupLoadBalancer(ctx context.Context, lb *models.LoadBalancer) error {
	// Delete the load balancer
	if _, err := lb.Delete(ctx, r.db, false); err != nil {
		r.logger.Error("failed to delete load balancer", zap.Error(err))
		return err
	}

	// Delete frontends assigned to the load balancer
	if _, err := models.Frontends(qm.Where(models.FrontendColumns.LoadBalancerID+" = ?", lb.LoadBalancerID)).DeleteAll(ctx, r.db, false); err != nil {
		r.logger.Error("failed to delete frontends", zap.Error(err))
		return err
	}

	return nil
}
