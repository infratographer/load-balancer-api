// Copyright Infratographer, Inc. and/or licensed to Infratographer, Inc. under one
// or more contributor license agreements. Licensed under the Elastic License 2.0;
// you may not use this file except in compliance with the Elastic License 2.0.
//
// Code generated by entc, DO NOT EDIT.

package generated

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"go.infratographer.com/load-balancer-api/internal/ent/generated/pool"
	"go.infratographer.com/load-balancer-api/internal/ent/generated/port"
	"go.infratographer.com/load-balancer-api/internal/ent/generated/predicate"
	"go.infratographer.com/x/gidx"
)

// PortUpdate is the builder for updating Port entities.
type PortUpdate struct {
	config
	hooks    []Hook
	mutation *PortMutation
}

// Where appends a list predicates to the PortUpdate builder.
func (pu *PortUpdate) Where(ps ...predicate.Port) *PortUpdate {
	pu.mutation.Where(ps...)
	return pu
}

// SetNumber sets the "number" field.
func (pu *PortUpdate) SetNumber(i int) *PortUpdate {
	pu.mutation.ResetNumber()
	pu.mutation.SetNumber(i)
	return pu
}

// AddNumber adds i to the "number" field.
func (pu *PortUpdate) AddNumber(i int) *PortUpdate {
	pu.mutation.AddNumber(i)
	return pu
}

// SetName sets the "name" field.
func (pu *PortUpdate) SetName(s string) *PortUpdate {
	pu.mutation.SetName(s)
	return pu
}

// AddPoolIDs adds the "pools" edge to the Pool entity by IDs.
func (pu *PortUpdate) AddPoolIDs(ids ...gidx.PrefixedID) *PortUpdate {
	pu.mutation.AddPoolIDs(ids...)
	return pu
}

