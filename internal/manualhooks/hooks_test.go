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
	defaultTimeout = 5 * time.Second
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
	expectedAdditionalSubjectIDs := []gidx.PrefixedID{lb.ID, lb.OwnerID, lb.LocationID, lb.ProviderID}
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
	expectedAdditionalSubjectIDs := []gidx.PrefixedID{lb.ID, lb.OwnerID, lb.LocationID, lb.ProviderID}
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
	expectedAdditionalSubjectIDs := []gidx.PrefixedID{lb.OwnerID, lb.LocationID, lb.ProviderID}
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
	expectedAdditionalSubjectIDs := []gidx.PrefixedID{pool.ID, pool.OwnerID, lb.ID, lb.LocationID, lb.ProviderID}
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
	expectedAdditionalSubjectIDs := []gidx.PrefixedID{pool.ID, pool.OwnerID, lb.ID, lb.LocationID, lb.ProviderID}
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
	expectedAdditionalSubjectIDs := []gidx.PrefixedID{pool.ID, pool.OwnerID, lb.ID, lb.LocationID, lb.ProviderID}
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
	expectedAdditionalSubjectIDs := []gidx.PrefixedID{pool.OwnerID, lb.ID, lb.LocationID, lb.ProviderID, origin.ID, port.ID}
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
	expectedAdditionalSubjectIDs := []gidx.PrefixedID{pool.OwnerID, lb.ID, lb.LocationID, lb.ProviderID}
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

	testutils.EntClient.Port.Use(manualhooks.PortHooks()...)

	t.Run("with pool", func(t *testing.T) {
		// Act
		lb := (&testutils.LoadBalancerBuilder{}).MustNew(ctx)
		pool := (&testutils.PoolBuilder{OwnerID: lb.OwnerID}).MustNew(ctx)
		port := (&testutils.PortBuilder{PoolIDs: []gidx.PrefixedID{pool.ID}, LoadBalancerID: lb.ID}).MustNew(ctx)

		msg := testutils.ChannelReceiveWithTimeout[events.Message[events.ChangeMessage]](t, changesChannel, defaultTimeout)

		// Assert
		expectedAdditionalSubjectIDs := []gidx.PrefixedID{pool.ID, lb.ID, lb.LocationID, lb.ProviderID, lb.OwnerID}
		actualAdditionalSubjectIDs := msg.Message().AdditionalSubjectIDs

		assert.ElementsMatch(t, expectedAdditionalSubjectIDs, actualAdditionalSubjectIDs)
		assert.Equal(t, port.ID, msg.Message().SubjectID)
		assert.Equal(t, createEventType, msg.Message().EventType)
	})

	t.Run("with multiple pools", func(t *testing.T) {
		// Act
		lb := (&testutils.LoadBalancerBuilder{}).MustNew(ctx)
		pool := (&testutils.PoolBuilder{OwnerID: lb.OwnerID}).MustNew(ctx)
		pool2 := (&testutils.PoolBuilder{OwnerID: lb.OwnerID}).MustNew(ctx)
		port := (&testutils.PortBuilder{PoolIDs: []gidx.PrefixedID{pool.ID, pool2.ID}, LoadBalancerID: lb.ID}).MustNew(ctx)

		msg := testutils.ChannelReceiveWithTimeout[events.Message[events.ChangeMessage]](t, changesChannel, defaultTimeout)

		// Assert
		expectedAdditionalSubjectIDs := []gidx.PrefixedID{pool.ID, pool2.ID, lb.ID, lb.LocationID, lb.ProviderID, lb.OwnerID}
		actualAdditionalSubjectIDs := msg.Message().AdditionalSubjectIDs

		assert.ElementsMatch(t, expectedAdditionalSubjectIDs, actualAdditionalSubjectIDs)
		assert.Equal(t, port.ID, msg.Message().SubjectID)
		assert.Equal(t, createEventType, msg.Message().EventType)
	})

	t.Run("with no pool", func(t *testing.T) {
		// Act
		lb := (&testutils.LoadBalancerBuilder{}).MustNew(ctx)
		port := (&testutils.PortBuilder{PoolIDs: []gidx.PrefixedID{}, LoadBalancerID: lb.ID}).MustNew(ctx)

		msg := testutils.ChannelReceiveWithTimeout[events.Message[events.ChangeMessage]](t, changesChannel, defaultTimeout)

		// Assert
		expectedAdditionalSubjectIDs := []gidx.PrefixedID{lb.ID, lb.LocationID, lb.ProviderID, lb.OwnerID}
		actualAdditionalSubjectIDs := msg.Message().AdditionalSubjectIDs

		assert.ElementsMatch(t, expectedAdditionalSubjectIDs, actualAdditionalSubjectIDs)
		assert.Equal(t, port.ID, msg.Message().SubjectID)
		assert.Equal(t, createEventType, msg.Message().EventType)
	})
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

