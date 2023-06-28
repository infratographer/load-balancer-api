package graphapi_test

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.infratographer.com/permissions-api/pkg/permissions"
	"go.infratographer.com/x/gidx"

	ent "go.infratographer.com/load-balancer-api/internal/ent/generated"
	"go.infratographer.com/load-balancer-api/internal/graphclient"
)

func TestQuery_loadBalancerProvider(t *testing.T) {
	ctx := context.Background()

	// Permit request
	ctx = context.WithValue(ctx, permissions.CheckerCtxKey, permissions.DefaultAllowChecker)

	p1 := ProviderBuilder{}.MustNew(ctx)
	p2 := ProviderBuilder{}.MustNew(ctx)

	testCases := []struct {
		TestName          string
		QueryID           gidx.PrefixedID
		ExpectedPrvovider *ent.Provider
		errorMsg          string
	}{
		{
			TestName:          "Happy Path - p1",
			QueryID:           p1.ID,
			ExpectedPrvovider: p1,
		},
		{
			TestName:          "Happy Path - p2",
			QueryID:           p2.ID,
			ExpectedPrvovider: p2,
		},
		{
			TestName: "No load balancer provider found with ID",
			QueryID:  gidx.MustNewID("testing"),
			errorMsg: "provider not found",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.TestName, func(t *testing.T) {
			resp, err := graphTestClient().GetLoadBalancerProvider(ctx, tt.QueryID)

			if tt.errorMsg != "" {
				require.Error(t, err)
				assert.ErrorContains(t, err, tt.errorMsg)
				assert.Nil(t, resp)

				return
			}

			require.NoError(t, err)
			require.NotNil(t, resp)
			require.NotNil(t, resp.LoadBalancerProvider)
			assert.EqualValues(t, tt.ExpectedPrvovider.Name, resp.LoadBalancerProvider.Name)
		})
	}
}

func TestFullProviderLifecycle(t *testing.T) {
	ctx := context.Background()

	// Permit request
	ctx = context.WithValue(ctx, permissions.CheckerCtxKey, permissions.DefaultAllowChecker)

	ownerID := gidx.MustNewID(ownerPrefix)
	name := gofakeit.DomainName()

	// create the Provider
	createdResp, err := graphTestClient().LoadBalancerProviderCreate(ctx, graphclient.CreateLoadBalancerProviderInput{
		Name:    name,
		OwnerID: ownerID,
	})

	require.NoError(t, err)
	require.NotNil(t, createdResp)
	require.NotNil(t, createdResp.LoadBalancerProviderCreate.LoadBalancerProvider)

	createdProv := createdResp.LoadBalancerProviderCreate.LoadBalancerProvider
	require.NotNil(t, createdProv.ID)
	require.Equal(t, name, createdProv.Name)
	assert.Equal(t, "loadpvd", createdProv.ID.Prefix())
	assert.Equal(t, ownerID, createdProv.Owner.ID)

	// Update the Provider
	newName := gofakeit.DomainName()
	updatedLBResp, err := graphTestClient().LoadBalancerProviderUpdate(ctx, createdProv.ID, graphclient.UpdateLoadBalancerProviderInput{Name: &newName})

	require.NoError(t, err)
	require.NotNil(t, updatedLBResp)
	require.NotNil(t, updatedLBResp.LoadBalancerProviderUpdate.LoadBalancerProvider)

	updatedLB := updatedLBResp.LoadBalancerProviderUpdate.LoadBalancerProvider
	require.EqualValues(t, createdProv.ID, updatedLB.ID)
	require.Equal(t, newName, updatedLB.Name)

	// Query the Provider
	queryLB, err := graphTestClient().GetLoadBalancerProvider(ctx, createdProv.ID)
	require.NoError(t, err)
	require.NotNil(t, queryLB)
	require.NotNil(t, queryLB.LoadBalancerProvider)
	require.Equal(t, newName, queryLB.LoadBalancerProvider.Name)

	// Delete the Provider
	deletedResp, err := graphTestClient().LoadBalancerProviderDelete(ctx, createdProv.ID)
	require.NoError(t, err)
	require.NotNil(t, deletedResp)
	require.NotNil(t, deletedResp.LoadBalancerProviderDelete)
	require.EqualValues(t, createdProv.ID, deletedResp.LoadBalancerProviderDelete.DeletedID.String())

	// Query the Provider to ensure it's no longer available
	deletedLB, err := graphTestClient().GetLoadBalancerProvider(ctx, createdProv.ID)
	require.Error(t, err)
	require.Nil(t, deletedLB)
	require.ErrorContains(t, err, "provider not found")
}
