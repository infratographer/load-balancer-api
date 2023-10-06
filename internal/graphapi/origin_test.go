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

func TestQueryPoolOrigin(t *testing.T) {
	ctx := context.Background()
	perms := new(mockpermissions.MockPermissions)
	perms.On("CreateAuthRelationships", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	ctx = perms.ContextWithHandler(ctx)

	// Permit request
	ctx = context.WithValue(ctx, permissions.CheckerCtxKey, permissions.DefaultAllowChecker)

	pool1 := (&testutils.PoolBuilder{}).MustNew(ctx)
	origin1 := (&testutils.OriginBuilder{PoolID: pool1.ID}).MustNew(ctx)
	origin2 := (&testutils.OriginBuilder{PoolID: pool1.ID}).MustNew(ctx)

	testCases := []struct {
		TestName       string
		QueryID        gidx.PrefixedID
		ExpectedOrigin *ent.Origin
		errorMsg       string
	}{
		{
			TestName:       "get origin 1",
			QueryID:        origin1.ID,
			ExpectedOrigin: origin1,
		},
		{
			TestName:       "get origin 2",
			QueryID:        origin2.ID,
			ExpectedOrigin: origin2,
		},
		{
			TestName: "invalid origin query ID",
			QueryID:  "an invalid origin id",
			errorMsg: "invalid id",
		},
	}

	for _, tt := range testCases {
		// lint
		tt := tt

		t.Run(tt.TestName, func(t *testing.T) {
			poolOriginResp, err := graphTestClient().GetLoadBalancerPoolOrigin(ctx, pool1.ID, tt.QueryID)
			if tt.errorMsg != "" {
				require.Error(t, err)
				assert.ErrorContains(t, err, tt.errorMsg)
				assert.Nil(t, poolOriginResp)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, poolOriginResp)

			poolOrigins := poolOriginResp.LoadBalancerPool.Origins.Edges
			assert.Len(t, poolOrigins, 1)
			assert.Equal(t, tt.ExpectedOrigin.Name, poolOrigins[0].Node.Name)
			assert.Equal(t, tt.ExpectedOrigin.Target, poolOrigins[0].Node.Target)
			assert.Equal(t, tt.ExpectedOrigin.PortNumber, int(poolOrigins[0].Node.PortNumber))
			assert.Equal(t, tt.ExpectedOrigin.Active, poolOrigins[0].Node.Active)
			assert.Equal(t, tt.ExpectedOrigin.PoolID, poolOrigins[0].Node.PoolID)
		})
	}
}

func TestMutate_OriginCreate(t *testing.T) {
	ctx := context.Background()
	perms := new(mockpermissions.MockPermissions)
	perms.On("CreateAuthRelationships", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	ctx = perms.ContextWithHandler(ctx)

	// Permit request
	ctx = context.WithValue(ctx, permissions.CheckerCtxKey, permissions.DefaultAllowChecker)

	pool1 := (&testutils.PoolBuilder{}).MustNew(ctx)

	testCases := []struct {
		TestName       string
		Input          graphclient.CreateLoadBalancerOriginInput
		ExpectedOrigin ent.LoadBalancerOrigin
		errorMsg       string
	}{
		{
			TestName: "poolID not found",
			Input: graphclient.CreateLoadBalancerOriginInput{
				Name:       "original",
				Target:     "1.2.3.4",
				PortNumber: 22,
				PoolID:     "loadpol-does-not-exist",
			},
			errorMsg: "pool not found",
		},
		{
			TestName: "invalid poolID format",
			Input: graphclient.CreateLoadBalancerOriginInput{
				Name:       "original",
				Target:     "1.2.3.4",
				PortNumber: 22,
				PoolID:     "derp",
			},
			errorMsg: "invalid id",
		},
		{
			TestName: "creates pool origin - defaults active to true",
			Input: graphclient.CreateLoadBalancerOriginInput{
				Name:       "original",
				Target:     "1.2.3.4",
				PortNumber: 22,
				PoolID:     pool1.ID,
			},
			ExpectedOrigin: ent.LoadBalancerOrigin{
				Name:       "original",
				Target:     "1.2.3.4",
				PortNumber: 22,
				PoolID:     pool1.ID,
				Active:     true,
			},
		},
		{
			TestName: "creates pool origin - active is false",
			Input: graphclient.CreateLoadBalancerOriginInput{
				Name:       "original",
				Target:     "1.2.3.4",
				PortNumber: 22,
				PoolID:     pool1.ID,
				Active:     newBool(false),
			},
			ExpectedOrigin: ent.LoadBalancerOrigin{
				Name:       "original",
				Target:     "1.2.3.4",
				PortNumber: 22,
				PoolID:     pool1.ID,
				Active:     false,
			},
		},
		{
			TestName: "invalid target ip",
			Input: graphclient.CreateLoadBalancerOriginInput{
				Name:       "original",
				Target:     "not a valid target ip",
				PortNumber: 22,
				PoolID:     pool1.ID,
				Active:     newBool(false),
			},
			errorMsg: "invalid ip address",
		},
	}

	for _, tt := range testCases {
		// lint
		tt := tt

		t.Run(tt.TestName, func(t *testing.T) {
			createdOriginResp, err := graphTestClient().LoadBalancerOriginCreate(ctx, tt.Input)

			if tt.errorMsg != "" {
				require.Error(t, err)
				assert.ErrorContains(t, err, tt.errorMsg)
				assert.Nil(t, createdOriginResp)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, createdOriginResp)

			createdOrigin := createdOriginResp.LoadBalancerOriginCreate.LoadBalancerOrigin
			assert.Contains(t, createdOrigin.ID, "loadogn-")
			assert.Equal(t, tt.ExpectedOrigin.Name, createdOrigin.Name)
			assert.Equal(t, tt.ExpectedOrigin.Target, createdOrigin.Target)
			assert.Equal(t, tt.ExpectedOrigin.PortNumber, int(createdOrigin.PortNumber))
			assert.Equal(t, tt.ExpectedOrigin.PoolID, createdOrigin.PoolID)
			assert.Equal(t, tt.ExpectedOrigin.Active, createdOrigin.Active)
		})
	}

	assert.Nil(t, nil)
}

func TestMutate_OriginUpdate(t *testing.T) {
	ctx := context.Background()
	perms := new(mockpermissions.MockPermissions)
	perms.On("CreateAuthRelationships", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	ctx = perms.ContextWithHandler(ctx)

	// Permit request
	ctx = context.WithValue(ctx, permissions.CheckerCtxKey, permissions.DefaultAllowChecker)

	pool1 := (&testutils.PoolBuilder{}).MustNew(ctx)
	origin1 := (&testutils.OriginBuilder{PoolID: pool1.ID}).MustNew(ctx)

	testCases := []struct {
		TestName       string
		OriginID       gidx.PrefixedID
		Input          graphclient.UpdateLoadBalancerOriginInput
		ExpectedOrigin ent.LoadBalancerOrigin
		errorMsg       string
	}{
		{
			TestName: "originID not found",
			OriginID: gidx.MustNewID("testing"),
			Input: graphclient.UpdateLoadBalancerOriginInput{
				Name:       newString("originator"),
				Target:     newString("5.6.7.8"),
				PortNumber: newInt64(222),
				Active:     newBool(false),
			},
			errorMsg: "origin not found",
		},
		{
			TestName: "invalid originID format",
			OriginID: "derp",
			Input: graphclient.UpdateLoadBalancerOriginInput{
				Name:       newString("originator"),
				Target:     newString("5.6.7.8"),
				PortNumber: newInt64(222),
				Active:     newBool(false),
			},
			errorMsg: "invalid id",
		},
		{
			TestName: "updates origin",
			OriginID: origin1.ID,
			Input: graphclient.UpdateLoadBalancerOriginInput{
				Name:       newString("originator"),
				Target:     newString("5.6.7.8"),
				PortNumber: newInt64(222),
				Active:     newBool(true),
			},
			ExpectedOrigin: ent.LoadBalancerOrigin{
				ID:         origin1.ID,
				Name:       "originator",
				Target:     "5.6.7.8",
				PortNumber: 222,
				Active:     true,
				PoolID:     pool1.ID,
			},
		},
	}

	for _, tt := range testCases {
		// lint
		tt := tt

		t.Run(tt.TestName, func(t *testing.T) {
			updatedOriginResp, err := graphTestClient().LoadBalancerOriginUpdate(ctx, tt.OriginID, tt.Input)

			if tt.errorMsg != "" {
				require.Error(t, err)
				assert.ErrorContains(t, err, tt.errorMsg)
				assert.Nil(t, updatedOriginResp)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, updatedOriginResp)

			updatedOrigin := updatedOriginResp.LoadBalancerOriginUpdate.LoadBalancerOrigin
			assert.Equal(t, tt.ExpectedOrigin.ID, updatedOrigin.ID)
			assert.Equal(t, tt.ExpectedOrigin.Name, updatedOrigin.Name)
			assert.Equal(t, tt.ExpectedOrigin.Target, updatedOrigin.Target)
			assert.Equal(t, tt.ExpectedOrigin.PortNumber, int(updatedOrigin.PortNumber))
			assert.Equal(t, tt.ExpectedOrigin.PoolID, updatedOrigin.PoolID)
			assert.Equal(t, tt.ExpectedOrigin.Active, updatedOrigin.Active)
		})
	}
}

func TestMutate_OriginDelete(t *testing.T) {
	ctx := context.Background()
	perms := new(mockpermissions.MockPermissions)
	perms.On("CreateAuthRelationships", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	perms.On("DeleteAuthRelationships", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	ctx = perms.ContextWithHandler(ctx)

	// Permit request
	ctx = context.WithValue(ctx, permissions.CheckerCtxKey, permissions.DefaultAllowChecker)

	pool1 := (&testutils.PoolBuilder{}).MustNew(ctx)
	origin1 := (&testutils.OriginBuilder{PoolID: pool1.ID}).MustNew(ctx)

	testCases := []struct {
		TestName string
		OriginID gidx.PrefixedID
		errorMsg string
	}{
		{
			TestName: "originID not found",
			OriginID: gidx.MustNewID("testing"),
			errorMsg: "origin not found",
		},
		{
			TestName: "invalid originID format",
			OriginID: "derp",
			errorMsg: "invalid id",
		},
		{
			TestName: "deletes origin",
			OriginID: origin1.ID,
		},
	}

	for _, tt := range testCases {
		// lint
		tt := tt

		t.Run(tt.TestName, func(t *testing.T) {
			deletedOriginResp, err := graphTestClient().LoadBalancerOriginDelete(ctx, tt.OriginID)

			if tt.errorMsg != "" {
				require.Error(t, err)
				assert.ErrorContains(t, err, tt.errorMsg)
				assert.Nil(t, deletedOriginResp)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, deletedOriginResp)

			deletedOriginID := deletedOriginResp.LoadBalancerOriginDelete.DeletedID
			assert.Equal(t, tt.OriginID, deletedOriginID)
		})
	}
}
