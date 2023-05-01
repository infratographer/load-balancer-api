package graphapi

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.30

import (
	"context"

	"go.infratographer.com/load-balancer-api/internal/ent/generated"
	"go.infratographer.com/x/gidx"
)

// CreateLoadBalancerProvider is the resolver for the createLoadBalancerProvider field.
func (r *mutationResolver) CreateLoadBalancerProvider(ctx context.Context, input generated.CreateLoadBalancerProviderInput) (*generated.Provider, error) {
	// TODO: authz check here
	return r.client.Provider.Create().SetInput(input).Save(ctx)
}

// UpdateLoadBalancerProvider is the resolver for the updateLoadBalancerProvider field.
func (r *mutationResolver) UpdateLoadBalancerProvider(ctx context.Context, id gidx.PrefixedID, input generated.UpdateLoadBalancerProviderInput) (*generated.Provider, error) {
	// TODO: authz check here
	p, err := r.client.Provider.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return p.Update().SetInput(input).Save(ctx)
}
