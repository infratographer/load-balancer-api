package api

import (
	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"github.com/labstack/echo/v4"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"go.infratographer.com/loadbalancerapi/internal/models"
)

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
		r.logger.Errorw("failed to begin transaction", "error", err)
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
			r.logger.Errorw("failed to insert frontend", "error", err)

			if err := tx.Rollback(); err != nil {
				r.logger.Errorw("failed to rollback transaction", "error", err)
			}

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
