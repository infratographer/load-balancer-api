package graphapi_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.infratographer.com/x/gidx"

	ent "go.infratographer.com/load-balancer-api/internal/ent/generated"
	"go.infratographer.com/load-balancer-api/internal/graphclient"
)

func TestTenantLoadBalancersResolver(t *testing.T) {
	ctx := context.Background()
	tenantID := gidx.MustNewID("testtnt")
	lb1 := (&LoadBalancerBuilder{TenantID: tenantID, LocationID: "testloc-CCCafdsaf", Name: "lb-a"}).MustNew(ctx)
	lb2 := (&LoadBalancerBuilder{TenantID: tenantID, LocationID: "testloc-AAAfasdf", Name: "lb-c"}).MustNew(ctx)
	lb3 := (&LoadBalancerBuilder{TenantID: tenantID, LocationID: "testloc-BBBasdfa", Name: "lb-1"}).MustNew(ctx)
	(&LoadBalancerBuilder{}).MustNew(ctx)
	// Update LB1 so it's updated at is most recent
	lb1.Update().SaveX(ctx)

	testCases := []struct {
		TestName      string
		OrderBy       *graphclient.LoadBalancerOrder
		TenantID      gidx.PrefixedID
		ResponseOrder []*ent.LoadBalancer
		errorMsg      string
	}{
		{
			TestName:      "Get Tenant LoadBalancers - Ordered By NAME ASC",
			OrderBy:       &graphclient.LoadBalancerOrder{Field: "NAME", Direction: "ASC"},
			TenantID:      tenantID,
			ResponseOrder: []*ent.LoadBalancer{lb3, lb1, lb2},
		},
		{
			TestName:      "Get Tenant LoadBalancers - Ordered By NAME DESC",
			OrderBy:       &graphclient.LoadBalancerOrder{Field: "NAME", Direction: "DESC"},
			TenantID:      tenantID,
			ResponseOrder: []*ent.LoadBalancer{lb2, lb1, lb3},
		},
		{
			TestName:      "Get Tenant LoadBalancers - Ordered By CREATED_AT ASC",
			OrderBy:       &graphclient.LoadBalancerOrder{Field: "CREATED_AT", Direction: "ASC"},
			TenantID:      tenantID,
			ResponseOrder: []*ent.LoadBalancer{lb1, lb2, lb3},
		},
		{
			TestName:      "Get Tenant LoadBalancers - Ordered By CREATED_AT DESC",
			OrderBy:       &graphclient.LoadBalancerOrder{Field: "CREATED_AT", Direction: "DESC"},
			TenantID:      tenantID,
			ResponseOrder: []*ent.LoadBalancer{lb3, lb2, lb1},
		},
		{
			TestName:      "Get Tenant LoadBalancers - Ordered By UPDATED_AT ASC",
			OrderBy:       &graphclient.LoadBalancerOrder{Field: "UPDATED_AT", Direction: "ASC"},
			TenantID:      tenantID,
			ResponseOrder: []*ent.LoadBalancer{lb2, lb3, lb1},
		},
		{
			TestName:      "Get Tenant LoadBalancers - Ordered By UPDATED_AT DESC",
			OrderBy:       &graphclient.LoadBalancerOrder{Field: "UPDATED_AT", Direction: "DESC"},
			TenantID:      tenantID,
			ResponseOrder: []*ent.LoadBalancer{lb1, lb3, lb2},
		},
		{
			TestName:      "Get Tenant LoadBalancers - No LBs for Tenant",
			TenantID:      gidx.MustNewID(tenantPrefix),
			ResponseOrder: []*ent.LoadBalancer{},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.TestName, func(t *testing.T) {
			resp, err := newGraphTestClient().GetTenantLoadBalancers(ctx, tt.TenantID, tt.OrderBy)

			if tt.errorMsg != "" {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tt.errorMsg)

				return
			}

			require.Len(t, resp.Entities[0].LoadBalancers.Edges, len(tt.ResponseOrder))
			for i, lb := range tt.ResponseOrder {
				respLB := resp.Entities[0].LoadBalancers.Edges[i].Node
				require.Equal(t, lb.ID, respLB.ID)
				require.Equal(t, lb.Name, respLB.Name)
			}
		})
	}
}
