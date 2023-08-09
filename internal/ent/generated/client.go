// Copyright 2023 The Infratographer Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Code generated by entc, DO NOT EDIT.

package generated

import (
	"context"
	"errors"
	"fmt"
	"log"

	"go.infratographer.com/load-balancer-api/internal/ent/generated/migrate"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"go.infratographer.com/load-balancer-api/internal/ent/generated/loadbalancer"
	"go.infratographer.com/load-balancer-api/internal/ent/generated/origin"
	"go.infratographer.com/load-balancer-api/internal/ent/generated/pool"
	"go.infratographer.com/load-balancer-api/internal/ent/generated/port"
	"go.infratographer.com/load-balancer-api/internal/ent/generated/provider"
	"go.infratographer.com/x/events"
	"go.infratographer.com/x/gidx"
)

// Client is the client that holds all ent builders.
type Client struct {
	config
	// Schema is the client for creating, migrating and dropping schema.
	Schema *migrate.Schema
	// LoadBalancer is the client for interacting with the LoadBalancer builders.
	LoadBalancer *LoadBalancerClient
	// Origin is the client for interacting with the Origin builders.
	Origin *OriginClient
	// Pool is the client for interacting with the Pool builders.
	Pool *PoolClient
	// Port is the client for interacting with the Port builders.
	Port *PortClient
	// Provider is the client for interacting with the Provider builders.
	Provider *ProviderClient
}

// NewClient creates a new client configured with the given options.
func NewClient(opts ...Option) *Client {
	cfg := config{log: log.Println, hooks: &hooks{}, inters: &inters{}}
	cfg.options(opts...)
	client := &Client{config: cfg}
	client.init()
	return client
}

func (c *Client) init() {
	c.Schema = migrate.NewSchema(c.driver)
	c.LoadBalancer = NewLoadBalancerClient(c.config)
	c.Origin = NewOriginClient(c.config)
	c.Pool = NewPoolClient(c.config)
	c.Port = NewPortClient(c.config)
	c.Provider = NewProviderClient(c.config)
}

type (
	// config is the configuration for the client and its builder.
	config struct {
		// driver used for executing database requests.
		driver dialect.Driver
		// debug enable a debug logging.
		debug bool
		// log used for logging on debug mode.
		log func(...any)
		// hooks to execute on mutations.
		hooks *hooks
		// interceptors to execute on queries.
		inters          *inters
		EventsPublisher events.Connection
	}
	// Option function to configure the client.
	Option func(*config)
)

// options applies the options on the config object.
func (c *config) options(opts ...Option) {
	for _, opt := range opts {
		opt(c)
	}
	if c.debug {
		c.driver = dialect.Debug(c.driver, c.log)
	}
}

// Debug enables debug logging on the ent.Driver.
func Debug() Option {
	return func(c *config) {
		c.debug = true
	}
}

// Log sets the logging function for debug mode.
func Log(fn func(...any)) Option {
	return func(c *config) {
		c.log = fn
	}
}

// Driver configures the client driver.
func Driver(driver dialect.Driver) Option {
	return func(c *config) {
		c.driver = driver
	}
}

// EventsPublisher configures the EventsPublisher.
func EventsPublisher(v events.Connection) Option {
	return func(c *config) {
		c.EventsPublisher = v
	}
}

// Open opens a database/sql.DB specified by the driver name and
// the data source name, and returns a new client attached to it.
// Optional parameters can be added for configuring the client.
func Open(driverName, dataSourceName string, options ...Option) (*Client, error) {
	switch driverName {
	case dialect.MySQL, dialect.Postgres, dialect.SQLite:
		drv, err := sql.Open(driverName, dataSourceName)
		if err != nil {
			return nil, err
		}
		return NewClient(append(options, Driver(drv))...), nil
	default:
		return nil, fmt.Errorf("unsupported driver: %q", driverName)
	}
}

// Tx returns a new transactional client. The provided context
// is used until the transaction is committed or rolled back.
func (c *Client) Tx(ctx context.Context) (*Tx, error) {
	if _, ok := c.driver.(*txDriver); ok {
		return nil, errors.New("generated: cannot start a transaction within a transaction")
	}
	tx, err := newTx(ctx, c.driver)
	if err != nil {
		return nil, fmt.Errorf("generated: starting a transaction: %w", err)
	}
	cfg := c.config
	cfg.driver = tx
	return &Tx{
		ctx:          ctx,
		config:       cfg,
		LoadBalancer: NewLoadBalancerClient(cfg),
		Origin:       NewOriginClient(cfg),
		Pool:         NewPoolClient(cfg),
		Port:         NewPortClient(cfg),
		Provider:     NewProviderClient(cfg),
	}, nil
}

