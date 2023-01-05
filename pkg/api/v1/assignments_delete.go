package api

import (
	"github.com/labstack/echo/v4"
	"go.infratographer.com/load-balancer-api/internal/models"
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

		return v1DeletedResponse(c)

	default:
		return v1BadRequestResponse(c, ErrAmbiguous)
	}
}
