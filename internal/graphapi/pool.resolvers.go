package graphapi

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.38

import (
	"context"
	"database/sql"

	"go.infratographer.com/load-balancer-api/internal/ent/generated"
	"go.infratographer.com/load-balancer-api/internal/ent/generated/loadbalancer"
	"go.infratographer.com/load-balancer-api/internal/ent/generated/origin"
	"go.infratographer.com/load-balancer-api/internal/ent/generated/pool"
	"go.infratographer.com/load-balancer-api/internal/ent/generated/port"
	"go.infratographer.com/load-balancer-api/internal/ent/generated/predicate"
	"go.infratographer.com/load-balancer-api/pkg/metadata"
	"go.infratographer.com/permissions-api/pkg/permissions"
	"go.infratographer.com/x/gidx"
	"golang.org/x/exp/slices"
)

// LoadBalancerPoolCreate is the resolver for the LoadBalancerPoolCreate field.
func (r *mutationResolver) LoadBalancerPoolCreate(ctx context.Context, input generated.CreateLoadBalancerPoolInput) (*LoadBalancerPoolCreatePayload, error) {
	logger := r.logger.With("ownerID", input.OwnerID)

	// check gidx owner format
	if _, err := gidx.Parse(input.OwnerID.String()); err != nil {
		return nil, err
	}

	if err := permissions.CheckAccess(ctx, input.OwnerID, actionLoadBalancerPoolCreate); err != nil {
		return nil, err
	}

	ports, err := r.client.Port.Query().Where(port.HasLoadBalancerWith(loadbalancer.OwnerIDEQ(input.OwnerID))).Where(port.IDIn(input.PortIDs...)).All(ctx)
	if err != nil {
		logger.Errorw("failed to query input ports", "error", err)
		return nil, ErrInternalServerError
	}

	if len(ports) < len(input.PortIDs) {
		return nil, ErrPortNotFound
	}

	for _, portId := range input.PortIDs {
		if err := permissions.CheckAccess(ctx, portId, actionLoadBalancerGet); err != nil {
			logger.Errorw("failed to check access", "error", err, "loadbalancerPortID", portId)
			return nil, err
		}
	}

	pool, err := r.client.Pool.Create().SetInput(input).Save(ctx)
	if err != nil {
		if generated.IsValidationError(err) {
			return nil, err
		}

		logger.Errorw("failed to create loadbalancer pool", "error", err)
		return nil, ErrInternalServerError
	}

	// if there are multiple loadbalancer ports with the same loadbalancer id, ensure the slice unique
	p := slices.CompactFunc(ports, func(x, y *generated.Port) bool {
		return x.LoadBalancerID == y.LoadBalancerID
	})

	// update metadata status for the port loadbalancer
	for _, port := range p {
		status := &metadata.LoadBalancerStatus{State: metadata.LoadBalancerStateUpdating}
		if err := r.LoadBalancerStatusUpdate(ctx, port.LoadBalancerID, status); err != nil {
			logger.Errorw("failed to update loadbalancer metadata status", "error", err, "loadbalancerID", port.LoadBalancerID)
			return nil, ErrInternalServerError
		}
	}

	return &LoadBalancerPoolCreatePayload{LoadBalancerPool: pool}, nil
}

// LoadBalancerPoolUpdate is the resolver for the LoadBalancerPoolUpdate field.
func (r *mutationResolver) LoadBalancerPoolUpdate(ctx context.Context, id gidx.PrefixedID, input generated.UpdateLoadBalancerPoolInput) (*LoadBalancerPoolUpdatePayload, error) {
	logger := r.logger.With("loadbalancerPoolID", id)

	// check gidx format
	if _, err := gidx.Parse(id.String()); err != nil {
		return nil, err
	}

	if err := permissions.CheckAccess(ctx, id, actionLoadBalancerPoolUpdate); err != nil {
		return nil, err
	}

	pool, err := r.client.Pool.Get(ctx, id)
	if err != nil {
		if generated.IsNotFound(err) {
			return nil, err
		}

		logger.Errorw("failed to get loadbalancer pool", "error", err)
		return nil, ErrInternalServerError
	}

	ports, err := r.client.Port.Query().Where(port.HasLoadBalancerWith(loadbalancer.OwnerIDEQ(pool.OwnerID))).Where(port.IDIn(input.AddPortIDs...)).All(ctx)
	if err != nil {
		logger.Errorw("failed to query input ports", "error", err)
		return nil, ErrInternalServerError
	}

	if len(ports) < len(input.AddPortIDs) {
		return nil, ErrPortNotFound
	}

	for _, portId := range input.AddPortIDs {
		if err := permissions.CheckAccess(ctx, portId, actionLoadBalancerGet); err != nil {
			logger.Errorw("failed to check access", "error", err, "loadbalancerPortID", portId)
			return nil, err
		}
	}

	pool, err = pool.Update().SetInput(input).Save(ctx)
	if err != nil {
		if generated.IsValidationError(err) {
			return nil, err
		}

		logger.Errorw("failed to update loadbalancer pool", "error", err)
		return nil, ErrInternalServerError
	}

	// if there are multiple loadbalancer ports with the same loadbalancer id, ensure the slice unique
	p := slices.CompactFunc(ports, func(x, y *generated.Port) bool {
		return x.LoadBalancerID == y.LoadBalancerID
	})

	// update metadata status for the port loadbalancer
	for _, port := range p {
		status := &metadata.LoadBalancerStatus{State: metadata.LoadBalancerStateUpdating}
		if err := r.LoadBalancerStatusUpdate(ctx, port.LoadBalancerID, status); err != nil {
			logger.Errorw("failed to update loadbalancer metadata status", "error", err, "loadbalancerID", port.LoadBalancerID)
			return nil, ErrInternalServerError
		}
	}

	return &LoadBalancerPoolUpdatePayload{LoadBalancerPool: pool}, nil
}

