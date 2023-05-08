package graphapi_test

import (
	"context"

	"github.com/brianvoe/gofakeit/v6"
	"go.infratographer.com/x/gidx"

	ent "go.infratographer.com/load-balancer-api/internal/ent/generated"
	"go.infratographer.com/load-balancer-api/internal/ent/generated/pool"
)

type ProviderBuilder struct {
	Name     string
	TenantID gidx.PrefixedID
}

func (p ProviderBuilder) MustNew(ctx context.Context) *ent.Provider {
	if p.Name == "" {
		p.Name = gofakeit.JobTitle()
	}

	if p.TenantID == "" {
		p.TenantID = gidx.MustNewID(tenantPrefix)
	}

	return EntClient.Provider.Create().SetName(p.Name).SetTenantID(p.TenantID).SaveX(ctx)
}

type LoadBalancerBuilder struct {
	Name       string
	TenantID   gidx.PrefixedID
	LocationID gidx.PrefixedID
	Provider   *ent.Provider
}

func (b LoadBalancerBuilder) MustNew(ctx context.Context) *ent.LoadBalancer {
	if b.Provider == nil {
		pb := &ProviderBuilder{TenantID: b.TenantID}
		b.Provider = pb.MustNew(ctx)
	}

	if b.Name == "" {
		b.Name = gofakeit.AppName()
	}

	if b.TenantID == "" {
		b.TenantID = b.Provider.TenantID
	}

	if b.LocationID == "" {
		b.LocationID = gidx.MustNewID(locationPrefix)
	}

	return EntClient.LoadBalancer.Create().SetName(b.Name).SetTenantID(b.TenantID).SetLocationID(b.LocationID).SetProvider(b.Provider).SaveX(ctx)
}

type PortBuilder struct {
	Name           string
	LoadBalancerID gidx.PrefixedID
	Number         int
}

func (p PortBuilder) MustNew(ctx context.Context) *ent.Port {
	if p.Name == "" {
		p.Name = gofakeit.AppName()
	}

	if p.LoadBalancerID == "" {
		p.LoadBalancerID = gidx.MustNewID(lbPrefix)
	}

	if p.Number == 0 {
		p.Number = gofakeit.Number(1, 65535)
	}

	return EntClient.Port.Create().SetName(p.Name).SetLoadBalancerID(p.LoadBalancerID).SetNumber(p.Number).SaveX(ctx)
}

type PoolBuilder struct {
	Name     string
	TenantID gidx.PrefixedID
	Protocol pool.Protocol
}

func (p *PoolBuilder) MustNew(ctx context.Context) *ent.Pool {
	if p.Name == "" {
		p.Name = gofakeit.AppName()
	}

	if p.TenantID == "" {
		p.TenantID = gidx.MustNewID(tenantPrefix)
	}

	if p.Protocol == "" {
		p.Protocol = pool.Protocol(gofakeit.RandomString([]string{"tcp", "udp"}))
	}

	return EntClient.Pool.Create().SetName(p.Name).SetTenantID(p.TenantID).SetProtocol(p.Protocol).SaveX(ctx)
}

type OriginBuilder struct {
	Name       string
	Target     string
	PortNumber int
	Active     bool
	PoolID     gidx.PrefixedID
}

func (o *OriginBuilder) MustNew(ctx context.Context) *ent.Origin {
	if o.Name == "" {
		o.Name = gofakeit.AppName()
	}

	if o.Target == "" {
		o.Target = gofakeit.IPv4Address()
	}

	if o.PortNumber == 0 {
		o.PortNumber = gofakeit.Number(1, 65535)
	}

	if o.PoolID == "" {
		pb := &PoolBuilder{}
		o.PoolID = pb.MustNew(ctx).ID
	}

	return EntClient.Origin.Create().SetName(o.Name).SetTarget(o.Target).SetPortNumber(o.PortNumber).SetActive(o.Active).SetPoolID(o.PoolID).SaveX(ctx)
}
