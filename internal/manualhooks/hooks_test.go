package manualhooks_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.infratographer.com/x/events"
	"go.infratographer.com/x/gidx"

	"go.infratographer.com/load-balancer-api/internal/manualhooks"
	"go.infratographer.com/load-balancer-api/internal/testutils"
)

const (
	ownerPrefix    = "testown"
	locationPrefix = "testloc"
	defaultTimeout = 2 * time.Second
)

var (
	createEventType = string(events.CreateChangeType)
	updateEventType = string(events.UpdateChangeType)
	deleteEventType = string(events.DeleteChangeType)
)

func TestMain(m *testing.M) {
	// setup the database
	testutils.SetupDB()

	// run the tests
	code := m.Run()

	// teardown the database
	testutils.TeardownDB()

	// return the test response code
	os.Exit(code)
}

func Test_LoadbalancerCreateHook(t *testing.T) {
	// Arrange
	ctx := testutils.MockPermissions(context.Background())

	changesChannel, err := testutils.EventsConn.SubscribeChanges(ctx, "create.load-balancer")
	require.NoError(t, err, "failed to subscribe to changes")

	testutils.EntClient.LoadBalancer.Use(manualhooks.LoadBalancerHooks()...)

	// Act
	lb := (&testutils.LoadBalancerBuilder{}).MustNew(ctx)

	msg := testutils.ChannelReceiveWithTimeout[events.Message[events.ChangeMessage]](t, changesChannel, defaultTimeout)

	// Assert
	expectedAdditionalSubjectIDs := []gidx.PrefixedID{lb.ID, lb.OwnerID, lb.LocationID}
	actualAdditionalSubjectIDs := msg.Message().AdditionalSubjectIDs

	assert.ElementsMatch(t, expectedAdditionalSubjectIDs, actualAdditionalSubjectIDs)
	assert.Equal(t, lb.ID, msg.Message().SubjectID)
	assert.Equal(t, createEventType, msg.Message().EventType)
}

func Test_LoadbalancerUpdateHook(t *testing.T) {
	// Arrange
	ctx := testutils.MockPermissions(context.Background())

	changesChannel, err := testutils.EventsConn.SubscribeChanges(ctx, "update.load-balancer")
	require.NoError(t, err, "failed to subscribe to changes")

	lb := (&testutils.LoadBalancerBuilder{}).MustNew(ctx)

	testutils.EntClient.LoadBalancer.Use(manualhooks.LoadBalancerHooks()...)

	// Act
	testutils.EntClient.LoadBalancer.UpdateOne(lb).SetName(("other-lb-name")).ExecX(ctx)

	msg := testutils.ChannelReceiveWithTimeout[events.Message[events.ChangeMessage]](t, changesChannel, defaultTimeout)

	// Assert
	expectedAdditionalSubjectIDs := []gidx.PrefixedID{lb.ID, lb.OwnerID, lb.LocationID}
	actualAdditionalSubjectIDs := msg.Message().AdditionalSubjectIDs

	assert.ElementsMatch(t, expectedAdditionalSubjectIDs, actualAdditionalSubjectIDs)
	assert.Equal(t, lb.ID, msg.Message().SubjectID)
	assert.Equal(t, updateEventType, msg.Message().EventType)
}

func Test_LoadbalancerDeleteHook(t *testing.T) {
	// Arrange
	ctx := testutils.MockPermissions(context.Background())

	changesChannel, err := testutils.EventsConn.SubscribeChanges(ctx, "delete.load-balancer")
	require.NoError(t, err, "failed to subscribe to changes")

	lb := (&testutils.LoadBalancerBuilder{}).MustNew(ctx)

	testutils.EntClient.LoadBalancer.Use(manualhooks.LoadBalancerHooks()...)

	// Act
	testutils.EntClient.LoadBalancer.DeleteOneID(lb.ID).ExecX(ctx)

	msg := testutils.ChannelReceiveWithTimeout[events.Message[events.ChangeMessage]](t, changesChannel, defaultTimeout)

	// Assert
	expectedAdditionalSubjectIDs := []gidx.PrefixedID{lb.OwnerID, lb.LocationID}
	actualAdditionalSubjectIDs := msg.Message().AdditionalSubjectIDs

	assert.ElementsMatch(t, expectedAdditionalSubjectIDs, actualAdditionalSubjectIDs)
	assert.Equal(t, lb.ID, msg.Message().SubjectID)
	assert.Equal(t, deleteEventType, msg.Message().EventType)
}

