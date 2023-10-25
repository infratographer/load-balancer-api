package testutils

import (
	"context"

	"github.com/brianvoe/gofakeit/v6"
	"go.infratographer.com/x/gidx"

	ent "go.infratographer.com/load-balancer-api/internal/ent/generated"
	"go.infratographer.com/load-balancer-api/internal/ent/generated/pool"
)

const (
	ownerPrefix    = "testown"
	locationPrefix = "testloc"
	lbPrefix       = "loadbal"
	minPortNum     = 1
	maxPortNum     = 65535
)

// ProviderBuilder is a provider-like struct for use in generating a provider using the ent client
type ProviderBuilder struct {
	Name    string
	OwnerID gidx.PrefixedID
}

// MustNew creates a provider from the receiver
func (p *ProviderBuilder) MustNew(ctx context.Context) *ent.Provider {
	if p.Name == "" {
		p.Name = gofakeit.JobTitle()
	}

	if p.OwnerID == "" {
		p.OwnerID = gidx.MustNewID(ownerPrefix)
	}

	return EntClient.Provider.Create().SetName(p.Name).SetOwnerID(p.OwnerID).SaveX(ctx)
}

// LoadBalancerBuilder is a loadbalancer-like struct for use in generating a loadbalancer using the ent client
type LoadBalancerBuilder struct {
	Name       string
	OwnerID    gidx.PrefixedID
	LocationID gidx.PrefixedID
	Provider   *ent.Provider
}

// MustNew creates a loadbalancer from the receiver
func (b *LoadBalancerBuilder) MustNew(ctx context.Context) *ent.LoadBalancer {
	if b.Provider == nil {
		pb := &ProviderBuilder{OwnerID: b.OwnerID}
		b.Provider = pb.MustNew(ctx)
	}

	if b.Name == "" {
		b.Name = gofakeit.AppName()
	}

	if b.OwnerID == "" {
		b.OwnerID = b.Provider.OwnerID
	}

	if b.LocationID == "" {
		b.LocationID = gidx.MustNewID(locationPrefix)
	}

	return EntClient.LoadBalancer.Create().SetName(b.Name).SetOwnerID(b.OwnerID).SetLocationID(b.LocationID).SetProvider(b.Provider).SaveX(ctx)
}

// PortBuilder is a port-like struct for use in generating a port using the ent client
type PortBuilder struct {
	Name           string
	LoadBalancerID gidx.PrefixedID
	Number         int
	PoolIDs        []gidx.PrefixedID
}

// MustNew creates a port from the receiver
func (p *PortBuilder) MustNew(ctx context.Context) *ent.Port {
	if p.Name == "" {
		p.Name = gofakeit.AppName()
	}

	if p.LoadBalancerID == "" {
		p.LoadBalancerID = gidx.MustNewID(lbPrefix)
	}

	if p.Number == 0 {
		p.Number = gofakeit.Number(minPortNum, maxPortNum)
	}

	return EntClient.Port.Create().SetName(p.Name).SetLoadBalancerID(p.LoadBalancerID).SetNumber(p.Number).AddPoolIDs(p.PoolIDs...).SaveX(ctx)
}

// PoolBuilder is a pool-like struct for use in generating a pool using the ent client
type PoolBuilder struct {
	Name     string
	OwnerID  gidx.PrefixedID
	Protocol pool.Protocol
}

// MustNew creates a pool from the receiver
func (p *PoolBuilder) MustNew(ctx context.Context) *ent.Pool {
	if p.Name == "" {
		p.Name = gofakeit.AppName()
	}

	if p.OwnerID == "" {
		p.OwnerID = gidx.MustNewID(ownerPrefix)
	}

	if p.Protocol == "" {
		p.Protocol = pool.Protocol(gofakeit.RandomString([]string{"tcp", "udp"}))
	}

	return EntClient.Pool.Create().SetName(p.Name).SetOwnerID(p.OwnerID).SetProtocol(p.Protocol).SaveX(ctx)
}

// OriginBuilder is an origin-like struct for use in generating an origin using the ent client
type OriginBuilder struct {
	Name       string
	Target     string
	PortNumber int
	Active     bool
	PoolID     gidx.PrefixedID
}

// MustNew creates an origin from the receiver
func (o *OriginBuilder) MustNew(ctx context.Context) *ent.Origin {
	if o.Name == "" {
		o.Name = gofakeit.AppName()
	}

	if o.Target == "" {
		o.Target = gofakeit.IPv4Address()
	}

	if o.PortNumber == 0 {
		o.PortNumber = gofakeit.Number(minPortNum, maxPortNum)
	}

	if o.PoolID == "" {
		pb := &PoolBuilder{}
		o.PoolID = pb.MustNew(ctx).ID
	}

	return EntClient.Origin.Create().SetName(o.Name).SetTarget(o.Target).SetPortNumber(o.PortNumber).SetActive(o.Active).SetPoolID(o.PoolID).SaveX(ctx)
}
