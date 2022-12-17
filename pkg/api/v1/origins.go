package api

import (
	"github.com/labstack/echo/v4"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"go.infratographer.com/loadbalancerapi/internal/models"
)

func (r *Router) originsParamsBinding(c echo.Context) ([]qm.QueryMod, error) {
	mods := []qm.QueryMod{}

	tenantID, err := r.parseTenantID(c)
	if err != nil {
		return nil, err
	}

	mods = append(mods, models.OriginWhere.TenantID.EQ(tenantID))

	originID := c.Param("origin_id")
	if originID != "" {
		mods = append(mods, models.OriginWhere.OriginID.EQ(originID))
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

func (r *Router) originsGet(c echo.Context) error {
	ctx := c.Request().Context()

	mods, err := r.originsParamsBinding(c)
	if err != nil {
		r.logger.Errorw("error parsing query params", "error", err)
		return v1BadRequestResponse(c, err)
	}

	os, err := models.Origins(mods...).All(ctx, r.db)
	if err != nil {
		r.logger.Errorw("error getting origins", "error", err)
		return v1InternalServerErrorResponse(c, err)
	}

	switch len(os) {
	case 0:
		return v1NotFoundResponse(c)
	default:
		return v1OriginsResponse(c, os)
	}

}
