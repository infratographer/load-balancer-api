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
	pool "go.infratographer.com/load-balancer-api/internal/ent/generated/pool"
	"go.infratographer.com/load-balancer-api/internal/graphclient"
)

func TestQueryPool(t *testing.T) {
	ctx := context.Background()
	perms := new(mockpermissions.MockPermissions)
	perms.On("CreateAuthRelationships", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	ctx = perms.ContextWithHandler(ctx)

	// Permit request
	ctx = context.WithValue(ctx, permissions.CheckerCtxKey, permissions.DefaultAllowChecker)

	pool1 := (&PoolBuilder{}).MustNew(ctx)
	pool2 := (&PoolBuilder{}).MustNew(ctx)

	testCases := []struct {
		TestName     string
		QueryID      gidx.PrefixedID
		ExpectedPool *ent.Pool
		errorMsg     string
	}{
		{
			TestName:     "get pool 1",
			QueryID:      pool1.ID,
			ExpectedPool: pool1,
		},
		{
			TestName:     "get pool 2",
			QueryID:      pool2.ID,
			ExpectedPool: pool2,
		},
		{
			TestName: "pool not found",
			QueryID:  gidx.MustNewID("testing"),
			errorMsg: "not found",
		},
		{
			TestName: "invalid pool query ID",
			QueryID:  "an invalid pool id",
			errorMsg: "invalid id",
		},
	}

	for _, tt := range testCases {
		// lint
		tt := tt

		t.Run(tt.TestName, func(t *testing.T) {
			resp, err := graphTestClient().GetLoadBalancerPool(ctx, tt.QueryID)
			if tt.errorMsg != "" {
				require.Error(t, err)
				assert.ErrorContains(t, err, tt.errorMsg)
				assert.Nil(t, resp)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, resp)
			assert.EqualValues(t, tt.ExpectedPool.Name, resp.LoadBalancerPool.Name)
		})
	}
}

func TestMutate_PoolCreate(t *testing.T) {
	ctx := context.Background()
	perms := new(mockpermissions.MockPermissions)
	perms.On("CreateAuthRelationships", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	ctx = perms.ContextWithHandler(ctx)

	// Permit request
	ctx = context.WithValue(ctx, permissions.CheckerCtxKey, permissions.DefaultAllowChecker)

	ownerID := gidx.MustNewID(ownerPrefix)

	lb := (&LoadBalancerBuilder{}).MustNew(ctx)
	port := (&PortBuilder{Name: "port80", LoadBalancerID: lb.ID, Number: 80}).MustNew(ctx)

	testCases := []struct {
		TestName     string
		Input        graphclient.CreateLoadBalancerPoolInput
		ExpectedPool ent.LoadBalancerPool
		errorMsg     string
	}{
		{
			TestName: "create pool",
			Input: graphclient.CreateLoadBalancerPoolInput{
				Name:     "pooly",
				Protocol: graphclient.LoadBalancerPoolProtocolTCP,
				OwnerID:  ownerID,
			},
			ExpectedPool: ent.LoadBalancerPool{
				Name:     "pooly",
				Protocol: pool.ProtocolTCP,
				OwnerID:  ownerID,
			},
		},
		{
			TestName: "invalid owner ID",
			Input: graphclient.CreateLoadBalancerPoolInput{
				Name:     "pooly",
				Protocol: graphclient.LoadBalancerPoolProtocolTCP,
				OwnerID:  "not a valid ID",
			},
			errorMsg: "invalid id",
		},
		{
			TestName: "invalid protocol",
			Input: graphclient.CreateLoadBalancerPoolInput{
				Name:     "pooly",
				Protocol: "invalid",
				OwnerID:  ownerID,
			},
			errorMsg: "not a valid LoadBalancerPoolProtocol",
		},
		{
			TestName: "empty name",
			Input: graphclient.CreateLoadBalancerPoolInput{
				Name:     "",
				Protocol: graphclient.LoadBalancerPoolProtocolUDP,
				OwnerID:  ownerID,
			},
			errorMsg: "validator failed",
		},
		{
			TestName: "port with conflicting OwnerID",
			Input: graphclient.CreateLoadBalancerPoolInput{
				Name:     "",
				Protocol: graphclient.LoadBalancerPoolProtocolUDP,
				OwnerID:  ownerID,
				PortIDs:  []gidx.PrefixedID{port.ID},
			},
			errorMsg: "conflicting ownership",
		},
	}

	for _, tt := range testCases {
		// lint
		tt := tt

		t.Run(tt.TestName, func(t *testing.T) {
			createdPoolResp, err := graphTestClient().LoadBalancerPoolCreate(ctx, tt.Input)
			if tt.errorMsg != "" {
				require.Error(t, err)
				assert.ErrorContains(t, err, tt.errorMsg)
				assert.Nil(t, createdPoolResp)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, createdPoolResp)

			createdPool := createdPoolResp.LoadBalancerPoolCreate.LoadBalancerPool

			assert.Equal(t, "loadpol", createdPool.ID.Prefix())
			assert.Equal(t, tt.ExpectedPool.Name, createdPool.Name)
			assert.Equal(t, tt.ExpectedPool.Protocol.String(), createdPool.Protocol.String())
			assert.Equal(t, tt.ExpectedPool.OwnerID, createdPool.OwnerID)
		})
	}
}

func TestMutate_PoolUpdate(t *testing.T) {
	ctx := context.Background()
	perms := new(mockpermissions.MockPermissions)
	perms.On("CreateAuthRelationships", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	ctx = perms.ContextWithHandler(ctx)

	// Permit request
	ctx = context.WithValue(ctx, permissions.CheckerCtxKey, permissions.DefaultAllowChecker)

	pool1 := (&PoolBuilder{Protocol: "tcp"}).MustNew(ctx)
	pool2 := (&PoolBuilder{Protocol: "tcp"}).MustNew(ctx)
	lb := (&LoadBalancerBuilder{}).MustNew(ctx)
	port := (&PortBuilder{LoadBalancerID: lb.ID}).MustNew(ctx)
	origin := (&OriginBuilder{PoolID: pool2.ID}).MustNew(ctx)
	updateProtocolUDP := graphclient.LoadBalancerPoolProtocolUDP

	testCases := []struct {
		TestName     string
		Input        graphclient.UpdateLoadBalancerPoolInput
		ExpectedPool ent.LoadBalancerPool
		errorMsg     string
	}{
		{
			TestName: "successfully updates name",
			Input: graphclient.UpdateLoadBalancerPoolInput{
				Name: newString("ImaPool"),
			},
			ExpectedPool: ent.LoadBalancerPool{
				Name:     "ImaPool",
				Protocol: pool.ProtocolTCP,
				OwnerID:  pool1.OwnerID,
			},
		},
		{
			TestName: "successfully updates protocol",
			Input: graphclient.UpdateLoadBalancerPoolInput{
				Name:     newString("ImaPool"),
				Protocol: &updateProtocolUDP,
			},
			ExpectedPool: ent.LoadBalancerPool{
				Name:     "ImaPool",
				Protocol: pool.ProtocolUDP,
				OwnerID:  pool1.OwnerID,
			},
		},
		{
			TestName: "empty name",
			Input: graphclient.UpdateLoadBalancerPoolInput{
				Name: newString(""),
			},
			errorMsg: "validator failed",
		},
		{
			TestName: "fails to update pool with port with conflicting OwnerID",
			Input: graphclient.UpdateLoadBalancerPoolInput{
				Name:       newString("ImaPool"),
				AddPortIDs: []gidx.PrefixedID{port.ID},
			},
			errorMsg: "conflicting ownership",
		},
		{
			TestName: "fails to update pool with another pool's origin",
			Input: graphclient.UpdateLoadBalancerPoolInput{
				Name:         newString("ImaPool"),
				AddOriginIDs: []gidx.PrefixedID{origin.ID},
			},
			errorMsg: "is already connected to a different pool_id",
		},
	}

	for _, tt := range testCases {
		// lint
		tt := tt

		t.Run(tt.TestName, func(t *testing.T) {
			updatedPoolResp, err := graphTestClient().LoadBalancerPoolUpdate(ctx, pool1.ID, tt.Input)
			if tt.errorMsg != "" {
				require.Error(t, err)
				assert.ErrorContains(t, err, tt.errorMsg)
				assert.Nil(t, updatedPoolResp)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, updatedPoolResp)

			updatedPool := updatedPoolResp.LoadBalancerPoolUpdate.LoadBalancerPool
			assert.Equal(t, "loadpol", updatedPool.ID.Prefix())
			assert.Equal(t, tt.ExpectedPool.Name, updatedPool.Name)
			assert.Equal(t, tt.ExpectedPool.Protocol.String(), updatedPool.Protocol.String())
			assert.Equal(t, tt.ExpectedPool.OwnerID, updatedPool.OwnerID)
		})
	}
}

func TestMutate_PoolDelete(t *testing.T) {
	ctx := context.Background()
	perms := new(mockpermissions.MockPermissions)
	perms.On("CreateAuthRelationships", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	perms.On("DeleteAuthRelationships", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	ctx = perms.ContextWithHandler(ctx)

	// Permit request
	ctx = context.WithValue(ctx, permissions.CheckerCtxKey, permissions.DefaultAllowChecker)

	pool1 := (&PoolBuilder{Protocol: "tcp"}).MustNew(ctx)

	testCases := []struct {
		TestName string
		DeleteID gidx.PrefixedID
		errorMsg string
	}{
		{
			TestName: "successfully deletes pool",
			DeleteID: pool1.ID,
		},
		{
			TestName: "invalid ID",
			DeleteID: "not a valid ID",
			errorMsg: "invalid id",
		},
		{
			TestName: "non-existent ID",
			DeleteID: gidx.MustNewID(ownerPrefix),
			errorMsg: "not found",
		},
	}

	for _, tt := range testCases {
		// lint
		tt := tt

		t.Run(tt.TestName, func(t *testing.T) {
			poolDeleteResp, err := graphTestClient().LoadBalancerPoolDelete(ctx, tt.DeleteID)
			if tt.errorMsg != "" {
				require.Error(t, err)
				assert.Nil(t, poolDeleteResp)
				assert.ErrorContains(t, err, tt.errorMsg)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, "loadpol", poolDeleteResp.LoadBalancerPoolDelete.DeletedID.Prefix())
		})
	}
}
