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

func TestQueryPool(t *testing.T) {
	ctx := context.Background()
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
			resp, err := newGraphTestClient().GetLoadBalancerPool(ctx, tt.QueryID)
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
	tenantID := gidx.MustNewID(tenantPrefix)
	ctx := context.Background()

	testCases := []struct {
		TestName     string
		Input        graphclient.CreateLoadBalancerPoolInput
		ExpectedPool graphclient.LoadBalancerPool
		errorMsg     string
	}{
		{
			TestName: "create pool",
			Input: graphclient.CreateLoadBalancerPoolInput{
				Name:     "pooly",
				Protocol: graphclient.LoadBalancerPoolProtocolTCP,
				TenantID: tenantID,
			},
			ExpectedPool: graphclient.LoadBalancerPool{
				Name:     "pooly",
				Protocol: graphclient.LoadBalancerPoolProtocolTCP,
				TenantID: tenantID,
			},
		},
		{
			TestName: "invalid tenant ID",
			Input: graphclient.CreateLoadBalancerPoolInput{
				Name:     "pooly",
				Protocol: graphclient.LoadBalancerPoolProtocolTCP,
				TenantID: "not a valid ID",
			},
			errorMsg: "invalid id",
		},
		{
			TestName: "invalid protocol",
			Input: graphclient.CreateLoadBalancerPoolInput{
				Name:     "pooly",
				Protocol: "invalid",
				TenantID: tenantID,
			},
			errorMsg: "not a valid LoadBalancerPoolProtocol",
		},
		{
			TestName: "empty name",
			Input: graphclient.CreateLoadBalancerPoolInput{
				Name:     "",
				Protocol: graphclient.LoadBalancerPoolProtocolUDP,
				TenantID: tenantID,
			},
			errorMsg: "validator failed",
		},
	}

	for _, tt := range testCases {
		// lint
		tt := tt

		t.Run(tt.TestName, func(t *testing.T) {
			createdPoolResp, err := newGraphTestClient().LoadBalancerPoolCreate(ctx, tt.Input)
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
			assert.Equal(t, tt.ExpectedPool.Protocol, createdPool.Protocol)
			assert.Equal(t, tt.ExpectedPool.TenantID, createdPool.TenantID)
		})
	}
}

func TestMutate_PoolUpdate(t *testing.T) {
	ctx := context.Background()
	pool1 := (&PoolBuilder{Protocol: "tcp"}).MustNew(ctx)
	updateProtocolUnknown := graphclient.LoadBalancerPoolProtocol("invalid")
	updateProtocolUDP := graphclient.LoadBalancerPoolProtocolUDP

	testCases := []struct {
		TestName     string
		Input        graphclient.UpdateLoadBalancerPoolInput
		ExpectedPool graphclient.LoadBalancerPool
		errorMsg     string
	}{
		{
			TestName: "successfully updates name",
			Input: graphclient.UpdateLoadBalancerPoolInput{
				Name: newString("ImaPool"),
			},
			ExpectedPool: graphclient.LoadBalancerPool{
				Name:     "ImaPool",
				Protocol: graphclient.LoadBalancerPoolProtocolTCP,
				TenantID: pool1.TenantID,
			},
		},
		{
			TestName: "successfully updates protocol",
			Input: graphclient.UpdateLoadBalancerPoolInput{
				Name:     newString("ImaPool"),
				Protocol: &updateProtocolUDP,
			},
			ExpectedPool: graphclient.LoadBalancerPool{
				Name:     "ImaPool",
				Protocol: graphclient.LoadBalancerPoolProtocolUDP,
				TenantID: pool1.TenantID,
			},
		},
		{
			TestName: "invalid protocol",
			Input: graphclient.UpdateLoadBalancerPoolInput{
				Name:     newString("ImaPool"),
				Protocol: &updateProtocolUnknown,
			},
			errorMsg: "not a valid LoadBalancerPoolProtocol",
		},
		{
			TestName: "empty name",
			Input: graphclient.UpdateLoadBalancerPoolInput{
				Name: newString(""),
			},
			errorMsg: "validator failed",
		},
	}

	for _, tt := range testCases {
		// lint
		tt := tt

		t.Run(tt.TestName, func(t *testing.T) {
			updatedPoolResp, err := newGraphTestClient().LoadBalancerPoolUpdate(ctx, pool1.ID, tt.Input)
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
			assert.Equal(t, tt.ExpectedPool.Protocol, updatedPool.Protocol)
			assert.Equal(t, tt.ExpectedPool.TenantID, updatedPool.TenantID)
		})
	}
}

func TestMutate_PoolDelete(t *testing.T) {
	ctx := context.Background()
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
			DeleteID: gidx.MustNewID(tenantPrefix),
			errorMsg: "not found",
		},
	}

	for _, tt := range testCases {
		// lint
		tt := tt

		t.Run(tt.TestName, func(t *testing.T) {
			poolDeleteResp, err := newGraphTestClient().LoadBalancerPoolDelete(ctx, tt.DeleteID)
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
