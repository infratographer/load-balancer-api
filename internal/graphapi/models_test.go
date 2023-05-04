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

func (p *ProviderBuilder) prepare(_ context.Context) {
	if p.Name == "" {
		p.Name = gofakeit.JobTitle()
	}

	if p.TenantID == "" {
		p.TenantID = gidx.MustNewID(tenantPrefix)
	}
}

func (p *ProviderBuilder) MustNew(ctx context.Context) *ent.Provider {
	p.prepare(ctx)
	return EntClient.Provider.Create().SetName(p.Name).SetTenantID(p.TenantID).SaveX(ctx)
}

type LoadBalancerBuilder struct {
	Name       string
	TenantID   gidx.PrefixedID
	LocationID gidx.PrefixedID
	Provider   *ent.Provider
}

func (b *LoadBalancerBuilder) prepare(ctx context.Context) {
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
}

func (b *LoadBalancerBuilder) MustNew(ctx context.Context) *ent.LoadBalancer {
	b.prepare(ctx)
	return EntClient.LoadBalancer.Create().SetName(b.Name).SetTenantID(b.TenantID).SetLocationID(b.LocationID).SetProvider(b.Provider).SaveX(ctx)
}
