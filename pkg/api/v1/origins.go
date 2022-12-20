package api

import (
	"github.com/gosimple/slug"
	"github.com/labstack/echo/v4"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"go.infratographer.com/loadbalancerapi/internal/models"
)

func (r *Router) originsParamsBinding(c echo.Context) ([]qm.QueryMod, error) {
	mods := []qm.QueryMod{}

	tenantID, err := r.parseTenantID(c)
	if err != nil {
		return nil, err
	}

	mods = append(mods, models.OriginWhere.TenantID.EQ(tenantID))

	originID := c.Param("origin_id")
	if originID != "" {
		mods = append(mods, models.OriginWhere.OriginID.EQ(originID))
	}

	queryParams := []string{"slug", "target", "port"}

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

// originsGet returns a list of origins
func (r *Router) originsGet(c echo.Context) error {
	ctx := c.Request().Context()

	mods, err := r.originsParamsBinding(c)
	if err != nil {
		r.logger.Errorw("error parsing query params", "error", err)
		return v1BadRequestResponse(c, err)
	}

	os, err := models.Origins(mods...).All(ctx, r.db)
	if err != nil {
		r.logger.Errorw("error getting origins", "error", err)
		return v1InternalServerErrorResponse(c, err)
	}

	switch len(os) {
	case 0:
		return v1NotFoundResponse(c)
	default:
		return v1OriginsResponse(c, os)
	}
}

// originsPost creates a new origin
func (r *Router) originsPost(c echo.Context) error {
	ctx := c.Request().Context()
	payload := []struct {
		Disabled    bool   `json:"disabled"`
		DisplayName string `json:"display_name"`
		Target      string `json:"target"`
		Port        int    `json:"port"`
		PoolID      string `json:"pool_id"`
	}{}

	if err := c.Bind(&payload); err != nil {
		r.logger.Errorw("error binding payload", "error", err)
		return v1BadRequestResponse(c, err)
	}

	tenantID, err := r.parseTenantID(c)
	if err != nil {
		return v1BadRequestResponse(c, err)
	}

	os := models.OriginSlice{}

	tx, err := r.db.Begin()
	if err != nil {
		r.logger.Errorw("error starting transaction", "error", err)
		return v1InternalServerErrorResponse(c, err)
	}

	for _, p := range payload {
		origin := models.Origin{
			DisplayName:               p.DisplayName,
			OriginUserSettingDisabled: p.Disabled,
			OriginTarget:              p.Target,
			PoolID:                    p.PoolID,
			Port:                      int64(p.Port),
			TenantID:                  tenantID,
			Slug:                      slug.Make(p.DisplayName),
			CurrentState:              "configuring",
		}

		os = append(os, &origin)

		if err := validateOrigin(origin); err != nil {
			r.logger.Errorw("error validating origins", "error", err)
			return v1BadRequestResponse(c, err)
		}

		if err := origin.Insert(ctx, r.db, boil.Infer()); err != nil {
			_ = tx.Rollback()

			r.logger.Errorw("error inserting origins", "error", err,
				"origin", origin, "request-id", c.Response().Header().Get(echo.HeaderXRequestID))

			return v1InternalServerErrorResponse(c, err)
		}
	}

	switch len(os) {
	case 0:
		if err := tx.Rollback(); err != nil {
			r.logger.Errorw("error rolling back transaction", "error", err)
			return v1InternalServerErrorResponse(c, err)
		}

		return v1NotFoundResponse(c)
	default:
		if err := tx.Commit(); err != nil {
			r.logger.Errorw("error committing transaction", "error", err)
			return v1InternalServerErrorResponse(c, err)
		}

		return v1CreatedResponse(c)
	}
}

func validateOrigin(o models.Origin) error {
	if o.OriginTarget == "" {
		return ErrMissingOriginTarget
	}

	if o.PoolID == "" {
		return ErrMissingPoolID
	}

	return nil
}

// originsDelete deletes an origin
func (r *Router) originsDelete(c echo.Context) error {
	ctx := c.Request().Context()

	mods, err := r.originsParamsBinding(c)
	if err != nil {
		r.logger.Errorw("error parsing query params", "error", err)
		return v1BadRequestResponse(c, err)
	}

	os, err := models.Origins(mods...).All(ctx, r.db)
	if err != nil {
		r.logger.Errorw("error getting origins", "error", err)
		return v1InternalServerErrorResponse(c, err)
	}

	switch len(os) {
	case 0:
		return v1NotFoundResponse(c)
	case 1:
		if _, err := os[0].Delete(ctx, r.db, false); err != nil {
			r.logger.Errorw("error deleting origin", "error", err)
			return v1InternalServerErrorResponse(c, err)
		}

		return v1DeletedResponse(c)
	default:
		return v1BadRequestResponse(c, ErrAmbiguous)
	}
}

// addOriginsRoutes adds the origins routes to the router
func (r *Router) addOriginRoutes(g *echo.Group) {
	g.GET("/origins", r.originsGet)
	g.GET("/origins/:origin_id", r.originsGet)

	g.POST("/origins", r.originsPost)

	g.DELETE("/origins", r.originsDelete)
	g.DELETE("/origins/:origin_id", r.originsDelete)
}
