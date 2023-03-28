package api

import (
	"github.com/labstack/echo/v4"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"go.uber.org/zap"

	"go.infratographer.com/load-balancer-api/internal/models"
)

func (r *Router) assignmentParamsBinding(c echo.Context) ([]qm.QueryMod, error) {
	mods := []qm.QueryMod{}

	tenantID, err := r.parseUUID(c, "tenant_id")
	if err != nil {
		return nil, err
	}

	mods = append(mods, models.AssignmentWhere.TenantID.EQ(tenantID))

	queryParams := []string{"port_id", "pool_id"}

	qpb := echo.QueryParamsBinder(c)

	for _, qp := range queryParams {
		mods = queryParamsToQueryMods(qpb, qp, mods)

		if len(c.QueryParam(qp)) > 0 {
			r.logger.Debug("assignment query parameters", zap.String("query.key", qp), zap.String("query.value", c.QueryParam(qp)))
		}
	}

	if err = qpb.BindError(); err != nil {
		return nil, err
	}

	return mods, nil
}
