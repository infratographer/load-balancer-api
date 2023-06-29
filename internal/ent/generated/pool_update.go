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

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"go.infratographer.com/load-balancer-api/internal/ent/generated/origin"
	"go.infratographer.com/load-balancer-api/internal/ent/generated/pool"
	"go.infratographer.com/load-balancer-api/internal/ent/generated/port"
	"go.infratographer.com/load-balancer-api/internal/ent/generated/predicate"
	"go.infratographer.com/x/gidx"
)

// PoolUpdate is the builder for updating Pool entities.
type PoolUpdate struct {
	config
	hooks    []Hook
	mutation *PoolMutation
}

// Where appends a list predicates to the PoolUpdate builder.
func (pu *PoolUpdate) Where(ps ...predicate.Pool) *PoolUpdate {
	pu.mutation.Where(ps...)
	return pu
}

// SetName sets the "name" field.
func (pu *PoolUpdate) SetName(s string) *PoolUpdate {
	pu.mutation.SetName(s)
	return pu
}

// SetProtocol sets the "protocol" field.
func (pu *PoolUpdate) SetProtocol(po pool.Protocol) *PoolUpdate {
	pu.mutation.SetProtocol(po)
	return pu
}

// AddPortIDs adds the "ports" edge to the Port entity by IDs.
func (pu *PoolUpdate) AddPortIDs(ids ...gidx.PrefixedID) *PoolUpdate {
	pu.mutation.AddPortIDs(ids...)
	return pu
}