// LoadBalancerPoolDelete is the resolver for the loadBalancerPoolDelete field.
func (r *mutationResolver) LoadBalancerPoolDelete(ctx context.Context, id gidx.PrefixedID) (*LoadBalancerPoolDeletePayload, error) {
	logger := r.logger.With("loadbalancerPoolID", id)

	// check gidx format
	if _, err := gidx.Parse(id.String()); err != nil {
		return nil, err
	}

	if err := permissions.CheckAccess(ctx, id, actionLoadBalancerPoolDelete); err != nil {
		return nil, err
	}

	if _, err := r.client.Pool.Get(ctx, id); err != nil {
		if generated.IsNotFound(err) {
			return nil, err
		}

		logger.Errorw("failed to get loadbalancer pool", "error", err)
		return nil, ErrInternalServerError
	}

	tx, err := r.client.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		logger.Errorw("failed to begin transaction", "error", err)
		return nil, ErrInternalServerError
	}

	defer func() {
		if err != nil {
			logger.Debugw("rolling back transaction")
			if err := tx.Rollback(); err != nil {
				logger.Errorw("failed to rollback transaction", "error", err)
			}

			return
		}

		logger.Debugw("committing transaction")
		if err := tx.Commit(); err != nil {
			logger.Errorw("failed to commit transaction", "error", err)
		}
	}()

	// cleanup origins associated with pool
	origins, err := tx.Origin.Query().Where(predicate.Origin(origin.PoolIDEQ(id))).All(ctx)
	if err != nil {
		logger.Errorw("failed to query origins", "error", err)
		return nil, ErrInternalServerError
	}

	for _, o := range origins {
		if err = tx.Origin.DeleteOne(o).Exec(ctx); err != nil {
			logger.Errorw("failed to delete origin", "loadbalancerOriginID", o.ID, "error", err)
			return nil, ErrInternalServerError
		}
	}

	// delete pool
	if err := tx.Pool.DeleteOneID(id).Exec(ctx); err != nil {
		logger.Errorw("failed to delete loadbalancer pool", "error", err)
		return nil, ErrInternalServerError
	}

	// find loadbalancers associated with this pool to update loadbalancer metadata status
	ports, err := r.client.Port.Query().Where(port.HasPoolsWith(pool.IDEQ(id))).All(ctx)
	if err == nil {
		for _, p := range ports {
			status := &metadata.LoadBalancerStatus{State: metadata.LoadBalancerStateUpdating}
			if err := r.LoadBalancerStatusUpdate(ctx, p.LoadBalancerID, status); err != nil {
				logger.Errorw("failed to update loadbalancer metadata status", "error", err, "loadbalancerID", p.LoadBalancerID)
				return nil, ErrInternalServerError
			}
		}
	}

	return &LoadBalancerPoolDeletePayload{DeletedID: &id}, nil
}

// LoadBalancerPool is the resolver for the loadBalancerPool field.
func (r *queryResolver) LoadBalancerPool(ctx context.Context, id gidx.PrefixedID) (*generated.Pool, error) {
	logger := r.logger.With("loadbalancerPoolID", id.String())

	// check gidx format
	if _, err := gidx.Parse(id.String()); err != nil {
		return nil, err
	}

	if err := permissions.CheckAccess(ctx, id, actionLoadBalancerPoolGet); err != nil {
		return nil, err
	}

	pool, err := r.client.Pool.Get(ctx, id)
	if err != nil {
		if generated.IsNotFound(err) {
			return nil, err
		}

		logger.Errorw("failed to get loadbalancer pool", "error", err)
		return nil, ErrInternalServerError
	}

	return pool, nil
}
