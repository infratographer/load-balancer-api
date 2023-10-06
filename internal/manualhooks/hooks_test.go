package manualhooks_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.infratographer.com/x/events"
	"go.infratographer.com/x/gidx"

	"go.infratographer.com/load-balancer-api/internal/manualhooks"
	"go.infratographer.com/load-balancer-api/internal/testutils"
)

const (
	ownerPrefix    = "testown"
	locationPrefix = "testloc"
	defualtTimeout = 2 * time.Second
)

func TestManualHooks(t *testing.T) {
	testutils.SetupDB()

	defer testutils.TeardownDB()

	t.Run("LoadbalancerCreateUpdateHook", LoadbalancerCreateUpdateHookTest)
	t.Run("LoadbalancerDeleteHookTest", LoadbalancerDeleteHookTest)
	t.Run("OriginCreateUpdateHookTest", OriginCreateUpdateHookTest)
	t.Run("OriginDeleteHookTest", OriginDeleteHookTest)
	t.Run("PoolCreateUpdateHookTest", PoolCreateUpdateHookTest)
	t.Run("PoolDeleteHookTest", PoolDeleteHookTest)
	t.Run("PortCreateUpdateHookTest", PortCreateUpdateHookTest)
	t.Run("PortDeleteHookTest", PortDeleteHookTest)
}

func LoadbalancerCreateUpdateHookTest(t *testing.T) {
	// Arrange
	ctx := testutils.MockPermissions(context.Background())

	changesChannel, err := testutils.EventsConn.SubscribeChanges(ctx, "update.load-balancer")
	testutils.IfErrPanic("failed to subscribe to changes", err)

	lb := (&testutils.LoadBalancerBuilder{}).MustNew(ctx)

	testutils.EntClient.LoadBalancer.Use(manualhooks.LoadBalancerHooks()...)

	// Act
	testutils.EntClient.LoadBalancer.UpdateOne(lb).SetName(("other-lb-name")).ExecX(ctx)

	msg := testutils.ChannelReceiveWithTimeout[events.Message[events.ChangeMessage]](changesChannel, defualtTimeout)

	// Assert
	expectedAdditionalSubjectIDs := []gidx.PrefixedID{lb.ID, lb.OwnerID, lb.LocationID}
	actualAdditionalSubjectIDs := msg.Message().AdditionalSubjectIDs

	assert.ElementsMatch(t, expectedAdditionalSubjectIDs, actualAdditionalSubjectIDs)
}

func LoadbalancerDeleteHookTest(t *testing.T) {
	// Arrange
	ctx := testutils.MockPermissions(context.Background())

	changesChannel, err := testutils.EventsConn.SubscribeChanges(ctx, "delete.load-balancer")
	testutils.IfErrPanic("failed to subscribe to changes", err)

	lb := (&testutils.LoadBalancerBuilder{}).MustNew(ctx)

	testutils.EntClient.LoadBalancer.Use(manualhooks.LoadBalancerHooks()...)

	// Act
	testutils.EntClient.LoadBalancer.DeleteOneID(lb.ID).ExecX(ctx)

	msg := testutils.ChannelReceiveWithTimeout[events.Message[events.ChangeMessage]](changesChannel, defualtTimeout)

	// Assert
	expectedAdditionalSubjectIDs := []gidx.PrefixedID{lb.OwnerID, lb.LocationID}
	actualAdditionalSubjectIDs := msg.Message().AdditionalSubjectIDs

	assert.ElementsMatch(t, expectedAdditionalSubjectIDs, actualAdditionalSubjectIDs)
}

