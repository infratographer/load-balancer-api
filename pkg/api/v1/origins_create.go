package api

import (
	"github.com/gosimple/slug"
	"github.com/labstack/echo/v4"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"go.infratographer.com/load-balancer-api/internal/models"
)

// originsCreate creates a new origin
func (r *Router) originsCreate(c echo.Context) error {
	ctx := c.Request().Context()
	payload := struct {
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

	origin := models.Origin{
		DisplayName:               payload.DisplayName,
		OriginUserSettingDisabled: payload.Disabled,
		OriginTarget:              payload.Target,
		PoolID:                    payload.PoolID,
		Port:                      int64(payload.Port),
		TenantID:                  tenantID,
		Slug:                      slug.Make(payload.DisplayName),
		CurrentState:              "configuring",
	}

	if err := validateOrigin(origin); err != nil {
		r.logger.Errorw("error validating origins", "error", err)
		return v1BadRequestResponse(c, err)
	}

	if err := origin.Insert(ctx, r.db, boil.Infer()); err != nil {
		r.logger.Errorw("error inserting origins", "error", err,
			"origin", origin, "request-id", c.Response().Header().Get(echo.HeaderXRequestID))

		return v1InternalServerErrorResponse(c, err)
	}

	return v1OriginCreatedResponse(c, origin.OriginID)
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
