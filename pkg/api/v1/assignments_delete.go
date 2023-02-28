package api

import (
	"github.com/labstack/echo/v4"
	"go.infratographer.com/load-balancer-api/internal/models"
	"go.infratographer.com/load-balancer-api/internal/pubsub"
)

// assignmentsDelete handles the DELETE /assignments route
func (r *Router) assignmentsDelete(c echo.Context) error {
	ctx := c.Request().Context()

	mods, err := r.assignmentParamsBinding(c)
	if err != nil {
		r.logger.Errorw("error parsing query params", "error", err)
		return v1BadRequestResponse(c, err)
	}

	assignments, err := models.Assignments(mods...).All(ctx, r.db)
	if err != nil {
		r.logger.Errorw("error getting assignments", "error", err)
		return v1InternalServerErrorResponse(c, err)
	}

	switch len(assignments) {
	case 0:
		return v1NotFoundResponse(c)
	case 1:
		if _, err := assignments[0].Delete(ctx, r.db, false); err != nil {
			r.logger.Errorw("error deleting assignment", "error", err)
			return v1InternalServerErrorResponse(c, err)
		}

		feMods := models.FrontendWhere.FrontendID.EQ(assignments[0].FrontendID)

		feModel, err := models.Frontends(feMods).One(ctx, r.db)
		if err != nil {
			// TODO: add status to reconcile and requeue this
			r.logger.Errorw("error fetching frontend", "error", err)
		}

		msg, err := pubsub.NewAssignmentMessage(someTestJWTURN, "urn:infratographer:infratographer.com:tenant:"+assignments[0].TenantID, pubsub.NewAssignmentURN(assignments[0].AssignmentID), "urn:infratographer:infratographer.com:load-balancer:"+feModel.LoadBalancerID)
		if err != nil {
			// TODO: add status to reconcile and requeue this
			r.logger.Errorw("error creating message", "error", err)
		}

		if err := r.pubsub.PublishDelete(ctx, "assignment", "global", msg); err != nil {
			// TODO: add status to reconcile and requeue this
			r.logger.Errorw("error publishing event", "error", err)
		}

		return v1DeletedResponse(c)

	default:
		return v1BadRequestResponse(c, ErrAmbiguous)
	}
}