func OriginCreateUpdateHookTest(t *testing.T) {
	// Arrange
	ctx := testutils.MockPermissions(context.Background())

	changesChannel, err := testutils.EventsConn.SubscribeChanges(ctx, "update.load-balancer-origin")
	testutils.IfErrPanic("failed to subscribe to changes", err)

	lb := (&testutils.LoadBalancerBuilder{}).MustNew(ctx)
	pool := (&testutils.PoolBuilder{}).MustNew(ctx)
	(&testutils.PortBuilder{PoolIDs: []gidx.PrefixedID{pool.ID}, LoadBalancerID: lb.ID}).MustNew(ctx)
	origin := (&testutils.OriginBuilder{PoolID: pool.ID}).MustNew(ctx)

	testutils.EntClient.Origin.Use(manualhooks.OriginHooks()...)

	// Act
	testutils.EntClient.Origin.UpdateOne(origin).SetName("other-origin-name").ExecX(ctx)

	msg := testutils.ChannelReceiveWithTimeout[events.Message[events.ChangeMessage]](changesChannel, defualtTimeout)

	// Assert
	expectedAdditionalSubjectIDs := []gidx.PrefixedID{pool.ID, pool.OwnerID, lb.ID, lb.LocationID}
	actualAdditionalSubjectIDs := msg.Message().AdditionalSubjectIDs

	assert.ElementsMatch(t, expectedAdditionalSubjectIDs, actualAdditionalSubjectIDs)
}

func OriginDeleteHookTest(t *testing.T) {
	// Arrange
	ctx := testutils.MockPermissions(context.Background())

	changesChannel, err := testutils.EventsConn.SubscribeChanges(ctx, "delete.load-balancer-origin")
	testutils.IfErrPanic("failed to subscribe to changes", err)

	lb := (&testutils.LoadBalancerBuilder{}).MustNew(ctx)
	pool := (&testutils.PoolBuilder{}).MustNew(ctx)
	(&testutils.PortBuilder{PoolIDs: []gidx.PrefixedID{pool.ID}, LoadBalancerID: lb.ID}).MustNew(ctx)
	origin := (&testutils.OriginBuilder{PoolID: pool.ID}).MustNew(ctx)

	testutils.EntClient.Origin.Use(manualhooks.OriginHooks()...)

	// Act
	testutils.EntClient.Origin.DeleteOne(origin).ExecX(ctx)

	msg := testutils.ChannelReceiveWithTimeout[events.Message[events.ChangeMessage]](changesChannel, defualtTimeout)

	// Assert
	expectedAdditionalSubjectIDs := []gidx.PrefixedID{pool.ID, pool.OwnerID, lb.ID, lb.LocationID}
	actualAdditionalSubjectIDs := msg.Message().AdditionalSubjectIDs

	assert.ElementsMatch(t, expectedAdditionalSubjectIDs, actualAdditionalSubjectIDs)
}

func PoolCreateUpdateHookTest(t *testing.T) {
	// Arrange
	ctx := testutils.MockPermissions(context.Background())

	changesChannel, err := testutils.EventsConn.SubscribeChanges(ctx, "update.load-balancer-pool")
	testutils.IfErrPanic("failed to subscribe to changes", err)

	lb := (&testutils.LoadBalancerBuilder{}).MustNew(ctx)
	pool := (&testutils.PoolBuilder{}).MustNew(ctx)
	port := (&testutils.PortBuilder{PoolIDs: []gidx.PrefixedID{pool.ID}, LoadBalancerID: lb.ID}).MustNew(ctx)
	origin := (&testutils.OriginBuilder{PoolID: pool.ID}).MustNew(ctx)

	testutils.EntClient.Pool.Use(manualhooks.PoolHooks()...)

	// Act
	testutils.EntClient.Pool.UpdateOne(pool).SetName("other-pool-name").ExecX(ctx)

	msg := testutils.ChannelReceiveWithTimeout[events.Message[events.ChangeMessage]](changesChannel, defualtTimeout)

	// Assert
	expectedAdditionalSubjectIDs := []gidx.PrefixedID{pool.ID, pool.OwnerID, lb.ID, lb.LocationID, origin.ID, port.ID}
	actualAdditionalSubjectIDs := msg.Message().AdditionalSubjectIDs

	assert.ElementsMatch(t, expectedAdditionalSubjectIDs, actualAdditionalSubjectIDs)
}

