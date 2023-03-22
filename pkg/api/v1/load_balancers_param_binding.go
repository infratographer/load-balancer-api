package api

import (
	"errors"

	"github.com/dspinhirne/netaddr-go"
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

	// optional tenant_id in the request path
	if tenantID, err = r.parseUUID(c, "tenant_id"); err != nil {
		if !errors.Is(err, ErrUUIDNotFound) {
			return nil, err
		}
	} else {
		// found tenant_id in path so add to query mods
		mods = append(mods, models.LoadBalancerWhere.TenantID.EQ(tenantID))
		r.logger.Debugw("path param", "tenant_id", tenantID)
	}

	// optional load_balancer_id in the request path
	if loadBalancerID, err = r.parseUUID(c, "load_balancer_id"); err != nil {
		if !errors.Is(err, ErrUUIDNotFound) {
			return nil, err
		}
	} else {
		// found load_balancer_id in path so add to query mods
		mods = append(mods, models.LoadBalancerWhere.LoadBalancerID.EQ(loadBalancerID))
		r.logger.Debugw("path param", "load_balancer_id", loadBalancerID)
	}

	if tenantID == "" && loadBalancerID == "" {
		r.logger.Debugw("either tenantID or loadBalancerID required in the path")
		return nil, ErrIDRequired
	}
	// query params
	queryParams := []string{"load_balancer_size", "load_balancer_type", "ip_addr", "location_id", "slug", "name"}

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

	if lb.Name == "" {
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
