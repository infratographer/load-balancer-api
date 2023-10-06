package graphapi_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.infratographer.com/permissions-api/pkg/permissions"
	"go.infratographer.com/permissions-api/pkg/permissions/mockpermissions"
	"go.infratographer.com/x/gidx"

	ent "go.infratographer.com/load-balancer-api/internal/ent/generated"
	"go.infratographer.com/load-balancer-api/internal/graphclient"
	"go.infratographer.com/load-balancer-api/internal/testutils"
)

func TestOwnerLoadBalancersResolver(t *testing.T) {
	ctx := context.Background()
	perms := new(mockpermissions.MockPermissions)
	perms.On("CreateAuthRelationships", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	ctx = perms.ContextWithHandler(ctx)

	// Permit request
	ctx = context.WithValue(ctx, permissions.CheckerCtxKey, permissions.DefaultAllowChecker)

	ownerID := gidx.MustNewID("testtnt")
	lb1 := (&testutils.LoadBalancerBuilder{OwnerID: ownerID, LocationID: "testloc-CCCafdsaf", Name: "lb-a"}).MustNew(ctx)
	lb2 := (&testutils.LoadBalancerBuilder{OwnerID: ownerID, LocationID: "testloc-AAAfasdf", Name: "lb-c"}).MustNew(ctx)
	lb3 := (&testutils.LoadBalancerBuilder{OwnerID: ownerID, LocationID: "testloc-BBBasdfa", Name: "lb-1"}).MustNew(ctx)
	(&testutils.LoadBalancerBuilder{}).MustNew(ctx)
	// Update LB1 so it's updated at is most recent
	lb1.Update().SaveX(ctx)

	testCases := []struct {
		TestName      string
		OrderBy       *graphclient.LoadBalancerOrder
		OwnerID       gidx.PrefixedID
		ResponseOrder []*ent.LoadBalancer
		errorMsg      string
	}{
		{
			TestName:      "Get Owner LoadBalancers - Ordered By NAME ASC",
			OrderBy:       &graphclient.LoadBalancerOrder{Field: "NAME", Direction: "ASC"},
			OwnerID:       ownerID,
			ResponseOrder: []*ent.LoadBalancer{lb3, lb1, lb2},
		},
		{
			TestName:      "Get Owner LoadBalancers - Ordered By NAME DESC",
			OrderBy:       &graphclient.LoadBalancerOrder{Field: "NAME", Direction: "DESC"},
			OwnerID:       ownerID,
			ResponseOrder: []*ent.LoadBalancer{lb2, lb1, lb3},
		},
		{
			TestName:      "Get Owner LoadBalancers - Ordered By CREATED_AT ASC",
			OrderBy:       &graphclient.LoadBalancerOrder{Field: "CREATED_AT", Direction: "ASC"},
			OwnerID:       ownerID,
			ResponseOrder: []*ent.LoadBalancer{lb1, lb2, lb3},
		},
		{
			TestName:      "Get Owner LoadBalancers - Ordered By CREATED_AT DESC",
			OrderBy:       &graphclient.LoadBalancerOrder{Field: "CREATED_AT", Direction: "DESC"},
			OwnerID:       ownerID,
			ResponseOrder: []*ent.LoadBalancer{lb3, lb2, lb1},
		},
		{
			TestName:      "Get Owner LoadBalancers - Ordered By UPDATED_AT ASC",
			OrderBy:       &graphclient.LoadBalancerOrder{Field: "UPDATED_AT", Direction: "ASC"},
			OwnerID:       ownerID,
			ResponseOrder: []*ent.LoadBalancer{lb2, lb3, lb1},
		},
		{
			TestName:      "Get Owner LoadBalancers - Ordered By UPDATED_AT DESC",
			OrderBy:       &graphclient.LoadBalancerOrder{Field: "UPDATED_AT", Direction: "DESC"},
			OwnerID:       ownerID,
			ResponseOrder: []*ent.LoadBalancer{lb1, lb3, lb2},
		},
		{
			TestName:      "Get Owner LoadBalancers - No LBs for Owner",
			OwnerID:       gidx.MustNewID(ownerPrefix),
			ResponseOrder: []*ent.LoadBalancer{},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.TestName, func(t *testing.T) {
			resp, err := graphTestClient().GetOwnerLoadBalancers(ctx, tt.OwnerID, tt.OrderBy)

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
