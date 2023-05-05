package graphapi

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.30

import (
	"context"

	"entgo.io/contrib/entgql"
	"go.infratographer.com/load-balancer-api/internal/ent/generated"
	"go.infratographer.com/load-balancer-api/internal/ent/generated/loadbalancer"
	"go.infratographer.com/x/gidx"
)

// Location is the resolver for the location field.
func (r *loadBalancerResolver) Location(ctx context.Context, obj *generated.LoadBalancer) (*Location, error) {
	return &Location{ID: obj.LocationID, scopedToTenantID: obj.TenantID}, nil
}

// LoadBalancers is the resolver for the loadBalancers field.
func (r *locationResolver) LoadBalancers(ctx context.Context, obj *Location, after *entgql.Cursor[gidx.PrefixedID], first *int, before *entgql.Cursor[gidx.PrefixedID], last *int, orderBy *generated.LoadBalancerOrder, where *generated.LoadBalancerWhereInput) (*generated.LoadBalancerConnection, error) {
	return r.client.LoadBalancer.Query().Where(loadbalancer.TenantID(obj.scopedToTenantID)).Paginate(ctx, after, first, before, last, generated.WithLoadBalancerOrder(orderBy), generated.WithLoadBalancerFilter(where.Filter))
}

// Location returns LocationResolver implementation.
func (r *Resolver) Location() LocationResolver { return &locationResolver{r} }

type locationResolver struct{ *Resolver }
