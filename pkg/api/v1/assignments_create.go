package api

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"go.uber.org/zap"

	"go.infratographer.com/load-balancer-api/internal/models"
	"go.infratographer.com/load-balancer-api/internal/pubsub"
)

// assignmentsCreate handles the POST /assignments route
func (r *Router) assignmentsCreate(c echo.Context) error {
	ctx := c.Request().Context()

	payload := struct {
		PortID string `json:"port_id"`
		PoolID string `json:"pool_id"`
	}{}

	if err := c.Bind(&payload); err != nil {
		r.logger.Error("error binding payload", zap.Error(err))
		return v1BadRequestResponse(c, err)
	}

	tenantID, err := r.parseUUID(c, "tenant_id")
	if err != nil {
		return err
	}

	// validate port exists
	port, err := models.Ports(
		models.PortWhere.PortID.EQ(payload.PortID),
		qm.Load("LoadBalancer"),
	).One(ctx, r.db)
	if err != nil {
		r.logger.Error("error fetching port", zap.Error(err))
		return v1BadRequestResponse(c, err)
	}

	assignmentID, err := r.createAssignment(ctx, r.db, tenantID, payload.PoolID, port.PortID)
	if err != nil {
		r.logger.Error("failed to create assignment", zap.Error(err))
		return v1BadRequestResponse(c, err)
	}

	msg, err := pubsub.NewMessage(
		pubsub.NewTenantURN(tenantID),
		pubsub.WithActorURN(someTestJWTURN),
		pubsub.WithSubjectURN(
			pubsub.NewAssignmentURN(assignmentID),
		),
		pubsub.WithAdditionalSubjectURNs(
			pubsub.NewLoadBalancerURN(port.LoadBalancerID),
			pubsub.NewPoolURN(payload.PoolID),
		),
		pubsub.WithSubjectFields(map[string]string{"tenant_id": tenantID}),
	)
	if err != nil {
		// TODO: add status to reconcile and requeue this
		r.logger.Error("error creating assignment message", zap.Error(err))
	}

	if err := r.pubsub.PublishCreate(ctx, "load-balancer-assignment", "global", msg); err != nil {
		// TODO: add status to reconcile and requeue this
		r.logger.Error("error publishing assignment event", zap.Error(err))
	}

	return v1AssignmentsCreatedResponse(c, assignmentID)
}

func (r *Router) createAssignment(ctx context.Context, exec boil.ContextExecutor, tenantID, poolID, portID string) (string, error) {
	r.logger.Debug("creating assignment",
		zap.String("tenant.id", tenantID),
		zap.String("pool.id", poolID),
		zap.String("port.id", portID),
	)

	// validate pool exists
	pool, err := models.Pools(
		models.PoolWhere.PoolID.EQ(poolID),
		models.PoolWhere.TenantID.EQ(tenantID),
	).One(ctx, r.db)
	if err != nil {
		r.logger.Error("error fetching pool", zap.Error(err))
		return "", err
	}

	assignment := models.Assignment{
		TenantID: tenantID,
		PortID:   portID,
		PoolID:   pool.PoolID,
	}

	if err := assignment.Insert(ctx, exec, boil.Infer()); err != nil {
		r.logger.Error("error inserting assignment", zap.Error(err))
		return "", err
	}

	return assignment.AssignmentID, nil
}
