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
	"go.infratographer.com/load-balancer-api/internal/ent/generated/provider"
	"go.infratographer.com/x/gidx"
)

// Representation of a load balancer.
type LoadBalancer struct {
	config `json:"-"`
	// ID of the ent.
	// The ID for the load balancer.
	ID gidx.PrefixedID `json:"id,omitempty"`
	// CreatedAt holds the value of the "created_at" field.
	CreatedAt time.Time `json:"created_at,omitempty"`
	// UpdatedAt holds the value of the "updated_at" field.
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	// The name of the load balancer.
	Name string `json:"name,omitempty"`
	// The ID for the owner for this load balancer.
	OwnerID gidx.PrefixedID `json:"owner_id,omitempty"`
	// The ID for the location of this load balancer.
	LocationID gidx.PrefixedID `json:"location_id,omitempty"`
	// The ID for the load balancer provider for this load balancer.
	ProviderID gidx.PrefixedID `json:"provider_id,omitempty"`
	// The ID of the ip address for this load balancer.
	IPID gidx.PrefixedID `json:"ip_id,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the LoadBalancerQuery when eager-loading is set.
	Edges        LoadBalancerEdges `json:"edges"`
	selectValues sql.SelectValues
}

// LoadBalancerEdges holds the relations/edges for other nodes in the graph.
type LoadBalancerEdges struct {
	// Annotations for the load balancer.
	Annotations []*LoadBalancerAnnotation `json:"annotations,omitempty"`
	// Statuses for the load balancer.
	Statuses []*LoadBalancerStatus `json:"statuses,omitempty"`
	// Ports holds the value of the ports edge.
	Ports []*Port `json:"ports,omitempty"`
	// The load balancer provider for the load balancer.
	Provider *Provider `json:"provider,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [4]bool
	// totalCount holds the count of the edges above.
	totalCount [4]map[string]int

	namedAnnotations map[string][]*LoadBalancerAnnotation
	namedStatuses    map[string][]*LoadBalancerStatus
	namedPorts       map[string][]*Port
}

// AnnotationsOrErr returns the Annotations value or an error if the edge
// was not loaded in eager-loading.
func (e LoadBalancerEdges) AnnotationsOrErr() ([]*LoadBalancerAnnotation, error) {
	if e.loadedTypes[0] {
		return e.Annotations, nil
	}
	return nil, &NotLoadedError{edge: "annotations"}
}

// StatusesOrErr returns the Statuses value or an error if the edge
// was not loaded in eager-loading.
func (e LoadBalancerEdges) StatusesOrErr() ([]*LoadBalancerStatus, error) {
	if e.loadedTypes[1] {
		return e.Statuses, nil
	}
	return nil, &NotLoadedError{edge: "statuses"}
}

// PortsOrErr returns the Ports value or an error if the edge
// was not loaded in eager-loading.
func (e LoadBalancerEdges) PortsOrErr() ([]*Port, error) {
	if e.loadedTypes[2] {
		return e.Ports, nil
	}
	return nil, &NotLoadedError{edge: "ports"}
}

// ProviderOrErr returns the Provider value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e LoadBalancerEdges) ProviderOrErr() (*Provider, error) {
	if e.loadedTypes[3] {
		if e.Provider == nil {
			// Edge was loaded but was not found.
			return nil, &NotFoundError{label: provider.Label}
		}
		return e.Provider, nil
	}
	return nil, &NotLoadedError{edge: "provider"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*LoadBalancer) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case loadbalancer.FieldID, loadbalancer.FieldOwnerID, loadbalancer.FieldLocationID, loadbalancer.FieldProviderID, loadbalancer.FieldIPID:
			values[i] = new(gidx.PrefixedID)
		case loadbalancer.FieldName:
			values[i] = new(sql.NullString)
		case loadbalancer.FieldCreatedAt, loadbalancer.FieldUpdatedAt:
			values[i] = new(sql.NullTime)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the LoadBalancer fields.
