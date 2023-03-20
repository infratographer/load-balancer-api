package api

import (
	"errors"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"go.uber.org/zap"

	"go.infratographer.com/load-balancer-api/internal/models"
)

// poolsParamsBinding return a set of query mods based
// on the query parameters and path parameters
func (r *Router) poolsParamsBinding(c echo.Context, relationships ...string) ([]qm.QueryMod, error) {
	mods := []qm.QueryMod{}

	// optional tenant_id in the request path
	if tenantID, err := r.parseUUID(c, "tenant_id"); err != nil {
		if !errors.Is(err, ErrUUIDNotFound) {
			return nil, err
		}
	} else {
		// found tenant_id in path so add to query mods
		mods = append(mods, models.PoolWhere.TenantID.EQ(tenantID))
		r.logger.Debug("path param", zap.String("tenant_id", tenantID))
	}

	poolID := c.Param("pool_id")
	if poolID != "" {
		if _, err := uuid.Parse(poolID); err != nil {
			return nil, ErrInvalidUUID
		}

		mods = append(mods, models.PoolWhere.PoolID.EQ(poolID))
		r.logger.Debug("path param", zap.String("pool_id", poolID))
	}

	queryParams := []string{"slug", "protocol", "display_name"}

	qpb := echo.QueryParamsBinder(c)

	for _, qp := range queryParams {
		mods = queryParamsToQueryMods(qpb, qp, mods)

		if len(c.QueryParam(qp)) > 0 {
			r.logger.Debug("query param", zap.String("query_param", qp), zap.String("param_vale", c.QueryParam(qp)))
		}
	}

	if err := qpb.BindError(); err != nil {
		return nil, err
	}

	for _, rel := range relationships {
		r.logger.Debug("appending relationships to query", zap.String("relationship", rel))
		mods = append(mods, qm.Load(rel))
	}

	return mods, nil
}
