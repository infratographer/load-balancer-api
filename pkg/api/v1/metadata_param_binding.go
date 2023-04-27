package api

import (
	"errors"

	"github.com/labstack/echo/v4"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"go.uber.org/zap"

	"go.infratographer.com/load-balancer-api/internal/models"
)

func (r *Router) metadataParamsBinding(c echo.Context) ([]qm.QueryMod, error) {
	var (
		err            error
		tenantID       string
		loadBalancerID string
		metadataID     string
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
		r.logger.Debug("path param", zap.String("tenant_id", tenantID))
	}

	// optional load_balancer_id in the request path
	if loadBalancerID, err = r.parseUUID(c, "load_balancer_id"); err != nil {
		if !errors.Is(err, ErrUUIDNotFound) {
			return nil, err
		}
	} else {
		// found load_balancer_id in path so add to query mods
		mods = append(mods, models.LoadBalancerWhere.LoadBalancerID.EQ(loadBalancerID))
		r.logger.Debug("path param", zap.String("load_balancer_id", loadBalancerID))
	}

	// optional metadata_id in the request path
	if metadataID, err = r.parseUUID(c, "metadata_id"); err != nil {
		if !errors.Is(err, ErrUUIDNotFound) {
			return nil, err
		}
	} else {
		// found metadata_id in path so add to query mods
		mods = append(mods, models.LoadBalancerMetadatumWhere.MetadataID.EQ(metadataID))
		r.logger.Debug("path param", zap.String("metadata_id", metadataID))
	}

	if tenantID == "" && loadBalancerID == "" && metadataID == "" {
		r.logger.Debug("either tenantID or loadBalancerID required in the path")
		return nil, ErrIDRequired
	}
	// query params
	queryParams := []string{"load_balancer_id", "load_balancer_size", "load_balancer_type", "ip_address_id", "location_id", "slug", "name"}

	qpb := echo.QueryParamsBinder(c)

	for _, qp := range queryParams {
		mods = queryParamsToQueryMods(qpb, qp, mods)

		if len(c.QueryParam(qp)) > 0 {
			r.logger.Debug("load balancer query parameters", zap.String("query.key", qp), zap.String("query.value", c.QueryParam(qp)))
		}
	}

	if err = qpb.BindError(); err != nil {
		return nil, err
	}

	return mods, nil
}
