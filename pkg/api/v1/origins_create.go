package api

import (
	"context"

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
		Port     int64  `json:"port"`
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

	originID, err := r.createOrigin(ctx, r.db, pool.PoolID, payload.Name, payload.Target, payload.Port, payload.Disabled)
	if err != nil {
		r.logger.Error("failed to create origins", zap.Error(err))
		return v1BadRequestResponse(c, err)
	}

	msg, err := pubsub.NewOriginMessage(
		someTestJWTURN,
		pubsub.NewTenantURN(pool.TenantID),
		pubsub.NewOriginURN(originID),
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

	return v1OriginCreatedResponse(c, originID)
}

func validateOrigin(o *models.Origin) error {
	if o.OriginTarget == "" {
		return ErrMissingOriginTarget
	}

	if o.PoolID == "" {
		return ErrMissingPoolID
	}

	if o.Port == 0 {
		return ErrMissingOriginPort
	}

	return nil
}

func (r *Router) createOrigin(ctx context.Context, exec boil.ContextExecutor, poolID, name, target string, port int64, disabled bool) (string, error) {
	r.logger.Debug("creating pool origin",
		zap.String("pool.id", poolID),
		zap.String("origin.name", name),
		zap.String("origin.target", target),
		zap.Int64("origin.port", port),
		zap.Bool("origin.disabled", disabled),
	)

	origin := models.Origin{
		Name:                      name,
		OriginUserSettingDisabled: disabled,
		OriginTarget:              target,
		PoolID:                    poolID,
		Port:                      port,
		Slug:                      slug.Make(name),
		CurrentState:              "configuring",
	}

	if err := validateOrigin(&origin); err != nil {
		r.logger.Error("error validating origins", zap.Error(err))
		return "", err
	}

	if err := origin.Insert(ctx, exec, boil.Infer()); err != nil {
		r.logger.Error("error inserting origins", zap.Error(err))
		return "", err
	}

	return origin.OriginID, nil
}
