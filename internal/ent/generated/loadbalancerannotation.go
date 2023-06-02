// Copyright Infratographer, Inc. and/or licensed to Infratographer, Inc. under one
// or more contributor license agreements. Licensed under the Elastic License 2.0;
// you may not use this file except in compliance with the Elastic License 2.0.
//
// Code generated by entc, DO NOT EDIT.

package generated

import (
	"fmt"
	"strings"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"go.infratographer.com/load-balancer-api/internal/ent/generated/loadbalancer"
	"go.infratographer.com/load-balancer-api/internal/ent/generated/loadbalancerannotation"
	"go.infratographer.com/x/gidx"
)

// LoadBalancerAnnotation is the model entity for the LoadBalancerAnnotation schema.
type LoadBalancerAnnotation struct {
	config `json:"-"`
	// ID of the ent.
	ID gidx.PrefixedID `json:"id,omitempty"`
	// CreatedAt holds the value of the "created_at" field.
	CreatedAt time.Time `json:"created_at,omitempty"`
	// UpdatedAt holds the value of the "updated_at" field.
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	// LoadBalancerID holds the value of the "load_balancer_id" field.
	LoadBalancerID gidx.PrefixedID `json:"load_balancer_id,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the LoadBalancerAnnotationQuery when eager-loading is set.
	Edges        LoadBalancerAnnotationEdges `json:"edges"`
	selectValues sql.SelectValues
}

// LoadBalancerAnnotationEdges holds the relations/edges for other nodes in the graph.
type LoadBalancerAnnotationEdges struct {
	// LoadBalancer holds the value of the load_balancer edge.
	LoadBalancer *LoadBalancer `json:"load_balancer,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [1]bool
	// totalCount holds the count of the edges above.
	totalCount [1]map[string]int
}

// LoadBalancerOrErr returns the LoadBalancer value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e LoadBalancerAnnotationEdges) LoadBalancerOrErr() (*LoadBalancer, error) {
	if e.loadedTypes[0] {
		if e.LoadBalancer == nil {
			// Edge was loaded but was not found.
			return nil, &NotFoundError{label: loadbalancer.Label}
		}
		return e.LoadBalancer, nil
	}
	return nil, &NotLoadedError{edge: "load_balancer"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*LoadBalancerAnnotation) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case loadbalancerannotation.FieldID, loadbalancerannotation.FieldLoadBalancerID:
			values[i] = new(gidx.PrefixedID)
		case loadbalancerannotation.FieldCreatedAt, loadbalancerannotation.FieldUpdatedAt:
			values[i] = new(sql.NullTime)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the LoadBalancerAnnotation fields.
func (lba *LoadBalancerAnnotation) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case loadbalancerannotation.FieldID:
			if value, ok := values[i].(*gidx.PrefixedID); !ok {
				return fmt.Errorf("unexpected type %T for field id", values[i])
			} else if value != nil {
				lba.ID = *value
			}
		case loadbalancerannotation.FieldCreatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field created_at", values[i])
			} else if value.Valid {
				lba.CreatedAt = value.Time
			}
		case loadbalancerannotation.FieldUpdatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field updated_at", values[i])
			} else if value.Valid {
				lba.UpdatedAt = value.Time
			}
		case loadbalancerannotation.FieldLoadBalancerID:
			if value, ok := values[i].(*gidx.PrefixedID); !ok {
				return fmt.Errorf("unexpected type %T for field load_balancer_id", values[i])
			} else if value != nil {
				lba.LoadBalancerID = *value
			}
		default:
			lba.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the LoadBalancerAnnotation.
// This includes values selected through modifiers, order, etc.
func (lba *LoadBalancerAnnotation) Value(name string) (ent.Value, error) {
	return lba.selectValues.Get(name)
}

// QueryLoadBalancer queries the "load_balancer" edge of the LoadBalancerAnnotation entity.
func (lba *LoadBalancerAnnotation) QueryLoadBalancer() *LoadBalancerQuery {
	return NewLoadBalancerAnnotationClient(lba.config).QueryLoadBalancer(lba)
}

// Update returns a builder for updating this LoadBalancerAnnotation.
// Note that you need to call LoadBalancerAnnotation.Unwrap() before calling this method if this LoadBalancerAnnotation
// was returned from a transaction, and the transaction was committed or rolled back.
func (lba *LoadBalancerAnnotation) Update() *LoadBalancerAnnotationUpdateOne {
	return NewLoadBalancerAnnotationClient(lba.config).UpdateOne(lba)
}

// Unwrap unwraps the LoadBalancerAnnotation entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (lba *LoadBalancerAnnotation) Unwrap() *LoadBalancerAnnotation {
	_tx, ok := lba.config.driver.(*txDriver)
	if !ok {
		panic("generated: LoadBalancerAnnotation is not a transactional entity")
	}
	lba.config.driver = _tx.drv
	return lba
}

// String implements the fmt.Stringer.
func (lba *LoadBalancerAnnotation) String() string {
	var builder strings.Builder
	builder.WriteString("LoadBalancerAnnotation(")
	builder.WriteString(fmt.Sprintf("id=%v, ", lba.ID))
	builder.WriteString("created_at=")
	builder.WriteString(lba.CreatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("updated_at=")
	builder.WriteString(lba.UpdatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("load_balancer_id=")
	builder.WriteString(fmt.Sprintf("%v", lba.LoadBalancerID))
	builder.WriteByte(')')
	return builder.String()
}

// IsEntity implement fedruntime.Entity
func (lba LoadBalancerAnnotation) IsEntity() {}

// LoadBalancerAnnotations is a parsable slice of LoadBalancerAnnotation.
type LoadBalancerAnnotations []*LoadBalancerAnnotation