func Test_MultipleLoadbalancersSharedPoolAddOrigin(t *testing.T) {
	// Scenario: 2 loadbalancers in different locations, with the same owner, share a pool.
	// An origin is added to the shared pool.
	// Assert the owner, loadbalancers, pool, and locations are all included in the additionalSubject ID list

	// Arrange
	ctx := testutils.MockPermissions(context.Background())

	changesChannel, err := testutils.EventsConn.SubscribeChanges(ctx, "create.load-balancer-origin")
	require.NoError(t, err, "failed to subscribe to changes")

	// create 2 loadbalancers with a shared pool of origins
	prov := (&testutils.ProviderBuilder{}).MustNew(ctx)
	lb1 := (&testutils.LoadBalancerBuilder{OwnerID: "tnttent-testing", Provider: prov}).MustNew(ctx)
	lb2 := (&testutils.LoadBalancerBuilder{OwnerID: "tnttent-testing", Provider: prov}).MustNew(ctx)
	pool := (&testutils.PoolBuilder{OwnerID: "tnttent-testing"}).MustNew(ctx)
	_ = (&testutils.PortBuilder{PoolIDs: []gidx.PrefixedID{pool.ID}, LoadBalancerID: lb1.ID}).MustNew(ctx)
	_ = (&testutils.PortBuilder{PoolIDs: []gidx.PrefixedID{pool.ID}, LoadBalancerID: lb2.ID}).MustNew(ctx)

	testutils.EntClient.Origin.Use(manualhooks.OriginHooks()...)

	// Act - add another origin to the pool
	ogn := (&testutils.OriginBuilder{PoolID: pool.ID}).MustNew(ctx)

	msg := testutils.ChannelReceiveWithTimeout[events.Message[events.ChangeMessage]](t, changesChannel, defaultTimeout)

	// Assert
	expectedAdditionalSubjectIDs := []gidx.PrefixedID{
		prov.ID,
		lb1.OwnerID,
		lb1.ID,
		lb2.ID,
		lb1.LocationID,
		lb2.LocationID,
		pool.ID,
	}
	actualAdditionalSubjectIDs := msg.Message().AdditionalSubjectIDs

	assert.ElementsMatch(t, expectedAdditionalSubjectIDs, actualAdditionalSubjectIDs)
	assert.Equal(t, ogn.ID, msg.Message().SubjectID)
	assert.Equal(t, createEventType, msg.Message().EventType)
}

func Test_MultipleLoadbalancersSharedPoolDeleteOrigin(t *testing.T) {
	// Scenario: 2 loadbalancers in different locations, with the same owner, share a pool.
	// An origin is removed from the shared pool.
	// Assert the owner, loadbalancers, pool, and locations are all included in the additionalSubject ID list

	// Arrange
	ctx := testutils.MockPermissions(context.Background())

	changesChannel, err := testutils.EventsConn.SubscribeChanges(ctx, "delete.load-balancer-origin")
	require.NoError(t, err, "failed to subscribe to changes")

	// create 2 loadbalancers with a shared pool of origins
	prov := (&testutils.ProviderBuilder{}).MustNew(ctx)
	lb1 := (&testutils.LoadBalancerBuilder{OwnerID: "tnttent-testing", Provider: prov}).MustNew(ctx)
	lb2 := (&testutils.LoadBalancerBuilder{OwnerID: "tnttent-testing", Provider: prov}).MustNew(ctx)
	pool := (&testutils.PoolBuilder{OwnerID: "tnttent-testing"}).MustNew(ctx)
	_ = (&testutils.PortBuilder{PoolIDs: []gidx.PrefixedID{pool.ID}, LoadBalancerID: lb1.ID}).MustNew(ctx)
	_ = (&testutils.PortBuilder{PoolIDs: []gidx.PrefixedID{pool.ID}, LoadBalancerID: lb2.ID}).MustNew(ctx)
	_ = (&testutils.OriginBuilder{PoolID: pool.ID}).MustNew(ctx)
	ogn2 := (&testutils.OriginBuilder{PoolID: pool.ID}).MustNew(ctx)

	testutils.EntClient.Origin.Use(manualhooks.OriginHooks()...)

	// Act - update the pool to remove an origin
	testutils.EntClient.Origin.DeleteOne(ogn2).ExecX(ctx)

	msg := testutils.ChannelReceiveWithTimeout[events.Message[events.ChangeMessage]](t, changesChannel, defaultTimeout)

	// Assert
	expectedAdditionalSubjectIDs := []gidx.PrefixedID{
		prov.ID,
		lb1.OwnerID,
		lb1.ID,
		lb2.ID,
		lb1.LocationID,
		lb2.LocationID,
		pool.ID,
	}
	actualAdditionalSubjectIDs := msg.Message().AdditionalSubjectIDs

	assert.ElementsMatch(t, expectedAdditionalSubjectIDs, actualAdditionalSubjectIDs)
	assert.Equal(t, ogn2.ID, msg.Message().SubjectID)
	assert.Equal(t, deleteEventType, msg.Message().EventType)
}