func Test_OriginCreateHook(t *testing.T) {
	// Arrange
	ctx := testutils.MockPermissions(context.Background())

	changesChannel, err := testutils.EventsConn.SubscribeChanges(ctx, "create.load-balancer-origin")
	require.NoError(t, err, "failed to subscribe to changes")

	lb := (&testutils.LoadBalancerBuilder{}).MustNew(ctx)
	pool := (&testutils.PoolBuilder{}).MustNew(ctx)
	(&testutils.PortBuilder{PoolIDs: []gidx.PrefixedID{pool.ID}, LoadBalancerID: lb.ID}).MustNew(ctx)

	testutils.EntClient.Origin.Use(manualhooks.OriginHooks()...)

	// Act
	origin := (&testutils.OriginBuilder{PoolID: pool.ID}).MustNew(ctx)

	msg := testutils.ChannelReceiveWithTimeout[events.Message[events.ChangeMessage]](t, changesChannel, defaultTimeout)

	// Assert
	expectedAdditionalSubjectIDs := []gidx.PrefixedID{pool.ID, pool.OwnerID, lb.ID, lb.LocationID}
	actualAdditionalSubjectIDs := msg.Message().AdditionalSubjectIDs

	assert.ElementsMatch(t, expectedAdditionalSubjectIDs, actualAdditionalSubjectIDs)
	assert.Equal(t, origin.ID, msg.Message().SubjectID)
	assert.Equal(t, createEventType, msg.Message().EventType)
}

func Test_OriginUpdateHook(t *testing.T) {
	// Arrange
	ctx := testutils.MockPermissions(context.Background())

	changesChannel, err := testutils.EventsConn.SubscribeChanges(ctx, "update.load-balancer-origin")
	require.NoError(t, err, "failed to subscribe to changes")

	lb := (&testutils.LoadBalancerBuilder{}).MustNew(ctx)
	pool := (&testutils.PoolBuilder{}).MustNew(ctx)
	(&testutils.PortBuilder{PoolIDs: []gidx.PrefixedID{pool.ID}, LoadBalancerID: lb.ID}).MustNew(ctx)
	origin := (&testutils.OriginBuilder{PoolID: pool.ID}).MustNew(ctx)

	testutils.EntClient.Origin.Use(manualhooks.OriginHooks()...)

	// Act
	testutils.EntClient.Origin.UpdateOne(origin).SetName("other-origin-name").ExecX(ctx)

	msg := testutils.ChannelReceiveWithTimeout[events.Message[events.ChangeMessage]](t, changesChannel, defaultTimeout)

	// Assert
	expectedAdditionalSubjectIDs := []gidx.PrefixedID{pool.ID, pool.OwnerID, lb.ID, lb.LocationID}
	actualAdditionalSubjectIDs := msg.Message().AdditionalSubjectIDs

	assert.ElementsMatch(t, expectedAdditionalSubjectIDs, actualAdditionalSubjectIDs)
	assert.Equal(t, origin.ID, msg.Message().SubjectID)
	assert.Equal(t, updateEventType, msg.Message().EventType)
}

func Test_OriginDeleteHook(t *testing.T) {
	// Arrange
	ctx := testutils.MockPermissions(context.Background())

	changesChannel, err := testutils.EventsConn.SubscribeChanges(ctx, "delete.load-balancer-origin")
	require.NoError(t, err, "failed to subscribe to changes")

	lb := (&testutils.LoadBalancerBuilder{}).MustNew(ctx)
	pool := (&testutils.PoolBuilder{}).MustNew(ctx)
	(&testutils.PortBuilder{PoolIDs: []gidx.PrefixedID{pool.ID}, LoadBalancerID: lb.ID}).MustNew(ctx)
	origin := (&testutils.OriginBuilder{PoolID: pool.ID}).MustNew(ctx)

	testutils.EntClient.Origin.Use(manualhooks.OriginHooks()...)

	// Act
	testutils.EntClient.Origin.DeleteOne(origin).ExecX(ctx)

	msg := testutils.ChannelReceiveWithTimeout[events.Message[events.ChangeMessage]](t, changesChannel, defaultTimeout)

	// Assert
	expectedAdditionalSubjectIDs := []gidx.PrefixedID{pool.ID, pool.OwnerID, lb.ID, lb.LocationID}
	actualAdditionalSubjectIDs := msg.Message().AdditionalSubjectIDs

	assert.ElementsMatch(t, expectedAdditionalSubjectIDs, actualAdditionalSubjectIDs)
	assert.Equal(t, origin.ID, msg.Message().SubjectID)
	assert.Equal(t, deleteEventType, msg.Message().EventType)
}

