package api

import (
	"context"

	"github.com/dspinhirne/netaddr-go"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"go.infratographer.com/loadbalancerapi/internal/models"
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
	queryParams := []string{"load_balancer_size", "load_balancer_type", "ip_addr", "location_id", "display_name"}

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

// loadBalancerGet returns a load balancer for a tenant by ID
func (r *Router) loadBalancerGet(c echo.Context) error {
	ctx := c.Request().Context()

	mods, err := r.loadBalancerParamsBinding(c)
	if err != nil {
		r.logger.Errorw("failed to bind params", "error", err)
		return v1BadRequestResponse(c, err)
	}

	lbs, err := models.LoadBalancers(mods...).All(ctx, r.db)
	if err != nil {
		return v1InternalServerErrorResponse(c, err)
	}

	switch len(lbs) {
	case 0:
		return v1NotFoundResponse(c)
	default:
		return v1LoadBalancers(c, lbs)
	}
}

// loadBalancerCreate creates a new load balancer for a tenant
func (r *Router) loadBalancerCreate(c echo.Context) error {
	ctx := c.Request().Context()

	payload := []struct {
		DisplayName      string `json:"display_name"`
		LoadBalancerSize string `json:"load_balancer_size"`
		LoadBalancerType string `json:"load_balancer_type"`
		IPAddr           string `json:"ip_addr"`
		LocationID       string `json:"location_id"`
	}{}

	if err := c.Bind(&payload); err != nil {
		r.logger.Errorw("failed to bind load balancer input", "error", err)
		return v1BadRequestResponse(c, err)
	}

	// Ensure the tenant ID is a set from the path,this prevents
	// a tenant from creating a load balancer for another tenant
	tenantID, err := r.parseTenantID(c)
	if err != nil {
		r.logger.Errorw("bad request", "error", err)
		return v1BadRequestResponse(c, err)
	}

	lbs := models.LoadBalancerSlice{}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		r.logger.Errorw("failed to begin transaction", "error", err)
		return v1BadRequestResponse(c, err)
	}

	for _, p := range payload {
		lb := &models.LoadBalancer{
			TenantID:         tenantID,
			DisplayName:      p.DisplayName,
			LoadBalancerSize: p.LoadBalancerSize,
			LoadBalancerType: p.LoadBalancerType,
			IPAddr:           p.IPAddr,
			LocationID:       p.LocationID,
		}

		if err := validateLoadBalancer(lb); err != nil {
			_ = tx.Rollback()

			r.logger.Errorw("failed to validate load balancer", "error", err)

			return v1UnprocessableEntityResponse(c, err)
		}

		lbs = append(lbs, lb)

		err = lb.Insert(ctx, tx, boil.Infer())
		if err != nil {
			r.logger.Errorw("failed to create load balancer, rolling back transaction", "error", err)

			if err := tx.Rollback(); err != nil {
				r.logger.Errorw("failed to rollback transaction", "error", err)
				return v1InternalServerErrorResponse(c, err)
			}

			return v1InternalServerErrorResponse(c, err)
		}
	}

	switch len(lbs) {
	case 0:
		if err := tx.Rollback(); err != nil {
			r.logger.Errorw("failed to rollback transaction", "error", err)
			return v1BadRequestResponse(c, err)
		}

		return v1UnprocessableEntityResponse(c, ErrInvalidLoadBalancer)
	default:
		if err := tx.Commit(); err != nil {
			r.logger.Errorw("failed to commit transaction", "error", err)
			return v1BadRequestResponse(c, err)
		}

		return v1CreatedResponse(c)
	}
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

// cleanupLoadBalancer deletes all related objects for a load balancer
func (r *Router) cleanupLoadBalancer(ctx context.Context, lb *models.LoadBalancer) error {
	// Delete the load balancer
	if _, err := lb.Delete(ctx, r.db, false); err != nil {
		r.logger.Errorw("failed to delete load balancer", "error", err)
		return err
	}

	// Deelete frontends assigned to the load balancer
	if _, err := models.Frontends(qm.Where(models.FrontendColumns.LoadBalancerID+" = ?", lb.LoadBalancerID)).DeleteAll(ctx, r.db, false); err != nil {
		r.logger.Errorw("failed to delete frontends", "error", err)
		return err
	}

	return nil
}

// loadBalancerDelete deletes a load balancer for a tenant
func (r *Router) loadBalancerDelete(c echo.Context) error {
	ctx := c.Request().Context()

	// Look up the load balancer by ID from the path and IP address from the query param
	// this is a unique index in the database, so it will only return one load balancer
	mods, err := r.loadBalancerParamsBinding(c)
	if err != nil {
		r.logger.Errorw("bad request", "error", err)
		return v1BadRequestResponse(c, err)
	}

	lb, err := models.LoadBalancers(mods...).All(ctx, r.db)
	if err != nil {
		r.logger.Errorw("failed to delete load balancer", "error", err)
		return v1InternalServerErrorResponse(c, err)
	}

	switch len(lb) {
	case 0:
		return v1NotFoundResponse(c)
	case 1:
		tx, err := r.db.BeginTx(ctx, nil)
		if err != nil {
			r.logger.Errorw("failed to begin transaction", "error", err)
			return v1InternalServerErrorResponse(c, err)
		}

		if err := r.cleanupLoadBalancer(ctx, lb[0]); err != nil {
			return v1InternalServerErrorResponse(c, err)
		}

		if err := tx.Commit(); err != nil {
			r.logger.Errorw("failed to commit transaction", "error", err)
			return v1InternalServerErrorResponse(c, err)
		}

		return v1DeletedResponse(c)
	default:
		return v1UnprocessableEntityResponse(c, ErrAmbiguous)
	}
}

func (r *Router) addLoadBalancerRoutes(g *echo.Group) {
	g.GET("/loadbalancers", r.loadBalancerGet)
	g.GET("/loadbalancers/:load_balancer_id", r.loadBalancerGet)

	g.POST("/loadbalancers", r.loadBalancerCreate)

	g.DELETE("/loadbalancers", r.loadBalancerDelete)
	g.DELETE("/loadbalancers/:load_balancer_id", r.loadBalancerDelete)
}
