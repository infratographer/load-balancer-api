package api

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"go.infratographer.com/load-balancer-api/internal/models"
	"go.infratographer.com/load-balancer-api/internal/pubsub"
	"go.uber.org/zap"
)

// portUpdate updates a port
func (r *Router) portUpdate(c echo.Context) error {
	ctx := c.Request().Context()

	mods, err := r.portParamsBinding(c)
	if err != nil {
		r.logger.Error("failed to bind port params", zap.Error(err))
		return v1BadRequestResponse(c, err)
	}

	ports, err := models.Ports(mods...).All(ctx, r.db)
	if err != nil {
		r.logger.Error("failed to get port", zap.Error(err))
		return v1InternalServerErrorResponse(c, err)
	}

	if len(ports) == 0 {
		return v1NotFoundResponse(c)
	} else if len(ports) != 1 {
		return v1BadRequestResponse(c, ErrAmbiguous)
	}

	port := ports[0]

	// collect original pool ids for the port
	origPools := make([]string, len(port.R.Assignments))
	for i, p := range port.R.Assignments {
		origPools[i] = p.PoolID
	}

	payload := struct {
		Name  string   `json:"name"`
		Port  int64    `json:"port"`
		Pools []string `json:"pools"`
	}{}

	if err := c.Bind(&payload); err != nil {
		r.logger.Error("failed to bind port update input", zap.Error(err))
		return v1BadRequestResponse(c, err)
	}

	port.Name = payload.Name
	port.Port = payload.Port
	// TODO do we need to update a CurrentState here?

	portID, err := r.updatePort(c, port, origPools, payload.Pools)
	if err != nil {
		return err
	}

	return v1UpdatePortResponse(c, portID)
}

// portPatch patches a port
func (r *Router) portPatch(c echo.Context) error {
	ctx := c.Request().Context()

	mods, err := r.portParamsBinding(c)
	if err != nil {
		r.logger.Error("failed to bind port params", zap.Error(err))
		return v1BadRequestResponse(c, err)
	}

	ports, err := models.Ports(mods...).All(ctx, r.db)
	if err != nil {
		r.logger.Error("failed to get port", zap.Error(err))
		return v1InternalServerErrorResponse(c, err)
	}

	if len(ports) == 0 {
		return v1NotFoundResponse(c)
	} else if len(ports) != 1 {
		return v1BadRequestResponse(c, ErrAmbiguous)
	}

	port := ports[0]

	// collect original pool ids for the port
	origPools := make([]string, len(port.R.Assignments))
	for i, p := range port.R.Assignments {
		origPools[i] = p.PoolID
	}

	payload := struct {
		Name  *string  `json:"name"`
		Port  *int64   `json:"port"`
		Pools []string `json:"pools"`
	}{}

	if err := c.Bind(&payload); err != nil {
		r.logger.Error("failed to bind port update input", zap.Error(err))
		return v1BadRequestResponse(c, err)
	}

	if payload.Name != nil {
		port.Name = *payload.Name
	}

	if payload.Port != nil {
		port.Port = *payload.Port
	}

	if payload.Pools == nil {
		payload.Pools = origPools
	}

	// TODO do we need to update a CurrentState here?

	portID, err := r.updatePort(c, port, origPools, payload.Pools)
	if err != nil {
		return err
	}

	return v1UpdatePortResponse(c, portID)
}

func (r *Router) updatePort(c echo.Context, port *models.Port, origPools, newPools []string) (string, error) {
	ctx := c.Request().Context()

	if err := validatePort(port); err != nil {
		return "", v1BadRequestResponse(c, err)
	}

	// validate load balancer
	lb, err := models.LoadBalancers(
		models.LoadBalancerWhere.LoadBalancerID.EQ(port.LoadBalancerID),
	).One(ctx, r.db)
	if err != nil {
		r.logger.Error("error looking up load balancer for port", zap.Error(err))
		return "", v1BadRequestResponse(c, err)
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		r.logger.Error("failed to begin transaction", zap.Error(err))
		return "", v1InternalServerErrorResponse(c, err)
	}

	if _, err := port.Update(ctx, tx, boil.Infer()); err != nil {
		r.logger.Error("failed to update port", zap.Error(err))

		if err := tx.Rollback(); err != nil {
			r.logger.Error("error rolling back transaction", zap.Error(err))
			return "", v1InternalServerErrorResponse(c, err)
		}

		return "", v1InternalServerErrorResponse(c, err)
	}

	additionalURNs := []string{
		pubsub.NewLoadBalancerURN(lb.LoadBalancerID),
	}

	diff := sliceCompare(origPools, newPools)

	for poolID, v := range diff {
		switch {
		case v < 0:
			assignmentID, err := r.deleteAssignment(ctx, tx, lb.TenantID, poolID, port.PortID)
			if err != nil {
				r.logger.Error("failed to create port assignment, rolling back transaction", zap.Error(err))

				if err := tx.Rollback(); err != nil {
					r.logger.Error("error rolling back transaction", zap.Error(err))
					return "", v1InternalServerErrorResponse(c, err)
				}

				return "", v1BadRequestResponse(c, err)
			}

			additionalURNs = append(additionalURNs, pubsub.NewAssignmentURN(assignmentID))
		case v >= 1:
			if _, err := uuid.Parse(poolID); err != nil {
				r.logger.Error("invalid uuid in port payload", zap.Error(err))
				return "", v1BadRequestResponse(c, err)
			}

			assignmentID, err := r.createAssignment(ctx, tx, lb.TenantID, poolID, port.PortID)
			if err != nil {
				r.logger.Error("failed to create port assignment, rolling back transaction", zap.Error(err))

				if err := tx.Rollback(); err != nil {
					r.logger.Error("error rolling back transaction", zap.Error(err))
					return "", v1InternalServerErrorResponse(c, err)
				}

				return "", v1BadRequestResponse(c, err)
			}

			additionalURNs = append(additionalURNs, pubsub.NewAssignmentURN(assignmentID))
		default:
			r.logger.Debug("assignment for pool already exists, skipping", zap.String("pool.id", poolID), zap.String("port.id", port.PortID))
		}
	}

	if err := tx.Commit(); err != nil {
		r.logger.Error("failed to commit transaction", zap.Error(err))
		return "", v1InternalServerErrorResponse(c, err)
	}

	msg, err := pubsub.NewMessage(
		pubsub.NewTenantURN(lb.TenantID),
		pubsub.WithActorURN(someTestJWTURN),
		pubsub.WithSubjectURN(
			pubsub.NewPortURN(port.PortID),
		),
		pubsub.WithAdditionalSubjectURNs(
			additionalURNs...,
		),
		pubsub.WithSubjectFields(map[string]string{"tenant_id": lb.TenantID}),
	)
	if err != nil {
		// TODO: add status to reconcile and requeue this
		r.logger.Error("failed to create load balancer port message", zap.Error(err))
	}

	if err := r.pubsub.PublishUpdate(ctx, "load-balancer-port", "global", msg); err != nil {
		// TODO: add status to reconcile and requeue this
		r.logger.Error("failed to publish load balancer port message", zap.Error(err))
	}

	return port.PortID, nil
}
