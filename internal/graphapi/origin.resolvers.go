package graphapi

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.38

import (
	"context"

	"go.infratographer.com/load-balancer-api/internal/ent/generated"
	"go.infratographer.com/load-balancer-api/internal/ent/generated/origin"
	"go.infratographer.com/load-balancer-api/internal/ent/generated/pool"
	"go.infratographer.com/load-balancer-api/internal/ent/generated/port"
	_ "go.infratographer.com/load-balancer-api/internal/ent/generated/runtime"
	"go.infratographer.com/load-balancer-api/pkg/metadata"
	"go.infratographer.com/permissions-api/pkg/permissions"
	"go.infratographer.com/x/gidx"
)

// LoadBalancerOriginCreate is the resolver for the loadBalancerOriginCreate field.
func (r *mutationResolver) LoadBalancerOriginCreate(ctx context.Context, input generated.CreateLoadBalancerOriginInput) (*LoadBalancerOriginCreatePayload, error) {
	logger := r.logger.With("poolID", input.PoolID)

	// check gidx format
	if _, err := gidx.Parse(input.PoolID.String()); err != nil {
		return nil, err
	}

	if err := permissions.CheckAccess(ctx, input.PoolID, actionLoadBalancerPoolUpdate); err != nil {
		return nil, err
	}

	// check if pool exists
	_, err := r.client.Pool.Get(ctx, input.PoolID)
	if err != nil {
		if generated.IsNotFound(err) {
			return nil, err
		}

		logger.Errorw("failed to get pool", "error", err)
		return nil, ErrInternalServerError
	}

	ogn, err := r.client.Origin.Create().SetInput(input).Save(ctx)
	if err != nil {
		if generated.IsValidationError(err) {
			return nil, err
		}

		logger.Errorw("failed to create origin", "error", err)
		return nil, ErrInternalServerError
	}

	// find loadbalancers associated with this origin to update loadbalancer metadata status
	ports, err := r.client.Port.Query().WithPools().WithLoadBalancer().Where(port.HasPoolsWith(pool.IDEQ(ogn.PoolID))).All(ctx)
	if err == nil {
		status := &metadata.LoadBalancerStatus{State: metadata.LoadBalancerStateUpdating}
		for _, p := range ports {
			if err := r.LoadBalancerStatusUpdate(ctx, p.LoadBalancerID, status); err != nil {
				logger.Errorw("failed to update loadbalancer metadata status", "error", err, "loadbalancerID", p.LoadBalancerID)
			}
		}
	}

	return &LoadBalancerOriginCreatePayload{LoadBalancerOrigin: ogn}, nil
}

// LoadBalancerOriginUpdate is the resolver for the loadBalancerOriginUpdate field.
func (r *mutationResolver) LoadBalancerOriginUpdate(ctx context.Context, id gidx.PrefixedID, input generated.UpdateLoadBalancerOriginInput) (*LoadBalancerOriginUpdatePayload, error) {
	logger := r.logger.With("originID", id.String())

	// check gidx format
	if _, err := gidx.Parse(id.String()); err != nil {
		return nil, err
	}

	ogn, err := r.client.Origin.Query().WithPool().Where(origin.IDEQ(id)).Only(ctx)
	if err != nil {
		if generated.IsNotFound(err) {
			return nil, err
		}

		logger.Errorw("failed to get origin", "error", err)
		return nil, ErrInternalServerError
	}

	if err := permissions.CheckAccess(ctx, ogn.Edges.Pool.OwnerID, actionLoadBalancerPoolUpdate); err != nil {
		return nil, err
	}

	ogn, err = ogn.Update().SetInput(input).Save(ctx)
	if err != nil {
		if generated.IsValidationError(err) {
			return nil, err
		}

		logger.Errorw("failed to update origin", "error", err)
		return nil, ErrInternalServerError
	}

	// find loadbalancers associated with this origin to update loadbalancer metadata status
	ports, err := r.client.Port.Query().WithPools().WithLoadBalancer().Where(port.HasPoolsWith(pool.HasOriginsWith(origin.IDEQ(id)))).All(ctx)
	if err == nil {
		status := &metadata.LoadBalancerStatus{State: metadata.LoadBalancerStateUpdating}
		for _, p := range ports {
			if err := r.LoadBalancerStatusUpdate(ctx, p.LoadBalancerID, status); err != nil {
				logger.Errorw("failed to update loadbalancer metadata status", "error", err, "loadbalancerID", p.LoadBalancerID)
			}
		}
	}

	return &LoadBalancerOriginUpdatePayload{LoadBalancerOrigin: ogn}, nil
}

// LoadBalancerOriginDelete is the resolver for the loadBalancerOriginDelete field.
func (r *mutationResolver) LoadBalancerOriginDelete(ctx context.Context, id gidx.PrefixedID) (*LoadBalancerOriginDeletePayload, error) {
	logger := r.logger.With("originID", id.String())

	// check gidx format
	if _, err := gidx.Parse(id.String()); err != nil {
		return nil, err
	}

	ogn, err := r.client.Origin.Query().WithPool().Where(origin.IDEQ(id)).Only(ctx)
	if err != nil {
		if generated.IsNotFound(err) {
			return nil, err
		}

		logger.Errorw("failed to get origin", "error", err)
		return nil, ErrInternalServerError
	}

	if err := permissions.CheckAccess(ctx, ogn.Edges.Pool.OwnerID, actionLoadBalancerPoolUpdate); err != nil {
		return nil, err
	}

	if err := r.client.Origin.DeleteOneID(id).Exec(ctx); err != nil {
		logger.Errorw("failed to delete origin", "error", err)
		return nil, ErrInternalServerError
	}

	// find loadbalancers associated with this origin to update loadbalancer metadata status
	ports, err := r.client.Port.Query().WithPools().WithLoadBalancer().Where(port.HasPoolsWith(pool.HasOriginsWith(origin.IDEQ(id)))).All(ctx)
	if err == nil {
		status := &metadata.LoadBalancerStatus{State: metadata.LoadBalancerStateUpdating}
		for _, p := range ports {
			if err := r.LoadBalancerStatusUpdate(ctx, p.LoadBalancerID, status); err != nil {
				logger.Errorw("failed to update loadbalancer metadata status", "error", err, "loadbalancerID", p.LoadBalancerID)
			}
		}
	}

	return &LoadBalancerOriginDeletePayload{DeletedID: id}, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