// BeginTx returns a transactional client with specified options.
func (c *Client) BeginTx(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	if _, ok := c.driver.(*txDriver); ok {
		return nil, errors.New("ent: cannot start a transaction within a transaction")
	}
	tx, err := c.driver.(interface {
		BeginTx(context.Context, *sql.TxOptions) (dialect.Tx, error)
	}).BeginTx(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("ent: starting a transaction: %w", err)
	}
	cfg := c.config
	cfg.driver = &txDriver{tx: tx, drv: c.driver}
	return &Tx{
		ctx:          ctx,
		config:       cfg,
		LoadBalancer: NewLoadBalancerClient(cfg),
		Origin:       NewOriginClient(cfg),
		Pool:         NewPoolClient(cfg),
		Port:         NewPortClient(cfg),
		Provider:     NewProviderClient(cfg),
	}, nil
}

// Debug returns a new debug-client. It's used to get verbose logging on specific operations.
//
//	client.Debug().
//		LoadBalancer.
//		Query().
//		Count(ctx)
func (c *Client) Debug() *Client {
	if c.debug {
		return c
	}
	cfg := c.config
	cfg.driver = dialect.Debug(c.driver, c.log)
	client := &Client{config: cfg}
	client.init()
	return client
}

// Close closes the database connection and prevents new queries from starting.
func (c *Client) Close() error {
	return c.driver.Close()
}

// Use adds the mutation hooks to all the entity clients.
// In order to add hooks to a specific client, call: `client.Node.Use(...)`.
func (c *Client) Use(hooks ...Hook) {
	c.LoadBalancer.Use(hooks...)
	c.Origin.Use(hooks...)
	c.Pool.Use(hooks...)
	c.Port.Use(hooks...)
	c.Provider.Use(hooks...)
}

// Intercept adds the query interceptors to all the entity clients.
// In order to add interceptors to a specific client, call: `client.Node.Intercept(...)`.
func (c *Client) Intercept(interceptors ...Interceptor) {
	c.LoadBalancer.Intercept(interceptors...)
	c.Origin.Intercept(interceptors...)
	c.Pool.Intercept(interceptors...)
	c.Port.Intercept(interceptors...)
	c.Provider.Intercept(interceptors...)
}

// Mutate implements the ent.Mutator interface.
func (c *Client) Mutate(ctx context.Context, m Mutation) (Value, error) {
	switch m := m.(type) {
	case *LoadBalancerMutation:
		return c.LoadBalancer.mutate(ctx, m)
	case *OriginMutation:
		return c.Origin.mutate(ctx, m)
	case *PoolMutation:
		return c.Pool.mutate(ctx, m)
	case *PortMutation:
		return c.Port.mutate(ctx, m)
	case *ProviderMutation:
		return c.Provider.mutate(ctx, m)
	default:
		return nil, fmt.Errorf("generated: unknown mutation type %T", m)
	}
}

// LoadBalancerClient is a client for the LoadBalancer schema.
type LoadBalancerClient struct {
	config
}

