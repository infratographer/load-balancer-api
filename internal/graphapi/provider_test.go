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
	"go.infratographer.com/load-balancer-api/internal/testutils"
)

func TestQuery_loadBalancerProvider(t *testing.T) {
	ctx := context.Background()

	// Permit request
	ctx = context.WithValue(ctx, permissions.CheckerCtxKey, permissions.DefaultAllowChecker)

	p1 := (&testutils.ProviderBuilder{}).MustNew(ctx)
	p2 := (&testutils.ProviderBuilder{}).MustNew(ctx)

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
		{
			TestName: "Invalid load balancer provider ID",
			QueryID:  gidx.PrefixedID("invalid"),
			errorMsg: "invalid id",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.TestName, func(t *testing.T) {
			tt := tt
			t.Parallel()

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

func TestCreate_Provider(t *testing.T) {
	ctx := context.Background()

	// Permit request
	ctx = context.WithValue(ctx, permissions.CheckerCtxKey, permissions.DefaultAllowChecker)

	ownerID := gidx.MustNewID(ownerPrefix)
	name := gofakeit.DomainName()

	testCases := []struct {
		TestName   string
		Input      graphclient.CreateLoadBalancerProviderInput
		ExpectedLB *ent.LoadBalancerProvider
		errorMsg   string
	}{
		{
			TestName: "creates provider",
			Input:    graphclient.CreateLoadBalancerProviderInput{Name: name, OwnerID: ownerID},
			ExpectedLB: &ent.LoadBalancerProvider{
				Name:    name,
				OwnerID: ownerID,
			},
		},
		{
			TestName: "fails to create provider with empty name",
			Input:    graphclient.CreateLoadBalancerProviderInput{Name: "", OwnerID: ownerID},
			errorMsg: "value is less than the required length",
		},
		{
			TestName: "fails to create provider with empty ownerID",
			Input:    graphclient.CreateLoadBalancerProviderInput{Name: name, OwnerID: ""},
			errorMsg: "value is less than the required length",
		},
		{
			TestName: "fails to create provider with invalid ownerID",
			Input:    graphclient.CreateLoadBalancerProviderInput{Name: name, OwnerID: gidx.PrefixedID("invalid")},
			errorMsg: "invalid id",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.TestName, func(t *testing.T) {
			tt := tt
			t.Parallel()

			resp, err := graphTestClient().LoadBalancerProviderCreate(ctx, tt.Input)

			if tt.errorMsg != "" {
				require.Error(t, err)
				assert.ErrorContains(t, err, tt.errorMsg)
				assert.Nil(t, resp)

				return
			}

			require.NoError(t, err)
			require.NotNil(t, resp)
			require.NotNil(t, resp.LoadBalancerProviderCreate)

			createdProvider := resp.LoadBalancerProviderCreate.LoadBalancerProvider
			assert.Equal(t, tt.ExpectedLB.Name, createdProvider.Name)
			assert.Equal(t, "loadpvd", createdProvider.ID.Prefix())
			assert.Equal(t, ownerID, createdProvider.Owner.ID)
		})
	}
}

func TestUpdate_Provider(t *testing.T) {
	ctx := context.Background()

	// Permit request
	ctx = context.WithValue(ctx, permissions.CheckerCtxKey, permissions.DefaultAllowChecker)

	prov := (&testutils.ProviderBuilder{}).MustNew(ctx)
	updateName := gofakeit.DomainName()

	testCases := []struct {
		TestName         string
		ID               gidx.PrefixedID
		Input            graphclient.UpdateLoadBalancerProviderInput
		ExpectedProvider *ent.LoadBalancerProvider
		errorMsg         string
	}{
		{
			TestName: "updates provider",
			ID:       prov.ID,
			Input:    graphclient.UpdateLoadBalancerProviderInput{Name: &updateName},
			ExpectedProvider: &ent.LoadBalancerProvider{
				Name:    updateName,
				ID:      prov.ID,
				OwnerID: prov.OwnerID,
			},
		},
		{
			TestName: "fails to update name to empty",
			ID:       prov.ID,
			Input:    graphclient.UpdateLoadBalancerProviderInput{Name: newString("")},
			errorMsg: "value is less than the required length",
		},
		{
			TestName: "fails to update provider that does not exist",
			ID:       gidx.PrefixedID("loadpvd-dne"),
			Input:    graphclient.UpdateLoadBalancerProviderInput{},
			errorMsg: "provider not found",
		},
		{
			TestName: "fails to update provider with invalid id",
			ID:       gidx.PrefixedID("invalid"),
			Input:    graphclient.UpdateLoadBalancerProviderInput{},
			errorMsg: "invalid id",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.TestName, func(t *testing.T) {
			resp, err := graphTestClient().LoadBalancerProviderUpdate(ctx, tt.ID, tt.Input)

			if tt.errorMsg != "" {
				require.Error(t, err)
				assert.ErrorContains(t, err, tt.errorMsg)
				assert.Nil(t, resp)

				return
			}

			require.NoError(t, err)
			require.NotNil(t, resp)
			require.NotNil(t, resp.LoadBalancerProviderUpdate)

			updatedProvider := resp.LoadBalancerProviderUpdate.LoadBalancerProvider
			assert.Equal(t, tt.ExpectedProvider.Name, updatedProvider.Name)
			assert.Equal(t, prov.ID, updatedProvider.ID)
		})
	}
}

func TestDelete_Provider(t *testing.T) {
	ctx := context.Background()

	// Permit request
	ctx = context.WithValue(ctx, permissions.CheckerCtxKey, permissions.DefaultAllowChecker)

	prov := (&testutils.ProviderBuilder{}).MustNew(ctx)

	testCases := []struct {
		TestName   string
		Input      gidx.PrefixedID
		ExpectedID gidx.PrefixedID
		errorMsg   string
	}{
		{
			TestName:   "deletes provider",
			Input:      prov.ID,
			ExpectedID: prov.ID,
		},
		{
			TestName: "fails to delete provider that does not exist",
			Input:    gidx.PrefixedID("loadpvd-dne"),
			errorMsg: "provider not found",
		},
		{
			TestName: "fails to delete empty provider ID",
			Input:    gidx.PrefixedID(""),
			errorMsg: "provider not found",
		},
		{
			TestName: "fails to delete invalid gidx id",
			Input:    gidx.PrefixedID("not-a-valid-gidx-id"),
			errorMsg: "invalid id",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.TestName, func(t *testing.T) {
			resp, err := graphTestClient().LoadBalancerProviderDelete(ctx, tt.Input)

			if tt.errorMsg != "" {
				require.Error(t, err)
				assert.ErrorContains(t, err, tt.errorMsg)
				assert.Nil(t, resp)

				return
			}

			require.NoError(t, err)
			require.NotNil(t, resp)
			require.NotNil(t, resp.LoadBalancerProviderDelete)

			deletedProvider := resp.LoadBalancerProviderDelete
			assert.Equal(t, tt.Input, deletedProvider.DeletedID)
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
