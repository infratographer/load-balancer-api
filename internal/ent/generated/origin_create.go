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
	"go.infratographer.com/x/gidx"
)

// OriginCreate is the builder for creating a Origin entity.
type OriginCreate struct {
	config
	mutation *OriginMutation
	hooks    []Hook
}

// SetCreatedAt sets the "created_at" field.
func (oc *OriginCreate) SetCreatedAt(t time.Time) *OriginCreate {
	oc.mutation.SetCreatedAt(t)
	return oc
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (oc *OriginCreate) SetNillableCreatedAt(t *time.Time) *OriginCreate {
	if t != nil {
		oc.SetCreatedAt(*t)
	}
	return oc
}

// SetUpdatedAt sets the "updated_at" field.
func (oc *OriginCreate) SetUpdatedAt(t time.Time) *OriginCreate {
	oc.mutation.SetUpdatedAt(t)
	return oc
}

// SetNillableUpdatedAt sets the "updated_at" field if the given value is not nil.
func (oc *OriginCreate) SetNillableUpdatedAt(t *time.Time) *OriginCreate {
	if t != nil {
		oc.SetUpdatedAt(*t)
	}
	return oc
}

// SetName sets the "name" field.
func (oc *OriginCreate) SetName(s string) *OriginCreate {
	oc.mutation.SetName(s)
	return oc
}

// SetWeight sets the "weight" field.
func (oc *OriginCreate) SetWeight(i int) *OriginCreate {
	oc.mutation.SetWeight(i)
	return oc
}

// SetNillableWeight sets the "weight" field if the given value is not nil.
func (oc *OriginCreate) SetNillableWeight(i *int) *OriginCreate {
	if i != nil {
		oc.SetWeight(*i)
	}
	return oc
}

// SetTarget sets the "target" field.
func (oc *OriginCreate) SetTarget(s string) *OriginCreate {
	oc.mutation.SetTarget(s)
	return oc
}

// SetPortNumber sets the "port_number" field.
func (oc *OriginCreate) SetPortNumber(i int) *OriginCreate {
	oc.mutation.SetPortNumber(i)
	return oc
}

// SetActive sets the "active" field.
func (oc *OriginCreate) SetActive(b bool) *OriginCreate {
	oc.mutation.SetActive(b)
	return oc
}

// SetNillableActive sets the "active" field if the given value is not nil.
func (oc *OriginCreate) SetNillableActive(b *bool) *OriginCreate {
	if b != nil {
		oc.SetActive(*b)
	}
	return oc
}

// SetPoolID sets the "pool_id" field.
func (oc *OriginCreate) SetPoolID(gi gidx.PrefixedID) *OriginCreate {
	oc.mutation.SetPoolID(gi)
	return oc
}

// SetID sets the "id" field.
func (oc *OriginCreate) SetID(gi gidx.PrefixedID) *OriginCreate {
	oc.mutation.SetID(gi)
	return oc
}

// SetNillableID sets the "id" field if the given value is not nil.
func (oc *OriginCreate) SetNillableID(gi *gidx.PrefixedID) *OriginCreate {
	if gi != nil {
		oc.SetID(*gi)
	}
	return oc
}

// SetPool sets the "pool" edge to the Pool entity.
func (oc *OriginCreate) SetPool(p *Pool) *OriginCreate {
	return oc.SetPoolID(p.ID)
}

// Mutation returns the OriginMutation object of the builder.
func (oc *OriginCreate) Mutation() *OriginMutation {
	return oc.mutation
}

// Save creates the Origin in the database.
func (oc *OriginCreate) Save(ctx context.Context) (*Origin, error) {
	oc.defaults()
	return withHooks(ctx, oc.sqlSave, oc.mutation, oc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (oc *OriginCreate) SaveX(ctx context.Context) *Origin {
	v, err := oc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (oc *OriginCreate) Exec(ctx context.Context) error {
	_, err := oc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (oc *OriginCreate) ExecX(ctx context.Context) {
	if err := oc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (oc *OriginCreate) defaults() {
	if _, ok := oc.mutation.CreatedAt(); !ok {
		v := origin.DefaultCreatedAt()
		oc.mutation.SetCreatedAt(v)
	}
	if _, ok := oc.mutation.UpdatedAt(); !ok {
		v := origin.DefaultUpdatedAt()
		oc.mutation.SetUpdatedAt(v)
	}
	if _, ok := oc.mutation.Weight(); !ok {
		v := origin.DefaultWeight
		oc.mutation.SetWeight(v)
	}
	if _, ok := oc.mutation.Active(); !ok {
		v := origin.DefaultActive
		oc.mutation.SetActive(v)
	}
	if _, ok := oc.mutation.ID(); !ok {
		v := origin.DefaultID()
		oc.mutation.SetID(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (oc *OriginCreate) check() error {
	if _, ok := oc.mutation.CreatedAt(); !ok {
		return &ValidationError{Name: "created_at", err: errors.New(`generated: missing required field "Origin.created_at"`)}
	}
	if _, ok := oc.mutation.UpdatedAt(); !ok {
		return &ValidationError{Name: "updated_at", err: errors.New(`generated: missing required field "Origin.updated_at"`)}
	}
	if _, ok := oc.mutation.Name(); !ok {
		return &ValidationError{Name: "name", err: errors.New(`generated: missing required field "Origin.name"`)}
	}
	if v, ok := oc.mutation.Name(); ok {
		if err := origin.NameValidator(v); err != nil {
			return &ValidationError{Name: "name", err: fmt.Errorf(`generated: validator failed for field "Origin.name": %w`, err)}
		}
	}
	if _, ok := oc.mutation.Weight(); !ok {
		return &ValidationError{Name: "weight", err: errors.New(`generated: missing required field "Origin.weight"`)}
	}
	if _, ok := oc.mutation.Target(); !ok {
		return &ValidationError{Name: "target", err: errors.New(`generated: missing required field "Origin.target"`)}
	}
	if v, ok := oc.mutation.Target(); ok {
		if err := origin.TargetValidator(v); err != nil {
			return &ValidationError{Name: "target", err: fmt.Errorf(`generated: validator failed for field "Origin.target": %w`, err)}
		}
	}
	if _, ok := oc.mutation.PortNumber(); !ok {
		return &ValidationError{Name: "port_number", err: errors.New(`generated: missing required field "Origin.port_number"`)}
	}
	if v, ok := oc.mutation.PortNumber(); ok {
		if err := origin.PortNumberValidator(v); err != nil {
			return &ValidationError{Name: "port_number", err: fmt.Errorf(`generated: validator failed for field "Origin.port_number": %w`, err)}
		}
	}
	if _, ok := oc.mutation.Active(); !ok {
		return &ValidationError{Name: "active", err: errors.New(`generated: missing required field "Origin.active"`)}
	}
	if _, ok := oc.mutation.PoolID(); !ok {
		return &ValidationError{Name: "pool_id", err: errors.New(`generated: missing required field "Origin.pool_id"`)}
	}
	if v, ok := oc.mutation.PoolID(); ok {
		if err := origin.PoolIDValidator(string(v)); err != nil {
			return &ValidationError{Name: "pool_id", err: fmt.Errorf(`generated: validator failed for field "Origin.pool_id": %w`, err)}
		}
	}
	if _, ok := oc.mutation.PoolID(); !ok {
		return &ValidationError{Name: "pool", err: errors.New(`generated: missing required edge "Origin.pool"`)}
	}
	return nil
}

func (oc *OriginCreate) sqlSave(ctx context.Context) (*Origin, error) {
	if err := oc.check(); err != nil {
		return nil, err
	}
	_node, _spec := oc.createSpec()
	if err := sqlgraph.CreateNode(ctx, oc.driver, _spec); err != nil {
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
	oc.mutation.id = &_node.ID
	oc.mutation.done = true
	return _node, nil
}

func (oc *OriginCreate) createSpec() (*Origin, *sqlgraph.CreateSpec) {
	var (
		_node = &Origin{config: oc.config}
		_spec = sqlgraph.NewCreateSpec(origin.Table, sqlgraph.NewFieldSpec(origin.FieldID, field.TypeString))
	)
	if id, ok := oc.mutation.ID(); ok {
		_node.ID = id
		_spec.ID.Value = &id
	}
	if value, ok := oc.mutation.CreatedAt(); ok {
		_spec.SetField(origin.FieldCreatedAt, field.TypeTime, value)
		_node.CreatedAt = value
	}
	if value, ok := oc.mutation.UpdatedAt(); ok {
		_spec.SetField(origin.FieldUpdatedAt, field.TypeTime, value)
		_node.UpdatedAt = value
	}
	if value, ok := oc.mutation.Name(); ok {
		_spec.SetField(origin.FieldName, field.TypeString, value)
		_node.Name = value
	}
	if value, ok := oc.mutation.Weight(); ok {
		_spec.SetField(origin.FieldWeight, field.TypeInt, value)
		_node.Weight = value
	}
	if value, ok := oc.mutation.Target(); ok {
		_spec.SetField(origin.FieldTarget, field.TypeString, value)
		_node.Target = value
	}
	if value, ok := oc.mutation.PortNumber(); ok {
		_spec.SetField(origin.FieldPortNumber, field.TypeInt, value)
		_node.PortNumber = value
	}
	if value, ok := oc.mutation.Active(); ok {
		_spec.SetField(origin.FieldActive, field.TypeBool, value)
		_node.Active = value
	}
	if nodes := oc.mutation.PoolIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   origin.PoolTable,
			Columns: []string{origin.PoolColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(pool.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.PoolID = nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// OriginCreateBulk is the builder for creating many Origin entities in bulk.
type OriginCreateBulk struct {
	config
	err      error
	builders []*OriginCreate
}

// Save creates the Origin entities in the database.
func (ocb *OriginCreateBulk) Save(ctx context.Context) ([]*Origin, error) {
	if ocb.err != nil {
		return nil, ocb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(ocb.builders))
	nodes := make([]*Origin, len(ocb.builders))
	mutators := make([]Mutator, len(ocb.builders))
	for i := range ocb.builders {
		func(i int, root context.Context) {
			builder := ocb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*OriginMutation)
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
					_, err = mutators[i+1].Mutate(root, ocb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, ocb.driver, spec); err != nil {
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
		if _, err := mutators[0].Mutate(ctx, ocb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (ocb *OriginCreateBulk) SaveX(ctx context.Context) []*Origin {
	v, err := ocb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (ocb *OriginCreateBulk) Exec(ctx context.Context) error {
	_, err := ocb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ocb *OriginCreateBulk) ExecX(ctx context.Context) {
	if err := ocb.Exec(ctx); err != nil {
		panic(err)
	}
}
