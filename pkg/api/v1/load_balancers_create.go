package api

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"github.com/labstack/echo/v4"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"go.infratographer.com/load-balancer-api/internal/models"
	"go.infratographer.com/load-balancer-api/internal/pubsub"
	"go.uber.org/zap"
)

// loadBalancerCreate creates a new load balancer for a tenant
func (r *Router) loadBalancerCreate(c echo.Context) error {
	ctx := c.Request().Context()

	payload := struct {
		Name             string `json:"name"`
		LoadBalancerSize string `json:"load_balancer_size"`
		LoadBalancerType string `json:"load_balancer_type"`
		IPAddressID      string `json:"ip_address_id"`
		LocationID       string `json:"location_id"`
		Ports            []struct {
			Name  string   `json:"name"`
			Port  int64    `json:"port"`
			Pools []string `json:"pools"`
		} `json:"ports"`
	}{}

	if err := c.Bind(&payload); err != nil {
		r.logger.Error("failed to bind load balancer input", zap.Error(err))
		return v1BadRequestResponse(c, err)
	}

	// Ensure the tenant ID is a set from the path,this prevents
	// a tenant from creating a load balancer for another tenant
	tenantID, err := r.parseUUID(c, "tenant_id")
	if err != nil {
		r.logger.Error("bad request", zap.Error(err))
		return v1BadRequestResponse(c, err)
	}

	// TODO get/validate IP address uuid from IPAM - just mock it out for now
	if payload.IPAddressID != "" {
		if _, err := uuid.Parse(payload.IPAddressID); err != nil {
			r.logger.Error("bad ip address uuid in request", zap.Error(err))
			return v1BadRequestResponse(c, err)
		}
	} else {
		payload.IPAddressID = uuid.NewString()
	}

	lb := &models.LoadBalancer{
		TenantID:         tenantID,
		Name:             payload.Name,
		LoadBalancerSize: payload.LoadBalancerSize,
		LoadBalancerType: payload.LoadBalancerType,
		IPAddressID:      payload.IPAddressID,
		LocationID:       payload.LocationID,
		Slug:             slug.Make(payload.Name),
		CurrentState:     "provisioning",
	}

	if err := validateLoadBalancer(lb); err != nil {
		r.logger.Error("failed to validate load balancer", zap.Error(err))
		return v1BadRequestResponse(c, err)
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		r.logger.Error("failed to begin transaction", zap.Error(err))
		return v1InternalServerErrorResponse(c, err)
	}

	if err = lb.Insert(ctx, tx, boil.Infer()); err != nil {
		r.logger.Error("failed to create load balancer, rolling back transaction", zap.Error(err))

		if err := tx.Rollback(); err != nil {
			r.logger.Error("error rolling back transaction", zap.Error(err))
			return v1InternalServerErrorResponse(c, err)
		}

		return v1InternalServerErrorResponse(c, err)
	}

	r.logger.Info("created new load balancer", zap.Any("loadbalancer.id", lb.LoadBalancerID))

	additionalURNs := []string{}

	for _, p := range payload.Ports {
		portID, err := r.loadBalancerPortCreate(ctx, tx, lb.LoadBalancerID, p.Name, p.Port)
		if err != nil {
			r.logger.Error("failed to create load balancer port, rolling back transaction", zap.Error(err))

			if err := tx.Rollback(); err != nil {
				r.logger.Error("error rolling back transaction", zap.Error(err))
				return v1InternalServerErrorResponse(c, err)
			}

			return v1BadRequestResponse(c, err)
		}

		additionalURNs = append(additionalURNs, pubsub.NewPortURN(portID))

		if len(p.Pools) > 0 {
			for _, pool := range p.Pools {
				assignmentID, err := r.loadBalancerAssignmentCreate(ctx, tx, tenantID, lb.LoadBalancerID, pool, portID)
				if err != nil {
					r.logger.Error("failed to create load balancer assignment, rolling back transaction", zap.Error(err))

					if err := tx.Rollback(); err != nil {
						r.logger.Error("error rolling back transaction", zap.Error(err))
						return v1InternalServerErrorResponse(c, err)
					}

					return v1BadRequestResponse(c, err)
				}

				additionalURNs = append(additionalURNs, pubsub.NewAssignmentURN(assignmentID))
			}
		}
	}

	if err := tx.Commit(); err != nil {
		r.logger.Error("failed to commit transaction", zap.Error(err))
		return v1InternalServerErrorResponse(c, err)
	}

	msg, err := pubsub.NewLoadBalancerMessage(
		someTestJWTURN,
		pubsub.NewTenantURN(tenantID),
		pubsub.NewLoadBalancerURN(lb.LoadBalancerID),
		additionalURNs...,
	)
	if err != nil {
		// TODO: add status to reconcile and requeue this
		r.logger.Error("failed to create load balancer message", zap.Error(err))
	}

	if err := r.pubsub.PublishCreate(ctx, "load-balancer", "global", msg); err != nil {
		// TODO: add status to reconcile and requeue this
		r.logger.Error("failed to publish load balancer message", zap.Error(err))
	}

	return v1LoadBalancerCreatedResponse(c, lb.LoadBalancerID)
}

func (r *Router) loadBalancerPortCreate(ctx context.Context, tx *sql.Tx, loadBalancerID string, portName string, portNumber int64) (string, error) {
	r.logger.Debug("creating loadbalancer port",
		zap.String("loadbalancer.id", loadBalancerID),
		zap.String("port.name", portName),
		zap.Int64("port.number", portNumber),
	)

	port := models.Port{
		Name:           portName,
		Port:           portNumber,
		LoadBalancerID: loadBalancerID,
		Slug:           slug.Make(portName),
		CurrentState:   "pending",
	}

	if err := validatePort(&port); err != nil {
		r.logger.Error("failed to validate port", zap.Error(err))
		return "", err
	}

	if err := port.Insert(ctx, tx, boil.Infer()); err != nil {
		r.logger.Error("failed to insert port", zap.Error(err))
		return "", err
	}

	return port.PortID, nil
}

func (r *Router) loadBalancerAssignmentCreate(ctx context.Context, tx *sql.Tx, tenantID, loadBalancerID, poolID, portID string) (string, error) {
	r.logger.Debug("creating loadbalancer assignment",
		zap.String("tenant.id", tenantID),
		zap.String("loadbalancer.id", loadBalancerID),
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

	if err := assignment.Insert(ctx, tx, boil.Infer()); err != nil {
		r.logger.Error("error inserting assignment", zap.Error(err))
		return "", err
	}

	return assignment.AssignmentID, nil
}
