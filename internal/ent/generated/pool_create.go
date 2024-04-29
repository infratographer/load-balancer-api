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
	"time"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"go.infratographer.com/load-balancer-api/internal/ent/generated/origin"
	"go.infratographer.com/load-balancer-api/internal/ent/generated/pool"
	"go.infratographer.com/load-balancer-api/internal/ent/generated/port"
	"go.infratographer.com/x/gidx"
)

// PoolCreate is the builder for creating a Pool entity.
type PoolCreate struct {
	config
	mutation *PoolMutation
	hooks    []Hook
}

// SetCreatedAt sets the "created_at" field.
func (pc *PoolCreate) SetCreatedAt(t time.Time) *PoolCreate {
	pc.mutation.SetCreatedAt(t)
	return pc
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (pc *PoolCreate) SetNillableCreatedAt(t *time.Time) *PoolCreate {
	if t != nil {
		pc.SetCreatedAt(*t)
	}
	return pc
}

// SetUpdatedAt sets the "updated_at" field.
func (pc *PoolCreate) SetUpdatedAt(t time.Time) *PoolCreate {
	pc.mutation.SetUpdatedAt(t)
	return pc
}

// SetNillableUpdatedAt sets the "updated_at" field if the given value is not nil.
func (pc *PoolCreate) SetNillableUpdatedAt(t *time.Time) *PoolCreate {
	if t != nil {
		pc.SetUpdatedAt(*t)
	}
	return pc
}

// SetCreatedBy sets the "created_by" field.
func (pc *PoolCreate) SetCreatedBy(s string) *PoolCreate {
	pc.mutation.SetCreatedBy(s)
	return pc
}

// SetNillableCreatedBy sets the "created_by" field if the given value is not nil.
func (pc *PoolCreate) SetNillableCreatedBy(s *string) *PoolCreate {
	if s != nil {
		pc.SetCreatedBy(*s)
	}
	return pc
}

// SetUpdatedBy sets the "updated_by" field.
func (pc *PoolCreate) SetUpdatedBy(s string) *PoolCreate {
	pc.mutation.SetUpdatedBy(s)
	return pc
}

// SetNillableUpdatedBy sets the "updated_by" field if the given value is not nil.
func (pc *PoolCreate) SetNillableUpdatedBy(s *string) *PoolCreate {
	if s != nil {
		pc.SetUpdatedBy(*s)
	}
	return pc
}

// SetDeletedAt sets the "deleted_at" field.
func (pc *PoolCreate) SetDeletedAt(t time.Time) *PoolCreate {
	pc.mutation.SetDeletedAt(t)
	return pc
}

// SetNillableDeletedAt sets the "deleted_at" field if the given value is not nil.
func (pc *PoolCreate) SetNillableDeletedAt(t *time.Time) *PoolCreate {
	if t != nil {
		pc.SetDeletedAt(*t)
	}
	return pc
}

// SetDeletedBy sets the "deleted_by" field.
func (pc *PoolCreate) SetDeletedBy(s string) *PoolCreate {
	pc.mutation.SetDeletedBy(s)
	return pc
}

// SetNillableDeletedBy sets the "deleted_by" field if the given value is not nil.
func (pc *PoolCreate) SetNillableDeletedBy(s *string) *PoolCreate {
	if s != nil {
		pc.SetDeletedBy(*s)
	}
	return pc
}

// SetName sets the "name" field.
func (pc *PoolCreate) SetName(s string) *PoolCreate {
	pc.mutation.SetName(s)
	return pc
}

// SetProtocol sets the "protocol" field.
func (pc *PoolCreate) SetProtocol(po pool.Protocol) *PoolCreate {
	pc.mutation.SetProtocol(po)
	return pc
}

// SetOwnerID sets the "owner_id" field.
func (pc *PoolCreate) SetOwnerID(gi gidx.PrefixedID) *PoolCreate {
	pc.mutation.SetOwnerID(gi)
	return pc
}

// SetID sets the "id" field.
func (pc *PoolCreate) SetID(gi gidx.PrefixedID) *PoolCreate {
	pc.mutation.SetID(gi)
	return pc
}

// SetNillableID sets the "id" field if the given value is not nil.
func (pc *PoolCreate) SetNillableID(gi *gidx.PrefixedID) *PoolCreate {
	if gi != nil {
		pc.SetID(*gi)
	}
	return pc
}

// AddPortIDs adds the "ports" edge to the Port entity by IDs.
func (pc *PoolCreate) AddPortIDs(ids ...gidx.PrefixedID) *PoolCreate {
	pc.mutation.AddPortIDs(ids...)
	return pc
}

// AddPorts adds the "ports" edges to the Port entity.
func (pc *PoolCreate) AddPorts(p ...*Port) *PoolCreate {
	ids := make([]gidx.PrefixedID, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return pc.AddPortIDs(ids...)
}

// AddOriginIDs adds the "origins" edge to the Origin entity by IDs.
func (pc *PoolCreate) AddOriginIDs(ids ...gidx.PrefixedID) *PoolCreate {
	pc.mutation.AddOriginIDs(ids...)
	return pc
}

// AddOrigins adds the "origins" edges to the Origin entity.
func (pc *PoolCreate) AddOrigins(o ...*Origin) *PoolCreate {
	ids := make([]gidx.PrefixedID, len(o))
	for i := range o {
		ids[i] = o[i].ID
	}
	return pc.AddOriginIDs(ids...)
}

// Mutation returns the PoolMutation object of the builder.
func (pc *PoolCreate) Mutation() *PoolMutation {
	return pc.mutation
}

// Save creates the Pool in the database.
func (pc *PoolCreate) Save(ctx context.Context) (*Pool, error) {
	if err := pc.defaults(); err != nil {
		return nil, err
	}
	return withHooks(ctx, pc.sqlSave, pc.mutation, pc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (pc *PoolCreate) SaveX(ctx context.Context) *Pool {
	v, err := pc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (pc *PoolCreate) Exec(ctx context.Context) error {
	_, err := pc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (pc *PoolCreate) ExecX(ctx context.Context) {
	if err := pc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (pc *PoolCreate) defaults() error {
	if _, ok := pc.mutation.CreatedAt(); !ok {
		if pool.DefaultCreatedAt == nil {
			return fmt.Errorf("generated: uninitialized pool.DefaultCreatedAt (forgotten import generated/runtime?)")
		}
		v := pool.DefaultCreatedAt()
		pc.mutation.SetCreatedAt(v)
	}
	if _, ok := pc.mutation.UpdatedAt(); !ok {
		if pool.DefaultUpdatedAt == nil {
			return fmt.Errorf("generated: uninitialized pool.DefaultUpdatedAt (forgotten import generated/runtime?)")
		}
		v := pool.DefaultUpdatedAt()
		pc.mutation.SetUpdatedAt(v)
	}
	if _, ok := pc.mutation.ID(); !ok {
		if pool.DefaultID == nil {
			return fmt.Errorf("generated: uninitialized pool.DefaultID (forgotten import generated/runtime?)")
		}
		v := pool.DefaultID()
		pc.mutation.SetID(v)
	}
	return nil
}

// check runs all checks and user-defined validators on the builder.
func (pc *PoolCreate) check() error {
	if _, ok := pc.mutation.CreatedAt(); !ok {
		return &ValidationError{Name: "created_at", err: errors.New(`generated: missing required field "Pool.created_at"`)}
	}
	if _, ok := pc.mutation.UpdatedAt(); !ok {
		return &ValidationError{Name: "updated_at", err: errors.New(`generated: missing required field "Pool.updated_at"`)}
	}
	if _, ok := pc.mutation.Name(); !ok {
		return &ValidationError{Name: "name", err: errors.New(`generated: missing required field "Pool.name"`)}
	}
	if v, ok := pc.mutation.Name(); ok {
		if err := pool.NameValidator(v); err != nil {
			return &ValidationError{Name: "name", err: fmt.Errorf(`generated: validator failed for field "Pool.name": %w`, err)}
		}
	}
	if _, ok := pc.mutation.Protocol(); !ok {
		return &ValidationError{Name: "protocol", err: errors.New(`generated: missing required field "Pool.protocol"`)}
	}
	if v, ok := pc.mutation.Protocol(); ok {
		if err := pool.ProtocolValidator(v); err != nil {
			return &ValidationError{Name: "protocol", err: fmt.Errorf(`generated: validator failed for field "Pool.protocol": %w`, err)}
		}
	}
	if _, ok := pc.mutation.OwnerID(); !ok {
		return &ValidationError{Name: "owner_id", err: errors.New(`generated: missing required field "Pool.owner_id"`)}
	}
	if v, ok := pc.mutation.OwnerID(); ok {
		if err := pool.OwnerIDValidator(string(v)); err != nil {
			return &ValidationError{Name: "owner_id", err: fmt.Errorf(`generated: validator failed for field "Pool.owner_id": %w`, err)}
		}
	}
	if v, ok := pc.mutation.ID(); ok {
		if err := v.Validate(); err != nil {
			return &ValidationError{Name: "id", err: fmt.Errorf(`generated: validator failed for field "Pool.id": %w`, err)}
		}
	}
	return nil
}

func (pc *PoolCreate) sqlSave(ctx context.Context) (*Pool, error) {
	if err := pc.check(); err != nil {
		return nil, err
	}
	_node, _spec := pc.createSpec()
	if err := sqlgraph.CreateNode(ctx, pc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	if _spec.ID.Value != nil {
		if id, ok := _spec.ID.Value.(*gidx.PrefixedID); ok {
			_node.ID = *id
		} else if err := _node.ID.Scan(_spec.ID.Value); err != nil {
			return nil, err
		}
	}
	pc.mutation.id = &_node.ID
	pc.mutation.done = true
	return _node, nil
}

func (pc *PoolCreate) createSpec() (*Pool, *sqlgraph.CreateSpec) {
	var (
		_node = &Pool{config: pc.config}
		_spec = sqlgraph.NewCreateSpec(pool.Table, sqlgraph.NewFieldSpec(pool.FieldID, field.TypeString))
	)
	if id, ok := pc.mutation.ID(); ok {
		_node.ID = id
		_spec.ID.Value = &id
	}
	if value, ok := pc.mutation.CreatedAt(); ok {
		_spec.SetField(pool.FieldCreatedAt, field.TypeTime, value)
		_node.CreatedAt = value
	}
	if value, ok := pc.mutation.UpdatedAt(); ok {
		_spec.SetField(pool.FieldUpdatedAt, field.TypeTime, value)
		_node.UpdatedAt = value
	}
	if value, ok := pc.mutation.CreatedBy(); ok {
		_spec.SetField(pool.FieldCreatedBy, field.TypeString, value)
		_node.CreatedBy = value
	}
	if value, ok := pc.mutation.UpdatedBy(); ok {
		_spec.SetField(pool.FieldUpdatedBy, field.TypeString, value)
		_node.UpdatedBy = value
	}
	if value, ok := pc.mutation.DeletedAt(); ok {
		_spec.SetField(pool.FieldDeletedAt, field.TypeTime, value)
		_node.DeletedAt = value
	}
	if value, ok := pc.mutation.DeletedBy(); ok {
		_spec.SetField(pool.FieldDeletedBy, field.TypeString, value)
		_node.DeletedBy = value
	}
	if value, ok := pc.mutation.Name(); ok {
		_spec.SetField(pool.FieldName, field.TypeString, value)
		_node.Name = value
	}
	if value, ok := pc.mutation.Protocol(); ok {
		_spec.SetField(pool.FieldProtocol, field.TypeEnum, value)
		_node.Protocol = value
	}
	if value, ok := pc.mutation.OwnerID(); ok {
		_spec.SetField(pool.FieldOwnerID, field.TypeString, value)
		_node.OwnerID = value
	}
	if nodes := pc.mutation.PortsIDs(); len(nodes) > 0 {
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := pc.mutation.OriginsIDs(); len(nodes) > 0 {
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// PoolCreateBulk is the builder for creating many Pool entities in bulk.
type PoolCreateBulk struct {
	config
	err      error
	builders []*PoolCreate
}

// Save creates the Pool entities in the database.
func (pcb *PoolCreateBulk) Save(ctx context.Context) ([]*Pool, error) {
	if pcb.err != nil {
		return nil, pcb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(pcb.builders))
	nodes := make([]*Pool, len(pcb.builders))
	mutators := make([]Mutator, len(pcb.builders))
	for i := range pcb.builders {
		func(i int, root context.Context) {
			builder := pcb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*PoolMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				var err error
				nodes[i], specs[i] = builder.createSpec()
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, pcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, pcb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				mutation.done = true
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, pcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (pcb *PoolCreateBulk) SaveX(ctx context.Context) []*Pool {
	v, err := pcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (pcb *PoolCreateBulk) Exec(ctx context.Context) error {
	_, err := pcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (pcb *PoolCreateBulk) ExecX(ctx context.Context) {
	if err := pcb.Exec(ctx); err != nil {
		panic(err)
	}
}
