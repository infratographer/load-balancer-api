package graphapi

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.38

import (
	"context"
	"fmt"

	"entgo.io/contrib/entgql"
	"go.infratographer.com/load-balancer-api/internal/ent/generated"
	"go.infratographer.com/load-balancer-api/internal/ent/generated/loadbalancer"
	"go.infratographer.com/load-balancer-api/internal/ent/generated/pool"
	_ "go.infratographer.com/load-balancer-api/internal/ent/generated/runtime"
	"go.infratographer.com/permissions-api/pkg/permissions"
	"go.infratographer.com/x/gidx"
)

// Owner is the resolver for the owner field.
func (r *loadBalancerResolver) Owner(ctx context.Context, obj *generated.LoadBalancer) (*ResourceOwner, error) {
	return &ResourceOwner{ID: obj.OwnerID}, nil
}

// Owner is the resolver for the owner field.
func (r *loadBalancerPoolResolver) Owner(ctx context.Context, obj *generated.Pool) (*ResourceOwner, error) {
	return &ResourceOwner{ID: obj.OwnerID}, nil
}

// Owner is the resolver for the owner field.
func (r *loadBalancerProviderResolver) Owner(ctx context.Context, obj *generated.Provider) (*ResourceOwner, error) {
	return &ResourceOwner{ID: obj.OwnerID}, nil
}

// LoadBalancers is the resolver for the loadBalancers field.
func (r *resourceOwnerResolver) LoadBalancers(ctx context.Context, obj *ResourceOwner, after *entgql.Cursor[gidx.PrefixedID], first *int, before *entgql.Cursor[gidx.PrefixedID], last *int, orderBy *generated.LoadBalancerOrder, where *generated.LoadBalancerWhereInput) (*generated.LoadBalancerConnection, error) {
	if err := permissions.CheckAccess(ctx, obj.ID, actionLoadBalancerGet); err != nil {
		return nil, err
	}

	return r.client.LoadBalancer.Query().Where(loadbalancer.OwnerID(obj.ID)).Paginate(ctx, after, first, before, last, generated.WithLoadBalancerOrder(orderBy), generated.WithLoadBalancerFilter(where.Filter))
}

// LoadBalancerPools is the resolver for the loadBalancerPools field.
func (r *resourceOwnerResolver) LoadBalancerPools(ctx context.Context, obj *ResourceOwner, after *entgql.Cursor[gidx.PrefixedID], first *int, before *entgql.Cursor[gidx.PrefixedID], last *int, orderBy *generated.LoadBalancerPoolOrder, where *generated.LoadBalancerPoolWhereInput) (*generated.LoadBalancerPoolConnection, error) {
	if err := permissions.CheckAccess(ctx, obj.ID, actionLoadBalancerPoolGet); err != nil {
		return nil, err
	}

	return r.client.Pool.Query().Where(pool.OwnerID(obj.ID)).Paginate(ctx, after, first, before, last, generated.WithLoadBalancerPoolOrder(orderBy), generated.WithLoadBalancerPoolFilter(where.Filter))
}

// LoadBalancersProviders is the resolver for the loadBalancersProviders field.
func (r *resourceOwnerResolver) LoadBalancersProviders(ctx context.Context, obj *ResourceOwner, after *entgql.Cursor[gidx.PrefixedID], first *int, before *entgql.Cursor[gidx.PrefixedID], last *int, orderBy *generated.LoadBalancerOrder, where *generated.LoadBalancerProviderWhereInput) (*generated.LoadBalancerProviderConnection, error) {
	panic(fmt.Errorf("not implemented: LoadBalancersProviders - loadBalancersProviders"))
}

// ResourceOwner returns ResourceOwnerResolver implementation.
func (r *Resolver) ResourceOwner() ResourceOwnerResolver { return &resourceOwnerResolver{r} }

type resourceOwnerResolver struct{ *Resolver }
