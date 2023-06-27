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
	"time"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"go.infratographer.com/load-balancer-api/internal/ent/generated/loadbalancer"
	"go.infratographer.com/load-balancer-api/internal/ent/generated/loadbalancerannotation"
	"go.infratographer.com/load-balancer-api/internal/ent/generated/loadbalancerstatus"
	"go.infratographer.com/load-balancer-api/internal/ent/generated/port"
	"go.infratographer.com/load-balancer-api/internal/ent/generated/provider"
	"go.infratographer.com/x/gidx"
)

// LoadBalancerCreate is the builder for creating a LoadBalancer entity.
type LoadBalancerCreate struct {
	config
	mutation *LoadBalancerMutation
	hooks    []Hook
}

// SetCreatedAt sets the "created_at" field.
func (lbc *LoadBalancerCreate) SetCreatedAt(t time.Time) *LoadBalancerCreate {
	lbc.mutation.SetCreatedAt(t)
	return lbc
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (lbc *LoadBalancerCreate) SetNillableCreatedAt(t *time.Time) *LoadBalancerCreate {
	if t != nil {
		lbc.SetCreatedAt(*t)
	}
	return lbc
}

// SetUpdatedAt sets the "updated_at" field.
func (lbc *LoadBalancerCreate) SetUpdatedAt(t time.Time) *LoadBalancerCreate {
	lbc.mutation.SetUpdatedAt(t)
	return lbc
}

// SetNillableUpdatedAt sets the "updated_at" field if the given value is not nil.
func (lbc *LoadBalancerCreate) SetNillableUpdatedAt(t *time.Time) *LoadBalancerCreate {
	if t != nil {
		lbc.SetUpdatedAt(*t)
	}
	return lbc
}

// SetName sets the "name" field.
func (lbc *LoadBalancerCreate) SetName(s string) *LoadBalancerCreate {
	lbc.mutation.SetName(s)
	return lbc
}

// SetOwnerID sets the "owner_id" field.
func (lbc *LoadBalancerCreate) SetOwnerID(gi gidx.PrefixedID) *LoadBalancerCreate {
	lbc.mutation.SetOwnerID(gi)
	return lbc
}

// SetLocationID sets the "location_id" field.
func (lbc *LoadBalancerCreate) SetLocationID(gi gidx.PrefixedID) *LoadBalancerCreate {
	lbc.mutation.SetLocationID(gi)
	return lbc
}

// SetProviderID sets the "provider_id" field.
func (lbc *LoadBalancerCreate) SetProviderID(gi gidx.PrefixedID) *LoadBalancerCreate {
	lbc.mutation.SetProviderID(gi)
	return lbc
}

// SetIPID sets the "ip_id" field.
func (lbc *LoadBalancerCreate) SetIPID(gi gidx.PrefixedID) *LoadBalancerCreate {
	lbc.mutation.SetIPID(gi)
	return lbc
}

// SetNillableIPID sets the "ip_id" field if the given value is not nil.
func (lbc *LoadBalancerCreate) SetNillableIPID(gi *gidx.PrefixedID) *LoadBalancerCreate {
	if gi != nil {
		lbc.SetIPID(*gi)
	}
	return lbc
}

// SetID sets the "id" field.
func (lbc *LoadBalancerCreate) SetID(gi gidx.PrefixedID) *LoadBalancerCreate {
	lbc.mutation.SetID(gi)
	return lbc
}

// SetNillableID sets the "id" field if the given value is not nil.
func (lbc *LoadBalancerCreate) SetNillableID(gi *gidx.PrefixedID) *LoadBalancerCreate {
	if gi != nil {
		lbc.SetID(*gi)
	}
	return lbc
}

// AddAnnotationIDs adds the "annotations" edge to the LoadBalancerAnnotation entity by IDs.
func (lbc *LoadBalancerCreate) AddAnnotationIDs(ids ...gidx.PrefixedID) *LoadBalancerCreate {
	lbc.mutation.AddAnnotationIDs(ids...)
	return lbc
}

// AddAnnotations adds the "annotations" edges to the LoadBalancerAnnotation entity.
func (lbc *LoadBalancerCreate) AddAnnotations(l ...*LoadBalancerAnnotation) *LoadBalancerCreate {
	ids := make([]gidx.PrefixedID, len(l))
	for i := range l {
		ids[i] = l[i].ID
	}
	return lbc.AddAnnotationIDs(ids...)
}

// AddStatusIDs adds the "statuses" edge to the LoadBalancerStatus entity by IDs.
func (lbc *LoadBalancerCreate) AddStatusIDs(ids ...gidx.PrefixedID) *LoadBalancerCreate {
	lbc.mutation.AddStatusIDs(ids...)
	return lbc
}

// AddStatuses adds the "statuses" edges to the LoadBalancerStatus entity.
func (lbc *LoadBalancerCreate) AddStatuses(l ...*LoadBalancerStatus) *LoadBalancerCreate {
	ids := make([]gidx.PrefixedID, len(l))
	for i := range l {
		ids[i] = l[i].ID
	}
	return lbc.AddStatusIDs(ids...)
}

// AddPortIDs adds the "ports" edge to the Port entity by IDs.
func (lbc *LoadBalancerCreate) AddPortIDs(ids ...gidx.PrefixedID) *LoadBalancerCreate {
	lbc.mutation.AddPortIDs(ids...)
	return lbc
}

// AddPorts adds the "ports" edges to the Port entity.
func (lbc *LoadBalancerCreate) AddPorts(p ...*Port) *LoadBalancerCreate {
	ids := make([]gidx.PrefixedID, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return lbc.AddPortIDs(ids...)
}

// SetProvider sets the "provider" edge to the Provider entity.
func (lbc *LoadBalancerCreate) SetProvider(p *Provider) *LoadBalancerCreate {
	return lbc.SetProviderID(p.ID)
}

// Mutation returns the LoadBalancerMutation object of the builder.
func (lbc *LoadBalancerCreate) Mutation() *LoadBalancerMutation {
	return lbc.mutation
}

// Save creates the LoadBalancer in the database.
func (lbc *LoadBalancerCreate) Save(ctx context.Context) (*LoadBalancer, error) {
	lbc.defaults()
	return withHooks(ctx, lbc.sqlSave, lbc.mutation, lbc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (lbc *LoadBalancerCreate) SaveX(ctx context.Context) *LoadBalancer {
	v, err := lbc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (lbc *LoadBalancerCreate) Exec(ctx context.Context) error {
	_, err := lbc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (lbc *LoadBalancerCreate) ExecX(ctx context.Context) {
	if err := lbc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (lbc *LoadBalancerCreate) defaults() {
	if _, ok := lbc.mutation.CreatedAt(); !ok {
		v := loadbalancer.DefaultCreatedAt()
		lbc.mutation.SetCreatedAt(v)
	}
	if _, ok := lbc.mutation.UpdatedAt(); !ok {
		v := loadbalancer.DefaultUpdatedAt()
		lbc.mutation.SetUpdatedAt(v)
	}
	if _, ok := lbc.mutation.ID(); !ok {
		v := loadbalancer.DefaultID()
		lbc.mutation.SetID(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (lbc *LoadBalancerCreate) check() error {
	if _, ok := lbc.mutation.CreatedAt(); !ok {
		return &ValidationError{Name: "created_at", err: errors.New(`generated: missing required field "LoadBalancer.created_at"`)}
	}
	if _, ok := lbc.mutation.UpdatedAt(); !ok {
		return &ValidationError{Name: "updated_at", err: errors.New(`generated: missing required field "LoadBalancer.updated_at"`)}
	}
	if _, ok := lbc.mutation.Name(); !ok {
		return &ValidationError{Name: "name", err: errors.New(`generated: missing required field "LoadBalancer.name"`)}
	}
	if v, ok := lbc.mutation.Name(); ok {
		if err := loadbalancer.NameValidator(v); err != nil {
			return &ValidationError{Name: "name", err: fmt.Errorf(`generated: validator failed for field "LoadBalancer.name": %w`, err)}
		}
	}
	if _, ok := lbc.mutation.OwnerID(); !ok {
		return &ValidationError{Name: "owner_id", err: errors.New(`generated: missing required field "LoadBalancer.owner_id"`)}
	}
	if _, ok := lbc.mutation.LocationID(); !ok {
		return &ValidationError{Name: "location_id", err: errors.New(`generated: missing required field "LoadBalancer.location_id"`)}
	}
	if v, ok := lbc.mutation.LocationID(); ok {
		if err := loadbalancer.LocationIDValidator(string(v)); err != nil {
			return &ValidationError{Name: "location_id", err: fmt.Errorf(`generated: validator failed for field "LoadBalancer.location_id": %w`, err)}
		}
	}
	if _, ok := lbc.mutation.ProviderID(); !ok {
		return &ValidationError{Name: "provider_id", err: errors.New(`generated: missing required field "LoadBalancer.provider_id"`)}
	}
	if v, ok := lbc.mutation.ProviderID(); ok {
		if err := loadbalancer.ProviderIDValidator(string(v)); err != nil {
			return &ValidationError{Name: "provider_id", err: fmt.Errorf(`generated: validator failed for field "LoadBalancer.provider_id": %w`, err)}
		}
	}
	if _, ok := lbc.mutation.ProviderID(); !ok {
		return &ValidationError{Name: "provider", err: errors.New(`generated: missing required edge "LoadBalancer.provider"`)}
	}
	return nil
}

func (lbc *LoadBalancerCreate) sqlSave(ctx context.Context) (*LoadBalancer, error) {
	if err := lbc.check(); err != nil {
		return nil, err
	}
	_node, _spec := lbc.createSpec()
	if err := sqlgraph.CreateNode(ctx, lbc.driver, _spec); err != nil {
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
	lbc.mutation.id = &_node.ID
	lbc.mutation.done = true
	return _node, nil
}

func (lbc *LoadBalancerCreate) createSpec() (*LoadBalancer, *sqlgraph.CreateSpec) {
	var (
		_node = &LoadBalancer{config: lbc.config}
		_spec = sqlgraph.NewCreateSpec(loadbalancer.Table, sqlgraph.NewFieldSpec(loadbalancer.FieldID, field.TypeString))
	)
	if id, ok := lbc.mutation.ID(); ok {
		_node.ID = id
		_spec.ID.Value = &id
	}
	if value, ok := lbc.mutation.CreatedAt(); ok {
		_spec.SetField(loadbalancer.FieldCreatedAt, field.TypeTime, value)
		_node.CreatedAt = value
	}
	if value, ok := lbc.mutation.UpdatedAt(); ok {
		_spec.SetField(loadbalancer.FieldUpdatedAt, field.TypeTime, value)
		_node.UpdatedAt = value
	}
	if value, ok := lbc.mutation.Name(); ok {
		_spec.SetField(loadbalancer.FieldName, field.TypeString, value)
		_node.Name = value
	}
	if value, ok := lbc.mutation.OwnerID(); ok {
		_spec.SetField(loadbalancer.FieldOwnerID, field.TypeString, value)
		_node.OwnerID = value
	}
	if value, ok := lbc.mutation.LocationID(); ok {
		_spec.SetField(loadbalancer.FieldLocationID, field.TypeString, value)
		_node.LocationID = value
	}
	if value, ok := lbc.mutation.IPID(); ok {
		_spec.SetField(loadbalancer.FieldIPID, field.TypeString, value)
		_node.IPID = value
	}
	if nodes := lbc.mutation.AnnotationsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   loadbalancer.AnnotationsTable,
			Columns: []string{loadbalancer.AnnotationsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(loadbalancerannotation.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := lbc.mutation.StatusesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   loadbalancer.StatusesTable,
			Columns: []string{loadbalancer.StatusesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(loadbalancerstatus.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := lbc.mutation.PortsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   loadbalancer.PortsTable,
			Columns: []string{loadbalancer.PortsColumn},
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
	if nodes := lbc.mutation.ProviderIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   loadbalancer.ProviderTable,
			Columns: []string{loadbalancer.ProviderColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(provider.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.ProviderID = nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// LoadBalancerCreateBulk is the builder for creating many LoadBalancer entities in bulk.
type LoadBalancerCreateBulk struct {
	config
	builders []*LoadBalancerCreate
}

// Save creates the LoadBalancer entities in the database.
func (lbcb *LoadBalancerCreateBulk) Save(ctx context.Context) ([]*LoadBalancer, error) {
	specs := make([]*sqlgraph.CreateSpec, len(lbcb.builders))
	nodes := make([]*LoadBalancer, len(lbcb.builders))
	mutators := make([]Mutator, len(lbcb.builders))
	for i := range lbcb.builders {
		func(i int, root context.Context) {
			builder := lbcb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*LoadBalancerMutation)
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
					_, err = mutators[i+1].Mutate(root, lbcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, lbcb.driver, spec); err != nil {
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
		if _, err := mutators[0].Mutate(ctx, lbcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (lbcb *LoadBalancerCreateBulk) SaveX(ctx context.Context) []*LoadBalancer {
	v, err := lbcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (lbcb *LoadBalancerCreateBulk) Exec(ctx context.Context) error {
	_, err := lbcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (lbcb *LoadBalancerCreateBulk) ExecX(ctx context.Context) {
	if err := lbcb.Exec(ctx); err != nil {
		panic(err)
	}
}
