package api

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"go.infratographer.com/loadbalancerapi/internal/models"
	"go.opentelemetry.io/otel/attribute"
)

// loadBalancerParamsBinding binds the request path and query params to a slice of query mods
// for use with sqlboiler. It returns an error if the tenant_id is not present in the request
// path or an invalid uuid is provided. It also returns an error if an invalid uuid is provided
// for the load_balancer_id in the request path. It also iterates the expected query params
// and appemds them to the slice of query mods if they are present in the request.
func (r *Router) loadBalancerParamsBinding(c echo.Context) ([]qm.QueryMod, error) {
	var (
		tenantID       string
		loadBalancerID string
	)

	mods := []qm.QueryMod{}
	ppb := echo.PathParamsBinder(c)

	// require tenant_id in the request path
	tenantID, err := parseTenantID(c)
	if err != nil {
		return nil, err
	}

	mods = append(mods, models.LoadBalancerWhere.TenantID.EQ(tenantID))
	r.logger.Debugw("path param", "tenant_id", tenantID)

	// optional load_balancer_id in the request path
	err = ppb.String("load_balancer_id", &tenantID).BindError()
	if err != nil {
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
		mods = queryParamsToQueryMods(qpb, models.TableNames.LoadBalancers, qp, mods)

		if len(c.QueryParam(qp)) > 0 {
			r.logger.Debugw("query param", "query_param", qp, "param_vale", c.QueryParam(qp))
		}
	}

	err = qpb.BindError()
	if err != nil {
		return nil, err
	}

	return mods, nil
}

// loadBalancerGet returns a load balancer for a tenant by ID
func (r *Router) loadBalancerGet(c echo.Context) error {
	ctx, span := tracer.Start(c.Request().Context(), "loadBalancerGet")
	defer span.End()

	span.SetAttributes(attribute.String("route", "loadBalancerGet"))

	mods, err := r.loadBalancerParamsBinding(c)
	if err != nil {
		r.logger.Errorw("bad request", "error", err)
		return c.JSON(http.StatusBadRequest, v1BadRequestResponse(err))
	}

	lbs, err := models.LoadBalancers(mods...).All(ctx, r.db)
	if err != nil {
		return err
	}

	switch len(lbs) {
	case 0:
		return c.JSON(http.StatusNotFound, v1NotFoundResponse())
	case 1:
		return c.JSON(http.StatusOK, v1LoadBalancer(lbs[0]))
	default:
		return c.JSON(http.StatusOK, v1LoadBalancerSlice(lbs))
	}
}

// loadBalancerGetByID returns a load balancer for a tenant by ID
func (r *Router) loadBalancerGetByID(c echo.Context) error {
	ctx, span := tracer.Start(c.Request().Context(), "loadBalancerGetByID")
	defer span.End()

	span.SetAttributes(attribute.String("route", "loadBalancerGetByID"))

	// Look up the load balancer by ID from the path
	lb, err := models.FindLoadBalancer(ctx, r.db, c.Param("load_balancer_id"))
	if err != nil {
		return err
	}

	tenantID, err := parseTenantID(c)
	if err != nil {
		r.logger.Errorw("bad request", "error", err)
		return c.JSON(http.StatusBadRequest, v1BadRequestResponse(err))
	}

	// Ensure the tenant ID matches the tenant ID in the path
	// If not, return a 404. This prevents a tenant from accessing another tenant's load balancer
	if lb.TenantID != tenantID {
		r.logger.Errorw("not found", "db_tenant_id", lb.TenantID, "path_tenant_id", tenantID)
		return c.JSON(http.StatusNotFound, v1NotFoundResponse())
	}

	return c.JSON(http.StatusOK, v1LoadBalancer(lb))
}

