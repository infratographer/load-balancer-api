package api

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"go.infratographer.com/loadbalancerapi/internal/models"
)

// poolsParamsBinding return a set of query mods based
// on the query parameters and path parameters
func (r *Router) poolsParamsBinding(c echo.Context) ([]qm.QueryMod, error) {
	mods := []qm.QueryMod{}

	tenantID, err := r.parseTenantID(c)
	if err != nil {
		return nil, err
	}

	mods = append(mods, models.PoolWhere.TenantID.EQ(tenantID))

	poolID := c.Param("pool_id")
	if poolID != "" {
		mods = append(mods, models.PoolWhere.PoolID.EQ(poolID))
	}

	queryParams := []string{"display_name", "protocol"}

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

// poolsGet returns a list of pools
func (r *Router) poolsGet(c echo.Context) error {
	ctx := c.Request().Context()

	mods, err := r.poolsParamsBinding(c)
	if err != nil {
		r.logger.Errorw("error parsing query params", "error", err)
		return v1BadRequestResponse(c, err)
	}

	ps, err := models.Pools(mods...).All(ctx, r.db)
	if err != nil {
		r.logger.Errorw("error getting pools", "error", err)
		return v1InternalServerErrorResponse(c, err)
	}

	switch len(ps) {
	case 0:
		return v1NotFoundResponse(c)
	default:
		return v1PoolsResponse(c, ps)
	}
}

// poolCreate creates a new pool
func (r *Router) poolCreate(c echo.Context) error {
	ctx := c.Request().Context()
	payload := []struct {
		DisplayName string `json:"display_name"`
		Protocol    string `json:"protocol"`
	}{}

	if err := c.Bind(&payload); err != nil {
		r.logger.Errorw("error binding payload", "error", err)
		return v1BadRequestResponse(c, err)
	}

	tenantID, err := r.parseTenantID(c)
	if err != nil {
		r.logger.Errorw("error parsing tenant id", "error", err)
		return v1BadRequestResponse(c, err)
	}

	ps := models.PoolSlice{}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		r.logger.Errorw("error starting transaction", "error", err)
		return v1InternalServerErrorResponse(c, err)
	}

	for _, p := range payload {
		pool := &models.Pool{
			DisplayName: p.DisplayName,
			Protocol:    p.Protocol,
			TenantID:    tenantID,
		}

		if err := validatePool(pool); err != nil {
			if err := tx.Rollback(); err != nil {
				r.logger.Errorw("error rolling back transaction", "error", err)
				return v1InternalServerErrorResponse(c, err)
			}

			r.logger.Errorw("error validating pool", "error", err)

			return v1BadRequestResponse(c, err)
		}

		ps = append(ps, pool)

		if err := pool.Insert(ctx, tx, boil.Infer()); err != nil {
			_ = tx.Rollback()

			r.logger.Errorw("error inserting pool", "error", err)

			return v1InternalServerErrorResponse(c, err)
		}
	}

	switch len(ps) {
	case 0:
		if err := tx.Rollback(); err != nil {
			r.logger.Errorw("error rolling back transaction", "error", err)
			return v1InternalServerErrorResponse(c, err)
		}

		return v1NotFoundResponse(c)
	default:
		if err := tx.Commit(); err != nil {
			r.logger.Errorw("error committing transaction", "error", err)
			return v1BadRequestResponse(c, err)
		}

		return v1CreatedResponse(c)
	}
}

// poolDelete deletes a pool
func (r *Router) poolDelete(c echo.Context) error {
	ctx := c.Request().Context()

	mods, err := r.poolsParamsBinding(c)
	if err != nil {
		r.logger.Errorw("error parsing query params", "error", err)
		return v1BadRequestResponse(c, err)
	}

	pool, err := models.Pools(mods...).All(ctx, r.db)
	if err != nil {
		r.logger.Errorw("error getting pool", "error", err)
		return v1InternalServerErrorResponse(c, err)
	}

	switch len(pool) {
	case 0:
		return v1NotFoundResponse(c)
	default:
		tx, err := r.db.BeginTx(ctx, nil)
		if err != nil {
			r.logger.Errorw("error starting transaction", "error", err)
			return v1InternalServerErrorResponse(c, err)
		}

		if err := r.cleanUpPool(ctx, pool[0]); err != nil {
			if err := tx.Rollback(); err != nil {
				r.logger.Errorw("error rolling back transaction", "error", err)
				return v1InternalServerErrorResponse(c, err)
			}

			r.logger.Errorw("error cleaning up pool", "error", err)

			return v1InternalServerErrorResponse(c, err)
		}

		if err := tx.Commit(); err != nil {
			r.logger.Errorw("failed to commit transaction", "error", err)
			return v1InternalServerErrorResponse(c, err)
		}

		return v1DeletedResponse(c)
	}
}

func (r *Router) cleanUpPool(ctx context.Context, pool *models.Pool) error {
	// delete all the pool members
	if _, err := pool.Delete(ctx, r.db, false); err != nil {
		return err
	}

	// delete origins
	//TODO

	return nil
}

// validatePool validates a pool
func validatePool(p *models.Pool) error {
	if p.DisplayName == "" {
		return ErrDisplayNameMissing
	}

	if p.Protocol == "" {
		p.Protocol = "tcp"
	}

	if p.Protocol != "tcp" {
		return ErrPoolProtocolInvalid
	}

	return nil
}

// addPoolsRoutes adds the routes for the pools API
func (r *Router) addPoolsRoutes(g *echo.Group) {
	g.GET("/pools", r.poolsGet)
	g.GET("/pools/:pool_id", r.poolsGet)

	g.POST("/pools", r.poolCreate)

	g.DELETE("/pools", r.poolDelete)
	g.DELETE("/pools/:pool_id", r.poolDelete)
}