func (lb *LoadBalancer) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case loadbalancer.FieldID:
			if value, ok := values[i].(*gidx.PrefixedID); !ok {
				return fmt.Errorf("unexpected type %T for field id", values[i])
			} else if value != nil {
				lb.ID = *value
			}
		case loadbalancer.FieldCreatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field created_at", values[i])
			} else if value.Valid {
				lb.CreatedAt = value.Time
			}
		case loadbalancer.FieldUpdatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field updated_at", values[i])
			} else if value.Valid {
				lb.UpdatedAt = value.Time
			}
		case loadbalancer.FieldName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field name", values[i])
			} else if value.Valid {
				lb.Name = value.String
			}
		case loadbalancer.FieldOwnerID:
			if value, ok := values[i].(*gidx.PrefixedID); !ok {
				return fmt.Errorf("unexpected type %T for field owner_id", values[i])
			} else if value != nil {
				lb.OwnerID = *value
			}
		case loadbalancer.FieldLocationID:
			if value, ok := values[i].(*gidx.PrefixedID); !ok {
				return fmt.Errorf("unexpected type %T for field location_id", values[i])
			} else if value != nil {
				lb.LocationID = *value
			}
		case loadbalancer.FieldProviderID:
			if value, ok := values[i].(*gidx.PrefixedID); !ok {
				return fmt.Errorf("unexpected type %T for field provider_id", values[i])
			} else if value != nil {
				lb.ProviderID = *value
			}
		case loadbalancer.FieldIPID:
			if value, ok := values[i].(*gidx.PrefixedID); !ok {
				return fmt.Errorf("unexpected type %T for field ip_id", values[i])
			} else if value != nil {
				lb.IPID = *value
			}
		default:
			lb.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the LoadBalancer.
// This includes values selected through modifiers, order, etc.
func (lb *LoadBalancer) Value(name string) (ent.Value, error) {
	return lb.selectValues.Get(name)
}

// QueryAnnotations queries the "annotations" edge of the LoadBalancer entity.
func (lb *LoadBalancer) QueryAnnotations() *LoadBalancerAnnotationQuery {
	return NewLoadBalancerClient(lb.config).QueryAnnotations(lb)
}

// QueryStatuses queries the "statuses" edge of the LoadBalancer entity.
func (lb *LoadBalancer) QueryStatuses() *LoadBalancerStatusQuery {
	return NewLoadBalancerClient(lb.config).QueryStatuses(lb)
}

// QueryPorts queries the "ports" edge of the LoadBalancer entity.
func (lb *LoadBalancer) QueryPorts() *PortQuery {
	return NewLoadBalancerClient(lb.config).QueryPorts(lb)
}

// QueryProvider queries the "provider" edge of the LoadBalancer entity.
func (lb *LoadBalancer) QueryProvider() *ProviderQuery {
	return NewLoadBalancerClient(lb.config).QueryProvider(lb)
}

// Update returns a builder for updating this LoadBalancer.
// Note that you need to call LoadBalancer.Unwrap() before calling this method if this LoadBalancer
// was returned from a transaction, and the transaction was committed or rolled back.
func (lb *LoadBalancer) Update() *LoadBalancerUpdateOne {
	return NewLoadBalancerClient(lb.config).UpdateOne(lb)
}

// Unwrap unwraps the LoadBalancer entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (lb *LoadBalancer) Unwrap() *LoadBalancer {
	_tx, ok := lb.config.driver.(*txDriver)
	if !ok {
		panic("generated: LoadBalancer is not a transactional entity")
	}
	lb.config.driver = _tx.drv
	return lb
}

