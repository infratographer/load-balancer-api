package api

import (
	"errors"

	"github.com/labstack/echo/v4"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"go.uber.org/zap"

	"go.infratographer.com/load-balancer-api/internal/models"
)

// frontendParamsBinding binds the request path and query params to a slice of query mods
// for use with sqlboiler. It returns an error if an invalid uuid is provided
// for the load_balancer_id or frontend_id in the request path. It also iterates the
// expected query params and appends them to the slice of query mods if they are present
// in the request.
func (r *Router) frontendParamsBinding(c echo.Context) ([]qm.QueryMod, error) {
	var (
		err            error
		loadBalancerID string
		frontendID     string
	)

	mods := []qm.QueryMod{}

	// optional load_balancer_id in the request path
	if loadBalancerID, err = r.parseUUID(c, "load_balancer_id"); err != nil {
		if !errors.Is(err, ErrUUIDNotFound) {
			return nil, err
		}
	} else {
		// found load_balancer_id in path so add to query mods
		mods = append(mods, models.FrontendWhere.LoadBalancerID.EQ(loadBalancerID))
		r.logger.Debug("path param", zap.String("load_balancer_id", loadBalancerID))
	}

	// optional frontend_id in the request path
	if frontendID, err = r.parseUUID(c, "frontend_id"); err != nil {
		if !errors.Is(err, ErrUUIDNotFound) {
			return nil, err
		}
	} else {
		// found frontend_id in path so add to query mods
		mods = append(mods, models.FrontendWhere.FrontendID.EQ(frontendID))
		r.logger.Debug("path param", zap.String("frontend_id", frontendID))
	}

	// query params
	queryParams := []string{"port", "load_balancer_id", "slug", "af_inet"}

	qpb := echo.QueryParamsBinder(c)

	for _, qp := range queryParams {
		mods = queryParamsToQueryMods(qpb, qp, mods)

		if len(c.QueryParam(qp)) > 0 {
			r.logger.Debug("query param", zap.String("query_param", qp), zap.String("param_vale", c.QueryParam(qp)))
		}
	}

	return mods, nil
}
