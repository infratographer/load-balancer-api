package manualhooks_test

import (
	"context"
	"log"
	"testing"
	"time"

	"entgo.io/ent/dialect"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.infratographer.com/permissions-api/pkg/permissions/mockpermissions"
	"go.infratographer.com/x/events"
	"go.infratographer.com/x/gidx"
	"go.infratographer.com/x/goosex"
	"go.infratographer.com/x/testing/eventtools"

	"go.infratographer.com/load-balancer-api/db"
	ent "go.infratographer.com/load-balancer-api/internal/ent/generated"
	"go.infratographer.com/load-balancer-api/internal/ent/generated/pool"
	"go.infratographer.com/load-balancer-api/internal/manualhooks"
	"go.infratographer.com/load-balancer-api/internal/testutils"
	"go.infratographer.com/load-balancer-api/x/testcontainersx"
)

const (
	ownerPrefix    = "testown"
	locationPrefix = "testloc"
	defualtTimeout = 2 * time.Second
)

var (
	EventsConn  events.Connection
	EntClient   *ent.Client
	DBContainer *testcontainersx.DBContainer
)

func TestManualHooks(t *testing.T) {
	setup()
	defer teardown()

	t.Run("LoadbalancerCreateUpdateHook", LoadbalancerCreateUpdateHookTest)
	t.Run("LoadbalancerDeleteHookTest", LoadbalancerDeleteHookTest)
	t.Run("OriginCreateUpdateHookTest", OriginCreateUpdateHookTest)
	t.Run("OriginDeleteHookTest", OriginDeleteHookTest)
	t.Run("PoolCreateUpdateHookTest", PoolCreateUpdateHookTest)
	t.Run("PoolDeleteHookTest", PoolDeleteHookTest)
	t.Run("PortCreateUpdateHookTest", PortCreateUpdateHookTest)
	t.Run("PortDeleteHookTest", PortDeleteHookTest)
}

func setup() {
	ctx := context.Background()

	// NATS setup
	nats, err := eventtools.NewNatsServer()
	testutils.IfErrPanic("failed to start nats server", err)

	conn, err := events.NewConnection(nats.Config)
	testutils.IfErrPanic("failed to create events connection", err)

	// DB and EntClient setup
	dia, uri, cntr := testutils.ParseDBURI(ctx)

	c, err := ent.Open(dia, uri, ent.Debug(), ent.EventsPublisher(conn))
	if err != nil {
		log.Println(err)
		testutils.IfErrPanic("failed terminating test db container after failing to connect to the db", cntr.Container.Terminate(ctx))
		testutils.IfErrPanic("failed opening connection to database:", err)
	}

	switch dia {
	case dialect.SQLite:
		// Run automatic migrations for SQLite
		testutils.IfErrPanic("failed creating db schema", c.Schema.Create(ctx))
	case dialect.Postgres:
		log.Println("Running database migrations")
		goosex.MigrateUp(uri, db.Migrations)
	}

	EventsConn = conn
	EntClient = c
	DBContainer = cntr
}

func teardown() {
	ctx := context.Background()

	if EntClient != nil {
		testutils.IfErrPanic("teardown failed to close database connection", EntClient.Close())
	}

	if DBContainer != nil && DBContainer.Container.IsRunning() {
		testutils.IfErrPanic("teardown failed to terminate test db container", DBContainer.Container.Terminate(ctx))
	}
}

