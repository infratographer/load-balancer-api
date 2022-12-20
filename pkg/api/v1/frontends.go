package api

import (
	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"github.com/labstack/echo/v4"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"go.infratographer.com/loadbalancerapi/internal/models"
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

// frontendGet returns a list of frontends for a given load balancer
func (r *Router) frontendGet(c echo.Context) error {
	ctx := c.Request().Context()

	mods, err := r.frontendParamsBinding(c)
	if err != nil {
		r.logger.Errorw("failed to bind frontend params", "error", err)
		return v1BadRequestResponse(c, err)
	}

	frontends, err := models.Frontends(mods...).All(ctx, r.db)
	if err != nil {
		return err
	}

	switch len(frontends) {
	case 0:
		return v1NotFoundResponse(c)
	default:
		return v1Frontends(c, frontends)
	}
}

// frontendDelete deletes a frontend
func (r *Router) frontendDelete(c echo.Context) error {
	ctx := c.Request().Context()

	mods, err := r.frontendParamsBinding(c)
	if err != nil {
		r.logger.Errorw("failed to bind frontend params", "error", err)
		return v1BadRequestResponse(c, err)
	}

	frontends, err := models.Frontends(mods...).All(ctx, r.db)
	if err != nil {
		return v1InternalServerErrorResponse(c, err)
	}

	switch len(frontends) {
	case 0:
		return v1NotFoundResponse(c)
	case 1:
		if _, err := frontends[0].Delete(ctx, r.db, false); err != nil {
			return v1InternalServerErrorResponse(c, err)
		}

		return v1DeletedResponse(c)
	default:
		return v1BadRequestResponse(c, ErrAmbiguous)
	}
}

// frontendCreate creates a new frontend
func (r *Router) frontendCreate(c echo.Context) error {
	ctx := c.Request().Context()

	payload := []struct {
		DisplayName    string `json:"display_name"`
		Port           int64  `json:"port"`
		LoadBalancerID string `json:"load_balancer_id"`
	}{}
	if err := c.Bind(&payload); err != nil {
		r.logger.Errorw("failed to bind frontend create input", "error", err)
		return v1BadRequestResponse(c, err)
	}

	tenantID, err := r.parseTenantID(c)
	if err != nil {
		return v1BadRequestResponse(c, err)
	}

	frontends := models.FrontendSlice{}

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	for _, p := range payload {
		frontend := models.Frontend{
			DisplayName:    p.DisplayName,
			Port:           p.Port,
			LoadBalancerID: p.LoadBalancerID,
			TenantID:       tenantID,
			Slug:           slug.Make(p.DisplayName),
			CurrentState:   "pending",
		}

		if err := validateFrontend(&frontend); err != nil {
			_ = tx.Rollback()
			return v1BadRequestResponse(c, err)
		}

		frontends = append(frontends, &frontend)

		if err := frontend.Insert(ctx, tx, boil.Infer()); err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	switch len(frontends) {
	case 0:
		_ = tx.Rollback()
		return v1UnprocessableEntityResponse(c, ErrEmptyPayload)
	default:
		if err := tx.Commit(); err != nil {
			return err
		}

		return v1CreatedResponse(c)
	}
}

// validateFrontend validates a frontend
func validateFrontend(frontend *models.Frontend) error {
	if frontend.Port < 1 || frontend.Port > 65535 {
		return ErrPortOutOfRange
	}

	if frontend.LoadBalancerID == "" {
		return ErrLoadBalancerIDMissing
	}

	if _, err := uuid.Parse(frontend.LoadBalancerID); err != nil {
		return ErrInvalidUUID
	}

	if frontend.DisplayName == "" {
		// TODO: generate a display name
		return ErrDisplayNameMissing
	}

	return nil
}

// addFrontendRoutes adds the frontend routes to the router
func (r *Router) addFrontendRoutes(rg *echo.Group) {
	rg.GET("/frontends", r.frontendGet)
	rg.GET("/frontends/:frontend_id", r.frontendGet)
	rg.GET("/loadbalancers/:load_balancer_id/frontends", r.frontendGet)

	rg.POST("/frontends", r.frontendCreate)

	rg.DELETE("/frontends", r.frontendDelete)
	rg.DELETE("/frontends/:frontend_id", r.frontendDelete)
	rg.DELETE("/loadbalancers/:load_balancer_id/frontends", r.frontendDelete)
}
