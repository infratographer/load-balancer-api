package graphapi_test

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.infratographer.com/x/gidx"

	ent "go.infratographer.com/load-balancer-api/internal/ent/generated"
	"go.infratographer.com/load-balancer-api/internal/graphclient"
)

func TestQuery_loadBalancer(t *testing.T) {
	ctx := context.Background()
	lb1 := (&LoadBalancerBuilder{}).MustNew(ctx)
	lb2 := (&LoadBalancerBuilder{}).MustNew(ctx)

	testCases := []struct {
		TestName   string
		QueryID    gidx.PrefixedID
		ExpectedLB *ent.LoadBalancer
		errorMsg   string
	}{
		{
			TestName:   "Happy Path - lb1",
			QueryID:    lb1.ID,
			ExpectedLB: lb1,
		},
		{
			TestName:   "Happy Path - lb2",
			QueryID:    lb2.ID,
			ExpectedLB: lb2,
		},
		{
			TestName: "No load balancer found with ID",
			QueryID:  gidx.MustNewID("testing"),
			errorMsg: "load_balancer not found",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.TestName, func(t *testing.T) {
			resp, err := newGraphTestClient().GetLoadBalancer(ctx, tt.QueryID)

			if tt.errorMsg != "" {
				require.Error(t, err)
				assert.ErrorContains(t, err, tt.errorMsg)
				assert.Nil(t, resp)

				return
			}

			require.NoError(t, err)
			require.NotNil(t, resp)
			require.NotNil(t, resp.LoadBalancer)
			assert.EqualValues(t, tt.ExpectedLB.Name, resp.LoadBalancer.Name)
		})
	}
}

func TestFullLoadBalancerLifecycle(t *testing.T) {
	ctx := context.Background()
	prov := (&ProviderBuilder{}).MustNew(ctx)
	tenantID := gidx.MustNewID(tenantPrefix)
	locationID := gidx.MustNewID(locationPrefix)
	name := gofakeit.DomainName()

	// create the LB
	createdLBResp, err := newGraphTestClient().LoadBalancerCreate(ctx, graphclient.CreateLoadBalancerInput{
		Name:       name,
		ProviderID: prov.ID,
		TenantID:   tenantID,
		LocationID: locationID,
	})

	require.NoError(t, err)
	require.NotNil(t, createdLBResp)
	require.NotNil(t, createdLBResp.LoadBalancerCreate.LoadBalancer)

	createdLB := createdLBResp.LoadBalancerCreate.LoadBalancer
	require.NotNil(t, createdLB.ID)
	require.Equal(t, name, createdLB.Name)
	assert.Equal(t, "loadbal", createdLB.ID.Prefix())
	assert.Equal(t, prov.ID, createdLB.LoadBalancerProvider.ID)
	assert.Equal(t, locationID, createdLB.Location.ID)
	assert.Equal(t, tenantID, createdLB.Tenant.ID)

	// Update the LB
	newName := gofakeit.DomainName()
	updatedLBResp, err := newGraphTestClient().LoadBalancerUpdate(ctx, createdLB.ID, graphclient.UpdateLoadBalancerInput{Name: &newName})

	require.NoError(t, err)
	require.NotNil(t, updatedLBResp)
	require.NotNil(t, updatedLBResp.LoadBalancerUpdate.LoadBalancer)

	updatedLB := updatedLBResp.LoadBalancerUpdate.LoadBalancer
	require.EqualValues(t, createdLB.ID, updatedLB.ID)
	require.Equal(t, newName, updatedLB.Name)

	// Query the LB
	queryLB, err := newGraphTestClient().GetLoadBalancer(ctx, createdLB.ID)
	require.NoError(t, err)
	require.NotNil(t, queryLB)
	require.NotNil(t, queryLB.LoadBalancer)
	require.Equal(t, newName, queryLB.LoadBalancer.Name)

	// Delete the LB
	deletedResp, err := newGraphTestClient().LoadBalancerDelete(ctx, createdLB.ID)
	require.NoError(t, err)
	require.NotNil(t, deletedResp)
	require.NotNil(t, deletedResp.LoadBalancerDelete)
	require.EqualValues(t, createdLB.ID, deletedResp.LoadBalancerDelete.DeletedID.String())

	// Query the LB to ensure it's no longer available
	deletedLB, err := newGraphTestClient().GetLoadBalancer(ctx, createdLB.ID)
	require.Error(t, err)
	require.Nil(t, deletedLB)
	require.ErrorContains(t, err, "load_balancer not found")
}
