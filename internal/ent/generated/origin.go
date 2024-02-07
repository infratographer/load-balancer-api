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
	"fmt"
	"strings"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"go.infratographer.com/load-balancer-api/internal/ent/generated/origin"
	"go.infratographer.com/load-balancer-api/internal/ent/generated/pool"
	"go.infratographer.com/x/gidx"
)

// Origin is the model entity for the Origin schema.
type Origin struct {
	config `json:"-"`
	// ID of the ent.
	ID gidx.PrefixedID `json:"id,omitempty"`
	// CreatedAt holds the value of the "created_at" field.
	CreatedAt time.Time `json:"created_at,omitempty"`
	// UpdatedAt holds the value of the "updated_at" field.
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	// CreatedBy holds the value of the "created_by" field.
	CreatedBy string `json:"created_by,omitempty"`
	// UpdatedBy holds the value of the "updated_by" field.
	UpdatedBy string `json:"updated_by,omitempty"`
	// Name holds the value of the "name" field.
	Name string `json:"name,omitempty"`
	// Weight holds the value of the "weight" field.
	Weight int32 `json:"weight,omitempty"`
	// Target holds the value of the "target" field.
	Target string `json:"target,omitempty"`
	// PortNumber holds the value of the "port_number" field.
	PortNumber int `json:"port_number,omitempty"`
	// Active holds the value of the "active" field.
	Active bool `json:"active,omitempty"`
	// PoolID holds the value of the "pool_id" field.
	PoolID gidx.PrefixedID `json:"pool_id,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the OriginQuery when eager-loading is set.
	Edges        OriginEdges `json:"edges"`
	selectValues sql.SelectValues
}

// OriginEdges holds the relations/edges for other nodes in the graph.
type OriginEdges struct {
	// Pool holds the value of the pool edge.
	Pool *Pool `json:"pool,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [1]bool
	// totalCount holds the count of the edges above.
	totalCount [1]map[string]int
}

// PoolOrErr returns the Pool value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e OriginEdges) PoolOrErr() (*Pool, error) {
	if e.loadedTypes[0] {
		if e.Pool == nil {
			// Edge was loaded but was not found.
			return nil, &NotFoundError{label: pool.Label}
		}
		return e.Pool, nil
	}
	return nil, &NotLoadedError{edge: "pool"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Origin) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case origin.FieldID, origin.FieldPoolID:
			values[i] = new(gidx.PrefixedID)
		case origin.FieldActive:
			values[i] = new(sql.NullBool)
		case origin.FieldWeight, origin.FieldPortNumber:
			values[i] = new(sql.NullInt64)
		case origin.FieldCreatedBy, origin.FieldUpdatedBy, origin.FieldName, origin.FieldTarget:
			values[i] = new(sql.NullString)
		case origin.FieldCreatedAt, origin.FieldUpdatedAt:
			values[i] = new(sql.NullTime)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Origin fields.
func (o *Origin) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case origin.FieldID:
			if value, ok := values[i].(*gidx.PrefixedID); !ok {
				return fmt.Errorf("unexpected type %T for field id", values[i])
			} else if value != nil {
				o.ID = *value
			}
		case origin.FieldCreatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field created_at", values[i])
			} else if value.Valid {
				o.CreatedAt = value.Time
			}
		case origin.FieldUpdatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field updated_at", values[i])
			} else if value.Valid {
				o.UpdatedAt = value.Time
			}
		case origin.FieldCreatedBy:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field created_by", values[i])
			} else if value.Valid {
				o.CreatedBy = value.String
			}
		case origin.FieldUpdatedBy:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field updated_by", values[i])
			} else if value.Valid {
				o.UpdatedBy = value.String
			}
		case origin.FieldName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field name", values[i])
			} else if value.Valid {
				o.Name = value.String
			}
		case origin.FieldWeight:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field weight", values[i])
			} else if value.Valid {
				o.Weight = int32(value.Int64)
			}
		case origin.FieldTarget:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field target", values[i])
			} else if value.Valid {
				o.Target = value.String
			}
		case origin.FieldPortNumber:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field port_number", values[i])
			} else if value.Valid {
				o.PortNumber = int(value.Int64)
			}
		case origin.FieldActive:
			if value, ok := values[i].(*sql.NullBool); !ok {
				return fmt.Errorf("unexpected type %T for field active", values[i])
			} else if value.Valid {
				o.Active = value.Bool
			}
		case origin.FieldPoolID:
			if value, ok := values[i].(*gidx.PrefixedID); !ok {
				return fmt.Errorf("unexpected type %T for field pool_id", values[i])
			} else if value != nil {
				o.PoolID = *value
			}
		default:
			o.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the Origin.
// This includes values selected through modifiers, order, etc.
func (o *Origin) Value(name string) (ent.Value, error) {
	return o.selectValues.Get(name)
}

// QueryPool queries the "pool" edge of the Origin entity.
func (o *Origin) QueryPool() *PoolQuery {
	return NewOriginClient(o.config).QueryPool(o)
}

// Update returns a builder for updating this Origin.
// Note that you need to call Origin.Unwrap() before calling this method if this Origin
// was returned from a transaction, and the transaction was committed or rolled back.
func (o *Origin) Update() *OriginUpdateOne {
	return NewOriginClient(o.config).UpdateOne(o)
}

// Unwrap unwraps the Origin entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (o *Origin) Unwrap() *Origin {
	_tx, ok := o.config.driver.(*txDriver)
	if !ok {
		panic("generated: Origin is not a transactional entity")
	}
	o.config.driver = _tx.drv
	return o
}

// String implements the fmt.Stringer.
func (o *Origin) String() string {
	var builder strings.Builder
	builder.WriteString("Origin(")
	builder.WriteString(fmt.Sprintf("id=%v, ", o.ID))
	builder.WriteString("created_at=")
	builder.WriteString(o.CreatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("updated_at=")
	builder.WriteString(o.UpdatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("created_by=")
	builder.WriteString(o.CreatedBy)
	builder.WriteString(", ")
	builder.WriteString("updated_by=")
	builder.WriteString(o.UpdatedBy)
	builder.WriteString(", ")
	builder.WriteString("name=")
	builder.WriteString(o.Name)
	builder.WriteString(", ")
	builder.WriteString("weight=")
	builder.WriteString(fmt.Sprintf("%v", o.Weight))
	builder.WriteString(", ")
	builder.WriteString("target=")
	builder.WriteString(o.Target)
	builder.WriteString(", ")
	builder.WriteString("port_number=")
	builder.WriteString(fmt.Sprintf("%v", o.PortNumber))
	builder.WriteString(", ")
	builder.WriteString("active=")
	builder.WriteString(fmt.Sprintf("%v", o.Active))
	builder.WriteString(", ")
	builder.WriteString("pool_id=")
	builder.WriteString(fmt.Sprintf("%v", o.PoolID))
	builder.WriteByte(')')
	return builder.String()
}

// IsEntity implement fedruntime.Entity
func (o Origin) IsEntity() {}

// Origins is a parsable slice of Origin.
type Origins []*Origin
