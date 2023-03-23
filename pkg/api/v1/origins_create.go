package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"go.infratographer.com/load-balancer-api/internal/models"
	"go.infratographer.com/load-balancer-api/internal/pubsub"
)

// originsCreate creates a new origin
func (r *Router) originsCreate(c *gin.Context) {
	ctx := c.Request.Context()
	payload := struct {
		Disabled bool   `json:"disabled"`
		Name     string `json:"name"`
		Target   string `json:"target"`
		Port     int    `json:"port"`
	}{}

	if err := c.Bind(&payload); err != nil {
		r.logger.Errorw("error binding payload", "error", err)
		v1BadRequestResponse(c, err)

		return
	}

	poolID, err := r.parsePoolID(c)
	if err != nil {
		v1BadRequestResponse(c, err)

		return
	}

	// validate pool exists
	pool, err := models.Pools(
		models.PoolWhere.PoolID.EQ(poolID),
	).One(ctx, r.db)
	if err != nil {
		r.logger.Errorw("error fetching pool", "error", err)
		v1BadRequestResponse(c, err)

		return
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
		r.logger.Errorw("error validating origins", "error", err)
		v1BadRequestResponse(c, err)

		return
	}

	if err := origin.Insert(ctx, r.db, boil.Infer()); err != nil {
		r.logger.Errorw("error inserting origin", "error", err)
		v1InternalServerErrorResponse(c, err)

		return
	}

	msg, err := pubsub.NewOriginMessage(
		someTestJWTURN,
		pubsub.NewTenantURN(pool.TenantID),
		pubsub.NewOriginURN(origin.OriginID),
		pubsub.NewPoolURN(pool.PoolID),
	)
	if err != nil {
		// TODO: add status to reconcile and requeue this
		r.logger.Errorw("error creating origin message", "error", err)
	}

	if err := r.pubsub.PublishCreate(ctx, "load-balancer-origin", "global", msg); err != nil {
		// TODO: add status to reconcile and requeue this
		r.logger.Errorw("error publishing origin event", "error", err)
	}

	v1OriginCreatedResponse(c, origin.OriginID)
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
