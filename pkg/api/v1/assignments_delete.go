package api

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"go.infratographer.com/load-balancer-api/internal/models"
	"go.infratographer.com/load-balancer-api/internal/pubsub"
	"go.uber.org/zap"
)

// assignmentsDelete handles the DELETE /assignments route
func (r *Router) assignmentsDelete(c echo.Context) error {
	ctx := c.Request().Context()

	mods, err := r.assignmentParamsBinding(c)
	if err != nil {
		r.logger.Error("error parsing query params", zap.Error(err))
		return v1BadRequestResponse(c, err)
	}

	assignments, err := models.Assignments(mods...).All(ctx, r.db)
	if err != nil {
		r.logger.Error("error getting assignments", zap.Error(err))
		return v1InternalServerErrorResponse(c, err)
	}

	switch len(assignments) {
	case 0:
		return v1NotFoundResponse(c)
	case 1:
		if _, err := assignments[0].Delete(ctx, r.db, false); err != nil {
			r.logger.Error("error deleting assignment", zap.Error(err))
			return v1InternalServerErrorResponse(c, err)
		}

		feMods := models.PortWhere.PortID.EQ(assignments[0].PortID)

		feModel, err := models.Ports(feMods).One(ctx, r.db)
		if err != nil {
			// TODO: add status to reconcile and requeue this
			r.logger.Error("error fetching port", zap.Error(err))
		}

		tenantID := assignments[0].TenantID
		assignmentID := assignments[0].AssignmentID

		msg, err := pubsub.NewMessage(
			pubsub.NewTenantURN(tenantID),
			pubsub.WithActorURN(someTestJWTURN),
			pubsub.WithSubjectURN(
				pubsub.NewAssignmentURN(assignmentID),
			),
			pubsub.WithAdditionalSubjectURNs(
				pubsub.NewLoadBalancerURN(feModel.LoadBalancerID),
			),
			pubsub.WithSubjectFields(
				map[string]string{
					"tenant_id":  tenantID,
					"tenant_urn": pubsub.NewTenantURN(tenantID),
				},
			),
		)
		if err != nil {
			// TODO: add status to reconcile and requeue this
			r.logger.Error("error creating message", zap.Error(err))
		}

		if err := r.pubsub.PublishDelete(ctx, "assignment", "global", msg); err != nil {
			// TODO: add status to reconcile and requeue this
			r.logger.Error("error publishing event", zap.Error(err))
		}

		return v1DeletedResponse(c)

	default:
		return v1BadRequestResponse(c, ErrAmbiguous)
	}
}

func (r *Router) deleteAssignment(ctx context.Context, _ boil.ContextExecutor, tenantID, poolID, portID string) (string, error) {
	r.logger.Debug("deleting assignment",
		zap.String("tenant.id", tenantID),
		zap.String("pool.id", poolID),
		zap.String("port.id", portID),
	)

	assignment, err := models.Assignments(
		models.AssignmentWhere.TenantID.EQ(tenantID),
		models.AssignmentWhere.PoolID.EQ(poolID),
		models.AssignmentWhere.PortID.EQ(portID),
	).One(ctx, r.db)
	if err != nil {
		r.logger.Error("error fetching assignment", zap.Error(err))
		return "", err
	}

	if _, err := assignment.Delete(ctx, r.db, false); err != nil {
		r.logger.Error("error deleting assignment", zap.Error(err))
		return "", err
	}

	return assignment.AssignmentID, nil
}
