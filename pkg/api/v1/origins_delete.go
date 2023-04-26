package api

import (
	"github.com/labstack/echo/v4"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"go.infratographer.com/load-balancer-api/internal/models"
	"go.infratographer.com/load-balancer-api/internal/pubsub"
	"go.uber.org/zap"
)

// originsDelete deletes an origin
func (r *Router) originsDelete(c echo.Context) error {
	ctx := c.Request().Context()

	mods, err := r.originsParamsBinding(c)
	if err != nil {
		r.logger.Error("error parsing query params", zap.Error(err))
		return v1BadRequestResponse(c, err)
	}

	mods = append(mods, qm.Load("Pool"))

	os, err := models.Origins(mods...).All(ctx, r.db)
	if err != nil {
		r.logger.Error("error getting origins", zap.Error(err))
		return v1InternalServerErrorResponse(c, err)
	}

	if len(os) == 0 {
		return v1NotFoundResponse(c)
	} else if len(os) != 1 {
		return v1BadRequestResponse(c, ErrAmbiguous)
	}

	origin := os[0]

	if _, err := origin.Delete(ctx, r.db, false); err != nil {
		r.logger.Error("error deleting origin", zap.Error(err))
		return v1InternalServerErrorResponse(c, err)
	}

	msg, err := pubsub.NewMessage(
		pubsub.NewTenantURN(origin.R.Pool.TenantID),
		pubsub.WithActorURN(someTestJWTURN),
		pubsub.WithSubjectURN(
			pubsub.NewOriginURN(origin.OriginID),
		),
		pubsub.WithAdditionalSubjectURNs(
			pubsub.NewPoolURN(origin.R.Pool.PoolID),
		),
		pubsub.WithSubjectFields(
			map[string]string{
				"tenant_id":  origin.R.Pool.TenantID,
				"tenant_urn": pubsub.NewTenantURN(origin.R.Pool.TenantID),
			},
		),
	)
	if err != nil {
		// TODO: add status to reconcile and requeue this
		r.logger.Error("error creating origin message", zap.Error(err))
	}

	if err := r.pubsub.PublishDelete(ctx, "load-balancer-origin", "global", msg); err != nil {
		// TODO: add status to reconcile and requeue this
		r.logger.Error("error publishing origin event", zap.Error(err))
	}

	return v1DeletedResponse(c)
}
