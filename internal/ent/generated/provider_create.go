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
	"go.infratographer.com/load-balancer-api/internal/ent/generated/loadbalancer"
	"go.infratographer.com/load-balancer-api/internal/ent/generated/provider"
	"go.infratographer.com/x/gidx"
)

// ProviderCreate is the builder for creating a Provider entity.
type ProviderCreate struct {
	config
	mutation *ProviderMutation
	hooks    []Hook
}

// SetCreatedAt sets the "created_at" field.
func (pc *ProviderCreate) SetCreatedAt(t time.Time) *ProviderCreate {
	pc.mutation.SetCreatedAt(t)
	return pc
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (pc *ProviderCreate) SetNillableCreatedAt(t *time.Time) *ProviderCreate {
	if t != nil {
		pc.SetCreatedAt(*t)
	}
	return pc
}

// SetUpdatedAt sets the "updated_at" field.
func (pc *ProviderCreate) SetUpdatedAt(t time.Time) *ProviderCreate {
	pc.mutation.SetUpdatedAt(t)
	return pc
}

// SetNillableUpdatedAt sets the "updated_at" field if the given value is not nil.
func (pc *ProviderCreate) SetNillableUpdatedAt(t *time.Time) *ProviderCreate {
	if t != nil {
		pc.SetUpdatedAt(*t)
	}
	return pc
}

// SetDeletedAt sets the "deleted_at" field.
func (pc *ProviderCreate) SetDeletedAt(t time.Time) *ProviderCreate {
	pc.mutation.SetDeletedAt(t)
	return pc
}

// SetNillableDeletedAt sets the "deleted_at" field if the given value is not nil.
func (pc *ProviderCreate) SetNillableDeletedAt(t *time.Time) *ProviderCreate {
	if t != nil {
		pc.SetDeletedAt(*t)
	}
	return pc
}

// SetDeletedBy sets the "deleted_by" field.
func (pc *ProviderCreate) SetDeletedBy(s string) *ProviderCreate {
	pc.mutation.SetDeletedBy(s)
	return pc
}

// SetNillableDeletedBy sets the "deleted_by" field if the given value is not nil.
func (pc *ProviderCreate) SetNillableDeletedBy(s *string) *ProviderCreate {
	if s != nil {
		pc.SetDeletedBy(*s)
	}
	return pc
}

// SetCreatedBy sets the "created_by" field.
func (pc *ProviderCreate) SetCreatedBy(s string) *ProviderCreate {
	pc.mutation.SetCreatedBy(s)
	return pc
}

// SetNillableCreatedBy sets the "created_by" field if the given value is not nil.
func (pc *ProviderCreate) SetNillableCreatedBy(s *string) *ProviderCreate {
	if s != nil {
		pc.SetCreatedBy(*s)
	}
	return pc
}

// SetUpdatedBy sets the "updated_by" field.
func (pc *ProviderCreate) SetUpdatedBy(s string) *ProviderCreate {
	pc.mutation.SetUpdatedBy(s)
	return pc
}

// SetNillableUpdatedBy sets the "updated_by" field if the given value is not nil.
func (pc *ProviderCreate) SetNillableUpdatedBy(s *string) *ProviderCreate {
	if s != nil {
		pc.SetUpdatedBy(*s)
	}
	return pc
}

// SetName sets the "name" field.
func (pc *ProviderCreate) SetName(s string) *ProviderCreate {
	pc.mutation.SetName(s)
	return pc
}

// SetOwnerID sets the "owner_id" field.
func (pc *ProviderCreate) SetOwnerID(gi gidx.PrefixedID) *ProviderCreate {
	pc.mutation.SetOwnerID(gi)
	return pc
}

// SetID sets the "id" field.
func (pc *ProviderCreate) SetID(gi gidx.PrefixedID) *ProviderCreate {
	pc.mutation.SetID(gi)
	return pc
}

// SetNillableID sets the "id" field if the given value is not nil.
func (pc *ProviderCreate) SetNillableID(gi *gidx.PrefixedID) *ProviderCreate {
	if gi != nil {
		pc.SetID(*gi)
	}
	return pc
}

// AddLoadBalancerIDs adds the "load_balancers" edge to the LoadBalancer entity by IDs.
func (pc *ProviderCreate) AddLoadBalancerIDs(ids ...gidx.PrefixedID) *ProviderCreate {
	pc.mutation.AddLoadBalancerIDs(ids...)
	return pc
}

// AddLoadBalancers adds the "load_balancers" edges to the LoadBalancer entity.
func (pc *ProviderCreate) AddLoadBalancers(l ...*LoadBalancer) *ProviderCreate {
	ids := make([]gidx.PrefixedID, len(l))
	for i := range l {
		ids[i] = l[i].ID
	}
	return pc.AddLoadBalancerIDs(ids...)
}

// Mutation returns the ProviderMutation object of the builder.
func (pc *ProviderCreate) Mutation() *ProviderMutation {
	return pc.mutation
}

// Save creates the Provider in the database.
func (pc *ProviderCreate) Save(ctx context.Context) (*Provider, error) {
	if err := pc.defaults(); err != nil {
		return nil, err
	}
	return withHooks(ctx, pc.sqlSave, pc.mutation, pc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (pc *ProviderCreate) SaveX(ctx context.Context) *Provider {
	v, err := pc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (pc *ProviderCreate) Exec(ctx context.Context) error {
	_, err := pc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (pc *ProviderCreate) ExecX(ctx context.Context) {
	if err := pc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (pc *ProviderCreate) defaults() error {
	if _, ok := pc.mutation.CreatedAt(); !ok {
		if provider.DefaultCreatedAt == nil {
			return fmt.Errorf("generated: uninitialized provider.DefaultCreatedAt (forgotten import generated/runtime?)")
		}
		v := provider.DefaultCreatedAt()
		pc.mutation.SetCreatedAt(v)
	}
	if _, ok := pc.mutation.UpdatedAt(); !ok {
		if provider.DefaultUpdatedAt == nil {
			return fmt.Errorf("generated: uninitialized provider.DefaultUpdatedAt (forgotten import generated/runtime?)")
		}
		v := provider.DefaultUpdatedAt()
		pc.mutation.SetUpdatedAt(v)
	}
	if _, ok := pc.mutation.ID(); !ok {
		if provider.DefaultID == nil {
			return fmt.Errorf("generated: uninitialized provider.DefaultID (forgotten import generated/runtime?)")
		}
		v := provider.DefaultID()
		pc.mutation.SetID(v)
	}
	return nil
}

// check runs all checks and user-defined validators on the builder.
func (pc *ProviderCreate) check() error {
	if _, ok := pc.mutation.CreatedAt(); !ok {
		return &ValidationError{Name: "created_at", err: errors.New(`generated: missing required field "Provider.created_at"`)}
	}
	if _, ok := pc.mutation.UpdatedAt(); !ok {
		return &ValidationError{Name: "updated_at", err: errors.New(`generated: missing required field "Provider.updated_at"`)}
	}
	if _, ok := pc.mutation.Name(); !ok {
		return &ValidationError{Name: "name", err: errors.New(`generated: missing required field "Provider.name"`)}
	}
	if v, ok := pc.mutation.Name(); ok {
		if err := provider.NameValidator(v); err != nil {
			return &ValidationError{Name: "name", err: fmt.Errorf(`generated: validator failed for field "Provider.name": %w`, err)}
		}
	}
	if _, ok := pc.mutation.OwnerID(); !ok {
		return &ValidationError{Name: "owner_id", err: errors.New(`generated: missing required field "Provider.owner_id"`)}
	}
	if v, ok := pc.mutation.OwnerID(); ok {
		if err := provider.OwnerIDValidator(string(v)); err != nil {
			return &ValidationError{Name: "owner_id", err: fmt.Errorf(`generated: validator failed for field "Provider.owner_id": %w`, err)}
		}
	}
	return nil
}

func (pc *ProviderCreate) sqlSave(ctx context.Context) (*Provider, error) {
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

func (pc *ProviderCreate) createSpec() (*Provider, *sqlgraph.CreateSpec) {
	var (
		_node = &Provider{config: pc.config}
		_spec = sqlgraph.NewCreateSpec(provider.Table, sqlgraph.NewFieldSpec(provider.FieldID, field.TypeString))
	)
	if id, ok := pc.mutation.ID(); ok {
		_node.ID = id
		_spec.ID.Value = &id
	}
	if value, ok := pc.mutation.CreatedAt(); ok {
		_spec.SetField(provider.FieldCreatedAt, field.TypeTime, value)
		_node.CreatedAt = value
	}
	if value, ok := pc.mutation.UpdatedAt(); ok {
		_spec.SetField(provider.FieldUpdatedAt, field.TypeTime, value)
		_node.UpdatedAt = value
	}
	if value, ok := pc.mutation.DeletedAt(); ok {
		_spec.SetField(provider.FieldDeletedAt, field.TypeTime, value)
		_node.DeletedAt = value
	}
	if value, ok := pc.mutation.DeletedBy(); ok {
		_spec.SetField(provider.FieldDeletedBy, field.TypeString, value)
		_node.DeletedBy = value
	}
	if value, ok := pc.mutation.CreatedBy(); ok {
		_spec.SetField(provider.FieldCreatedBy, field.TypeString, value)
		_node.CreatedBy = value
	}
	if value, ok := pc.mutation.UpdatedBy(); ok {
		_spec.SetField(provider.FieldUpdatedBy, field.TypeString, value)
		_node.UpdatedBy = value
	}
	if value, ok := pc.mutation.Name(); ok {
		_spec.SetField(provider.FieldName, field.TypeString, value)
		_node.Name = value
	}
	if value, ok := pc.mutation.OwnerID(); ok {
		_spec.SetField(provider.FieldOwnerID, field.TypeString, value)
		_node.OwnerID = value
	}
	if nodes := pc.mutation.LoadBalancersIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   provider.LoadBalancersTable,
			Columns: []string{provider.LoadBalancersColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(loadbalancer.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// ProviderCreateBulk is the builder for creating many Provider entities in bulk.
type ProviderCreateBulk struct {
	config
	err      error
	builders []*ProviderCreate
}

// Save creates the Provider entities in the database.
func (pcb *ProviderCreateBulk) Save(ctx context.Context) ([]*Provider, error) {
	if pcb.err != nil {
		return nil, pcb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(pcb.builders))
	nodes := make([]*Provider, len(pcb.builders))
	mutators := make([]Mutator, len(pcb.builders))
	for i := range pcb.builders {
		func(i int, root context.Context) {
			builder := pcb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*ProviderMutation)
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
func (pcb *ProviderCreateBulk) SaveX(ctx context.Context) []*Provider {
	v, err := pcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (pcb *ProviderCreateBulk) Exec(ctx context.Context) error {
	_, err := pcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (pcb *ProviderCreateBulk) ExecX(ctx context.Context) {
	if err := pcb.Exec(ctx); err != nil {
		panic(err)
	}
}
