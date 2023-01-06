package api

import (
	"github.com/gosimple/slug"
	"github.com/labstack/echo/v4"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"go.infratographer.com/load-balancer-api/internal/models"
)

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
			Slug:             slug.Make(p.DisplayName),
			CurrentState:     "provisioning",
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
