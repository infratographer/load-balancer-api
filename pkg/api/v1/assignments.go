package api

import (
	"github.com/labstack/echo/v4"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"go.infratographer.com/loadbalancerapi/internal/models"
)

func (r *Router) assignmentParamsBinding(c echo.Context) ([]qm.QueryMod, error) {
	mods := []qm.QueryMod{}

	tenantID, err := r.parseTenantID(c)
	if err != nil {
		return nil, err
	}

	mods = append(mods, models.AssignmentWhere.TenantID.EQ(tenantID))

	queryParams := []string{"frontend_id", "load_balancer_id", "pool_id"}

	qpb := echo.QueryParamsBinder(c)

	for _, qp := range queryParams {
		mods = queryParamsToQueryMods(qpb, qp, mods)

		if len(c.QueryParam(qp)) > 0 {
			r.logger.Debugw("query param", "query_param", qp, "param_vale", c.QueryParam(qp))
		}
	}

	if err = qpb.BindError(); err != nil {
		return nil, err
	}

	return mods, nil
}

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

// assignmentsPost handles the POST /assignments route
func (r *Router) assignmentsPost(c echo.Context) error {
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

// addAssignRoutes adds the assignment routes to the router
func (r *Router) addAssignRoutes(g *echo.Group) {
	g.GET("/assignments", r.assignmentsGet)
	g.POST("/assignments", r.assignmentsPost)
	g.DELETE("/assignments", r.assignmentsDelete)
}