// NewLoadBalancerClient returns a client for the LoadBalancer from the given config.
func NewLoadBalancerClient(c config) *LoadBalancerClient {
	return &LoadBalancerClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `loadbalancer.Hooks(f(g(h())))`.
func (c *LoadBalancerClient) Use(hooks ...Hook) {
	c.hooks.LoadBalancer = append(c.hooks.LoadBalancer, hooks...)
}

// Intercept adds a list of query interceptors to the interceptors stack.
// A call to `Intercept(f, g, h)` equals to `loadbalancer.Intercept(f(g(h())))`.
func (c *LoadBalancerClient) Intercept(interceptors ...Interceptor) {
	c.inters.LoadBalancer = append(c.inters.LoadBalancer, interceptors...)
}

// Create returns a builder for creating a LoadBalancer entity.
func (c *LoadBalancerClient) Create() *LoadBalancerCreate {
	mutation := newLoadBalancerMutation(c.config, OpCreate)
	return &LoadBalancerCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of LoadBalancer entities.
func (c *LoadBalancerClient) CreateBulk(builders ...*LoadBalancerCreate) *LoadBalancerCreateBulk {
	return &LoadBalancerCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for LoadBalancer.
func (c *LoadBalancerClient) Update() *LoadBalancerUpdate {
	mutation := newLoadBalancerMutation(c.config, OpUpdate)
	return &LoadBalancerUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *LoadBalancerClient) UpdateOne(lb *LoadBalancer) *LoadBalancerUpdateOne {
	mutation := newLoadBalancerMutation(c.config, OpUpdateOne, withLoadBalancer(lb))
	return &LoadBalancerUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *LoadBalancerClient) UpdateOneID(id gidx.PrefixedID) *LoadBalancerUpdateOne {
	mutation := newLoadBalancerMutation(c.config, OpUpdateOne, withLoadBalancerID(id))
	return &LoadBalancerUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for LoadBalancer.
func (c *LoadBalancerClient) Delete() *LoadBalancerDelete {
	mutation := newLoadBalancerMutation(c.config, OpDelete)
	return &LoadBalancerDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a builder for deleting the given entity.
func (c *LoadBalancerClient) DeleteOne(lb *LoadBalancer) *LoadBalancerDeleteOne {
	return c.DeleteOneID(lb.ID)
}

// DeleteOneID returns a builder for deleting the given entity by its id.
func (c *LoadBalancerClient) DeleteOneID(id gidx.PrefixedID) *LoadBalancerDeleteOne {
	builder := c.Delete().Where(loadbalancer.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &LoadBalancerDeleteOne{builder}
}

// Query returns a query builder for LoadBalancer.
func (c *LoadBalancerClient) Query() *LoadBalancerQuery {
	return &LoadBalancerQuery{
		config: c.config,
		ctx:    &QueryContext{Type: TypeLoadBalancer},
		inters: c.Interceptors(),
	}
}

// Get returns a LoadBalancer entity by its id.
func (c *LoadBalancerClient) Get(ctx context.Context, id gidx.PrefixedID) (*LoadBalancer, error) {
	return c.Query().Where(loadbalancer.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *LoadBalancerClient) GetX(ctx context.Context, id gidx.PrefixedID) *LoadBalancer {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// QueryPorts queries the ports edge of a LoadBalancer.
func (c *LoadBalancerClient) QueryPorts(lb *LoadBalancer) *PortQuery {
	query := (&PortClient{config: c.config}).Query()
	query.path = func(context.Context) (fromV *sql.Selector, _ error) {
		id := lb.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(loadbalancer.Table, loadbalancer.FieldID, id),
			sqlgraph.To(port.Table, port.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, true, loadbalancer.PortsTable, loadbalancer.PortsColumn),
		)
		fromV = sqlgraph.Neighbors(lb.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryProvider queries the provider edge of a LoadBalancer.
func (c *LoadBalancerClient) QueryProvider(lb *LoadBalancer) *ProviderQuery {
	query := (&ProviderClient{config: c.config}).Query()
	query.path = func(context.Context) (fromV *sql.Selector, _ error) {
		id := lb.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(loadbalancer.Table, loadbalancer.FieldID, id),
			sqlgraph.To(provider.Table, provider.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, loadbalancer.ProviderTable, loadbalancer.ProviderColumn),
		)
		fromV = sqlgraph.Neighbors(lb.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *LoadBalancerClient) Hooks() []Hook {
	return c.hooks.LoadBalancer
}

// Interceptors returns the client interceptors.
func (c *LoadBalancerClient) Interceptors() []Interceptor {
	return c.inters.LoadBalancer
}

func (c *LoadBalancerClient) mutate(ctx context.Context, m *LoadBalancerMutation) (Value, error) {
	switch m.Op() {
	case OpCreate:
		return (&LoadBalancerCreate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdate:
		return (&LoadBalancerUpdate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdateOne:
		return (&LoadBalancerUpdateOne{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpDelete, OpDeleteOne:
		return (&LoadBalancerDelete{config: c.config, hooks: c.Hooks(), mutation: m}).Exec(ctx)
	default:
		return nil, fmt.Errorf("generated: unknown LoadBalancer mutation op: %q", m.Op())
	}
}

// OriginClient is a client for the Origin schema.
type OriginClient struct {
	config
}

// NewOriginClient returns a client for the Origin from the given config.
func NewOriginClient(c config) *OriginClient {
	return &OriginClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `origin.Hooks(f(g(h())))`.
func (c *OriginClient) Use(hooks ...Hook) {
	c.hooks.Origin = append(c.hooks.Origin, hooks...)
}

// Intercept adds a list of query interceptors to the interceptors stack.
// A call to `Intercept(f, g, h)` equals to `origin.Intercept(f(g(h())))`.
func (c *OriginClient) Intercept(interceptors ...Interceptor) {
	c.inters.Origin = append(c.inters.Origin, interceptors...)
}

// Create returns a builder for creating a Origin entity.
func (c *OriginClient) Create() *OriginCreate {
	mutation := newOriginMutation(c.config, OpCreate)
	return &OriginCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of Origin entities.
func (c *OriginClient) CreateBulk(builders ...*OriginCreate) *OriginCreateBulk {
	return &OriginCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for Origin.
func (c *OriginClient) Update() *OriginUpdate {
	mutation := newOriginMutation(c.config, OpUpdate)
	return &OriginUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *OriginClient) UpdateOne(o *Origin) *OriginUpdateOne {
	mutation := newOriginMutation(c.config, OpUpdateOne, withOrigin(o))
	return &OriginUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *OriginClient) UpdateOneID(id gidx.PrefixedID) *OriginUpdateOne {
	mutation := newOriginMutation(c.config, OpUpdateOne, withOriginID(id))
	return &OriginUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for Origin.
func (c *OriginClient) Delete() *OriginDelete {
	mutation := newOriginMutation(c.config, OpDelete)
	return &OriginDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a builder for deleting the given entity.
func (c *OriginClient) DeleteOne(o *Origin) *OriginDeleteOne {
	return c.DeleteOneID(o.ID)
}

// DeleteOneID returns a builder for deleting the given entity by its id.
func (c *OriginClient) DeleteOneID(id gidx.PrefixedID) *OriginDeleteOne {
	builder := c.Delete().Where(origin.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &OriginDeleteOne{builder}
}

// Query returns a query builder for Origin.
func (c *OriginClient) Query() *OriginQuery {
	return &OriginQuery{
		config: c.config,
		ctx:    &QueryContext{Type: TypeOrigin},
		inters: c.Interceptors(),
	}
}

// Get returns a Origin entity by its id.
func (c *OriginClient) Get(ctx context.Context, id gidx.PrefixedID) (*Origin, error) {
	return c.Query().Where(origin.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *OriginClient) GetX(ctx context.Context, id gidx.PrefixedID) *Origin {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// QueryPool queries the pool edge of a Origin.
func (c *OriginClient) QueryPool(o *Origin) *PoolQuery {
	query := (&PoolClient{config: c.config}).Query()
	query.path = func(context.Context) (fromV *sql.Selector, _ error) {
		id := o.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(origin.Table, origin.FieldID, id),
			sqlgraph.To(pool.Table, pool.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, origin.PoolTable, origin.PoolColumn),
		)
		fromV = sqlgraph.Neighbors(o.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *OriginClient) Hooks() []Hook {
	return c.hooks.Origin
}

// Interceptors returns the client interceptors.
func (c *OriginClient) Interceptors() []Interceptor {
	return c.inters.Origin
}

func (c *OriginClient) mutate(ctx context.Context, m *OriginMutation) (Value, error) {
	switch m.Op() {
	case OpCreate:
		return (&OriginCreate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdate:
		return (&OriginUpdate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdateOne:
		return (&OriginUpdateOne{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpDelete, OpDeleteOne:
		return (&OriginDelete{config: c.config, hooks: c.Hooks(), mutation: m}).Exec(ctx)
	default:
		return nil, fmt.Errorf("generated: unknown Origin mutation op: %q", m.Op())
	}
}

// PoolClient is a client for the Pool schema.
type PoolClient struct {
	config
}

// NewPoolClient returns a client for the Pool from the given config.
func NewPoolClient(c config) *PoolClient {
	return &PoolClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `pool.Hooks(f(g(h())))`.
func (c *PoolClient) Use(hooks ...Hook) {
	c.hooks.Pool = append(c.hooks.Pool, hooks...)
}

// Intercept adds a list of query interceptors to the interceptors stack.
// A call to `Intercept(f, g, h)` equals to `pool.Intercept(f(g(h())))`.
func (c *PoolClient) Intercept(interceptors ...Interceptor) {
	c.inters.Pool = append(c.inters.Pool, interceptors...)
}

// Create returns a builder for creating a Pool entity.
func (c *PoolClient) Create() *PoolCreate {
	mutation := newPoolMutation(c.config, OpCreate)
	return &PoolCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of Pool entities.
func (c *PoolClient) CreateBulk(builders ...*PoolCreate) *PoolCreateBulk {
	return &PoolCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for Pool.
func (c *PoolClient) Update() *PoolUpdate {
	mutation := newPoolMutation(c.config, OpUpdate)
	return &PoolUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *PoolClient) UpdateOne(po *Pool) *PoolUpdateOne {
	mutation := newPoolMutation(c.config, OpUpdateOne, withPool(po))
	return &PoolUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *PoolClient) UpdateOneID(id gidx.PrefixedID) *PoolUpdateOne {
	mutation := newPoolMutation(c.config, OpUpdateOne, withPoolID(id))
	return &PoolUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for Pool.
func (c *PoolClient) Delete() *PoolDelete {
	mutation := newPoolMutation(c.config, OpDelete)
	return &PoolDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a builder for deleting the given entity.
func (c *PoolClient) DeleteOne(po *Pool) *PoolDeleteOne {
	return c.DeleteOneID(po.ID)
}

// DeleteOneID returns a builder for deleting the given entity by its id.
func (c *PoolClient) DeleteOneID(id gidx.PrefixedID) *PoolDeleteOne {
	builder := c.Delete().Where(pool.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &PoolDeleteOne{builder}
}

// Query returns a query builder for Pool.
func (c *PoolClient) Query() *PoolQuery {
	return &PoolQuery{
		config: c.config,
		ctx:    &QueryContext{Type: TypePool},
		inters: c.Interceptors(),
	}
}

// Get returns a Pool entity by its id.
func (c *PoolClient) Get(ctx context.Context, id gidx.PrefixedID) (*Pool, error) {
	return c.Query().Where(pool.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *PoolClient) GetX(ctx context.Context, id gidx.PrefixedID) *Pool {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// QueryPorts queries the ports edge of a Pool.
func (c *PoolClient) QueryPorts(po *Pool) *PortQuery {
	query := (&PortClient{config: c.config}).Query()
	query.path = func(context.Context) (fromV *sql.Selector, _ error) {
		id := po.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(pool.Table, pool.FieldID, id),
			sqlgraph.To(port.Table, port.FieldID),
			sqlgraph.Edge(sqlgraph.M2M, false, pool.PortsTable, pool.PortsPrimaryKey...),
		)
		fromV = sqlgraph.Neighbors(po.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryOrigins queries the origins edge of a Pool.
func (c *PoolClient) QueryOrigins(po *Pool) *OriginQuery {
	query := (&OriginClient{config: c.config}).Query()
	query.path = func(context.Context) (fromV *sql.Selector, _ error) {
		id := po.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(pool.Table, pool.FieldID, id),
			sqlgraph.To(origin.Table, origin.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, true, pool.OriginsTable, pool.OriginsColumn),
		)
		fromV = sqlgraph.Neighbors(po.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *PoolClient) Hooks() []Hook {
	return c.hooks.Pool
}

// Interceptors returns the client interceptors.
func (c *PoolClient) Interceptors() []Interceptor {
	return c.inters.Pool
}

func (c *PoolClient) mutate(ctx context.Context, m *PoolMutation) (Value, error) {
	switch m.Op() {
	case OpCreate:
		return (&PoolCreate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdate:
		return (&PoolUpdate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdateOne:
		return (&PoolUpdateOne{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpDelete, OpDeleteOne:
		return (&PoolDelete{config: c.config, hooks: c.Hooks(), mutation: m}).Exec(ctx)
	default:
		return nil, fmt.Errorf("generated: unknown Pool mutation op: %q", m.Op())
	}
}

// PortClient is a client for the Port schema.
type PortClient struct {
	config
}

// NewPortClient returns a client for the Port from the given config.
func NewPortClient(c config) *PortClient {
	return &PortClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `port.Hooks(f(g(h())))`.
func (c *PortClient) Use(hooks ...Hook) {
	c.hooks.Port = append(c.hooks.Port, hooks...)
}

// Intercept adds a list of query interceptors to the interceptors stack.
// A call to `Intercept(f, g, h)` equals to `port.Intercept(f(g(h())))`.
func (c *PortClient) Intercept(interceptors ...Interceptor) {
	c.inters.Port = append(c.inters.Port, interceptors...)
}

// Create returns a builder for creating a Port entity.
func (c *PortClient) Create() *PortCreate {
	mutation := newPortMutation(c.config, OpCreate)
	return &PortCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of Port entities.
func (c *PortClient) CreateBulk(builders ...*PortCreate) *PortCreateBulk {
	return &PortCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for Port.
func (c *PortClient) Update() *PortUpdate {
	mutation := newPortMutation(c.config, OpUpdate)
	return &PortUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *PortClient) UpdateOne(po *Port) *PortUpdateOne {
	mutation := newPortMutation(c.config, OpUpdateOne, withPort(po))
	return &PortUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *PortClient) UpdateOneID(id gidx.PrefixedID) *PortUpdateOne {
	mutation := newPortMutation(c.config, OpUpdateOne, withPortID(id))
	return &PortUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for Port.
func (c *PortClient) Delete() *PortDelete {
	mutation := newPortMutation(c.config, OpDelete)
	return &PortDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a builder for deleting the given entity.
func (c *PortClient) DeleteOne(po *Port) *PortDeleteOne {
	return c.DeleteOneID(po.ID)
}

// DeleteOneID returns a builder for deleting the given entity by its id.
func (c *PortClient) DeleteOneID(id gidx.PrefixedID) *PortDeleteOne {
	builder := c.Delete().Where(port.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &PortDeleteOne{builder}
}

// Query returns a query builder for Port.
func (c *PortClient) Query() *PortQuery {
	return &PortQuery{
		config: c.config,
		ctx:    &QueryContext{Type: TypePort},
		inters: c.Interceptors(),
	}
}

// Get returns a Port entity by its id.
func (c *PortClient) Get(ctx context.Context, id gidx.PrefixedID) (*Port, error) {
	return c.Query().Where(port.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *PortClient) GetX(ctx context.Context, id gidx.PrefixedID) *Port {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// QueryPools queries the pools edge of a Port.
func (c *PortClient) QueryPools(po *Port) *PoolQuery {
	query := (&PoolClient{config: c.config}).Query()
	query.path = func(context.Context) (fromV *sql.Selector, _ error) {
		id := po.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(port.Table, port.FieldID, id),
			sqlgraph.To(pool.Table, pool.FieldID),
			sqlgraph.Edge(sqlgraph.M2M, true, port.PoolsTable, port.PoolsPrimaryKey...),
		)
		fromV = sqlgraph.Neighbors(po.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryLoadBalancer queries the load_balancer edge of a Port.
func (c *PortClient) QueryLoadBalancer(po *Port) *LoadBalancerQuery {
	query := (&LoadBalancerClient{config: c.config}).Query()
	query.path = func(context.Context) (fromV *sql.Selector, _ error) {
		id := po.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(port.Table, port.FieldID, id),
			sqlgraph.To(loadbalancer.Table, loadbalancer.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, port.LoadBalancerTable, port.LoadBalancerColumn),
		)
		fromV = sqlgraph.Neighbors(po.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *PortClient) Hooks() []Hook {
	return c.hooks.Port
}

// Interceptors returns the client interceptors.
func (c *PortClient) Interceptors() []Interceptor {
	return c.inters.Port
}

func (c *PortClient) mutate(ctx context.Context, m *PortMutation) (Value, error) {
	switch m.Op() {
	case OpCreate:
		return (&PortCreate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdate:
		return (&PortUpdate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdateOne:
		return (&PortUpdateOne{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpDelete, OpDeleteOne:
		return (&PortDelete{config: c.config, hooks: c.Hooks(), mutation: m}).Exec(ctx)
	default:
		return nil, fmt.Errorf("generated: unknown Port mutation op: %q", m.Op())
	}
}

// ProviderClient is a client for the Provider schema.
type ProviderClient struct {
	config
}

// NewProviderClient returns a client for the Provider from the given config.
func NewProviderClient(c config) *ProviderClient {
	return &ProviderClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `provider.Hooks(f(g(h())))`.
func (c *ProviderClient) Use(hooks ...Hook) {
	c.hooks.Provider = append(c.hooks.Provider, hooks...)
}

// Intercept adds a list of query interceptors to the interceptors stack.
// A call to `Intercept(f, g, h)` equals to `provider.Intercept(f(g(h())))`.
func (c *ProviderClient) Intercept(interceptors ...Interceptor) {
	c.inters.Provider = append(c.inters.Provider, interceptors...)
}

// Create returns a builder for creating a Provider entity.
func (c *ProviderClient) Create() *ProviderCreate {
	mutation := newProviderMutation(c.config, OpCreate)
	return &ProviderCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of Provider entities.
func (c *ProviderClient) CreateBulk(builders ...*ProviderCreate) *ProviderCreateBulk {
	return &ProviderCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for Provider.
func (c *ProviderClient) Update() *ProviderUpdate {
	mutation := newProviderMutation(c.config, OpUpdate)
	return &ProviderUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *ProviderClient) UpdateOne(pr *Provider) *ProviderUpdateOne {
	mutation := newProviderMutation(c.config, OpUpdateOne, withProvider(pr))
	return &ProviderUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *ProviderClient) UpdateOneID(id gidx.PrefixedID) *ProviderUpdateOne {
	mutation := newProviderMutation(c.config, OpUpdateOne, withProviderID(id))
	return &ProviderUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for Provider.
func (c *ProviderClient) Delete() *ProviderDelete {
	mutation := newProviderMutation(c.config, OpDelete)
	return &ProviderDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a builder for deleting the given entity.
func (c *ProviderClient) DeleteOne(pr *Provider) *ProviderDeleteOne {
	return c.DeleteOneID(pr.ID)
}

// DeleteOneID returns a builder for deleting the given entity by its id.
func (c *ProviderClient) DeleteOneID(id gidx.PrefixedID) *ProviderDeleteOne {
	builder := c.Delete().Where(provider.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &ProviderDeleteOne{builder}
}

// Query returns a query builder for Provider.
func (c *ProviderClient) Query() *ProviderQuery {
	return &ProviderQuery{
		config: c.config,
		ctx:    &QueryContext{Type: TypeProvider},
		inters: c.Interceptors(),
	}
}

// Get returns a Provider entity by its id.
func (c *ProviderClient) Get(ctx context.Context, id gidx.PrefixedID) (*Provider, error) {
	return c.Query().Where(provider.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *ProviderClient) GetX(ctx context.Context, id gidx.PrefixedID) *Provider {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// QueryLoadBalancers queries the load_balancers edge of a Provider.
func (c *ProviderClient) QueryLoadBalancers(pr *Provider) *LoadBalancerQuery {
	query := (&LoadBalancerClient{config: c.config}).Query()
	query.path = func(context.Context) (fromV *sql.Selector, _ error) {
		id := pr.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(provider.Table, provider.FieldID, id),
			sqlgraph.To(loadbalancer.Table, loadbalancer.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, true, provider.LoadBalancersTable, provider.LoadBalancersColumn),
		)
		fromV = sqlgraph.Neighbors(pr.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *ProviderClient) Hooks() []Hook {
	return c.hooks.Provider
}

// Interceptors returns the client interceptors.
func (c *ProviderClient) Interceptors() []Interceptor {
	return c.inters.Provider
}

func (c *ProviderClient) mutate(ctx context.Context, m *ProviderMutation) (Value, error) {
	switch m.Op() {
	case OpCreate:
		return (&ProviderCreate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdate:
		return (&ProviderUpdate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdateOne:
		return (&ProviderUpdateOne{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpDelete, OpDeleteOne:
		return (&ProviderDelete{config: c.config, hooks: c.Hooks(), mutation: m}).Exec(ctx)
	default:
		return nil, fmt.Errorf("generated: unknown Provider mutation op: %q", m.Op())
	}
}

// hooks and interceptors per client, for fast access.
type (
	hooks struct {
		LoadBalancer, Origin, Pool, Port, Provider []ent.Hook
	}
	inters struct {
		LoadBalancer, Origin, Pool, Port, Provider []ent.Interceptor
	}
)
