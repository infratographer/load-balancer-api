package api

import (
	"github.com/gosimple/slug"
	"github.com/labstack/echo/v4"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"go.infratographer.com/loadbalancerapi/internal/models"
)

// originsCreate creates a new origin
func (r *Router) originsCreate(c echo.Context) error {
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
