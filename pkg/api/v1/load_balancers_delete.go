package api

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"go.infratographer.com/load-balancer-api/internal/models"
	"go.infratographer.com/load-balancer-api/internal/pubsub"
)

// loadBalancerDelete deletes a load balancer for a tenant
func (r *Router) loadBalancerDelete(c echo.Context) error {
	ctx := c.Request().Context()

	// Look up the load balancer by ID from the path and IP address from the query param
	// this is a unique index in the database, so it will only return one load balancer
	mods, err := r.loadBalancerParamsBinding(c)
	if err != nil {
		r.logger.Errorw("bad request", "error", err)
		return v1BadRequestResponse(c, err)
	}

	lb, err := models.LoadBalancers(mods...).All(ctx, r.db)
	if err != nil {
		r.logger.Errorw("failed to delete load balancer", "error", err)
		return v1InternalServerErrorResponse(c, err)
	}

	switch len(lb) {
	case 0:
		return v1NotFoundResponse(c)
	case 1:
		tx, err := r.db.BeginTx(ctx, nil)
		if err != nil {
			r.logger.Errorw("failed to begin transaction", "error", err)
			return v1InternalServerErrorResponse(c, err)
		}

		if err := r.cleanupLoadBalancer(ctx, lb[0]); err != nil {
			return v1InternalServerErrorResponse(c, err)
		}

		if err := tx.Commit(); err != nil {
			r.logger.Errorw("failed to commit transaction", "error", err)
			return v1InternalServerErrorResponse(c, err)
		}

		msg, err := pubsub.NewLoadBalancerMessage(someTestJWTURN, "urn:infratographer:tenant:"+lb[0].TenantID, pubsub.NewLoadBalancerURN(lb[0].LoadBalancerID))
		if err != nil {
			// TODO: add status to reconcile and requeue this
			r.logger.Errorw("failed to create load balancer message", "error", err)
		}

		if err := r.pubsub.PublishDelete(ctx, "load-balancer", "global", msg); err != nil {
			// TODO: add status to reconcile and requeue this
			r.logger.Errorw("failed to publish load balancer message", "error", err)
		}

		return v1DeletedResponse(c)
	default:
		return v1UnprocessableEntityResponse(c, ErrAmbiguous)
	}
}

// cleanupLoadBalancer deletes all related objects for a load balancer
func (r *Router) cleanupLoadBalancer(ctx context.Context, lb *models.LoadBalancer) error {
	// Delete the load balancer
	if _, err := lb.Delete(ctx, r.db, false); err != nil {
		r.logger.Errorw("failed to delete load balancer", "error", err)
		return err
	}

	// Delete frontends assigned to the load balancer
	if _, err := models.Frontends(qm.Where(models.FrontendColumns.LoadBalancerID+" = ?", lb.LoadBalancerID)).DeleteAll(ctx, r.db, false); err != nil {
		r.logger.Errorw("failed to delete frontends", "error", err)
		return err
	}

	return nil
}
