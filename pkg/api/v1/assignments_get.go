package api

import (
	"github.com/gin-gonic/gin"
	"go.infratographer.com/load-balancer-api/internal/models"
)

// assignmentsGet handles the GET /assignments route
func (r *Router) assignmentsGet(c *gin.Context) {
	ctx := c.Request.Context()

	mods, err := r.assignmentParamsBinding(c)
	if err != nil {
		r.logger.Errorw("error parsing query params", "error", err)
		v1BadRequestResponse(c, err)

		return
	}

	assignments, err := models.Assignments(mods...).All(ctx, r.db)
	if err != nil {
		r.logger.Errorw("error getting assignments", "error", err)
		v1InternalServerErrorResponse(c, err)

		return
	}

	switch len(assignments) {
	case 0:
		v1NotFoundResponse(c)
	default:
		v1Assignments(c, assignments)
	}
}
