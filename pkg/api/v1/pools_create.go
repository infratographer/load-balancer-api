package api

import (
	"github.com/gosimple/slug"
	"github.com/labstack/echo/v4"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"go.infratographer.com/loadbalancerapi/internal/models"
)

// poolCreate creates a new pool
func (r *Router) poolCreate(c echo.Context) error {
	ctx := c.Request().Context()
	payload := []struct {
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

	ps := models.PoolSlice{}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		r.logger.Errorw("error starting transaction", "error", err)
		return v1InternalServerErrorResponse(c, err)
	}

	for _, p := range payload {
		pool := &models.Pool{
			DisplayName: p.DisplayName,
			Protocol:    p.Protocol,
			TenantID:    tenantID,
			Slug:        slug.Make(p.DisplayName),
		}

		if err := validatePool(pool); err != nil {
			if err := tx.Rollback(); err != nil {
				r.logger.Errorw("error rolling back transaction", "error", err)
				return v1InternalServerErrorResponse(c, err)
			}

			r.logger.Errorw("error validating pool", "error", err)

			return v1BadRequestResponse(c, err)
		}

		ps = append(ps, pool)

		if err := pool.Insert(ctx, tx, boil.Infer()); err != nil {
			_ = tx.Rollback()

			r.logger.Errorw("error inserting pool", "error", err)

			return v1InternalServerErrorResponse(c, err)
		}
	}

	switch len(ps) {
	case 0:
		if err := tx.Rollback(); err != nil {
			r.logger.Errorw("error rolling back transaction", "error", err)
			return v1InternalServerErrorResponse(c, err)
		}

		return v1NotFoundResponse(c)
	default:
		if err := tx.Commit(); err != nil {
			r.logger.Errorw("error committing transaction", "error", err)
			return v1BadRequestResponse(c, err)
		}

		return v1CreatedResponse(c)
	}
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
