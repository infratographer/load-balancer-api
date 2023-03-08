package api

import (
	"github.com/labstack/echo/v4"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"go.infratographer.com/load-balancer-api/internal/models"
	"go.infratographer.com/load-balancer-api/internal/pubsub"
)

// assignmentsCreate handles the POST /assignments route
func (r *Router) assignmentsCreate(c echo.Context) error {
	ctx := c.Request().Context()

	payload := struct {
		FrontendID string `json:"frontend_id"`
		PoolID     string `json:"pool_id"`
	}{}

	if err := c.Bind(&payload); err != nil {
		r.logger.Errorw("error binding payload", "error", err)
		return v1BadRequestResponse(c, err)
	}

	tenantID, err := r.parseUUID(c, "tenant_id")
	if err != nil {
		return err
	}

	// validate frontend exists
	frontend, err := models.Frontends(
		models.FrontendWhere.FrontendID.EQ(payload.FrontendID),
		qm.Load("LoadBalancer"),
	).One(ctx, r.db)
	if err != nil {
		r.logger.Errorw("error fetching frontend", "error", err)
		return v1BadRequestResponse(c, err)
	}

	// validate pool exists
	pool, err := models.Pools(
		models.PoolWhere.PoolID.EQ(payload.PoolID),
	).One(ctx, r.db)
	if err != nil {
		r.logger.Errorw("error fetching pool", "error", err)
		return v1BadRequestResponse(c, err)
	}

	assignment := models.Assignment{
		TenantID:   tenantID,
		FrontendID: frontend.FrontendID,
		PoolID:     pool.PoolID,
	}

	if err := assignment.Insert(ctx, r.db, boil.Infer()); err != nil {
		r.logger.Errorw("error inserting assignment", "error", err)
		return v1InternalServerErrorResponse(c, err)
	}

	msg, err := pubsub.NewAssignmentMessage(
		someTestJWTURN,
		pubsub.NewTenantURN(tenantID),
		pubsub.NewAssignmentURN(assignment.AssignmentID),
		pubsub.NewLoadBalancerURN(frontend.LoadBalancerID),
		pubsub.NewPoolURN(pool.PoolID),
	)
	if err != nil {
		// TODO: add status to reconcile and requeue this
		r.logger.Errorw("error creating assignment message", "error", err)
	}

	if err := r.pubsub.PublishCreate(ctx, "load-balancer-assignment", "global", msg); err != nil {
		// TODO: add status to reconcile and requeue this
		r.logger.Errorw("error publishing assignment event", "error", err)
	}

	return v1AssignmentsCreatedResponse(c, assignment.AssignmentID)
}
