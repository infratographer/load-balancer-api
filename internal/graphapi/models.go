package graphapi

import (
	"go.infratographer.com/x/gidx"

	"go.infratographer.com/load-balancer-api/internal/ent/generated"
)

// Location represents a Location in the graph for the bits load-balancer-api is able to return
type Location struct {
	ID              gidx.PrefixedID                   `json:"id"`
	LoadBalancers   *generated.LoadBalancerConnection `json:"loadBalancers"`
	scopedToOwnerID gidx.PrefixedID                   `json:"-"`
}

// IsEntity ensures the entity interface is met
func (Location) IsEntity() {}

// Owner represents a Location in the graph for the bits load-balancer-api is able to return
type Owner struct {
	ID                gidx.PrefixedID                       `json:"id"`
	LoadBalancers     *generated.LoadBalancerConnection     `json:"loadBalancers"`
	LoadBalancerPools *generated.LoadBalancerPoolConnection `json:"loadBalancerPools"`
}

// IsEntity ensures the entity interface is met
func (Owner) IsEntity() {}
