package graphapi

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.38

import (
	"context"
	"database/sql"

	"go.infratographer.com/permissions-api/pkg/permissions"
	"go.infratographer.com/x/events"
	"go.infratographer.com/x/gidx"

	"go.infratographer.com/load-balancer-api/internal/ent/generated"
	"go.infratographer.com/load-balancer-api/internal/ent/generated/port"
	"go.infratographer.com/load-balancer-api/internal/ent/generated/predicate"
)

// LoadBalancerCreate is the resolver for the loadBalancerCreate field.
func (r *mutationResolver) LoadBalancerCreate(ctx context.Context, input generated.CreateLoadBalancerInput) (*LoadBalancerCreatePayload, error) {
	if err := permissions.CheckAccess(ctx, input.OwnerID, actionLoadBalancerCreate); err != nil {
		return nil, err
	}

	lb, err := r.client.LoadBalancer.Create().SetInput(input).Save(ctx)
	if err != nil {
		if generated.IsValidationError(err) {
			return nil, err
		}

		r.logger.Errorw("failed to create loadbalancer", "error", err)
		return nil, ErrInternalServerError
	}

	return &LoadBalancerCreatePayload{LoadBalancer: lb}, nil
}

// LoadBalancerUpdate is the resolver for the loadBalancerUpdate field.
func (r *mutationResolver) LoadBalancerUpdate(ctx context.Context, id gidx.PrefixedID, input generated.UpdateLoadBalancerInput) (*LoadBalancerUpdatePayload, error) {
	logger := r.logger.With("loadbalancerID", id.String())

	// check gidx format
	if _, err := gidx.Parse(id.String()); err != nil {
		return nil, err
	}

	if err := permissions.CheckAccess(ctx, id, actionLoadBalancerUpdate); err != nil {
		return nil, err
	}

	lb, err := r.client.LoadBalancer.Get(ctx, id)
	if err != nil {
		if generated.IsNotFound(err) {
			return nil, err
		}

		logger.Errorw("failed to get loadbalancer", "error", err)
		return nil, ErrInternalServerError
	}

	lb, err = lb.Update().SetInput(input).Save(ctx)
	if err != nil {
		if generated.IsValidationError(err) {
			return nil, err
		}

		logger.Errorw("failed to update loadbalancer", "error", err)
		return nil, ErrInternalServerError
	}

	return &LoadBalancerUpdatePayload{LoadBalancer: lb}, nil
}

// LoadBalancerDelete is the resolver for the loadBalancerDelete field.
func (r *mutationResolver) LoadBalancerDelete(ctx context.Context, id gidx.PrefixedID) (*LoadBalancerDeletePayload, error) {
	logger := r.logger.With("loadbalancerID", id.String())

	// check gidx format
	if _, err := gidx.Parse(id.String()); err != nil {
		return nil, err
	}

	if err := permissions.CheckAccess(ctx, id, actionLoadBalancerDelete); err != nil {
		return nil, err
	}

	lb, err := r.client.LoadBalancer.Get(ctx, id)
	if err != nil {
		if generated.IsNotFound(err) {
			return nil, err
		}

		logger.Errorw("failed to get loadbalancer", "error", err)
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

	// cleanup ports associated with loadbalancer
	ports, err := tx.Port.Query().Where(predicate.Port(port.LoadBalancerIDEQ(id))).All(ctx)
	if err != nil {
		logger.Errorw("failed to query ports", "error", err)
		return nil, ErrInternalServerError
	}

	for _, p := range ports {
		if err = tx.Port.DeleteOne(p).Exec(ctx); err != nil {
			logger.Errorw("failed to delete port", "port", p.ID, "error", err)
			return nil, ErrInternalServerError
		}
	}

	// delete loadbalancer
	if err = tx.LoadBalancer.DeleteOneID(id).Exec(ctx); err != nil {
		logger.Errorw("failed to delete loadbalancer", "error", err)
		return nil, ErrInternalServerError
	}

	// delete auth relationship
	relationship := events.AuthRelationshipRelation{
		Relation:  "owner",
		SubjectID: lb.OwnerID,
	}

	// Strip cancellation from context so the auth-relationship delete fully succeeds or fails due something other than cancellation
	noCancelCtx := context.WithoutCancel(ctx)
	if err := permissions.DeleteAuthRelationships(noCancelCtx, "load-balancer", id, relationship); err != nil {
		logger.Errorw("failed to delete auth relationship", "error", err)
		return nil, ErrInternalServerError
	}

	return &LoadBalancerDeletePayload{DeletedID: id}, nil
}

// LoadBalancer is the resolver for the loadBalancer field.
func (r *queryResolver) LoadBalancer(ctx context.Context, id gidx.PrefixedID) (*generated.LoadBalancer, error) {
	logger := r.logger.With("loadbalancerID", id.String())

	// check gidx format
	if _, err := gidx.Parse(id.String()); err != nil {
		return nil, err
	}

	if err := permissions.CheckAccess(ctx, id, actionLoadBalancerGet); err != nil {
		return nil, err
	}

	lb, err := r.client.LoadBalancer.Get(ctx, id)
	if err != nil {
		if generated.IsNotFound(err) {
			return nil, err
		}

		logger.Errorw("failed to get loadbalancer", "error", err)
		return nil, ErrInternalServerError
	}

	return lb, nil
}