func PoolDeleteHookTest(t *testing.T) {
	// Arrange
	ctx := testutils.MockPermissions(context.Background())

	changesChannel, err := testutils.EventsConn.SubscribeChanges(ctx, "delete.load-balancer-pool")
	testutils.IfErrPanic("failed to subscribe to changes", err)

	lb := (&testutils.LoadBalancerBuilder{}).MustNew(ctx)
	pool := (&testutils.PoolBuilder{}).MustNew(ctx)
	(&testutils.PortBuilder{PoolIDs: []gidx.PrefixedID{pool.ID}, LoadBalancerID: lb.ID}).MustNew(ctx)

	testutils.EntClient.Pool.Use(manualhooks.PoolHooks()...)

	// Act
	testutils.EntClient.Pool.DeleteOne(pool).ExecX(ctx)

	msg := testutils.ChannelReceiveWithTimeout[events.Message[events.ChangeMessage]](changesChannel, defualtTimeout)

	// Assert
	expectedAdditionalSubjectIDs := []gidx.PrefixedID{pool.OwnerID, lb.ID, lb.LocationID}
	actualAdditionalSubjectIDs := msg.Message().AdditionalSubjectIDs

	assert.ElementsMatch(t, expectedAdditionalSubjectIDs, actualAdditionalSubjectIDs)
}

func PortCreateUpdateHookTest(t *testing.T) {
	// Arrange
	ctx := testutils.MockPermissions(context.Background())

	changesChannel, err := testutils.EventsConn.SubscribeChanges(ctx, "update.load-balancer-port")
	testutils.IfErrPanic("failed to subscribe to changes", err)

	lb := (&testutils.LoadBalancerBuilder{}).MustNew(ctx)
	pool := (&testutils.PoolBuilder{}).MustNew(ctx)
	port := (&testutils.PortBuilder{PoolIDs: []gidx.PrefixedID{pool.ID}, LoadBalancerID: lb.ID}).MustNew(ctx)

	testutils.EntClient.Port.Use(manualhooks.PortHooks()...)

	// Act
	testutils.EntClient.Port.UpdateOne(port).SetName("other-port-name").ExecX(ctx)

	msg := testutils.ChannelReceiveWithTimeout[events.Message[events.ChangeMessage]](changesChannel, defualtTimeout)

	// Assert
	expectedAdditionalSubjectIDs := []gidx.PrefixedID{pool.ID, pool.OwnerID, lb.ID, lb.LocationID, lb.ProviderID, lb.OwnerID}
	actualAdditionalSubjectIDs := msg.Message().AdditionalSubjectIDs

	assert.ElementsMatch(t, expectedAdditionalSubjectIDs, actualAdditionalSubjectIDs)
}

func PortDeleteHookTest(t *testing.T) {
	// Arrange
	ctx := testutils.MockPermissions(context.Background())

	changesChannel, err := testutils.EventsConn.SubscribeChanges(ctx, "delete.load-balancer-port")
	testutils.IfErrPanic("failed to subscribe to changes", err)

	lb := (&testutils.LoadBalancerBuilder{}).MustNew(ctx)
	pool := (&testutils.PoolBuilder{}).MustNew(ctx)
	port := (&testutils.PortBuilder{PoolIDs: []gidx.PrefixedID{pool.ID}, LoadBalancerID: lb.ID}).MustNew(ctx)

	testutils.EntClient.Port.Use(manualhooks.PortHooks()...)

	// Act
	testutils.EntClient.Port.DeleteOne(port).ExecX(ctx)

	msg := testutils.ChannelReceiveWithTimeout[events.Message[events.ChangeMessage]](changesChannel, defualtTimeout)

	// Assert
	expectedAdditionalSubjectIDs := []gidx.PrefixedID{lb.OwnerID, lb.ID, lb.LocationID, lb.ProviderID}
	actualAdditionalSubjectIDs := msg.Message().AdditionalSubjectIDs

	assert.ElementsMatch(t, expectedAdditionalSubjectIDs, actualAdditionalSubjectIDs)
}