func Test_PoolCreateHook(t *testing.T) {
	// Arrange
	ctx := testutils.MockPermissions(context.Background())

	changesChannel, err := testutils.EventsConn.SubscribeChanges(ctx, "create.load-balancer-pool")
	require.NoError(t, err, "failed to subscribe to changes")

	testutils.EntClient.Pool.Use(manualhooks.PoolHooks()...)

	// Act
	pool := (&testutils.PoolBuilder{}).MustNew(ctx)

	msg := testutils.ChannelReceiveWithTimeout[events.Message[events.ChangeMessage]](t, changesChannel, defaultTimeout)

	// Assert
	expectedAdditionalSubjectIDs := []gidx.PrefixedID{pool.OwnerID}
	actualAdditionalSubjectIDs := msg.Message().AdditionalSubjectIDs

	assert.ElementsMatch(t, expectedAdditionalSubjectIDs, actualAdditionalSubjectIDs)
	assert.Equal(t, pool.ID, msg.Message().SubjectID)
	assert.Equal(t, createEventType, msg.Message().EventType)
}

func Test_PoolUpdateHook(t *testing.T) {
	// Arrange
	ctx := testutils.MockPermissions(context.Background())

	changesChannel, err := testutils.EventsConn.SubscribeChanges(ctx, "update.load-balancer-pool")
	require.NoError(t, err, "failed to subscribe to changes")

	lb := (&testutils.LoadBalancerBuilder{}).MustNew(ctx)
	pool := (&testutils.PoolBuilder{}).MustNew(ctx)
	port := (&testutils.PortBuilder{PoolIDs: []gidx.PrefixedID{pool.ID}, LoadBalancerID: lb.ID}).MustNew(ctx)
	origin := (&testutils.OriginBuilder{PoolID: pool.ID}).MustNew(ctx)

	testutils.EntClient.Pool.Use(manualhooks.PoolHooks()...)

	// Act
	testutils.EntClient.Pool.UpdateOne(pool).SetName("other-pool-name").ExecX(ctx)

	msg := testutils.ChannelReceiveWithTimeout[events.Message[events.ChangeMessage]](t, changesChannel, defaultTimeout)

	// Assert
	expectedAdditionalSubjectIDs := []gidx.PrefixedID{pool.ID, pool.OwnerID, lb.ID, lb.LocationID, origin.ID, port.ID}
	actualAdditionalSubjectIDs := msg.Message().AdditionalSubjectIDs

	assert.ElementsMatch(t, expectedAdditionalSubjectIDs, actualAdditionalSubjectIDs)
	assert.Equal(t, pool.ID, msg.Message().SubjectID)
	assert.Equal(t, updateEventType, msg.Message().EventType)
}

func Test_PoolDeleteHook(t *testing.T) {
	// Arrange
	ctx := testutils.MockPermissions(context.Background())

	changesChannel, err := testutils.EventsConn.SubscribeChanges(ctx, "delete.load-balancer-pool")
	require.NoError(t, err, "failed to subscribe to changes")

	lb := (&testutils.LoadBalancerBuilder{}).MustNew(ctx)
	pool := (&testutils.PoolBuilder{}).MustNew(ctx)
	(&testutils.PortBuilder{PoolIDs: []gidx.PrefixedID{pool.ID}, LoadBalancerID: lb.ID}).MustNew(ctx)

	testutils.EntClient.Pool.Use(manualhooks.PoolHooks()...)

	// Act
	testutils.EntClient.Pool.DeleteOne(pool).ExecX(ctx)

	msg := testutils.ChannelReceiveWithTimeout[events.Message[events.ChangeMessage]](t, changesChannel, defaultTimeout)

	// Assert
	expectedAdditionalSubjectIDs := []gidx.PrefixedID{pool.OwnerID, lb.ID, lb.LocationID}
	actualAdditionalSubjectIDs := msg.Message().AdditionalSubjectIDs

	assert.ElementsMatch(t, expectedAdditionalSubjectIDs, actualAdditionalSubjectIDs)
	assert.Equal(t, pool.ID, msg.Message().SubjectID)
	assert.Equal(t, deleteEventType, msg.Message().EventType)
}

