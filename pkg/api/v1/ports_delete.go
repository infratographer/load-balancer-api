package api

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
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

	mods = append(
		mods,
		qm.Load("Assignments"),
	)

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

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		r.logger.Error("error starting transaction", zap.Error(err))
		return v1InternalServerErrorResponse(c, err)
	}

	assignments, err := r.cleanupPortAssignments(ctx, tx, port)
	if err != nil {
		r.logger.Error("error cleaning up assignments for pool, rolling back", zap.Error(err))

		if err := tx.Rollback(); err != nil {
			r.logger.Error("error rolling back transaction", zap.Error(err))
			return v1InternalServerErrorResponse(c, err)
		}

		return v1InternalServerErrorResponse(c, err)
	}

	if _, err := port.Delete(ctx, tx, false); err != nil {
		r.logger.Error("error deleting port, rolling back", zap.Error(err))

		if err := tx.Rollback(); err != nil {
			r.logger.Error("error rolling back transaction", zap.Error(err))
			return v1InternalServerErrorResponse(c, err)
		}

		return v1InternalServerErrorResponse(c, err)
	}

	if err := tx.Commit(); err != nil {
		r.logger.Error("failed to commit transaction", zap.Error(err))
		return v1InternalServerErrorResponse(c, err)
	}

	msg, err := pubsub.NewMessage(
		pubsub.NewTenantURN(loadBalancer.TenantID),
		pubsub.WithActorURN(someTestJWTURN),
		pubsub.WithSubjectURN(
			pubsub.NewPortURN(port.PortID),
		),
		pubsub.WithAdditionalSubjectURNs(
			append(assignments, pubsub.NewLoadBalancerURN(loadBalancer.LoadBalancerID))...,
		),
		pubsub.WithSubjectFields(map[string]string{"tenant_id": loadBalancer.TenantID}),
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

func (r *Router) cleanupPortAssignments(ctx context.Context, exec boil.ContextExecutor, port *models.Port) ([]string, error) {
	assignmentUrns := []string{}

	// delete assignments
	for _, assignment := range port.R.Assignments {
		r.logger.Debug("deleting assignment for port",
			zap.String("port.id", port.PortID),
			zap.String("assignment.id", assignment.AssignmentID),
		)

		if _, err := assignment.Delete(ctx, exec, false); err != nil {
			return nil, err
		}

		assignmentUrns = append(assignmentUrns, pubsub.NewAssignmentURN(assignment.AssignmentID))
	}

	return assignmentUrns, nil
}
