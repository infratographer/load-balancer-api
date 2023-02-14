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

	tenantID, err := r.parseTenantID(c)
	if err != nil {
		return err
	}

	assignment := models.Assignment{
		TenantID:   tenantID,
		FrontendID: payload.FrontendID,
		PoolID:     payload.PoolID,
	}

	if err := assignment.Insert(ctx, r.db, boil.Infer()); err != nil {
		r.logger.Errorw("error inserting assignment", "error", err)
		return v1InternalServerErrorResponse(c, err)
	}

	aMods := []qm.QueryMod{
		models.AssignmentWhere.TenantID.EQ(tenantID),
		models.AssignmentWhere.FrontendID.EQ(payload.FrontendID),
		models.AssignmentWhere.PoolID.EQ(payload.PoolID),
	}

	aModel, err := models.Assignments(aMods...).One(ctx, r.db)
	if err != nil {
		return v1InternalServerErrorResponse(c, err)
	}

	feMods := models.FrontendWhere.FrontendID.EQ(payload.FrontendID)

	feModel, err := models.Frontends(feMods).One(ctx, r.db)
	if err != nil {
		// TODO: add status to reconcile and requeue this
		r.logger.Errorw("error fetching frontend", "error", err)
	}

	msg, err := pubsub.NewAssignmentMessage(someTestJWTURN, "urn:infratographer:infratographer.com:tenant:"+tenantID, pubsub.NewAssignmentURN(aModel.AssignmentID), "urn:infratographer:infratographer.com:load-balancer:"+feModel.LoadBalancerID)
	if err != nil {
		// TODO: add status to reconcile and requeue this
		r.logger.Errorw("error creating message", "error", err)
	}

	if err := pubsub.PublishCreate(ctx, r.events, "assignment", "global", msg); err != nil {
		// TODO: add status to reconcile and requeue this
		r.logger.Errorw("error publishing event", "error", err)
	}

	return v1AssignmentsCreatedResponse(c, assignment.AssignmentID)
}
