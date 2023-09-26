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

var (
	ownerPrefix    = "testown"
	locationPrefix = "testloc"
	defualtTimeout = 2 * time.Second
)

func setup(subChangesTopic string) (context.Context, *ent.Client, <-chan events.Message[events.ChangeMessage], func()) {
	ctx := context.Background()

	// NATS setup
	nats, err := eventtools.NewNatsServer()
	testutils.IfErrPanic("failed to start nats server", err)

	conn, err := events.NewConnection(nats.Config)
	testutils.IfErrPanic("failed to create events connection", err)

	subChan, err := conn.SubscribeChanges(ctx, subChangesTopic)
	testutils.IfErrPanic("failed to subscribe to changes", err)

	// DB and EntClient setup
	dia, uri, cntr := testutils.ParseDBURI(ctx)

	entClient, err := ent.Open(dia, uri, ent.Debug(), ent.EventsPublisher(conn))
	if err != nil {
		log.Println(err)
		testutils.IfErrPanic("failed terminating test db container after failing to connect to the db", cntr.Container.Terminate(ctx))
		testutils.IfErrPanic("failed opening connection to database:", err)
	}

	switch dia {
	case dialect.SQLite:
		// Run automatic migrations for SQLite
		testutils.IfErrPanic("failed creating db schema", entClient.Schema.Create(ctx))
	case dialect.Postgres:
		log.Println("Running database migrations")
		goosex.MigrateUp(uri, db.Migrations)
	}

	// mock permissions
	perms := new(mockpermissions.MockPermissions)
	perms.On("CreateAuthRelationships", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	perms.On("DeleteAuthRelationships", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	ctx = perms.ContextWithHandler(ctx)

	teardown := func() { teardown(entClient, cntr) }

	return ctx, entClient, subChan, teardown
}

func teardown(entClient *ent.Client, dbc *testcontainersx.DBContainer) {
	ctx := context.Background()

	if entClient != nil {
		testutils.IfErrPanic("teardown failed to close database connection", entClient.Close())
	}

	if dbc != nil {
		testutils.IfErrPanic("teardown failed to terminate test db container", dbc.Container.Terminate(ctx))
	}
}

func Test_LoadbalancerCreateUpdateHook(t *testing.T) {
	// Arrange
	ctx, entClient, changesChannel, teardown := setup("update.load-balancer")
	defer teardown()

	ownerId := gidx.MustNewID(ownerPrefix)
	provider := entClient.Provider.Create().SetName("provider-name").SetOwnerID(ownerId).SaveX(ctx)
	lb := entClient.LoadBalancer.Create().SetName("lb-name").SetProvider(provider).SetOwnerID(ownerId).SetLocationID(gidx.MustNewID(locationPrefix)).SaveX(ctx)

	entClient.LoadBalancer.Use(manualhooks.LoadBalancerHooks()...)

	// Act
	entClient.LoadBalancer.UpdateOne(lb).SetName(("other-lb-name")).ExecX(ctx)

	msg := testutils.ChannelReceiveWithTimeout[events.Message[events.ChangeMessage]](changesChannel, defualtTimeout)

	// Assert
	expectedAdditionalSubjectIDs := []gidx.PrefixedID{lb.ID, lb.OwnerID, lb.LocationID}
	actualAdditionalSubjectIDs := msg.Message().AdditionalSubjectIDs

	assert.ElementsMatch(t, expectedAdditionalSubjectIDs, actualAdditionalSubjectIDs)
}

func Test_LoadbalancerDeleteHook(t *testing.T) {
	// Arrange
	ctx, entClient, changesChannel, teardown := setup("delete.load-balancer")
	defer teardown()

	ownerId := gidx.MustNewID(ownerPrefix)
	provider := entClient.Provider.Create().SetName("provider-name").SetOwnerID(ownerId).SaveX(ctx)
	lb := entClient.LoadBalancer.Create().SetName("lb-name").SetProvider(provider).SetOwnerID(ownerId).SetLocationID(gidx.MustNewID(locationPrefix)).SaveX(ctx)

	entClient.LoadBalancer.Use(manualhooks.LoadBalancerHooks()...)

	// Act
	entClient.LoadBalancer.DeleteOneID(lb.ID).ExecX(ctx)

	msg := testutils.ChannelReceiveWithTimeout[events.Message[events.ChangeMessage]](changesChannel, defualtTimeout)

	// Assert
	expectedAdditionalSubjectIDs := []gidx.PrefixedID{lb.OwnerID, lb.LocationID}
	actualAdditionalSubjectIDs := msg.Message().AdditionalSubjectIDs

	assert.ElementsMatch(t, expectedAdditionalSubjectIDs, actualAdditionalSubjectIDs)
}

func Test_OriginCreateUpdateHook(t *testing.T) {
	// Arrange
	ctx, entClient, changesChannel, teardown := setup("update.load-balancer-origin")
	defer teardown()

	ownerId := gidx.MustNewID(ownerPrefix)

	provider := entClient.Provider.Create().SetName("provider-name").SetOwnerID(ownerId).SaveX(ctx)
	lb := entClient.LoadBalancer.Create().SetName("lb-name").SetProvider(provider).SetOwnerID(ownerId).SetLocationID(gidx.MustNewID(locationPrefix)).SaveX(ctx)
	pool := entClient.Pool.Create().SetName("pool-name").SetOwnerID(ownerId).SetProtocol(pool.ProtocolTCP).SaveX(ctx)
	entClient.Port.Create().SetName("port-name").AddPoolIDs(pool.ID).SetLoadBalancer(lb).SetNumber(11).SaveX(ctx)
	origin := entClient.Origin.Create().SetName("origin-name").SetPool(pool).SetTarget("127.0.0.1").SetPortNumber(12).SaveX(ctx)

	entClient.Origin.Use(manualhooks.OriginHooks()...)

	// Act
	entClient.Origin.UpdateOne(origin).SetName("other-origin-name").ExecX(ctx)

	msg := testutils.ChannelReceiveWithTimeout[events.Message[events.ChangeMessage]](changesChannel, defualtTimeout)

	// Assert
	expectedAdditionalSubjectIDs := []gidx.PrefixedID{pool.ID, ownerId, lb.ID, lb.LocationID}
	actualAdditionalSubjectIDs := msg.Message().AdditionalSubjectIDs

	assert.ElementsMatch(t, expectedAdditionalSubjectIDs, actualAdditionalSubjectIDs)
}

func Test_OriginDeleteHook(t *testing.T) {
	// Arrange
	ctx, entClient, changesChannel, teardown := setup("delete.load-balancer-origin")
	defer teardown()

	ownerId := gidx.MustNewID(ownerPrefix)

	provider := entClient.Provider.Create().SetName("provider-name").SetOwnerID(ownerId).SaveX(ctx)
	lb := entClient.LoadBalancer.Create().SetName("lb-name").SetProvider(provider).SetOwnerID(ownerId).SetLocationID(gidx.MustNewID(locationPrefix)).SaveX(ctx)
	pool := entClient.Pool.Create().SetName("pool-name").SetOwnerID(ownerId).SetProtocol(pool.ProtocolTCP).SaveX(ctx)
	entClient.Port.Create().SetName("port-name").AddPoolIDs(pool.ID).SetLoadBalancer(lb).SetNumber(11).SaveX(ctx)
	origin := entClient.Origin.Create().SetName("origin-name").SetPool(pool).SetTarget("127.0.0.1").SetPortNumber(12).SaveX(ctx)

	entClient.Origin.Use(manualhooks.OriginHooks()...)

	// Act
	entClient.Origin.DeleteOne(origin).ExecX(ctx)

	msg := testutils.ChannelReceiveWithTimeout[events.Message[events.ChangeMessage]](changesChannel, defualtTimeout)

	// Assert
	expectedAdditionalSubjectIDs := []gidx.PrefixedID{pool.ID, ownerId, lb.ID, lb.LocationID}
	actualAdditionalSubjectIDs := msg.Message().AdditionalSubjectIDs

	assert.ElementsMatch(t, expectedAdditionalSubjectIDs, actualAdditionalSubjectIDs)
}

func Test_PoolCreateUpdateHook(t *testing.T) {
	// Arrange
	ctx, entClient, changesChannel, teardown := setup("update.load-balancer-pool")
	defer teardown()

	ownerId := gidx.MustNewID(ownerPrefix)

	provider := entClient.Provider.Create().SetName("provider-name").SetOwnerID(ownerId).SaveX(ctx)
	lb := entClient.LoadBalancer.Create().SetName("lb-name").SetProvider(provider).SetOwnerID(ownerId).SetLocationID(gidx.MustNewID(locationPrefix)).SaveX(ctx)
	pool := entClient.Pool.Create().SetName("pool-name").SetOwnerID(ownerId).SetProtocol(pool.ProtocolTCP).SaveX(ctx)
	port := entClient.Port.Create().SetName("port-name").AddPoolIDs(pool.ID).SetLoadBalancer(lb).SetNumber(11).SaveX(ctx)
	origin := entClient.Origin.Create().SetName("origin-name").SetPool(pool).SetTarget("127.0.0.1").SetPortNumber(12).SaveX(ctx)

	entClient.Pool.Use(manualhooks.PoolHooks()...)

	// Act
	entClient.Pool.UpdateOne(pool).SetName("other-pool-name").ExecX(ctx)

	msg := testutils.ChannelReceiveWithTimeout[events.Message[events.ChangeMessage]](changesChannel, defualtTimeout)

	// Assert
	expectedAdditionalSubjectIDs := []gidx.PrefixedID{pool.ID, ownerId, lb.ID, lb.LocationID, origin.ID, port.ID}
	actualAdditionalSubjectIDs := msg.Message().AdditionalSubjectIDs

	assert.ElementsMatch(t, expectedAdditionalSubjectIDs, actualAdditionalSubjectIDs)
}

func Test_PoolDeleteHook(t *testing.T) {
	// Arrange
	ctx, entClient, changesChannel, teardown := setup("delete.load-balancer-pool")
	defer teardown()

	ownerId := gidx.MustNewID(ownerPrefix)

	provider := entClient.Provider.Create().SetName("provider-name").SetOwnerID(ownerId).SaveX(ctx)
	lb := entClient.LoadBalancer.Create().SetName("lb-name").SetProvider(provider).SetOwnerID(ownerId).SetLocationID(gidx.MustNewID(locationPrefix)).SaveX(ctx)
	pool := entClient.Pool.Create().SetName("pool-name").SetOwnerID(ownerId).SetProtocol(pool.ProtocolTCP).SaveX(ctx)
	entClient.Port.Create().AddPoolIDs(pool.ID).SetLoadBalancer(lb).SetName("port-name").SetNumber(11).SaveX(ctx)

	entClient.Pool.Use(manualhooks.PoolHooks()...)

	// Act
	entClient.Pool.DeleteOne(pool).ExecX(ctx)

	msg := testutils.ChannelReceiveWithTimeout[events.Message[events.ChangeMessage]](changesChannel, defualtTimeout)

	// Assert
	expectedAdditionalSubjectIDs := []gidx.PrefixedID{ownerId, lb.ID, lb.LocationID}
	actualAdditionalSubjectIDs := msg.Message().AdditionalSubjectIDs

	assert.ElementsMatch(t, expectedAdditionalSubjectIDs, actualAdditionalSubjectIDs)
}

func Test_PortCreateUpdateHook(t *testing.T) {
	// Arrange
	ctx, entClient, changesChannel, teardown := setup("update.load-balancer-port")
	defer teardown()

	ownerId := gidx.MustNewID(ownerPrefix)

	provider := entClient.Provider.Create().SetName("provider-name").SetOwnerID(ownerId).SaveX(ctx)
	lb := entClient.LoadBalancer.Create().SetName("lb-name").SetProvider(provider).SetOwnerID(ownerId).SetLocationID(gidx.MustNewID(locationPrefix)).SaveX(ctx)
	pool := entClient.Pool.Create().SetName("pool-name").SetOwnerID(ownerId).SetProtocol(pool.ProtocolTCP).SaveX(ctx)
	port := entClient.Port.Create().AddPoolIDs(pool.ID).SetLoadBalancer(lb).SetName("port-name").SetNumber(11).SaveX(ctx)

	entClient.Port.Use(manualhooks.PortHooks()...)

	// Act
	entClient.Port.UpdateOne(port).SetName("other-port-name").ExecX(ctx)

	msg := testutils.ChannelReceiveWithTimeout[events.Message[events.ChangeMessage]](changesChannel, defualtTimeout)

	// Assert
	expectedAdditionalSubjectIDs := []gidx.PrefixedID{ownerId, lb.ID, lb.LocationID, provider.ID, pool.ID}
	actualAdditionalSubjectIDs := msg.Message().AdditionalSubjectIDs

	assert.ElementsMatch(t, expectedAdditionalSubjectIDs, actualAdditionalSubjectIDs)
}

func Test_PortDeleteHook(t *testing.T) {
	// Arrange
	ctx, entClient, changesChannel, teardown := setup("delete.load-balancer-port")
	defer teardown()

	ownerId := gidx.MustNewID(ownerPrefix)

	provider := entClient.Provider.Create().SetName("provider-name").SetOwnerID(ownerId).SaveX(ctx)
	lb := entClient.LoadBalancer.Create().SetName("lb-name").SetProvider(provider).SetOwnerID(ownerId).SetLocationID(gidx.MustNewID(locationPrefix)).SaveX(ctx)
	pool := entClient.Pool.Create().SetName("pool-name").SetOwnerID(ownerId).SetProtocol(pool.ProtocolTCP).SaveX(ctx)
	port := entClient.Port.Create().AddPoolIDs(pool.ID).SetLoadBalancer(lb).SetName("port-name").SetNumber(11).SaveX(ctx)

	entClient.Port.Use(manualhooks.PortHooks()...)

	// Act
	entClient.Port.DeleteOne(port).ExecX(ctx)

	msg := testutils.ChannelReceiveWithTimeout[events.Message[events.ChangeMessage]](changesChannel, defualtTimeout)

	// Assert
	expectedAdditionalSubjectIDs := []gidx.PrefixedID{ownerId, lb.ID, lb.LocationID, provider.ID}
	actualAdditionalSubjectIDs := msg.Message().AdditionalSubjectIDs

	assert.ElementsMatch(t, expectedAdditionalSubjectIDs, actualAdditionalSubjectIDs)
}
