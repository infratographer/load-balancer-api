package graphapi_test

import (
	"context"

	"github.com/brianvoe/gofakeit/v6"
	"go.infratographer.com/x/gidx"

	ent "go.infratographer.com/load-balancer-api/internal/ent/generated"
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
