package api

import (
	"errors"

	"github.com/labstack/echo/v4"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"go.uber.org/zap"

	"go.infratographer.com/load-balancer-api/internal/models"
)

// portParamsBinding binds the request path and query params to a slice of query mods
// for use with sqlboiler. It returns an error if an invalid uuid is provided
// for the load_balancer_id or port_id in the request path. It also iterates the
// expected query params and appends them to the slice of query mods if they are present
// in the request.
func (r *Router) portParamsBinding(c echo.Context) ([]qm.QueryMod, error) {
	var (
		err            error
		loadBalancerID string
		portID         string
	)

	mods := []qm.QueryMod{}

	// optional load_balancer_id in the request path
	if loadBalancerID, err = r.parseUUID(c, "load_balancer_id"); err != nil {
		if !errors.Is(err, ErrUUIDNotFound) {
			return nil, err
		}
	} else {
		// found load_balancer_id in path so add to query mods
		mods = append(mods, models.PortWhere.LoadBalancerID.EQ(loadBalancerID))
		r.logger.Debug("path param", zap.String("load_balancer_id", loadBalancerID))
	}

	// optional port_id in the request path
	if portID, err = r.parseUUID(c, "port_id"); err != nil {
		if !errors.Is(err, ErrUUIDNotFound) {
			return nil, err
		}
	} else {
		// found port_id in path so add to query mods
		mods = append(mods, models.PortWhere.PortID.EQ(portID))
		r.logger.Debug("path param", zap.String("port_id", portID))
	}

	// query params
	queryParams := []string{"port", "load_balancer_id", "slug", "af_inet", "port_id", "name"}

	qpb := echo.QueryParamsBinder(c)

	for _, qp := range queryParams {
		mods = queryParamsToQueryMods(qpb, qp, mods)

		if len(c.QueryParam(qp)) > 0 {
			r.logger.Debug("query param", zap.String("query_param", qp), zap.String("param_value", c.QueryParam(qp)))
		}
	}

	return mods, nil
}
