package api

import (
	"github.com/labstack/echo/v4"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"go.infratographer.com/loadbalancerapi/internal/models"
)

// assignmentsCreate handles the POST /assignments route
func (r *Router) assignmentsCreate(c echo.Context) error {
	ctx := c.Request().Context()

	payload := []struct {
		FrontendID     string `json:"frontend_id"`
		LoadBalancerID string `json:"load_balancer_id"`
		PoolID         string `json:"pool_id"`
	}{}

	if err := c.Bind(&payload); err != nil {
		r.logger.Errorw("error binding payload", "error", err)
		return v1BadRequestResponse(c, err)
	}

	tenantID, err := r.parseTenantID(c)
	if err != nil {
		return err
	}

	assignments := models.AssignmentSlice{}

	tx, err := r.db.Begin()
	if err != nil {
		r.logger.Errorw("error starting transaction", "error", err)
		return v1InternalServerErrorResponse(c, err)
	}

	for _, p := range payload {
		assignment := models.Assignment{
			TenantID:       tenantID,
			FrontendID:     p.FrontendID,
			LoadBalancerID: p.LoadBalancerID,
			PoolID:         p.PoolID,
		}

		assignments = append(assignments, &assignment)

		if err := assignment.Insert(ctx, tx, boil.Infer()); err != nil {
			r.logger.Errorw("error inserting assignment", "error", err)

			if err := tx.Rollback(); err != nil {
				r.logger.Errorw("error rolling back transaction", "error", err)
				return v1InternalServerErrorResponse(c, err)
			}

			return v1InternalServerErrorResponse(c, err)
		}
	}

	switch len(assignments) {
	case 0:
		_ = tx.Rollback()
		return v1NotFoundResponse(c)
	default:
		if err := tx.Commit(); err != nil {
			r.logger.Errorw("error committing transaction", "error", err)
			return v1InternalServerErrorResponse(c, err)
		}

		return v1CreatedResponse(c)
	}
}
