package graphapi_test

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.infratographer.com/permissions-api/pkg/permissions"
	"go.infratographer.com/permissions-api/pkg/permissions/mockpermissions"
	"go.infratographer.com/x/gidx"

	"go.infratographer.com/load-balancer-api/internal/config"
	ent "go.infratographer.com/load-balancer-api/internal/ent/generated"
	"go.infratographer.com/load-balancer-api/internal/graphclient"
	"go.infratographer.com/load-balancer-api/internal/testutils"
)

func TestQuery_loadBalancer(t *testing.T) {
	ctx := context.Background()
	perms := new(mockpermissions.MockPermissions)
	perms.On("CreateAuthRelationships", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	ctx = perms.ContextWithHandler(ctx)

	// Permit request
	ctx = context.WithValue(ctx, permissions.CheckerCtxKey, permissions.DefaultAllowChecker)

	lb1 := (&testutils.LoadBalancerBuilder{}).MustNew(ctx)
	lb2 := (&testutils.LoadBalancerBuilder{}).MustNew(ctx)

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
		{
			TestName: "invalid gidx format",
			QueryID:  "test-invalid-id",
			errorMsg: "invalid id",
		},
		{
			TestName: "empty loadbalancer id",
			QueryID:  "test-invalid-id",
			errorMsg: "invalid id",
		},
		{
			TestName: "whitespace loadbalancer id",
			QueryID:  "test-invalid-id",
			errorMsg: "invalid id",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.TestName, func(t *testing.T) {
			tt := tt

			t.Parallel()

			resp, err := graphTestClient().GetLoadBalancer(ctx, tt.QueryID)

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

func TestCreate_loadBalancer(t *testing.T) {
	ctx := context.Background()
	perms := new(mockpermissions.MockPermissions)
	perms.On("CreateAuthRelationships", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	ctx = perms.ContextWithHandler(ctx)

	// Permit request
	ctx = context.WithValue(ctx, permissions.CheckerCtxKey, permissions.DefaultAllowChecker)

	prov := (&testutils.ProviderBuilder{}).MustNew(ctx)
	ownerID := gidx.MustNewID(ownerPrefix)
	locationID := gidx.MustNewID(locationPrefix)
	name := gofakeit.DomainName()

	testCases := []struct {
		TestName   string
		Input      graphclient.CreateLoadBalancerInput
		ExpectedLB *ent.LoadBalancer
		errorMsg   string
	}{
		{
			TestName: "creates loadbalancer",
			Input:    graphclient.CreateLoadBalancerInput{Name: name, ProviderID: prov.ID, OwnerID: ownerID, LocationID: locationID},
			ExpectedLB: &ent.LoadBalancer{
				Name:       name,
				ProviderID: prov.ID,
				OwnerID:    ownerID,
				LocationID: locationID,
			},
		},
		{
			TestName: "fails to create loadbalancer with empty name",
			Input:    graphclient.CreateLoadBalancerInput{Name: "", ProviderID: prov.ID, OwnerID: ownerID, LocationID: locationID},
			errorMsg: "must not be empty",
		},
		{
			TestName: "fails to create loadbalancer with empty ownerID",
			Input:    graphclient.CreateLoadBalancerInput{Name: name, ProviderID: prov.ID, OwnerID: "", LocationID: locationID},
			errorMsg: "must not be empty",
		},
		{
			TestName: "fails to create loadbalancer with empty locationID",
			Input:    graphclient.CreateLoadBalancerInput{Name: name, ProviderID: prov.ID, OwnerID: ownerID, LocationID: ""},
			errorMsg: "must not be empty",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.TestName, func(t *testing.T) {
			tt := tt

			t.Parallel()

			resp, err := graphTestClient().LoadBalancerCreate(ctx, tt.Input)

			if tt.errorMsg != "" {
				require.Error(t, err)
				assert.ErrorContains(t, err, tt.errorMsg)
				assert.Nil(t, resp)

				return
			}

			require.NoError(t, err)
			require.NotNil(t, resp)
			require.NotNil(t, resp.LoadBalancerCreate)

			createdLB := resp.LoadBalancerCreate.LoadBalancer
			assert.Equal(t, tt.ExpectedLB.Name, createdLB.Name)
			assert.Equal(t, "loadbal", createdLB.ID.Prefix())
			assert.Equal(t, prov.ID, createdLB.LoadBalancerProvider.ID)
			assert.Equal(t, locationID, createdLB.Location.ID)
			assert.Equal(t, ownerID, createdLB.Owner.ID)
		})
	}
}

func TestCreate_loadBalancer_limit(t *testing.T) {
	ctx := context.Background()
	perms := new(mockpermissions.MockPermissions)
	perms.On("CreateAuthRelationships", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	ctx = perms.ContextWithHandler(ctx)

	// Permit request
	ctx = context.WithValue(ctx, permissions.CheckerCtxKey, permissions.DefaultAllowChecker)

	prov := (&testutils.ProviderBuilder{}).MustNew(ctx)
	locationID := gidx.MustNewID(locationPrefix)
	name := gofakeit.DomainName()

	config.AppConfig.LoadBalancerLimit = 3

	testCases := []struct {
		TestName string
		lbCount  int
		Input    graphclient.CreateLoadBalancerInput
		errorMsg string
	}{
		{
			TestName: "creates loadbalancers - under limit",
			Input:    graphclient.CreateLoadBalancerInput{Name: name, ProviderID: prov.ID, OwnerID: gidx.MustNewID(ownerPrefix), LocationID: locationID},
			lbCount:  2,
		},
		{
			TestName: "fails to create loadbalancers - over limit",
			Input:    graphclient.CreateLoadBalancerInput{Name: name, ProviderID: prov.ID, OwnerID: gidx.MustNewID(ownerPrefix), LocationID: locationID},
			lbCount:  5,
			errorMsg: "load balancer limit reached",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.TestName, func(t *testing.T) {
			tt := tt

			t.Parallel()

			var err error

			for i := 1; i < tt.lbCount; i++ {
				_, err = graphTestClient().LoadBalancerCreate(ctx, tt.Input)
				if err != nil {
					return
				}
			}

			if tt.errorMsg != "" {
				require.Error(t, err)
				assert.ErrorContains(t, err, tt.errorMsg)

				return
			}

			require.NoError(t, err)
		})
	}
}

func TestUpdate_loadBalancer(t *testing.T) {
	ctx := context.Background()
	perms := new(mockpermissions.MockPermissions)
	perms.On("CreateAuthRelationships", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	ctx = perms.ContextWithHandler(ctx)

	// Permit request
	ctx = context.WithValue(ctx, permissions.CheckerCtxKey, permissions.DefaultAllowChecker)

	lb := (&testutils.LoadBalancerBuilder{}).MustNew(ctx)
	updateName := gofakeit.DomainName()

	testCases := []struct {
		TestName   string
		ID         gidx.PrefixedID
		Input      graphclient.UpdateLoadBalancerInput
		ExpectedLB *ent.LoadBalancer
		errorMsg   string
	}{
		{
			TestName: "updates loadbalancer",
			ID:       lb.ID,
			Input:    graphclient.UpdateLoadBalancerInput{Name: &updateName},
			ExpectedLB: &ent.LoadBalancer{
				Name:       updateName,
				ProviderID: lb.ProviderID,
				OwnerID:    lb.OwnerID,
				LocationID: lb.LocationID,
			},
		},
		{
			TestName: "fails to update name to empty",
			ID:       lb.ID,
			Input:    graphclient.UpdateLoadBalancerInput{Name: newString("")},
			errorMsg: "must not be empty",
		},
		{
			TestName: "fails to update name to whitespace",
			ID:       lb.ID,
			Input:    graphclient.UpdateLoadBalancerInput{Name: newString("   ")},
			errorMsg: "must not be empty",
		},
		{
			TestName: "fails to update loadbalancer that does not exist",
			ID:       gidx.PrefixedID("loadbal-dne"),
			Input:    graphclient.UpdateLoadBalancerInput{Name: newString("loadbal-dne")},
			errorMsg: "load_balancer not found",
		},
		{
			TestName: "fails with invalid gidx",
			ID:       "test-invalid-id",
			Input:    graphclient.UpdateLoadBalancerInput{Name: newString("loadbal-dne")},
			errorMsg: "invalid id",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.TestName, func(t *testing.T) {
			resp, err := graphTestClient().LoadBalancerUpdate(ctx, tt.ID, tt.Input)

			if tt.errorMsg != "" {
				require.Error(t, err)
				assert.ErrorContains(t, err, tt.errorMsg)
				assert.Nil(t, resp)

				return
			}

			require.NoError(t, err)
			require.NotNil(t, resp)
			require.NotNil(t, resp.LoadBalancerUpdate)

			updatedLB := resp.LoadBalancerUpdate.LoadBalancer
			assert.Equal(t, tt.ExpectedLB.Name, updatedLB.Name)
			assert.Equal(t, lb.ID, updatedLB.ID)
		})
	}
}

func TestDelete_loadBalancer(t *testing.T) {
	ctx := context.Background()
	perms := new(mockpermissions.MockPermissions)
	perms.On("CreateAuthRelationships", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	perms.On("DeleteAuthRelationships", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	ctx = perms.ContextWithHandler(ctx)

	// Permit request
	ctx = context.WithValue(ctx, permissions.CheckerCtxKey, permissions.DefaultAllowChecker)

	lb := (&testutils.LoadBalancerBuilder{}).MustNew(ctx)

	testCases := []struct {
		TestName   string
		Input      gidx.PrefixedID
		ExpectedID gidx.PrefixedID
		errorMsg   string
	}{
		{
			TestName:   "deletes loadbalancer",
			Input:      lb.ID,
			ExpectedID: lb.ID,
		},
		{
			TestName: "fails to delete loadbalancer that does not exist",
			Input:    gidx.PrefixedID("loadbal-dne"),
			errorMsg: "load_balancer not found",
		},
		{
			TestName: "fails to delete empty loadbalancer ID",
			Input:    gidx.PrefixedID(""),
			errorMsg: "must not be empty",
		},
		{
			TestName: "fails with invalid gidx",
			Input:    "test-invalid-id",
			errorMsg: "invalid id",
		},
		{
			TestName: "fails with invalid characters",
			Input:    gidx.PrefixedID("loadbal-!@#$%^&*()"),
			errorMsg: "valid characters are A-Z a-z 0-9 _ -",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.TestName, func(t *testing.T) {
			resp, err := graphTestClient().LoadBalancerDelete(ctx, tt.Input)

			if tt.errorMsg != "" {
				require.Error(t, err)
				assert.ErrorContains(t, err, tt.errorMsg)
				assert.Nil(t, resp)

				return
			}

			require.NoError(t, err)
			require.NotNil(t, resp)
			require.NotNil(t, resp.LoadBalancerDelete)

			deletedLB := resp.LoadBalancerDelete
			assert.EqualValues(t, tt.ExpectedID, deletedLB.DeletedID)
		})
	}
}

func TestFullLoadBalancerLifecycle(t *testing.T) {
	ctx := context.Background()
	perms := new(mockpermissions.MockPermissions)
	perms.On("CreateAuthRelationships", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	perms.On("DeleteAuthRelationships", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	ctx = perms.ContextWithHandler(ctx)

	// Permit request
	ctx = context.WithValue(ctx, permissions.CheckerCtxKey, permissions.DefaultAllowChecker)

	prov := (&testutils.ProviderBuilder{}).MustNew(ctx)
	ownerID := gidx.MustNewID(ownerPrefix)
	locationID := gidx.MustNewID(locationPrefix)
	name := gofakeit.DomainName()

	// create the LB
	createdLBResp, err := graphTestClient().LoadBalancerCreate(ctx, graphclient.CreateLoadBalancerInput{
		Name:       name,
		ProviderID: prov.ID,
		OwnerID:    ownerID,
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
	assert.Equal(t, ownerID, createdLB.Owner.ID)

	createdPortResp, err := graphTestClient().LoadBalancerPortCreate(ctx, graphclient.CreateLoadBalancerPortInput{
		Name:           newString(gofakeit.DomainName()),
		Number:         8080,
		LoadBalancerID: createdLB.ID,
	})

	require.NoError(t, err)
	require.NotNil(t, createdPortResp)
	require.NotNil(t, createdPortResp.LoadBalancerPortCreate.LoadBalancerPort)

	// Update the LB
	newName := gofakeit.DomainName()
	updatedLBResp, err := graphTestClient().LoadBalancerUpdate(ctx, createdLB.ID, graphclient.UpdateLoadBalancerInput{Name: &newName})

	require.NoError(t, err)
	require.NotNil(t, updatedLBResp)
	require.NotNil(t, updatedLBResp.LoadBalancerUpdate.LoadBalancer)

	updatedLB := updatedLBResp.LoadBalancerUpdate.LoadBalancer
	require.EqualValues(t, createdLB.ID, updatedLB.ID)
	require.Equal(t, newName, updatedLB.Name)

	// Query the LB
	queryLB, err := graphTestClient().GetLoadBalancer(ctx, createdLB.ID)
	require.NoError(t, err)
	require.NotNil(t, queryLB)
	require.NotNil(t, queryLB.LoadBalancer)
	require.Equal(t, newName, queryLB.LoadBalancer.Name)

	// Delete the LB
	deletedResp, err := graphTestClient().LoadBalancerDelete(ctx, createdLB.ID)
	require.NoError(t, err)
	require.NotNil(t, deletedResp)
	require.NotNil(t, deletedResp.LoadBalancerDelete)
	require.EqualValues(t, createdLB.ID, deletedResp.LoadBalancerDelete.DeletedID.String())

	// Query the LB to ensure it's no longer available
	deletedLB, err := graphTestClient().GetLoadBalancer(ctx, createdLB.ID)
	require.Error(t, err)
	require.Nil(t, deletedLB)
	require.ErrorContains(t, err, "load_balancer not found")
}
