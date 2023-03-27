package api

import (
	"github.com/gosimple/slug"
	"github.com/labstack/echo/v4"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"go.infratographer.com/load-balancer-api/internal/models"
	"go.infratographer.com/load-balancer-api/internal/pubsub"
	"go.uber.org/zap"
)

// originsCreate creates a new origin
func (r *Router) originsCreate(c echo.Context) error {
	ctx := c.Request().Context()
	payload := struct {
		Disabled bool   `json:"disabled"`
		Name     string `json:"name"`
		Target   string `json:"target"`
		Port     int    `json:"port"`
	}{}

	if err := c.Bind(&payload); err != nil {
		r.logger.Error("error binding payload", zap.Error(err))
		return v1BadRequestResponse(c, err)
	}

	poolID, err := r.parseUUID(c, "pool_id")
	if err != nil {
		return v1BadRequestResponse(c, err)
	}

	// validate pool exists
	pool, err := models.Pools(
		models.PoolWhere.PoolID.EQ(poolID),
	).One(ctx, r.db)
	if err != nil {
		r.logger.Error("error fetching pool", zap.Error(err))
		return v1BadRequestResponse(c, err)
	}

	origin := models.Origin{
		Name:                      payload.Name,
		OriginUserSettingDisabled: payload.Disabled,
		OriginTarget:              payload.Target,
		PoolID:                    pool.PoolID,
		Port:                      int64(payload.Port),
		Slug:                      slug.Make(payload.Name),
		CurrentState:              "configuring",
	}

	if err := validateOrigin(origin); err != nil {
		r.logger.Error("error validating origins", zap.Error(err))
		return v1BadRequestResponse(c, err)
	}

	if err := origin.Insert(ctx, r.db, boil.Infer()); err != nil {
		r.logger.Error("error inserting origins", zap.Error(err))

		return v1InternalServerErrorResponse(c, err)
	}

	msg, err := pubsub.NewOriginMessage(
		someTestJWTURN,
		pubsub.NewTenantURN(pool.TenantID),
		pubsub.NewOriginURN(origin.OriginID),
		pubsub.NewPoolURN(pool.PoolID),
	)
	if err != nil {
		// TODO: add status to reconcile and requeue this
		r.logger.Error("error creating origin message", zap.Error(err))
	}

	if err := r.pubsub.PublishCreate(ctx, "load-balancer-origin", "global", msg); err != nil {
		// TODO: add status to reconcile and requeue this
		r.logger.Error("error publishing origin event", zap.Error(err))
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
