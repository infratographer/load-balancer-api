package api

import (
	"github.com/gosimple/slug"
	"github.com/labstack/echo/v4"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"go.infratographer.com/load-balancer-api/internal/models"
)

// poolCreate creates a new pool
func (r *Router) poolCreate(c echo.Context) error {
	ctx := c.Request().Context()
	payload := struct {
		DisplayName string `json:"display_name"`
		Protocol    string `json:"protocol"`
	}{}

	if err := c.Bind(&payload); err != nil {
		r.logger.Errorw("error binding payload", "error", err)
		return v1BadRequestResponse(c, err)
	}

	tenantID, err := r.parseTenantID(c)
	if err != nil {
		r.logger.Errorw("error parsing tenant id", "error", err)
		return v1BadRequestResponse(c, err)
	}

	pool := &models.Pool{
		DisplayName: payload.DisplayName,
		Protocol:    payload.Protocol,
		TenantID:    tenantID,
		Slug:        slug.Make(payload.DisplayName),
	}

	if err := validatePool(pool); err != nil {
		r.logger.Errorw("error validating pool", "error", err)

		return v1BadRequestResponse(c, err)
	}

	if err := pool.Insert(ctx, r.db, boil.Infer()); err != nil {
		r.logger.Errorw("error inserting pool", "error", err)

		return v1InternalServerErrorResponse(c, err)
	}

	return v1PoolCreatedResponse(c, pool.PoolID)
}

// validatePool validates a pool
func validatePool(p *models.Pool) error {
	if p.DisplayName == "" {
		return ErrDisplayNameMissing
	}

	if p.Protocol == "" {
		p.Protocol = "tcp"
	}

	if p.Protocol != "tcp" {
		return ErrPoolProtocolInvalid
	}

	return nil
}
