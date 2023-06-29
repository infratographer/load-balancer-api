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
	"fmt"
	"math"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"go.infratographer.com/load-balancer-api/internal/ent/generated/origin"
	"go.infratographer.com/load-balancer-api/internal/ent/generated/pool"
	"go.infratographer.com/load-balancer-api/internal/ent/generated/predicate"
	"go.infratographer.com/x/gidx"
)

// OriginQuery is the builder for querying Origin entities.
type OriginQuery struct {
	config
	ctx        *QueryContext
	order      []origin.OrderOption
	inters     []Interceptor
	predicates []predicate.Origin
	withPool   *PoolQuery
	modifiers  []func(*sql.Selector)
	loadTotal  []func(context.Context, []*Origin) error
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the OriginQuery builder.
func (oq *OriginQuery) Where(ps ...predicate.Origin) *OriginQuery {
	oq.predicates = append(oq.predicates, ps...)
	return oq
}

// Limit the number of records to be returned by this query.
func (oq *OriginQuery) Limit(limit int) *OriginQuery {
	oq.ctx.Limit = &limit
	return oq
}

// Offset to start from.
func (oq *OriginQuery) Offset(offset int) *OriginQuery {
	oq.ctx.Offset = &offset
	return oq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (oq *OriginQuery) Unique(unique bool) *OriginQuery {
	oq.ctx.Unique = &unique
	return oq
}

// Order specifies how the records should be ordered.
func (oq *OriginQuery) Order(o ...origin.OrderOption) *OriginQuery {
	oq.order = append(oq.order, o...)
	return oq
}

// QueryPool chains the current query on the "pool" edge.
func (oq *OriginQuery) QueryPool() *PoolQuery {
	query := (&PoolClient{config: oq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := oq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := oq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(origin.Table, origin.FieldID, selector),
			sqlgraph.To(pool.Table, pool.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, origin.PoolTable, origin.PoolColumn),
		)
		fromU = sqlgraph.SetNeighbors(oq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first Origin entity from the query.
// Returns a *NotFoundError when no Origin was found.
func (oq *OriginQuery) First(ctx context.Context) (*Origin, error) {
	nodes, err := oq.Limit(1).All(setContextOp(ctx, oq.ctx, "First"))
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{origin.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (oq *OriginQuery) FirstX(ctx context.Context) *Origin {
	node, err := oq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first Origin ID from the query.
// Returns a *NotFoundError when no Origin ID was found.
func (oq *OriginQuery) FirstID(ctx context.Context) (id gidx.PrefixedID, err error) {
	var ids []gidx.PrefixedID
	if ids, err = oq.Limit(1).IDs(setContextOp(ctx, oq.ctx, "FirstID")); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{origin.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (oq *OriginQuery) FirstIDX(ctx context.Context) gidx.PrefixedID {
	id, err := oq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single Origin entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one Origin entity is found.
// Returns a *NotFoundError when no Origin entities are found.
func (oq *OriginQuery) Only(ctx context.Context) (*Origin, error) {
	nodes, err := oq.Limit(2).All(setContextOp(ctx, oq.ctx, "Only"))
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{origin.Label}
	default:
		return nil, &NotSingularError{origin.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (oq *OriginQuery) OnlyX(ctx context.Context) *Origin {
	node, err := oq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only Origin ID in the query.
// Returns a *NotSingularError when more than one Origin ID is found.
// Returns a *NotFoundError when no entities are found.
func (oq *OriginQuery) OnlyID(ctx context.Context) (id gidx.PrefixedID, err error) {
	var ids []gidx.PrefixedID
	if ids, err = oq.Limit(2).IDs(setContextOp(ctx, oq.ctx, "OnlyID")); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{origin.Label}
	default:
		err = &NotSingularError{origin.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (oq *OriginQuery) OnlyIDX(ctx context.Context) gidx.PrefixedID {
	id, err := oq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of Origins.
func (oq *OriginQuery) All(ctx context.Context) ([]*Origin, error) {
	ctx = setContextOp(ctx, oq.ctx, "All")
	if err := oq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	qr := querierAll[[]*Origin, *OriginQuery]()
	return withInterceptors[[]*Origin](ctx, oq, qr, oq.inters)
}

// AllX is like All, but panics if an error occurs.
func (oq *OriginQuery) AllX(ctx context.Context) []*Origin {
	nodes, err := oq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of Origin IDs.
func (oq *OriginQuery) IDs(ctx context.Context) (ids []gidx.PrefixedID, err error) {
	if oq.ctx.Unique == nil && oq.path != nil {
		oq.Unique(true)
	}
	ctx = setContextOp(ctx, oq.ctx, "IDs")
	if err = oq.Select(origin.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (oq *OriginQuery) IDsX(ctx context.Context) []gidx.PrefixedID {
	ids, err := oq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (oq *OriginQuery) Count(ctx context.Context) (int, error) {
	ctx = setContextOp(ctx, oq.ctx, "Count")
	if err := oq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return withInterceptors[int](ctx, oq, querierCount[*OriginQuery](), oq.inters)
}

// CountX is like Count, but panics if an error occurs.
func (oq *OriginQuery) CountX(ctx context.Context) int {
	count, err := oq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (oq *OriginQuery) Exist(ctx context.Context) (bool, error) {
	ctx = setContextOp(ctx, oq.ctx, "Exist")
	switch _, err := oq.FirstID(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("generated: check existence: %w", err)
	default:
		return true, nil
	}
}

// ExistX is like Exist, but panics if an error occurs.
func (oq *OriginQuery) ExistX(ctx context.Context) bool {
	exist, err := oq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the OriginQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (oq *OriginQuery) Clone() *OriginQuery {
	if oq == nil {
		return nil
	}
	return &OriginQuery{
		config:     oq.config,
		ctx:        oq.ctx.Clone(),
		order:      append([]origin.OrderOption{}, oq.order...),
		inters:     append([]Interceptor{}, oq.inters...),
		predicates: append([]predicate.Origin{}, oq.predicates...),
		withPool:   oq.withPool.Clone(),
		// clone intermediate query.
		sql:  oq.sql.Clone(),
		path: oq.path,
	}
}

// WithPool tells the query-builder to eager-load the nodes that are connected to
// the "pool" edge. The optional arguments are used to configure the query builder of the edge.
func (oq *OriginQuery) WithPool(opts ...func(*PoolQuery)) *OriginQuery {
	query := (&PoolClient{config: oq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	oq.withPool = query
	return oq
}

// GroupBy is used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
//
// Example:
//
//	var v []struct {
//		CreatedAt time.Time `json:"created_at,omitempty"`
//		Count int `json:"count,omitempty"`
//	}
//
//	client.Origin.Query().
//		GroupBy(origin.FieldCreatedAt).
//		Aggregate(generated.Count()).
//		Scan(ctx, &v)
func (oq *OriginQuery) GroupBy(field string, fields ...string) *OriginGroupBy {
	oq.ctx.Fields = append([]string{field}, fields...)
	grbuild := &OriginGroupBy{build: oq}
	grbuild.flds = &oq.ctx.Fields
	grbuild.label = origin.Label
	grbuild.scan = grbuild.Scan
	return grbuild
}

// Select allows the selection one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
//
// Example:
//
//	var v []struct {
//		CreatedAt time.Time `json:"created_at,omitempty"`
//	}
//
//	client.Origin.Query().
//		Select(origin.FieldCreatedAt).
//		Scan(ctx, &v)
func (oq *OriginQuery) Select(fields ...string) *OriginSelect {
	oq.ctx.Fields = append(oq.ctx.Fields, fields...)
	sbuild := &OriginSelect{OriginQuery: oq}
	sbuild.label = origin.Label
	sbuild.flds, sbuild.scan = &oq.ctx.Fields, sbuild.Scan
	return sbuild
}

// Aggregate returns a OriginSelect configured with the given aggregations.
func (oq *OriginQuery) Aggregate(fns ...AggregateFunc) *OriginSelect {
	return oq.Select().Aggregate(fns...)
}

func (oq *OriginQuery) prepareQuery(ctx context.Context) error {
	for _, inter := range oq.inters {
		if inter == nil {
			return fmt.Errorf("generated: uninitialized interceptor (forgotten import generated/runtime?)")
		}
		if trv, ok := inter.(Traverser); ok {
			if err := trv.Traverse(ctx, oq); err != nil {
				return err
			}
		}
	}
	for _, f := range oq.ctx.Fields {
		if !origin.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("generated: invalid field %q for query", f)}
		}
	}
	if oq.path != nil {
		prev, err := oq.path(ctx)
		if err != nil {
			return err
		}
		oq.sql = prev
	}
	return nil
}

func (oq *OriginQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*Origin, error) {
	var (
		nodes       = []*Origin{}
		_spec       = oq.querySpec()
		loadedTypes = [1]bool{
			oq.withPool != nil,
		}
	)
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*Origin).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &Origin{config: oq.config}
		nodes = append(nodes, node)
		node.Edges.loadedTypes = loadedTypes
		return node.assignValues(columns, values)
	}
	if len(oq.modifiers) > 0 {
		_spec.Modifiers = oq.modifiers
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, oq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	if query := oq.withPool; query != nil {
		if err := oq.loadPool(ctx, query, nodes, nil,
			func(n *Origin, e *Pool) { n.Edges.Pool = e }); err != nil {
			return nil, err
		}
	}
	for i := range oq.loadTotal {
		if err := oq.loadTotal[i](ctx, nodes); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

func (oq *OriginQuery) loadPool(ctx context.Context, query *PoolQuery, nodes []*Origin, init func(*Origin), assign func(*Origin, *Pool)) error {
	ids := make([]gidx.PrefixedID, 0, len(nodes))
	nodeids := make(map[gidx.PrefixedID][]*Origin)
	for i := range nodes {
		fk := nodes[i].PoolID
		if _, ok := nodeids[fk]; !ok {
			ids = append(ids, fk)
		}
		nodeids[fk] = append(nodeids[fk], nodes[i])
	}
	if len(ids) == 0 {
		return nil
	}
	query.Where(pool.IDIn(ids...))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		nodes, ok := nodeids[n.ID]
		if !ok {
			return fmt.Errorf(`unexpected foreign-key "pool_id" returned %v`, n.ID)
		}
		for i := range nodes {
			assign(nodes[i], n)
		}
	}
	return nil
}

func (oq *OriginQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := oq.querySpec()
	if len(oq.modifiers) > 0 {
		_spec.Modifiers = oq.modifiers
	}
	_spec.Node.Columns = oq.ctx.Fields
	if len(oq.ctx.Fields) > 0 {
		_spec.Unique = oq.ctx.Unique != nil && *oq.ctx.Unique
	}
	return sqlgraph.CountNodes(ctx, oq.driver, _spec)
}

func (oq *OriginQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := sqlgraph.NewQuerySpec(origin.Table, origin.Columns, sqlgraph.NewFieldSpec(origin.FieldID, field.TypeString))
	_spec.From = oq.sql
	if unique := oq.ctx.Unique; unique != nil {
		_spec.Unique = *unique
	} else if oq.path != nil {
		_spec.Unique = true
	}
	if fields := oq.ctx.Fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, origin.FieldID)
		for i := range fields {
			if fields[i] != origin.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
		if oq.withPool != nil {
			_spec.Node.AddColumnOnce(origin.FieldPoolID)
		}
	}
	if ps := oq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := oq.ctx.Limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := oq.ctx.Offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := oq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (oq *OriginQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(oq.driver.Dialect())
	t1 := builder.Table(origin.Table)
	columns := oq.ctx.Fields
	if len(columns) == 0 {
		columns = origin.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if oq.sql != nil {
		selector = oq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if oq.ctx.Unique != nil && *oq.ctx.Unique {
		selector.Distinct()
	}
	for _, p := range oq.predicates {
		p(selector)
	}
	for _, p := range oq.order {
		p(selector)
	}
	if offset := oq.ctx.Offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := oq.ctx.Limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// OriginGroupBy is the group-by builder for Origin entities.
type OriginGroupBy struct {
	selector
	build *OriginQuery
}

// Aggregate adds the given aggregation functions to the group-by query.
func (ogb *OriginGroupBy) Aggregate(fns ...AggregateFunc) *OriginGroupBy {
	ogb.fns = append(ogb.fns, fns...)
	return ogb
}

// Scan applies the selector query and scans the result into the given value.
func (ogb *OriginGroupBy) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, ogb.build.ctx, "GroupBy")
	if err := ogb.build.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*OriginQuery, *OriginGroupBy](ctx, ogb.build, ogb, ogb.build.inters, v)
}

func (ogb *OriginGroupBy) sqlScan(ctx context.Context, root *OriginQuery, v any) error {
	selector := root.sqlQuery(ctx).Select()
	aggregation := make([]string, 0, len(ogb.fns))
	for _, fn := range ogb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(*ogb.flds)+len(ogb.fns))
		for _, f := range *ogb.flds {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	selector.GroupBy(selector.Columns(*ogb.flds...)...)
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := ogb.build.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// OriginSelect is the builder for selecting fields of Origin entities.
type OriginSelect struct {
	*OriginQuery
	selector
}

// Aggregate adds the given aggregation functions to the selector query.
func (os *OriginSelect) Aggregate(fns ...AggregateFunc) *OriginSelect {
	os.fns = append(os.fns, fns...)
	return os
}

// Scan applies the selector query and scans the result into the given value.
func (os *OriginSelect) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, os.ctx, "Select")
	if err := os.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*OriginQuery, *OriginSelect](ctx, os.OriginQuery, os, os.inters, v)
}

func (os *OriginSelect) sqlScan(ctx context.Context, root *OriginQuery, v any) error {
	selector := root.sqlQuery(ctx)
	aggregation := make([]string, 0, len(os.fns))
	for _, fn := range os.fns {
		aggregation = append(aggregation, fn(selector))
	}
	switch n := len(*os.selector.flds); {
	case n == 0 && len(aggregation) > 0:
		selector.Select(aggregation...)
	case n != 0 && len(aggregation) > 0:
		selector.AppendSelect(aggregation...)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := os.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}