func mockPermissions(ctx context.Context) context.Context {
	// mock permissions
	perms := new(mockpermissions.MockPermissions)
	perms.On("CreateAuthRelationships", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	perms.On("DeleteAuthRelationships", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	ctx = perms.ContextWithHandler(ctx)

	return ctx
}

func LoadbalancerCreateUpdateHookTest(t *testing.T) {
	// Arrange
	ctx := mockPermissions(context.Background())

	changesChannel, err := EventsConn.SubscribeChanges(ctx, "update.load-balancer")
	testutils.IfErrPanic("failed to subscribe to changes", err)

	ownerId := gidx.MustNewID(ownerPrefix)
	provider := EntClient.Provider.Create().SetName("provider-name").SetOwnerID(ownerId).SaveX(ctx)
	lb := EntClient.LoadBalancer.Create().SetName("lb-name").SetProvider(provider).SetOwnerID(ownerId).SetLocationID(gidx.MustNewID(locationPrefix)).SaveX(ctx)

	EntClient.LoadBalancer.Use(manualhooks.LoadBalancerHooks()...)

	// Act
	EntClient.LoadBalancer.UpdateOne(lb).SetName(("other-lb-name")).ExecX(ctx)

	msg := testutils.ChannelReceiveWithTimeout[events.Message[events.ChangeMessage]](changesChannel, defualtTimeout)

	// Assert
	expectedAdditionalSubjectIDs := []gidx.PrefixedID{lb.ID, lb.OwnerID, lb.LocationID}
	actualAdditionalSubjectIDs := msg.Message().AdditionalSubjectIDs

	assert.ElementsMatch(t, expectedAdditionalSubjectIDs, actualAdditionalSubjectIDs)
}

func LoadbalancerDeleteHookTest(t *testing.T) {
	// Arrange
	ctx := mockPermissions(context.Background())

	changesChannel, err := EventsConn.SubscribeChanges(ctx, "delete.load-balancer")
	testutils.IfErrPanic("failed to subscribe to changes", err)

	ownerId := gidx.MustNewID(ownerPrefix)
	provider := EntClient.Provider.Create().SetName("provider-name").SetOwnerID(ownerId).SaveX(ctx)
	lb := EntClient.LoadBalancer.Create().SetName("lb-name").SetProvider(provider).SetOwnerID(ownerId).SetLocationID(gidx.MustNewID(locationPrefix)).SaveX(ctx)

	EntClient.LoadBalancer.Use(manualhooks.LoadBalancerHooks()...)

	// Act
	EntClient.LoadBalancer.DeleteOneID(lb.ID).ExecX(ctx)

	msg := testutils.ChannelReceiveWithTimeout[events.Message[events.ChangeMessage]](changesChannel, defualtTimeout)

	// Assert
	expectedAdditionalSubjectIDs := []gidx.PrefixedID{lb.OwnerID, lb.LocationID}
	actualAdditionalSubjectIDs := msg.Message().AdditionalSubjectIDs

	assert.ElementsMatch(t, expectedAdditionalSubjectIDs, actualAdditionalSubjectIDs)
}

func OriginCreateUpdateHookTest(t *testing.T) {
	// Arrange
	ctx := mockPermissions(context.Background())

	changesChannel, err := EventsConn.SubscribeChanges(ctx, "update.load-balancer-origin")
	testutils.IfErrPanic("failed to subscribe to changes", err)

	ownerId := gidx.MustNewID(ownerPrefix)

	provider := EntClient.Provider.Create().SetName("provider-name").SetOwnerID(ownerId).SaveX(ctx)
	lb := EntClient.LoadBalancer.Create().SetName("lb-name").SetProvider(provider).SetOwnerID(ownerId).SetLocationID(gidx.MustNewID(locationPrefix)).SaveX(ctx)
	pool := EntClient.Pool.Create().SetName("pool-name").SetOwnerID(ownerId).SetProtocol(pool.ProtocolTCP).SaveX(ctx)
	EntClient.Port.Create().SetName("port-name").AddPoolIDs(pool.ID).SetLoadBalancer(lb).SetNumber(11).SaveX(ctx)
	origin := EntClient.Origin.Create().SetName("origin-name").SetPool(pool).SetTarget("127.0.0.1").SetPortNumber(12).SaveX(ctx)

	EntClient.Origin.Use(manualhooks.OriginHooks()...)

	// Act
	EntClient.Origin.UpdateOne(origin).SetName("other-origin-name").ExecX(ctx)

	msg := testutils.ChannelReceiveWithTimeout[events.Message[events.ChangeMessage]](changesChannel, defualtTimeout)

	// Assert
	expectedAdditionalSubjectIDs := []gidx.PrefixedID{pool.ID, ownerId, lb.ID, lb.LocationID}
	actualAdditionalSubjectIDs := msg.Message().AdditionalSubjectIDs

	assert.ElementsMatch(t, expectedAdditionalSubjectIDs, actualAdditionalSubjectIDs)
}

func OriginDeleteHookTest(t *testing.T) {
	// Arrange
	ctx := mockPermissions(context.Background())

	changesChannel, err := EventsConn.SubscribeChanges(ctx, "delete.load-balancer-origin")
	testutils.IfErrPanic("failed to subscribe to changes", err)

	ownerId := gidx.MustNewID(ownerPrefix)

	provider := EntClient.Provider.Create().SetName("provider-name").SetOwnerID(ownerId).SaveX(ctx)
	lb := EntClient.LoadBalancer.Create().SetName("lb-name").SetProvider(provider).SetOwnerID(ownerId).SetLocationID(gidx.MustNewID(locationPrefix)).SaveX(ctx)
	pool := EntClient.Pool.Create().SetName("pool-name").SetOwnerID(ownerId).SetProtocol(pool.ProtocolTCP).SaveX(ctx)
	EntClient.Port.Create().SetName("port-name").AddPoolIDs(pool.ID).SetLoadBalancer(lb).SetNumber(11).SaveX(ctx)
	origin := EntClient.Origin.Create().SetName("origin-name").SetPool(pool).SetTarget("127.0.0.1").SetPortNumber(12).SaveX(ctx)

	EntClient.Origin.Use(manualhooks.OriginHooks()...)

	// Act
	EntClient.Origin.DeleteOne(origin).ExecX(ctx)

	msg := testutils.ChannelReceiveWithTimeout[events.Message[events.ChangeMessage]](changesChannel, defualtTimeout)

	// Assert
	expectedAdditionalSubjectIDs := []gidx.PrefixedID{pool.ID, ownerId, lb.ID, lb.LocationID}
	actualAdditionalSubjectIDs := msg.Message().AdditionalSubjectIDs

	assert.ElementsMatch(t, expectedAdditionalSubjectIDs, actualAdditionalSubjectIDs)
}

func PoolCreateUpdateHookTest(t *testing.T) {
	// Arrange
	ctx := mockPermissions(context.Background())

	changesChannel, err := EventsConn.SubscribeChanges(ctx, "update.load-balancer-pool")
	testutils.IfErrPanic("failed to subscribe to changes", err)

	ownerId := gidx.MustNewID(ownerPrefix)

	provider := EntClient.Provider.Create().SetName("provider-name").SetOwnerID(ownerId).SaveX(ctx)
	lb := EntClient.LoadBalancer.Create().SetName("lb-name").SetProvider(provider).SetOwnerID(ownerId).SetLocationID(gidx.MustNewID(locationPrefix)).SaveX(ctx)
	pool := EntClient.Pool.Create().SetName("pool-name").SetOwnerID(ownerId).SetProtocol(pool.ProtocolTCP).SaveX(ctx)
	port := EntClient.Port.Create().SetName("port-name").AddPoolIDs(pool.ID).SetLoadBalancer(lb).SetNumber(11).SaveX(ctx)
	origin := EntClient.Origin.Create().SetName("origin-name").SetPool(pool).SetTarget("127.0.0.1").SetPortNumber(12).SaveX(ctx)

	EntClient.Pool.Use(manualhooks.PoolHooks()...)

	// Act
	EntClient.Pool.UpdateOne(pool).SetName("other-pool-name").ExecX(ctx)

	msg := testutils.ChannelReceiveWithTimeout[events.Message[events.ChangeMessage]](changesChannel, defualtTimeout)

	// Assert
	expectedAdditionalSubjectIDs := []gidx.PrefixedID{pool.ID, ownerId, lb.ID, lb.LocationID, origin.ID, port.ID}
	actualAdditionalSubjectIDs := msg.Message().AdditionalSubjectIDs

	assert.ElementsMatch(t, expectedAdditionalSubjectIDs, actualAdditionalSubjectIDs)
}

func PoolDeleteHookTest(t *testing.T) {
	// Arrange
	ctx := mockPermissions(context.Background())

	changesChannel, err := EventsConn.SubscribeChanges(ctx, "delete.load-balancer-pool")
	testutils.IfErrPanic("failed to subscribe to changes", err)

	ownerId := gidx.MustNewID(ownerPrefix)

	provider := EntClient.Provider.Create().SetName("provider-name").SetOwnerID(ownerId).SaveX(ctx)
	lb := EntClient.LoadBalancer.Create().SetName("lb-name").SetProvider(provider).SetOwnerID(ownerId).SetLocationID(gidx.MustNewID(locationPrefix)).SaveX(ctx)
	pool := EntClient.Pool.Create().SetName("pool-name").SetOwnerID(ownerId).SetProtocol(pool.ProtocolTCP).SaveX(ctx)
	EntClient.Port.Create().AddPoolIDs(pool.ID).SetLoadBalancer(lb).SetName("port-name").SetNumber(11).SaveX(ctx)

	EntClient.Pool.Use(manualhooks.PoolHooks()...)

	// Act
	EntClient.Pool.DeleteOne(pool).ExecX(ctx)

	msg := testutils.ChannelReceiveWithTimeout[events.Message[events.ChangeMessage]](changesChannel, defualtTimeout)

	// Assert
	expectedAdditionalSubjectIDs := []gidx.PrefixedID{ownerId, lb.ID, lb.LocationID}
	actualAdditionalSubjectIDs := msg.Message().AdditionalSubjectIDs

	assert.ElementsMatch(t, expectedAdditionalSubjectIDs, actualAdditionalSubjectIDs)
}

func PortCreateUpdateHookTest(t *testing.T) {
	// Arrange
	ctx := mockPermissions(context.Background())

	changesChannel, err := EventsConn.SubscribeChanges(ctx, "update.load-balancer-port")
	testutils.IfErrPanic("failed to subscribe to changes", err)

	ownerId := gidx.MustNewID(ownerPrefix)

	provider := EntClient.Provider.Create().SetName("provider-name").SetOwnerID(ownerId).SaveX(ctx)
	lb := EntClient.LoadBalancer.Create().SetName("lb-name").SetProvider(provider).SetOwnerID(ownerId).SetLocationID(gidx.MustNewID(locationPrefix)).SaveX(ctx)
	pool := EntClient.Pool.Create().SetName("pool-name").SetOwnerID(ownerId).SetProtocol(pool.ProtocolTCP).SaveX(ctx)
	port := EntClient.Port.Create().AddPoolIDs(pool.ID).SetLoadBalancer(lb).SetName("port-name").SetNumber(11).SaveX(ctx)

	EntClient.Port.Use(manualhooks.PortHooks()...)

	// Act
	EntClient.Port.UpdateOne(port).SetName("other-port-name").ExecX(ctx)

	msg := testutils.ChannelReceiveWithTimeout[events.Message[events.ChangeMessage]](changesChannel, defualtTimeout)

	// Assert
	expectedAdditionalSubjectIDs := []gidx.PrefixedID{ownerId, lb.ID, lb.LocationID, provider.ID, pool.ID}
	actualAdditionalSubjectIDs := msg.Message().AdditionalSubjectIDs

	assert.ElementsMatch(t, expectedAdditionalSubjectIDs, actualAdditionalSubjectIDs)
}

func PortDeleteHookTest(t *testing.T) {
	// Arrange
	ctx := mockPermissions(context.Background())

	changesChannel, err := EventsConn.SubscribeChanges(ctx, "delete.load-balancer-port")
	testutils.IfErrPanic("failed to subscribe to changes", err)

	ownerId := gidx.MustNewID(ownerPrefix)

	provider := EntClient.Provider.Create().SetName("provider-name").SetOwnerID(ownerId).SaveX(ctx)
	lb := EntClient.LoadBalancer.Create().SetName("lb-name").SetProvider(provider).SetOwnerID(ownerId).SetLocationID(gidx.MustNewID(locationPrefix)).SaveX(ctx)
	pool := EntClient.Pool.Create().SetName("pool-name").SetOwnerID(ownerId).SetProtocol(pool.ProtocolTCP).SaveX(ctx)
	port := EntClient.Port.Create().AddPoolIDs(pool.ID).SetLoadBalancer(lb).SetName("port-name").SetNumber(11).SaveX(ctx)

	EntClient.Port.Use(manualhooks.PortHooks()...)

	// Act
	EntClient.Port.DeleteOne(port).ExecX(ctx)

	msg := testutils.ChannelReceiveWithTimeout[events.Message[events.ChangeMessage]](changesChannel, defualtTimeout)

	// Assert
	expectedAdditionalSubjectIDs := []gidx.PrefixedID{ownerId, lb.ID, lb.LocationID, provider.ID}
	actualAdditionalSubjectIDs := msg.Message().AdditionalSubjectIDs

	assert.ElementsMatch(t, expectedAdditionalSubjectIDs, actualAdditionalSubjectIDs)
}
