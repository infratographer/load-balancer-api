package api

import (
	"github.com/labstack/echo/v4"
	"go.infratographer.com/load-balancer-api/internal/models"
)

// assignmentsGet handles the GET /assignments route
func (r *Router) assignmentsGet(c echo.Context) error {
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
	default:
		return v1Assignments(c, assignments)
	}
}
