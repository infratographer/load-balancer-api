package api

import (
	"github.com/dspinhirne/netaddr-go"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"go.infratographer.com/load-balancer-api/internal/models"
)

// loadBalancerParamsBinding binds the request path and query params to a slice of query mods
// for use with sqlboiler. It returns an error if the tenant_id is not present in the request
// path or an invalid uuid is provided. It also returns an error if an invalid uuid is provided
// for the load_balancer_id in the request path. It also iterates the expected query params
// and appends them to the slice of query mods if they are present in the request.
func (r *Router) loadBalancerParamsBinding(c echo.Context) ([]qm.QueryMod, error) {
	var (
		err            error
		tenantID       string
		loadBalancerID string
	)

	mods := []qm.QueryMod{}
	ppb := echo.PathParamsBinder(c)

	// require tenant_id in the request path
	if tenantID, err = r.parseTenantID(c); err != nil {
		return nil, err
	}

	mods = append(mods, models.LoadBalancerWhere.TenantID.EQ(tenantID))
	r.logger.Debugw("path param", "tenant_id", tenantID)

	// optional load_balancer_id in the request path
	if err = ppb.String("load_balancer_id", &tenantID).BindError(); err != nil {
		return nil, err
	}

	if loadBalancerID != "" {
		if _, err := uuid.Parse(loadBalancerID); err != nil {
			return nil, ErrInvalidUUID
		}

		mods = append(mods, models.LoadBalancerWhere.LoadBalancerID.EQ(loadBalancerID))
		r.logger.Debugw("path param", "load_balancer_id", loadBalancerID)
	}
	// query params
	queryParams := []string{"load_balancer_size", "load_balancer_type", "ip_addr", "location_id", "slug"}

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

// validateLoadBalancer validates a load balancer
func validateLoadBalancer(lb *models.LoadBalancer) error {
	if lb.IPAddr == "" {
		return ErrLoadBalancerIPMissing
	}

	if _, err := netaddr.ParseIP(lb.IPAddr); err != nil {
		return ErrLoadBalancerIPInvalid
	}

	if lb.DisplayName == "" {
		return ErrDisplayNameMissing
	}

	if lb.LoadBalancerSize == "" {
		return ErrSizeRequired
	}

	if lb.LoadBalancerType != "layer-3" {
		return ErrTypeInvalid
	}

	if lb.LocationID == "" {
		return ErrLocationIDRequired
	}

	return nil
}
