package graphapi_test

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.infratographer.com/x/gidx"

	ent "go.infratographer.com/load-balancer-api/internal/ent/generated"
	"go.infratographer.com/load-balancer-api/internal/graphapi"
	"go.infratographer.com/load-balancer-api/internal/graphclient"
)

func TestQuery_loadBalancer(t *testing.T) {
	ctx := context.Background()
	graphClient := graphclient.New(graphapi.NewResolver(EntClient))

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
			resp, err := graphClient.QueryLoadBalancer(tt.QueryID)

			if tt.errorMsg != "" {
				require.Error(t, err)
				assert.ErrorContains(t, err, tt.errorMsg)
				assert.Nil(t, resp)

				return
			}

			require.NoError(t, err)
			require.NotNil(t, resp)
			assert.EqualValues(t, tt.ExpectedLB.Name, resp.Name)
		})
	}
}

func TestFullLoadBalancerLifecycle(t *testing.T) {
	ctx := context.Background()
	graphClient := graphclient.New(graphapi.NewResolver(EntClient))

	prov := (&ProviderBuilder{}).MustNew(ctx)
	tenantID := gidx.MustNewID(tenantPrefix)
	locationID := gidx.MustNewID(locationPrefix)
	name := gofakeit.DomainName()

	// create the LB
	createdLB, err := graphClient.LoadBalancerCreate(ent.CreateLoadBalancerInput{
		Name:       name,
		ProviderID: prov.ID,
		TenantID:   tenantID,
		LocationID: locationID,
	})

	require.NoError(t, err)
	require.NotNil(t, createdLB)
	require.NotNil(t, createdLB.ID)
	require.Equal(t, name, createdLB.Name)
	assert.Equal(t, "loadbal", createdLB.ID.Prefix())
	assert.Equal(t, prov.ID, createdLB.Provider.ID)
	assert.Equal(t, locationID, createdLB.LocationID)
	assert.Equal(t, tenantID, createdLB.TenantID)

	// Update the LB
	newName := gofakeit.DomainName()
	updatedLB, err := graphClient.LoadBalancerUpdate(createdLB.ID, ent.UpdateLoadBalancerInput{Name: &newName})

	require.NoError(t, err)
	require.NotNil(t, updatedLB)
	require.EqualValues(t, createdLB.ID, updatedLB.ID)
	require.Equal(t, newName, updatedLB.Name)

	// Query the LB
	queryLB, err := graphClient.QueryLoadBalancer(createdLB.ID)
	require.NoError(t, err)
	require.NotNil(t, queryLB)
	require.Equal(t, newName, queryLB.Name)

	// Delete the LB
	deletedID, err := graphClient.LoadBalancerDelete(createdLB.ID)
	require.NoError(t, err)
	require.NotNil(t, deletedID)
	require.Equal(t, createdLB.ID, deletedID)

	// Query the LB to ensure it's no longer available
	deletedLB, err := graphClient.QueryLoadBalancer(createdLB.ID)
	require.Error(t, err)
	require.Nil(t, deletedLB)
	require.ErrorContains(t, err, "load_balancer not found")
}
