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

func TestCreate_LoadbalancerPort(t *testing.T) {
	config.AppConfig.RestrictedPorts = []int{1234}

	ctx := context.Background()
	perms := new(mockpermissions.MockPermissions)
	perms.On("CreateAuthRelationships", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	ctx = perms.ContextWithHandler(ctx)

	// Permit request
	ctx = context.WithValue(ctx, permissions.CheckerCtxKey, permissions.DefaultAllowChecker)

	lb := (&testutils.LoadBalancerBuilder{}).MustNew(ctx)
	poolBad := (&testutils.PoolBuilder{}).MustNew(ctx)
	_ = (&testutils.PortBuilder{Name: "port80", LoadBalancerID: lb.ID, Number: 80}).MustNew(ctx)

	testCases := []struct {
		TestName string
		Input    graphclient.CreateLoadBalancerPortInput
		Expected *graphclient.LoadBalancerPort
		errorMsg string
	}{
		{
			TestName: "creates loadbalancer port",
			Input: graphclient.CreateLoadBalancerPortInput{
				Name:           newString("lb-port"),
				LoadBalancerID: lb.ID,
				Number:         22,
			},
			Expected: &graphclient.LoadBalancerPort{
				Name:   newString("lb-port"),
				Number: 22,
			},
		},
		{
			TestName: "succeeds in creating loadbalancer port with empty name",
			Input: graphclient.CreateLoadBalancerPortInput{
				Name:           newString(""),
				LoadBalancerID: lb.ID,
				Number:         23,
			},
			Expected: &graphclient.LoadBalancerPort{
				Name:   newString(""),
				Number: 23,
			},
		},
		{
			TestName: "succeeds in creating loadbalancer port with nil name",
			Input: graphclient.CreateLoadBalancerPortInput{
				Name:           nil,
				LoadBalancerID: lb.ID,
				Number:         24,
			},
			Expected: &graphclient.LoadBalancerPort{
				Name:   newString(""),
				Number: 24,
			},
		},
		{
			TestName: "fails to create loadbalancer port with empty loadbalancer id",
			Input: graphclient.CreateLoadBalancerPortInput{
				Name:           newString("lb-port"),
				LoadBalancerID: "",
				Number:         22,
			},
			errorMsg: "load_balancer not found",
		},
		{
			TestName: "fails to create loadbalancer port with number < min",
			Input: graphclient.CreateLoadBalancerPortInput{
				Name:           newString("lb-port"),
				LoadBalancerID: lb.ID,
				Number:         0,
			},
			errorMsg: "value out of range",
		},
		{
			TestName: "fails to create loadbalancer port with number > max",
			Input: graphclient.CreateLoadBalancerPortInput{
				Name:           newString("lb-port"),
				LoadBalancerID: lb.ID,
				Number:         65536,
			},
			errorMsg: "value out of range",
		},
		{
			TestName: "fails to create loadbalancer port with duplicate port number",
			Input: graphclient.CreateLoadBalancerPortInput{
				Name:           newString("lb-port"),
				LoadBalancerID: lb.ID,
				Number:         80,
			},
			errorMsg: "port number already in use",
		},
		{
			TestName: "fails to create loadbalancer port with restricted port number",
			Input: graphclient.CreateLoadBalancerPortInput{
				Name:           newString("lb-port"),
				LoadBalancerID: lb.ID,
				Number:         1234,
			},
			errorMsg: "port number restricted",
		},
		{
			TestName: "fails to create loadbalancer port with invalid pool id",
			Input: graphclient.CreateLoadBalancerPortInput{
				Name:           newString("lb-port"),
				LoadBalancerID: lb.ID,
				Number:         1234,
				PoolIDs:        []gidx.PrefixedID{"not-a-valid-pool-id"},
			},
			errorMsg: "invalid id",
		},
		{
			TestName: "fails to create loadbalancer port with pool with conflicting OwnerID",
			Input: graphclient.CreateLoadBalancerPortInput{
				Name:           newString("lb-port"),
				LoadBalancerID: lb.ID,
				Number:         1234,
				PoolIDs:        []gidx.PrefixedID{poolBad.ID},
			},
			errorMsg: "one or more pools not found",
		},
		{
			TestName: "fails to create port with long name",
			Input: graphclient.CreateLoadBalancerPortInput{
				Name:           &longName,
				LoadBalancerID: lb.ID,
				Number:         22,
			},
			errorMsg: "must not be longer than",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.TestName, func(t *testing.T) {
			tt := tt

			t.Parallel()

			resp, err := graphTestClient().LoadBalancerPortCreate(ctx, tt.Input)

			if tt.errorMsg != "" {
				require.Error(t, err)
				assert.ErrorContains(t, err, tt.errorMsg)
				assert.Nil(t, resp)

				return
			}

			require.NoError(t, err)
			require.NotNil(t, resp)
			require.NotNil(t, resp.LoadBalancerPortCreate)

			createdPort := resp.LoadBalancerPortCreate.LoadBalancerPort
			require.NotNil(t, createdPort.ID)
			require.Equal(t, tt.Expected.Name, createdPort.Name)
			require.Equal(t, tt.Expected.Number, createdPort.Number)
			require.Equal(t, "loadprt", createdPort.ID.Prefix())
			require.Equal(t, lb.ID, createdPort.LoadBalancer.ID)
		})
	}
}

func TestUpdate_LoadbalancerPort(t *testing.T) {
	config.AppConfig.RestrictedPorts = []int{1234}

	ctx := context.Background()
	perms := new(mockpermissions.MockPermissions)
	perms.On("CreateAuthRelationships", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	ctx = perms.ContextWithHandler(ctx)

	// Permit request
	ctx = context.WithValue(ctx, permissions.CheckerCtxKey, permissions.DefaultAllowChecker)

	lb := (&testutils.LoadBalancerBuilder{}).MustNew(ctx)
	port := (&testutils.PortBuilder{Name: "port80", LoadBalancerID: lb.ID, Number: 80}).MustNew(ctx)
	poolBad := (&testutils.PoolBuilder{}).MustNew(ctx)
	_ = (&testutils.PortBuilder{Name: "dupeport8080", LoadBalancerID: lb.ID, Number: 8080}).MustNew(ctx)

	testCases := []struct {
		TestName string
		Input    graphclient.UpdateLoadBalancerPortInput
		ID       gidx.PrefixedID
		Expected *graphclient.LoadBalancerPort
		errorMsg string
	}{
		{
			TestName: "fails to update loadbalancer port number to duplicate of another port",
			ID:       port.ID,
			Input: graphclient.UpdateLoadBalancerPortInput{
				Number: newInt64(8080),
			},
			errorMsg: "port number already in use",
		},
		{
			TestName: "fails to update loadbalancer port number to restricted port",
			ID:       port.ID,
			Input: graphclient.UpdateLoadBalancerPortInput{
				Number: newInt64(1234),
			},
			errorMsg: "port number restricted",
		},
		{
			TestName: "updates loadbalancer port name",
			ID:       port.ID,
			Input: graphclient.UpdateLoadBalancerPortInput{
				Name: newString("lb-port"),
			},
			Expected: &graphclient.LoadBalancerPort{
				Name:   newString("lb-port"),
				Number: 80,
			},
		},
		{
			TestName: "updates loadbalancer port number",
			ID:       port.ID,
			Input: graphclient.UpdateLoadBalancerPortInput{
				Number: newInt64(22),
			},
			Expected: &graphclient.LoadBalancerPort{
				Name:   newString("lb-port"),
				Number: 22,
			},
		},
		{
			TestName: "succeeds in updating loadbalancer port name to empty",
			ID:       port.ID,
			Input: graphclient.UpdateLoadBalancerPortInput{
				Name: newString(""),
			},
			Expected: &graphclient.LoadBalancerPort{
				Name:   newString(""),
				Number: 22,
			},
		},
		{
			TestName: "fails to update loadbalancer port number < min",
			ID:       port.ID,
			Input: graphclient.UpdateLoadBalancerPortInput{
				Number: newInt64(0),
			},
			errorMsg: "value out of range",
		},
		{
			TestName: "fails to update loadbalancer port number > max",
			ID:       port.ID,
			Input: graphclient.UpdateLoadBalancerPortInput{
				Number: newInt64(65536),
			},
			errorMsg: "value out of range",
		},
		{
			TestName: "fails to update port that does not exist",
			ID:       gidx.PrefixedID("loadprt-doesnotexist"),
			Input:    graphclient.UpdateLoadBalancerPortInput{},
			errorMsg: "not found",
		},
		{
			TestName: "fails to update port with invalid gidx",
			ID:       gidx.PrefixedID("not a valid gidx"),
			Input:    graphclient.UpdateLoadBalancerPortInput{},
			errorMsg: "invalid id",
		},
		{
			TestName: "fails to update loadbalancer port with pool with conflicting OwnerID",
			ID:       port.ID,
			Input: graphclient.UpdateLoadBalancerPortInput{
				AddPoolIDs: []gidx.PrefixedID{poolBad.ID},
			},
			errorMsg: "one or more pools not found",
		},
		{
			TestName: "fails to update port with long name",
			ID:       port.ID,
			Input:    graphclient.UpdateLoadBalancerPortInput{Name: &longName},
			errorMsg: "must not be longer than",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.TestName, func(t *testing.T) {
			resp, err := graphTestClient().LoadBalancerPortUpdate(ctx, tt.ID, tt.Input)

			if tt.errorMsg != "" {
				require.Error(t, err)
				assert.ErrorContains(t, err, tt.errorMsg)
				assert.Nil(t, resp)

				return
			}

			require.NoError(t, err)
			require.NotNil(t, resp)
			require.NotNil(t, resp.LoadBalancerPortUpdate)

			updatedPort := resp.LoadBalancerPortUpdate.LoadBalancerPort
			require.NotNil(t, updatedPort.ID)
			require.Equal(t, tt.Expected.Name, updatedPort.Name)
			require.Equal(t, tt.Expected.Number, updatedPort.Number)
			require.Equal(t, "loadprt", updatedPort.ID.Prefix())
		})
	}
}

func TestDelete_LoadbalancerPort(t *testing.T) {
	ctx := context.Background()
	perms := new(mockpermissions.MockPermissions)
	perms.On("CreateAuthRelationships", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	perms.On("DeleteAuthRelationships", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	ctx = perms.ContextWithHandler(ctx)

	// Permit request
	ctx = context.WithValue(ctx, permissions.CheckerCtxKey, permissions.DefaultAllowChecker)

	lb := (&testutils.LoadBalancerBuilder{}).MustNew(ctx)
	port := (&testutils.PortBuilder{Name: "port80", LoadBalancerID: lb.ID, Number: 80}).MustNew(ctx)

	testCases := []struct {
		TestName string
		Input    gidx.PrefixedID
		errorMsg string
	}{
		{
			TestName: "deletes loadbalancer port",
			Input:    port.ID,
		},
		{
			TestName: "fails to delete loadbalancer port that does not exist",
			Input:    gidx.PrefixedID("loadprt-dne"),
			errorMsg: "port not found",
		},
		{
			TestName: "fails to delete empty loadbalancer port ID",
			Input:    gidx.PrefixedID(""),
			errorMsg: "must not be empty",
		},
		{
			TestName: "fails to delete with invalid gidx port ID",
			Input:    gidx.PrefixedID("not-a-valid-gidx"),
			errorMsg: "invalid id",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.TestName, func(t *testing.T) {
			tt := tt

			t.Parallel()

			resp, err := graphTestClient().LoadBalancerPortDelete(ctx, tt.Input)

			if tt.errorMsg != "" {
				require.Error(t, err)
				assert.ErrorContains(t, err, tt.errorMsg)
				assert.Nil(t, resp)

				return
			}

			require.NoError(t, err)
			require.NotNil(t, resp)
			require.NotNil(t, resp.LoadBalancerPortDelete)

			deletedPortID := resp.LoadBalancerPortDelete.DeletedID
			require.NotNil(t, deletedPortID)
			require.Equal(t, tt.Input, deletedPortID)
		})
	}
}

func TestGet_LoadbalancerPort(t *testing.T) {
	ctx := context.Background()
	perms := new(mockpermissions.MockPermissions)
	perms.On("CreateAuthRelationships", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	ctx = perms.ContextWithHandler(ctx)

	// Permit request
	ctx = context.WithValue(ctx, permissions.CheckerCtxKey, permissions.DefaultAllowChecker)

	lb := (&testutils.LoadBalancerBuilder{}).MustNew(ctx)
	port := (&testutils.PortBuilder{Name: "", LoadBalancerID: lb.ID, Number: 80}).MustNew(ctx)

	testCases := []struct {
		TestName     string
		QueryID      gidx.PrefixedID
		ExpectedPort *ent.Port
		errorMsg     string
	}{
		{
			TestName:     "get port",
			QueryID:      port.ID,
			ExpectedPort: port,
		},
		{
			TestName: "port not found",
			QueryID:  gidx.MustNewID("loadprt"),
			errorMsg: "not found",
		},
		{
			TestName: "invalid port ID",
			QueryID:  "an invalid port id",
			errorMsg: "invalid id",
		},
		{
			TestName: "fails with invalid gidx",
			QueryID:  "test-invalid-id",
			errorMsg: "invalid id",
		},
		{
			TestName: "fails with invalid characters",
			QueryID:  gidx.PrefixedID("loadprt-!@#$%^&*()"),
			errorMsg: "valid characters are A-Z a-z 0-9 _ -",
		},
	}

	for _, tt := range testCases {
		// lint
		tt := tt

		t.Run(tt.TestName, func(t *testing.T) {
			resp, err := graphTestClient().GetLoadBalancerPort(ctx, tt.QueryID)
			if tt.errorMsg != "" {
				require.Error(t, err)
				assert.ErrorContains(t, err, tt.errorMsg)
				assert.Nil(t, resp)

				return
			}

			require.NoError(t, err)
			require.NotNil(t, resp)

			assert.Equal(t, port.ID, resp.LoadBalancerPort.ID)
			assert.EqualValues(t, port.Number, resp.LoadBalancerPort.Number)
			assert.Equal(t, "loadprt", resp.LoadBalancerPort.ID.Prefix())
		})
	}
}

func TestFullLoadBalancerPortLifecycle(t *testing.T) {
	ctx := context.Background()
	perms := new(mockpermissions.MockPermissions)
	perms.On("CreateAuthRelationships", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	perms.On("DeleteAuthRelationships", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	ctx = perms.ContextWithHandler(ctx)

	// Permit request
	ctx = context.WithValue(ctx, permissions.CheckerCtxKey, permissions.DefaultAllowChecker)

	lb := (&testutils.LoadBalancerBuilder{}).MustNew(ctx)
	name := newString(gofakeit.DomainName())

	createdPortResp, err := graphTestClient().LoadBalancerPortCreate(ctx, graphclient.CreateLoadBalancerPortInput{
		Name:           name,
		LoadBalancerID: lb.ID,
		Number:         22,
	})

	require.NoError(t, err)
	require.NotNil(t, createdPortResp)
	require.NotNil(t, createdPortResp.LoadBalancerPortCreate.LoadBalancerPort)

	createdPort := createdPortResp.LoadBalancerPortCreate.LoadBalancerPort
	require.NotNil(t, createdPort.ID)
	require.Equal(t, name, createdPort.Name)
	require.EqualValues(t, 22, createdPort.Number)
	assert.Equal(t, "loadprt", createdPort.ID.Prefix())
	assert.Equal(t, lb.ID, createdPort.LoadBalancer.ID)

	// Update the Port
	newPort := int64(gofakeit.Number(1, 65535))
	updatedPort, err := graphTestClient().LoadBalancerPortUpdate(ctx, createdPort.ID, graphclient.UpdateLoadBalancerPortInput{Number: &newPort})

	require.NoError(t, err)
	require.NotNil(t, updatedPort)
	require.EqualValues(t, createdPort.ID, updatedPort.LoadBalancerPortUpdate.LoadBalancerPort.ID)
	require.Equal(t, newPort, updatedPort.LoadBalancerPortUpdate.LoadBalancerPort.Number)

	// Query the Port
	queryPort, err := graphTestClient().GetPortByLoadBalancer(ctx, lb.ID, createdPort.ID)
	require.NoError(t, err)
	require.NotNil(t, queryPort)
	require.Len(t, queryPort.LoadBalancer.Ports.Edges, 1)
	require.Equal(t, newPort, queryPort.LoadBalancer.Ports.Edges[0].Node.Number)

	// Delete the Port
	deletedResp, err := graphTestClient().LoadBalancerPortDelete(ctx, createdPort.ID)
	require.NoError(t, err)
	require.NotNil(t, deletedResp)
	require.EqualValues(t, createdPort.ID, deletedResp.LoadBalancerPortDelete.DeletedID.String())

	// Query the Port
	queryPort, err = graphTestClient().GetPortByLoadBalancer(ctx, lb.ID, createdPort.ID)
	// The Load balancer still exists so this doesn't cause a failure
	require.NoError(t, err)
	require.Len(t, queryPort.LoadBalancer.Ports.Edges, 0)
}