// AddPorts adds the "ports" edges to the Port entity.
func (pu *PoolUpdate) AddPorts(p ...*Port) *PoolUpdate {
	ids := make([]gidx.PrefixedID, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return pu.AddPortIDs(ids...)
}

// AddOriginIDs adds the "origins" edge to the Origin entity by IDs.
func (pu *PoolUpdate) AddOriginIDs(ids ...gidx.PrefixedID) *PoolUpdate {
	pu.mutation.AddOriginIDs(ids...)
	return pu
}

// AddOrigins adds the "origins" edges to the Origin entity.
func (pu *PoolUpdate) AddOrigins(o ...*Origin) *PoolUpdate {
	ids := make([]gidx.PrefixedID, len(o))
	for i := range o {
		ids[i] = o[i].ID
	}
	return pu.AddOriginIDs(ids...)
}

// Mutation returns the PoolMutation object of the builder.
func (pu *PoolUpdate) Mutation() *PoolMutation {
	return pu.mutation
}

// ClearPorts clears all "ports" edges to the Port entity.
func (pu *PoolUpdate) ClearPorts() *PoolUpdate {
	pu.mutation.ClearPorts()
	return pu
}

// RemovePortIDs removes the "ports" edge to Port entities by IDs.
func (pu *PoolUpdate) RemovePortIDs(ids ...gidx.PrefixedID) *PoolUpdate {
	pu.mutation.RemovePortIDs(ids...)
	return pu
}

// RemovePorts removes "ports" edges to Port entities.
func (pu *PoolUpdate) RemovePorts(p ...*Port) *PoolUpdate {
	ids := make([]gidx.PrefixedID, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return pu.RemovePortIDs(ids...)
}

// ClearOrigins clears all "origins" edges to the Origin entity.
func (pu *PoolUpdate) ClearOrigins() *PoolUpdate {
	pu.mutation.ClearOrigins()
	return pu
}

// RemoveOriginIDs removes the "origins" edge to Origin entities by IDs.
func (pu *PoolUpdate) RemoveOriginIDs(ids ...gidx.PrefixedID) *PoolUpdate {
	pu.mutation.RemoveOriginIDs(ids...)
	return pu
}

// RemoveOrigins removes "origins" edges to Origin entities.
func (pu *PoolUpdate) RemoveOrigins(o ...*Origin) *PoolUpdate {
	ids := make([]gidx.PrefixedID, len(o))
	for i := range o {
		ids[i] = o[i].ID
	}
	return pu.RemoveOriginIDs(ids...)
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (pu *PoolUpdate) Save(ctx context.Context) (int, error) {
	pu.defaults()
	return withHooks(ctx, pu.sqlSave, pu.mutation, pu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (pu *PoolUpdate) SaveX(ctx context.Context) int {
	affected, err := pu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (pu *PoolUpdate) Exec(ctx context.Context) error {
	_, err := pu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (pu *PoolUpdate) ExecX(ctx context.Context) {
	if err := pu.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (pu *PoolUpdate) defaults() {
	if _, ok := pu.mutation.UpdatedAt(); !ok {
		v := pool.UpdateDefaultUpdatedAt()
		pu.mutation.SetUpdatedAt(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (pu *PoolUpdate) check() error {
	if v, ok := pu.mutation.Name(); ok {
		if err := pool.NameValidator(v); err != nil {
			return &ValidationError{Name: "name", err: fmt.Errorf(`generated: validator failed for field "Pool.name": %w`, err)}
		}
	}
	if v, ok := pu.mutation.Protocol(); ok {
		if err := pool.ProtocolValidator(v); err != nil {
			return &ValidationError{Name: "protocol", err: fmt.Errorf(`generated: validator failed for field "Pool.protocol": %w`, err)}
		}
	}
	return nil
}

func (pu *PoolUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := pu.check(); err != nil {
		return n, err
	}
	_spec := sqlgraph.NewUpdateSpec(pool.Table, pool.Columns, sqlgraph.NewFieldSpec(pool.FieldID, field.TypeString))
	if ps := pu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := pu.mutation.UpdatedAt(); ok {
		_spec.SetField(pool.FieldUpdatedAt, field.TypeTime, value)
	}
	if value, ok := pu.mutation.Name(); ok {
		_spec.SetField(pool.FieldName, field.TypeString, value)
	}
	if value, ok := pu.mutation.Protocol(); ok {
		_spec.SetField(pool.FieldProtocol, field.TypeEnum, value)
	}
	if pu.mutation.PortsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   pool.PortsTable,
			Columns: pool.PortsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(port.FieldID, field.TypeString),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := pu.mutation.RemovedPortsIDs(); len(nodes) > 0 && !pu.mutation.PortsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   pool.PortsTable,
			Columns: pool.PortsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(port.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := pu.mutation.PortsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   pool.PortsTable,
			Columns: pool.PortsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(port.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if pu.mutation.OriginsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   pool.OriginsTable,
			Columns: []string{pool.OriginsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(origin.FieldID, field.TypeString),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := pu.mutation.RemovedOriginsIDs(); len(nodes) > 0 && !pu.mutation.OriginsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   pool.OriginsTable,
			Columns: []string{pool.OriginsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(origin.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := pu.mutation.OriginsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   pool.OriginsTable,
			Columns: []string{pool.OriginsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(origin.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, pu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{pool.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	pu.mutation.done = true
	return n, nil
}

// PoolUpdateOne is the builder for updating a single Pool entity.
type PoolUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *PoolMutation
}

// SetName sets the "name" field.
func (puo *PoolUpdateOne) SetName(s string) *PoolUpdateOne {
	puo.mutation.SetName(s)
	return puo
}

// SetProtocol sets the "protocol" field.
func (puo *PoolUpdateOne) SetProtocol(po pool.Protocol) *PoolUpdateOne {
	puo.mutation.SetProtocol(po)
	return puo
}

// AddPortIDs adds the "ports" edge to the Port entity by IDs.
func (puo *PoolUpdateOne) AddPortIDs(ids ...gidx.PrefixedID) *PoolUpdateOne {
	puo.mutation.AddPortIDs(ids...)
	return puo
}

// AddPorts adds the "ports" edges to the Port entity.
func (puo *PoolUpdateOne) AddPorts(p ...*Port) *PoolUpdateOne {
	ids := make([]gidx.PrefixedID, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return puo.AddPortIDs(ids...)
}

// AddOriginIDs adds the "origins" edge to the Origin entity by IDs.
func (puo *PoolUpdateOne) AddOriginIDs(ids ...gidx.PrefixedID) *PoolUpdateOne {
	puo.mutation.AddOriginIDs(ids...)
	return puo
}

// AddOrigins adds the "origins" edges to the Origin entity.
func (puo *PoolUpdateOne) AddOrigins(o ...*Origin) *PoolUpdateOne {
	ids := make([]gidx.PrefixedID, len(o))
	for i := range o {
		ids[i] = o[i].ID
	}
	return puo.AddOriginIDs(ids...)
}

// Mutation returns the PoolMutation object of the builder.
func (puo *PoolUpdateOne) Mutation() *PoolMutation {
	return puo.mutation
}

// ClearPorts clears all "ports" edges to the Port entity.
func (puo *PoolUpdateOne) ClearPorts() *PoolUpdateOne {
	puo.mutation.ClearPorts()
	return puo
}

// RemovePortIDs removes the "ports" edge to Port entities by IDs.
func (puo *PoolUpdateOne) RemovePortIDs(ids ...gidx.PrefixedID) *PoolUpdateOne {
	puo.mutation.RemovePortIDs(ids...)
	return puo
}

// RemovePorts removes "ports" edges to Port entities.
func (puo *PoolUpdateOne) RemovePorts(p ...*Port) *PoolUpdateOne {
	ids := make([]gidx.PrefixedID, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return puo.RemovePortIDs(ids...)
}

// ClearOrigins clears all "origins" edges to the Origin entity.
func (puo *PoolUpdateOne) ClearOrigins() *PoolUpdateOne {
	puo.mutation.ClearOrigins()
	return puo
}

// RemoveOriginIDs removes the "origins" edge to Origin entities by IDs.
func (puo *PoolUpdateOne) RemoveOriginIDs(ids ...gidx.PrefixedID) *PoolUpdateOne {
	puo.mutation.RemoveOriginIDs(ids...)
	return puo
}

// RemoveOrigins removes "origins" edges to Origin entities.
func (puo *PoolUpdateOne) RemoveOrigins(o ...*Origin) *PoolUpdateOne {
	ids := make([]gidx.PrefixedID, len(o))
	for i := range o {
		ids[i] = o[i].ID
	}
	return puo.RemoveOriginIDs(ids...)
}

// Where appends a list predicates to the PoolUpdate builder.
func (puo *PoolUpdateOne) Where(ps ...predicate.Pool) *PoolUpdateOne {
	puo.mutation.Where(ps...)
	return puo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (puo *PoolUpdateOne) Select(field string, fields ...string) *PoolUpdateOne {
	puo.fields = append([]string{field}, fields...)
	return puo
}

// Save executes the query and returns the updated Pool entity.
func (puo *PoolUpdateOne) Save(ctx context.Context) (*Pool, error) {
	puo.defaults()
	return withHooks(ctx, puo.sqlSave, puo.mutation, puo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (puo *PoolUpdateOne) SaveX(ctx context.Context) *Pool {
	node, err := puo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (puo *PoolUpdateOne) Exec(ctx context.Context) error {
	_, err := puo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (puo *PoolUpdateOne) ExecX(ctx context.Context) {
	if err := puo.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (puo *PoolUpdateOne) defaults() {
	if _, ok := puo.mutation.UpdatedAt(); !ok {
		v := pool.UpdateDefaultUpdatedAt()
		puo.mutation.SetUpdatedAt(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (puo *PoolUpdateOne) check() error {
	if v, ok := puo.mutation.Name(); ok {
		if err := pool.NameValidator(v); err != nil {
			return &ValidationError{Name: "name", err: fmt.Errorf(`generated: validator failed for field "Pool.name": %w`, err)}
		}
	}
	if v, ok := puo.mutation.Protocol(); ok {
		if err := pool.ProtocolValidator(v); err != nil {
			return &ValidationError{Name: "protocol", err: fmt.Errorf(`generated: validator failed for field "Pool.protocol": %w`, err)}
		}
	}
	return nil
}

func (puo *PoolUpdateOne) sqlSave(ctx context.Context) (_node *Pool, err error) {
	if err := puo.check(); err != nil {
		return _node, err
	}
	_spec := sqlgraph.NewUpdateSpec(pool.Table, pool.Columns, sqlgraph.NewFieldSpec(pool.FieldID, field.TypeString))
	id, ok := puo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`generated: missing "Pool.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := puo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, pool.FieldID)
		for _, f := range fields {
			if !pool.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("generated: invalid field %q for query", f)}
			}
			if f != pool.FieldID {
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
		_spec.SetField(pool.FieldUpdatedAt, field.TypeTime, value)
	}
	if value, ok := puo.mutation.Name(); ok {
		_spec.SetField(pool.FieldName, field.TypeString, value)
	}
	if value, ok := puo.mutation.Protocol(); ok {
		_spec.SetField(pool.FieldProtocol, field.TypeEnum, value)
	}
	if puo.mutation.PortsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   pool.PortsTable,
			Columns: pool.PortsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(port.FieldID, field.TypeString),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := puo.mutation.RemovedPortsIDs(); len(nodes) > 0 && !puo.mutation.PortsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   pool.PortsTable,
			Columns: pool.PortsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(port.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := puo.mutation.PortsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   pool.PortsTable,
			Columns: pool.PortsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(port.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if puo.mutation.OriginsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   pool.OriginsTable,
			Columns: []string{pool.OriginsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(origin.FieldID, field.TypeString),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := puo.mutation.RemovedOriginsIDs(); len(nodes) > 0 && !puo.mutation.OriginsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   pool.OriginsTable,
			Columns: []string{pool.OriginsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(origin.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := puo.mutation.OriginsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   pool.OriginsTable,
			Columns: []string{pool.OriginsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(origin.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &Pool{config: puo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, puo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{pool.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	puo.mutation.done = true
	return _node, nil
}
