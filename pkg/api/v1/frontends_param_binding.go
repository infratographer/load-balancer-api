package api

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"go.infratographer.com/load-balancer-api/internal/models"
)

// frontendParamsBinding binds the request path and query params to a slice of query mods
// for use with sqlboiler. It returns an error if the tenant_id is not present in the request
// path or an invalid uuid is provided. It also returns an error if an invalid uuid is provided
// for the load_balancer_id in the request path. It also iterates the expected query params
// and appends them to the slice of query mods if they are present in the request.
func (r *Router) frontendParamsBinding(c echo.Context) ([]qm.QueryMod, error) {
	var (
		err      error
		tenantID string
		// loadBalancerID string
		frontendID string
	)

	mods := []qm.QueryMod{}
	ppb := echo.PathParamsBinder(c)

	if tenantID, err = r.parseTenantID(c); err != nil {
		return nil, err
	}

	mods = append(mods, models.FrontendWhere.TenantID.EQ(tenantID))
	r.logger.Debugw("path param", "tenant_id", tenantID)

	// optional frontend_id in the request path
	if err = ppb.String("frontend_id", &frontendID).BindError(); err != nil {
		return nil, err
	}

	if frontendID != "" {
		if _, err := uuid.Parse(frontendID); err != nil {
			return nil, ErrInvalidUUID
		}

		mods = append(mods, models.FrontendWhere.FrontendID.EQ(frontendID))
		r.logger.Debugw("path param", "frontend_id", frontendID)
	}

	// query params
	queryParams := []string{"port", "load_balancer_id", "slug", "af_inet"}

	qpb := echo.QueryParamsBinder(c)

	for _, qp := range queryParams {
		mods = queryParamsToQueryMods(qpb, qp, mods)

		if len(c.QueryParam(qp)) > 0 {
			r.logger.Debugw("query param", "query_param", qp, "param_vale", c.QueryParam(qp))
		}
	}

	return mods, nil
}