// AddPools adds the "pools" edges to the Pool entity.
func (pu *PortUpdate) AddPools(p ...*Pool) *PortUpdate {
	ids := make([]gidx.PrefixedID, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return pu.AddPoolIDs(ids...)
}

// Mutation returns the PortMutation object of the builder.
func (pu *PortUpdate) Mutation() *PortMutation {
	return pu.mutation
}

// ClearPools clears all "pools" edges to the Pool entity.
func (pu *PortUpdate) ClearPools() *PortUpdate {
	pu.mutation.ClearPools()
	return pu
}

// RemovePoolIDs removes the "pools" edge to Pool entities by IDs.
func (pu *PortUpdate) RemovePoolIDs(ids ...gidx.PrefixedID) *PortUpdate {
	pu.mutation.RemovePoolIDs(ids...)
	return pu
}

// RemovePools removes "pools" edges to Pool entities.
func (pu *PortUpdate) RemovePools(p ...*Pool) *PortUpdate {
	ids := make([]gidx.PrefixedID, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return pu.RemovePoolIDs(ids...)
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (pu *PortUpdate) Save(ctx context.Context) (int, error) {
	pu.defaults()
	return withHooks(ctx, pu.sqlSave, pu.mutation, pu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (pu *PortUpdate) SaveX(ctx context.Context) int {
	affected, err := pu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (pu *PortUpdate) Exec(ctx context.Context) error {
	_, err := pu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (pu *PortUpdate) ExecX(ctx context.Context) {
	if err := pu.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (pu *PortUpdate) defaults() {
	if _, ok := pu.mutation.UpdatedAt(); !ok {
		v := port.UpdateDefaultUpdatedAt()
		pu.mutation.SetUpdatedAt(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (pu *PortUpdate) check() error {
	if v, ok := pu.mutation.Number(); ok {
		if err := port.NumberValidator(v); err != nil {
			return &ValidationError{Name: "number", err: fmt.Errorf(`generated: validator failed for field "Port.number": %w`, err)}
		}
	}
	if v, ok := pu.mutation.Name(); ok {
		if err := port.NameValidator(v); err != nil {
			return &ValidationError{Name: "name", err: fmt.Errorf(`generated: validator failed for field "Port.name": %w`, err)}
		}
	}
	if _, ok := pu.mutation.LoadBalancerID(); pu.mutation.LoadBalancerCleared() && !ok {
		return errors.New(`generated: clearing a required unique edge "Port.load_balancer"`)
	}
	return nil
}

func (pu *PortUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := pu.check(); err != nil {
		return n, err
	}
	_spec := sqlgraph.NewUpdateSpec(port.Table, port.Columns, sqlgraph.NewFieldSpec(port.FieldID, field.TypeString))
	if ps := pu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := pu.mutation.UpdatedAt(); ok {
		_spec.SetField(port.FieldUpdatedAt, field.TypeTime, value)
	}
	if value, ok := pu.mutation.Number(); ok {
		_spec.SetField(port.FieldNumber, field.TypeInt, value)
	}
	if value, ok := pu.mutation.AddedNumber(); ok {
		_spec.AddField(port.FieldNumber, field.TypeInt, value)
	}
	if value, ok := pu.mutation.Name(); ok {
		_spec.SetField(port.FieldName, field.TypeString, value)
	}
	if pu.mutation.PoolsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   port.PoolsTable,
			Columns: port.PoolsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(pool.FieldID, field.TypeString),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := pu.mutation.RemovedPoolsIDs(); len(nodes) > 0 && !pu.mutation.PoolsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   port.PoolsTable,
			Columns: port.PoolsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(pool.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := pu.mutation.PoolsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   port.PoolsTable,
			Columns: port.PoolsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(pool.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, pu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{port.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	pu.mutation.done = true
	return n, nil
}

// PortUpdateOne is the builder for updating a single Port entity.
type PortUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *PortMutation
}

// SetNumber sets the "number" field.
func (puo *PortUpdateOne) SetNumber(i int) *PortUpdateOne {
	puo.mutation.ResetNumber()
	puo.mutation.SetNumber(i)
	return puo
}

// AddNumber adds i to the "number" field.
func (puo *PortUpdateOne) AddNumber(i int) *PortUpdateOne {
	puo.mutation.AddNumber(i)
	return puo
}

// SetName sets the "name" field.
func (puo *PortUpdateOne) SetName(s string) *PortUpdateOne {
	puo.mutation.SetName(s)
	return puo
}

// AddPoolIDs adds the "pools" edge to the Pool entity by IDs.
func (puo *PortUpdateOne) AddPoolIDs(ids ...gidx.PrefixedID) *PortUpdateOne {
	puo.mutation.AddPoolIDs(ids...)
	return puo
}

// AddPools adds the "pools" edges to the Pool entity.
func (puo *PortUpdateOne) AddPools(p ...*Pool) *PortUpdateOne {
	ids := make([]gidx.PrefixedID, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return puo.AddPoolIDs(ids...)
}

// Mutation returns the PortMutation object of the builder.
func (puo *PortUpdateOne) Mutation() *PortMutation {
	return puo.mutation
}

// ClearPools clears all "pools" edges to the Pool entity.
func (puo *PortUpdateOne) ClearPools() *PortUpdateOne {
	puo.mutation.ClearPools()
	return puo
}

// RemovePoolIDs removes the "pools" edge to Pool entities by IDs.
func (puo *PortUpdateOne) RemovePoolIDs(ids ...gidx.PrefixedID) *PortUpdateOne {
	puo.mutation.RemovePoolIDs(ids...)
	return puo
}

// RemovePools removes "pools" edges to Pool entities.
func (puo *PortUpdateOne) RemovePools(p ...*Pool) *PortUpdateOne {
	ids := make([]gidx.PrefixedID, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return puo.RemovePoolIDs(ids...)
}

// Where appends a list predicates to the PortUpdate builder.
func (puo *PortUpdateOne) Where(ps ...predicate.Port) *PortUpdateOne {
	puo.mutation.Where(ps...)
	return puo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (puo *PortUpdateOne) Select(field string, fields ...string) *PortUpdateOne {
	puo.fields = append([]string{field}, fields...)
	return puo
}

// Save executes the query and returns the updated Port entity.
func (puo *PortUpdateOne) Save(ctx context.Context) (*Port, error) {
	puo.defaults()
	return withHooks(ctx, puo.sqlSave, puo.mutation, puo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (puo *PortUpdateOne) SaveX(ctx context.Context) *Port {
	node, err := puo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (puo *PortUpdateOne) Exec(ctx context.Context) error {
	_, err := puo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (puo *PortUpdateOne) ExecX(ctx context.Context) {
	if err := puo.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (puo *PortUpdateOne) defaults() {
	if _, ok := puo.mutation.UpdatedAt(); !ok {
		v := port.UpdateDefaultUpdatedAt()
		puo.mutation.SetUpdatedAt(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (puo *PortUpdateOne) check() error {
	if v, ok := puo.mutation.Number(); ok {
		if err := port.NumberValidator(v); err != nil {
			return &ValidationError{Name: "number", err: fmt.Errorf(`generated: validator failed for field "Port.number": %w`, err)}
		}
	}
	if v, ok := puo.mutation.Name(); ok {
		if err := port.NameValidator(v); err != nil {
			return &ValidationError{Name: "name", err: fmt.Errorf(`generated: validator failed for field "Port.name": %w`, err)}
		}
	}
	if _, ok := puo.mutation.LoadBalancerID(); puo.mutation.LoadBalancerCleared() && !ok {
		return errors.New(`generated: clearing a required unique edge "Port.load_balancer"`)
	}
	return nil
}

func (puo *PortUpdateOne) sqlSave(ctx context.Context) (_node *Port, err error) {
	if err := puo.check(); err != nil {
		return _node, err
	}
	_spec := sqlgraph.NewUpdateSpec(port.Table, port.Columns, sqlgraph.NewFieldSpec(port.FieldID, field.TypeString))
	id, ok := puo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`generated: missing "Port.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := puo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, port.FieldID)
		for _, f := range fields {
			if !port.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("generated: invalid field %q for query", f)}
			}
			if f != port.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := puo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := puo.mutation.UpdatedAt(); ok {
		_spec.SetField(port.FieldUpdatedAt, field.TypeTime, value)
	}
	if value, ok := puo.mutation.Number(); ok {
		_spec.SetField(port.FieldNumber, field.TypeInt, value)
	}
	if value, ok := puo.mutation.AddedNumber(); ok {
		_spec.AddField(port.FieldNumber, field.TypeInt, value)
	}
	if value, ok := puo.mutation.Name(); ok {
		_spec.SetField(port.FieldName, field.TypeString, value)
	}
	if puo.mutation.PoolsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   port.PoolsTable,
			Columns: port.PoolsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(pool.FieldID, field.TypeString),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := puo.mutation.RemovedPoolsIDs(); len(nodes) > 0 && !puo.mutation.PoolsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   port.PoolsTable,
			Columns: port.PoolsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(pool.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := puo.mutation.PoolsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   port.PoolsTable,
			Columns: port.PoolsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(pool.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &Port{config: puo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, puo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{port.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	puo.mutation.done = true
	return _node, nil
}