func Test_PortCreateHook(t *testing.T) {
	// Arrange
	ctx := testutils.MockPermissions(context.Background())

	changesChannel, err := testutils.EventsConn.SubscribeChanges(ctx, "create.load-balancer-port")
	require.NoError(t, err, "failed to subscribe to changes")

	lb := (&testutils.LoadBalancerBuilder{}).MustNew(ctx)
	pool := (&testutils.PoolBuilder{}).MustNew(ctx)

	testutils.EntClient.Port.Use(manualhooks.PortHooks()...)

	// Act
	port := (&testutils.PortBuilder{PoolIDs: []gidx.PrefixedID{pool.ID}, LoadBalancerID: lb.ID}).MustNew(ctx)

	msg := testutils.ChannelReceiveWithTimeout[events.Message[events.ChangeMessage]](t, changesChannel, defaultTimeout)

	// Assert
	expectedAdditionalSubjectIDs := []gidx.PrefixedID{pool.ID, pool.OwnerID, lb.ID, lb.LocationID, lb.ProviderID, lb.OwnerID}
	actualAdditionalSubjectIDs := msg.Message().AdditionalSubjectIDs

	assert.ElementsMatch(t, expectedAdditionalSubjectIDs, actualAdditionalSubjectIDs)
	assert.Equal(t, port.ID, msg.Message().SubjectID)
	assert.Equal(t, createEventType, msg.Message().EventType)
}

func Test_PortUpdateHook(t *testing.T) {
	// Arrange
	ctx := testutils.MockPermissions(context.Background())

	changesChannel, err := testutils.EventsConn.SubscribeChanges(ctx, "update.load-balancer-port")
	require.NoError(t, err, "failed to subscribe to changes")

	lb := (&testutils.LoadBalancerBuilder{}).MustNew(ctx)
	pool := (&testutils.PoolBuilder{}).MustNew(ctx)
	port := (&testutils.PortBuilder{PoolIDs: []gidx.PrefixedID{pool.ID}, LoadBalancerID: lb.ID}).MustNew(ctx)

	testutils.EntClient.Port.Use(manualhooks.PortHooks()...)

	// Act
	testutils.EntClient.Port.UpdateOne(port).SetName("other-port-name").ExecX(ctx)

	msg := testutils.ChannelReceiveWithTimeout[events.Message[events.ChangeMessage]](t, changesChannel, defaultTimeout)

	// Assert
	expectedAdditionalSubjectIDs := []gidx.PrefixedID{pool.ID, pool.OwnerID, lb.ID, lb.LocationID, lb.ProviderID, lb.OwnerID}
	actualAdditionalSubjectIDs := msg.Message().AdditionalSubjectIDs

	assert.ElementsMatch(t, expectedAdditionalSubjectIDs, actualAdditionalSubjectIDs)
	assert.Equal(t, port.ID, msg.Message().SubjectID)
	assert.Equal(t, updateEventType, msg.Message().EventType)
}

func Test_PortDeleteHook(t *testing.T) {
	// Arrange
	ctx := testutils.MockPermissions(context.Background())

	changesChannel, err := testutils.EventsConn.SubscribeChanges(ctx, "delete.load-balancer-port")
	require.NoError(t, err, "failed to subscribe to changes")

	lb := (&testutils.LoadBalancerBuilder{}).MustNew(ctx)
	pool := (&testutils.PoolBuilder{}).MustNew(ctx)
	port := (&testutils.PortBuilder{PoolIDs: []gidx.PrefixedID{pool.ID}, LoadBalancerID: lb.ID}).MustNew(ctx)

	testutils.EntClient.Port.Use(manualhooks.PortHooks()...)

	// Act
	testutils.EntClient.Port.DeleteOne(port).ExecX(ctx)

	msg := testutils.ChannelReceiveWithTimeout[events.Message[events.ChangeMessage]](t, changesChannel, defaultTimeout)

	// Assert
	expectedAdditionalSubjectIDs := []gidx.PrefixedID{lb.OwnerID, lb.ID, lb.LocationID, lb.ProviderID}
	actualAdditionalSubjectIDs := msg.Message().AdditionalSubjectIDs

	assert.ElementsMatch(t, expectedAdditionalSubjectIDs, actualAdditionalSubjectIDs)
	assert.Equal(t, port.ID, msg.Message().SubjectID)
	assert.Equal(t, deleteEventType, msg.Message().EventType)
}