// loadBalancerCreate creates a new load balancer for a tenant
func (r *Router) loadBalancerCreate(c echo.Context) error {
	ctx, span := tracer.Start(c.Request().Context(), "loadBalancerCreate")
	defer span.End()

	span.SetAttributes(attribute.String("route", "loadBalancerCreate"))

	type input struct {
		DisplayName      string `json:"display_name"`
		LoadBalancerSize string `json:"load_balancer_size"`
		LoadBalancerType string `json:"load_balancer_type"`
		IPAddr           string `json:"ip_addr"`
		LocationID       string `json:"location_id"`
	}

	type inputSlice []input

	payload := inputSlice{}

	err := c.Bind(&payload)
	if err != nil {
		r.logger.Errorw("bad request", "error", err)
		return c.JSON(http.StatusBadRequest, v1BadRequestResponse(err))
	}

	// Ensure the tenant ID is a set from the path,this prevents
	// a tenant from creating a load balancer for another tenant
	tenantID, err := parseTenantID(c)
	if err != nil {
		r.logger.Errorw("bad request", "error", err)
		return c.JSON(http.StatusBadRequest, v1BadRequestResponse(err))
	}

	lbs := models.LoadBalancerSlice{}

	tx, err := r.db.Begin()
	if err != nil {
		r.logger.Errorw("failed to begin transaction", "error", err)
		return c.JSON(http.StatusBadRequest, v1BadRequestResponse(err))
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

		lbs = append(lbs, lb)

		err = lb.Insert(ctx, tx, boil.Infer())
		if err != nil {
			r.logger.Errorw("failed to create load balancer, rolling back transaction", "error", err)

			_ = tx.Rollback()

			return c.JSON(http.StatusBadRequest, v1BadRequestResponse(err))
		}
	}

	switch len(lbs) {
	case 0:
		_ = tx.Rollback()
		return c.JSON(http.StatusUnprocessableEntity, v1UnprocessableEntityResponse(ErrInvalidLoadBalancer))
	default:
		if err := tx.Commit(); err != nil {
			r.logger.Errorw("failed to commit transaction", "error", err)
			return c.JSON(http.StatusBadRequest, v1BadRequestResponse(err))
		}

		return c.JSON(http.StatusCreated, v1CreatedResponse())
	}
}

// loadBalancerDeleteByID deletes a load balancer for a tenant by ID
func (r *Router) loadBalancerDeleteByID(c echo.Context) error {
	ctx, span := tracer.Start(c.Request().Context(), "loadBalancerDeleteByID")
	defer span.End()

	span.SetAttributes(attribute.String("route", "loadBalancerDeleteByID"))

	tenantID, err := parseTenantID(c)
	if err != nil {
		r.logger.Errorw("bad request", "error", err)
		return c.JSON(http.StatusBadRequest, v1BadRequestResponse(err))
	}

	mods := []qm.QueryMod{
		qm.Where(models.LoadBalancerColumns.LoadBalancerID+" = ?", c.Param("load_balancer_id")),
		qm.Where(models.LoadBalancerColumns.TenantID+" = ?", tenantID),
	}
	// Look up the load balancer by ID from the path
	lb, err := models.LoadBalancers(mods...).One(ctx, r.db)
	if err != nil {
		return c.JSON(http.StatusNotFound, v1NotFoundResponse())
	}

	// Delete the load balancer
	_, err = lb.Delete(ctx, r.db, false)
	if err != nil {
		r.logger.Errorw("failed to delete load balancer", "error", err)
		return c.JSON(http.StatusInternalServerError, v1InternalServerErrorResponse(err))
	}

	return c.JSON(http.StatusNoContent, v1DeletedResponse())
}

// loadBalancerDelete deletes a load balancer for a tenant
func (r *Router) loadBalancerDelete(c echo.Context) error {
	ctx, span := tracer.Start(c.Request().Context(), "loadBalancerDelete")
	defer span.End()

	span.SetAttributes(attribute.String("route", "loadBalancerDelete"))

	// Look up the load balancer by ID from the path and IP address from the query param
	// this is a unique index in the database, so it will only return one load balancer
	mods, err := r.loadBalancerParamsBinding(c)
	if err != nil {
		r.logger.Errorw("bad request", "error", err)
		return c.JSON(http.StatusBadRequest, v1BadRequestResponse(err))
	}

	lb, err := models.LoadBalancers(mods...).All(ctx, r.db)
	if err != nil {
		r.logger.Errorw("failed to delete load balancer", "error", err)
		return err
	}

	switch len(lb) {
	case 0:
		return c.JSON(http.StatusNotFound, v1NotFoundResponse())
	case 1:
		_, err = lb[0].Delete(ctx, r.db, false)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusNoContent, v1DeletedResponse())
	default:
		return c.JSON(http.StatusUnprocessableEntity, v1UnprocessableEntityResponse(ErrAmbiguous))
	}
}

func (r *Router) addLoadBalancerRoutes(g *echo.Group) {
	g.GET("/tenant/:tenant_id/loadbalancers", r.loadBalancerGet)
	g.POST("/tenant/:tenant_id/loadbalancers", r.loadBalancerCreate)
	g.DELETE("/tenant/:tenant_id/loadbalancers", r.loadBalancerDelete)
	g.GET("/tenant/:tenant_id/loadbalancers/:load_balancer_id", r.loadBalancerGetByID)
	g.DELETE("/tenant/:tenant_id/loadbalancers/:load_balancer_id", r.loadBalancerDeleteByID)
}
