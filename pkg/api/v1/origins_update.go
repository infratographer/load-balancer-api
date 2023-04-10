package api

import (
	"github.com/labstack/echo/v4"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"go.infratographer.com/load-balancer-api/internal/models"
	"go.infratographer.com/load-balancer-api/internal/pubsub"
	"go.uber.org/zap"
)

// originUpdate updates an origin in a pool
func (r *Router) originUpdate(c echo.Context) error {
	ctx := c.Request().Context()

	mods, err := r.originsParamsBinding(c)
	if err != nil {
		r.logger.Error("failed to bind origin params", zap.Error(err))
		return v1BadRequestResponse(c, err)
	}

	mods = append(mods, qm.Load("Pool"))

	origins, err := models.Origins(mods...).All(ctx, r.db)
	if err != nil {
		r.logger.Error("failed to get port", zap.Error(err))
		return v1InternalServerErrorResponse(c, err)
	}

	if len(origins) == 0 {
		return v1NotFoundResponse(c)
	} else if len(origins) != 1 {
		r.logger.Warn("ambiguous query ", zap.Any("origins", origins))
		return v1BadRequestResponse(c, ErrAmbiguous)
	}

	origin := origins[0]

	payload := struct {
		Disabled bool   `json:"disabled"`
		Name     string `json:"name"`
		Target   string `json:"target"`
		Port     int64  `json:"port"`
	}{}

	if err := c.Bind(&payload); err != nil {
		return v1BadRequestResponse(c, err)
	}

	// update origin
	origin.OriginUserSettingDisabled = payload.Disabled
	origin.Name = payload.Name
	origin.OriginTarget = payload.Target
	origin.Port = payload.Port

	return r.updateOrigin(c, origin)
}

// originPatch patches an origin
func (r *Router) originPatch(c echo.Context) error {
	ctx := c.Request().Context()

	mods, err := r.originsParamsBinding(c)
	if err != nil {
		r.logger.Error("failed to bind origin params", zap.Error(err))
		return v1BadRequestResponse(c, err)
	}

	mods = append(mods, qm.Load("Pool"))

	origins, err := models.Origins(mods...).All(ctx, r.db)
	if err != nil {
		r.logger.Error("failed to get port", zap.Error(err))
		return v1InternalServerErrorResponse(c, err)
	}

	if len(origins) == 0 {
		return v1NotFoundResponse(c)
	} else if len(origins) != 1 {
		return v1BadRequestResponse(c, ErrAmbiguous)
	}

	origin := origins[0]

	payload := struct {
		Disabled *bool   `json:"disabled"`
		Name     *string `json:"name"`
		Target   *string `json:"target"`
		Port     *int64  `json:"port"`
	}{}

	if err := c.Bind(&payload); err != nil {
		r.logger.Error("failed to bind origin patch input", zap.Error(err))
		return v1BadRequestResponse(c, err)
	}

	if payload.Disabled != nil {
		origin.OriginUserSettingDisabled = *payload.Disabled
	}

	if payload.Name != nil {
		origin.Name = *payload.Name
	}

	if payload.Target != nil {
		origin.OriginTarget = *payload.Target
	}

	if payload.Port != nil {
		origin.Port = *payload.Port
	}

	return r.updateOrigin(c, origin)
}

func (r *Router) updateOrigin(c echo.Context, origin *models.Origin) error {
	ctx := c.Request().Context()

	if err := validateOrigin(origin); err != nil {
		return v1BadRequestResponse(c, err)
	}

	additionalURNs := []string{
		pubsub.NewPoolURN(origin.R.Pool.PoolID),
	}

	if _, err := origin.Update(ctx, r.db, boil.Infer()); err != nil {
		r.logger.Error("failed to update port", zap.Error(err))
		return v1InternalServerErrorResponse(c, err)
	}

	msg, err := pubsub.NewOriginMessage(
		someTestJWTURN,
		pubsub.NewTenantURN(origin.R.Pool.TenantID),
		pubsub.NewOriginURN(origin.OriginID),
		additionalURNs...,
	)
	if err != nil {
		// TODO: add status to reconcile and requeue this
		r.logger.Error("failed to create load balancer origin message", zap.Error(err))
	}

	if err := r.pubsub.PublishUpdate(ctx, "load-balancer-origin", "global", msg); err != nil {
		// TODO: add status to reconcile and requeue this
		r.logger.Error("failed to publish load balancer origin message", zap.Error(err))
	}

	return v1UpdateOriginResponse(c, origin.OriginID)
}
