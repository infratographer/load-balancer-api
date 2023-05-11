package graphapi_test

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.infratographer.com/load-balancer-api/internal/graphclient"
)

func TestFullLoadBalancerPortLifecycle(t *testing.T) {
	ctx := context.Background()
	lb := (&LoadBalancerBuilder{}).MustNew(ctx)
	name := gofakeit.DomainName()

	createdPortResp, err := newGraphTestClient().LoadBalancerPortCreate(ctx, graphclient.CreateLoadBalancerPortInput{
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
	updatedPort, err := newGraphTestClient().LoadBalancerPortUpdate(ctx, createdPort.ID, graphclient.UpdateLoadBalancerPortInput{Number: &newPort})

	require.NoError(t, err)
	require.NotNil(t, updatedPort)
	require.EqualValues(t, createdPort.ID, updatedPort.LoadBalancerPortUpdate.LoadBalancerPort.ID)
	require.Equal(t, newPort, updatedPort.LoadBalancerPortUpdate.LoadBalancerPort.Number)

	// Query the Port
	queryPort, err := newGraphTestClient().GetLoadBalancerPort(ctx, lb.ID, createdPort.ID)
	require.NoError(t, err)
	require.NotNil(t, queryPort)
	require.Len(t, queryPort.LoadBalancer.Ports.Edges, 1)
	require.Equal(t, newPort, queryPort.LoadBalancer.Ports.Edges[0].Node.Number)

	// Delete the Port
	deletedResp, err := newGraphTestClient().LoadBalancerPortDelete(ctx, createdPort.ID)
	require.NoError(t, err)
	require.NotNil(t, deletedResp)
	require.EqualValues(t, createdPort.ID, deletedResp.LoadBalancerPortDelete.DeletedID.String())

	// Query the Port
	queryPort, err = newGraphTestClient().GetLoadBalancerPort(ctx, lb.ID, createdPort.ID)
	// The Load balancer still exists so this doesn't cause a failure
	require.NoError(t, err)
	require.Len(t, queryPort.LoadBalancer.Ports.Edges, 0)
}
