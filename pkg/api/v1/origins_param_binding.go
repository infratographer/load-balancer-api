package api

import (
	"github.com/labstack/echo/v4"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"go.infratographer.com/load-balancer-api/internal/models"
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

	queryParams := []string{"slug", "target", "port"}

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
