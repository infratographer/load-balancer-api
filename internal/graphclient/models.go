package graphclient

import (
	"go.infratographer.com/x/gidx"
)

// Mutations represents the responses returned from mutation calls
type Mutations struct {
	LoadBalancerCreate LoadBalancerCreate `json:"loadBalancerCreate"`
	LoadBalancerDelete LoadBalancerDelete `json:"loadBalancerDelete"`
	LoadBalancerUpdate LoadBalancerUpdate `json:"loadBalancerUpdate"`

	PoolCreate PoolCreate `json:"poolCreate"`
	PoolDelete PoolDelete `json:"poolDelete"`
	PoolUpdate PoolUpdate `json:"poolUpdate"`
}

// LoadBalancerCreate response
type LoadBalancerCreate struct {
	LoadBalancer *LoadBalancer `json:"loadBalancer"`
}

// LoadBalancerDelete response
type LoadBalancerDelete struct {
	DeletedID gidx.PrefixedID `json:"deletedID"`
}

// LoadBalancerUpdate response
type LoadBalancerUpdate struct {
	LoadBalancer *LoadBalancer `json:"loadBalancer"`
}

// LoadBalancer represents a GraphQL LoadBalancer type
type LoadBalancer struct {
	ID         gidx.PrefixedID `json:"id"`
	Name       string          `json:"name"`
	CreatedAt  string          `json:"createdAt"`
	UpdatedAt  string          `json:"updatedAt"`
	Provider   *Provider       `json:"loadBalancerProvider"`
	LocationID gidx.PrefixedID `json:"locationID"`
	TenantID   gidx.PrefixedID `json:"tenantID"`
}

// Provider represents a GraphQL Provider type
type Provider struct {
	ID        gidx.PrefixedID `json:"id"`
	Name      string          `json:"name"`
	CreatedAt string          `json:"createdAt"`
	UpdatedAt string          `json:"updatedAt"`
}

/////////////////////////////////////////////////////////////////
// Pools
/////////////////////////////////////////////////////////////////

// PoolCreate response
type PoolCreate struct {
	Pool *Pool `json:"loadBalancerPool"`
}

// PoolDelete response
type PoolDelete struct {
	DeletedID gidx.PrefixedID `json:"deletedID"`
}

// PoolUpdate response
type PoolUpdate struct {
	Pool *Pool `json:"loadBalancerPool"`
}

// Pool represents a GraphQL Pool type
type Pool struct {
	ID        gidx.PrefixedID `json:"id"`
	Name      string          `json:"name"`
	CreatedAt string          `json:"createdAt"`
	UpdatedAt string          `json:"updatedAt"`
	Protocol  string          `json:"protocol"`
	TenantID  gidx.PrefixedID `json:"tenantID"`
}
