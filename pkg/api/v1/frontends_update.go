package api

import (
	"github.com/labstack/echo/v4"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"go.infratographer.com/load-balancer-api/internal/models"
)

// frontendUpdate updates a frontend
func (r *Router) frontendUpdate(c echo.Context) error {
	ctx := c.Request().Context()

	mods, err := r.frontendParamsBinding(c)
	if err != nil {
		r.logger.Errorw("failed to bind frontend params", "error", err)
		return v1BadRequestResponse(c, err)
	}

	frontend, err := models.Frontends(mods...).One(ctx, r.db)
	if err != nil {
		r.logger.Errorw("failed to get frontend", "error", err)
		return v1InternalServerErrorResponse(c, err)
	}

	payload := struct {
		DisplayName string `json:"display_name"`
		Port        int64  `json:"port"`
	}{}
	if err := c.Bind(&payload); err != nil {
		r.logger.Errorw("failed to bind frontend update input", "error", err)
		return v1BadRequestResponse(c, err)
	}

	frontend.DisplayName = payload.DisplayName
	frontend.Port = payload.Port
	// TODO do we need to update a CurrentState here?

	if err := validateFrontend(frontend); err != nil {
		return v1BadRequestResponse(c, err)
	}

	if _, err := frontend.Update(ctx, r.db, boil.Infer()); err != nil {
		r.logger.Errorw("failed to update frontend", "error", err)
		return v1InternalServerErrorResponse(c, err)
	}

	// TODO emit event that load balancers associated with frontend is updated

	return v1UpdateFrontendResponse(c, frontend.FrontendID)
}