// String implements the fmt.Stringer.
func (lb *LoadBalancer) String() string {
	var builder strings.Builder
	builder.WriteString("LoadBalancer(")
	builder.WriteString(fmt.Sprintf("id=%v, ", lb.ID))
	builder.WriteString("created_at=")
	builder.WriteString(lb.CreatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("updated_at=")
	builder.WriteString(lb.UpdatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("name=")
	builder.WriteString(lb.Name)
	builder.WriteString(", ")
	builder.WriteString("owner_id=")
	builder.WriteString(fmt.Sprintf("%v", lb.OwnerID))
	builder.WriteString(", ")
	builder.WriteString("location_id=")
	builder.WriteString(fmt.Sprintf("%v", lb.LocationID))
	builder.WriteString(", ")
	builder.WriteString("provider_id=")
	builder.WriteString(fmt.Sprintf("%v", lb.ProviderID))
	builder.WriteString(", ")
	builder.WriteString("ip_id=")
	builder.WriteString(fmt.Sprintf("%v", lb.IPID))
	builder.WriteByte(')')
	return builder.String()
}

// IsEntity implement fedruntime.Entity
func (lb LoadBalancer) IsEntity() {}

// IsIPv4Addressable implements interface for IPv4Addressable
func (lb LoadBalancer) IsIPv4Addressable() {}

// NamedAnnotations returns the Annotations named value or an error if the edge was not
// loaded in eager-loading with this name.
func (lb *LoadBalancer) NamedAnnotations(name string) ([]*LoadBalancerAnnotation, error) {
	if lb.Edges.namedAnnotations == nil {
		return nil, &NotLoadedError{edge: name}
	}
	nodes, ok := lb.Edges.namedAnnotations[name]
	if !ok {
		return nil, &NotLoadedError{edge: name}
	}
	return nodes, nil
}

func (lb *LoadBalancer) appendNamedAnnotations(name string, edges ...*LoadBalancerAnnotation) {
	if lb.Edges.namedAnnotations == nil {
		lb.Edges.namedAnnotations = make(map[string][]*LoadBalancerAnnotation)
	}
	if len(edges) == 0 {
		lb.Edges.namedAnnotations[name] = []*LoadBalancerAnnotation{}
	} else {
		lb.Edges.namedAnnotations[name] = append(lb.Edges.namedAnnotations[name], edges...)
	}
}

// NamedStatuses returns the Statuses named value or an error if the edge was not
// loaded in eager-loading with this name.
func (lb *LoadBalancer) NamedStatuses(name string) ([]*LoadBalancerStatus, error) {
	if lb.Edges.namedStatuses == nil {
		return nil, &NotLoadedError{edge: name}
	}
	nodes, ok := lb.Edges.namedStatuses[name]
	if !ok {
		return nil, &NotLoadedError{edge: name}
	}
	return nodes, nil
}

func (lb *LoadBalancer) appendNamedStatuses(name string, edges ...*LoadBalancerStatus) {
	if lb.Edges.namedStatuses == nil {
		lb.Edges.namedStatuses = make(map[string][]*LoadBalancerStatus)
	}
	if len(edges) == 0 {
		lb.Edges.namedStatuses[name] = []*LoadBalancerStatus{}
	} else {
		lb.Edges.namedStatuses[name] = append(lb.Edges.namedStatuses[name], edges...)
	}
}

// NamedPorts returns the Ports named value or an error if the edge was not
// loaded in eager-loading with this name.
func (lb *LoadBalancer) NamedPorts(name string) ([]*Port, error) {
	if lb.Edges.namedPorts == nil {
		return nil, &NotLoadedError{edge: name}
	}
	nodes, ok := lb.Edges.namedPorts[name]
	if !ok {
		return nil, &NotLoadedError{edge: name}
	}
	return nodes, nil
}

func (lb *LoadBalancer) appendNamedPorts(name string, edges ...*Port) {
	if lb.Edges.namedPorts == nil {
		lb.Edges.namedPorts = make(map[string][]*Port)
	}
	if len(edges) == 0 {
		lb.Edges.namedPorts[name] = []*Port{}
	} else {
		lb.Edges.namedPorts[name] = append(lb.Edges.namedPorts[name], edges...)
	}
}

// LoadBalancers is a parsable slice of LoadBalancer.
type LoadBalancers []*LoadBalancer
