// Copyright 2023 The Infratographer Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Code generated by entc, DO NOT EDIT.

package loadbalancerannotation

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/dialect/sql/sqljson"
	"go.infratographer.com/load-balancer-api/internal/ent/generated/predicate"
	"go.infratographer.com/x/gidx"
)

// ID filters vertices based on their ID field.
func ID(id gidx.PrefixedID) predicate.LoadBalancerAnnotation {
	return predicate.LoadBalancerAnnotation(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id gidx.PrefixedID) predicate.LoadBalancerAnnotation {
	return predicate.LoadBalancerAnnotation(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id gidx.PrefixedID) predicate.LoadBalancerAnnotation {
	return predicate.LoadBalancerAnnotation(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...gidx.PrefixedID) predicate.LoadBalancerAnnotation {
	return predicate.LoadBalancerAnnotation(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...gidx.PrefixedID) predicate.LoadBalancerAnnotation {
	return predicate.LoadBalancerAnnotation(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id gidx.PrefixedID) predicate.LoadBalancerAnnotation {
	return predicate.LoadBalancerAnnotation(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id gidx.PrefixedID) predicate.LoadBalancerAnnotation {
	return predicate.LoadBalancerAnnotation(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id gidx.PrefixedID) predicate.LoadBalancerAnnotation {
	return predicate.LoadBalancerAnnotation(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id gidx.PrefixedID) predicate.LoadBalancerAnnotation {
	return predicate.LoadBalancerAnnotation(sql.FieldLTE(FieldID, id))
}

// Namespace applies equality check predicate on the "namespace" field. It's identical to NamespaceEQ.
func Namespace(v string) predicate.LoadBalancerAnnotation {
	return predicate.LoadBalancerAnnotation(sql.FieldEQ(FieldNamespace, v))
}

// CreatedAt applies equality check predicate on the "created_at" field. It's identical to CreatedAtEQ.
func CreatedAt(v time.Time) predicate.LoadBalancerAnnotation {
	return predicate.LoadBalancerAnnotation(sql.FieldEQ(FieldCreatedAt, v))
}

// UpdatedAt applies equality check predicate on the "updated_at" field. It's identical to UpdatedAtEQ.
func UpdatedAt(v time.Time) predicate.LoadBalancerAnnotation {
	return predicate.LoadBalancerAnnotation(sql.FieldEQ(FieldUpdatedAt, v))
}

// LoadBalancerID applies equality check predicate on the "load_balancer_id" field. It's identical to LoadBalancerIDEQ.
func LoadBalancerID(v gidx.PrefixedID) predicate.LoadBalancerAnnotation {
	return predicate.LoadBalancerAnnotation(sql.FieldEQ(FieldLoadBalancerID, v))
}

// NamespaceEQ applies the EQ predicate on the "namespace" field.
func NamespaceEQ(v string) predicate.LoadBalancerAnnotation {
	return predicate.LoadBalancerAnnotation(sql.FieldEQ(FieldNamespace, v))
}

// NamespaceNEQ applies the NEQ predicate on the "namespace" field.
func NamespaceNEQ(v string) predicate.LoadBalancerAnnotation {
	return predicate.LoadBalancerAnnotation(sql.FieldNEQ(FieldNamespace, v))
}

// NamespaceIn applies the In predicate on the "namespace" field.
func NamespaceIn(vs ...string) predicate.LoadBalancerAnnotation {
	return predicate.LoadBalancerAnnotation(sql.FieldIn(FieldNamespace, vs...))
}

// NamespaceNotIn applies the NotIn predicate on the "namespace" field.
func NamespaceNotIn(vs ...string) predicate.LoadBalancerAnnotation {
	return predicate.LoadBalancerAnnotation(sql.FieldNotIn(FieldNamespace, vs...))
}

// NamespaceGT applies the GT predicate on the "namespace" field.
func NamespaceGT(v string) predicate.LoadBalancerAnnotation {
	return predicate.LoadBalancerAnnotation(sql.FieldGT(FieldNamespace, v))
}

// NamespaceGTE applies the GTE predicate on the "namespace" field.
func NamespaceGTE(v string) predicate.LoadBalancerAnnotation {
	return predicate.LoadBalancerAnnotation(sql.FieldGTE(FieldNamespace, v))
}

// NamespaceLT applies the LT predicate on the "namespace" field.
func NamespaceLT(v string) predicate.LoadBalancerAnnotation {
	return predicate.LoadBalancerAnnotation(sql.FieldLT(FieldNamespace, v))
}

// NamespaceLTE applies the LTE predicate on the "namespace" field.
func NamespaceLTE(v string) predicate.LoadBalancerAnnotation {
	return predicate.LoadBalancerAnnotation(sql.FieldLTE(FieldNamespace, v))
}

// NamespaceContains applies the Contains predicate on the "namespace" field.
func NamespaceContains(v string) predicate.LoadBalancerAnnotation {
	return predicate.LoadBalancerAnnotation(sql.FieldContains(FieldNamespace, v))
}

// NamespaceHasPrefix applies the HasPrefix predicate on the "namespace" field.
func NamespaceHasPrefix(v string) predicate.LoadBalancerAnnotation {
	return predicate.LoadBalancerAnnotation(sql.FieldHasPrefix(FieldNamespace, v))
}

// NamespaceHasSuffix applies the HasSuffix predicate on the "namespace" field.
func NamespaceHasSuffix(v string) predicate.LoadBalancerAnnotation {
	return predicate.LoadBalancerAnnotation(sql.FieldHasSuffix(FieldNamespace, v))
}

// NamespaceEqualFold applies the EqualFold predicate on the "namespace" field.
func NamespaceEqualFold(v string) predicate.LoadBalancerAnnotation {
	return predicate.LoadBalancerAnnotation(sql.FieldEqualFold(FieldNamespace, v))
}

// NamespaceContainsFold applies the ContainsFold predicate on the "namespace" field.
func NamespaceContainsFold(v string) predicate.LoadBalancerAnnotation {
	return predicate.LoadBalancerAnnotation(sql.FieldContainsFold(FieldNamespace, v))
}

// CreatedAtEQ applies the EQ predicate on the "created_at" field.
func CreatedAtEQ(v time.Time) predicate.LoadBalancerAnnotation {
	return predicate.LoadBalancerAnnotation(sql.FieldEQ(FieldCreatedAt, v))
}

// CreatedAtNEQ applies the NEQ predicate on the "created_at" field.
func CreatedAtNEQ(v time.Time) predicate.LoadBalancerAnnotation {
	return predicate.LoadBalancerAnnotation(sql.FieldNEQ(FieldCreatedAt, v))
}

// CreatedAtIn applies the In predicate on the "created_at" field.
func CreatedAtIn(vs ...time.Time) predicate.LoadBalancerAnnotation {
	return predicate.LoadBalancerAnnotation(sql.FieldIn(FieldCreatedAt, vs...))
}

// CreatedAtNotIn applies the NotIn predicate on the "created_at" field.
func CreatedAtNotIn(vs ...time.Time) predicate.LoadBalancerAnnotation {
	return predicate.LoadBalancerAnnotation(sql.FieldNotIn(FieldCreatedAt, vs...))
}

// CreatedAtGT applies the GT predicate on the "created_at" field.
func CreatedAtGT(v time.Time) predicate.LoadBalancerAnnotation {
	return predicate.LoadBalancerAnnotation(sql.FieldGT(FieldCreatedAt, v))
}

// CreatedAtGTE applies the GTE predicate on the "created_at" field.
func CreatedAtGTE(v time.Time) predicate.LoadBalancerAnnotation {
	return predicate.LoadBalancerAnnotation(sql.FieldGTE(FieldCreatedAt, v))
}

// CreatedAtLT applies the LT predicate on the "created_at" field.
func CreatedAtLT(v time.Time) predicate.LoadBalancerAnnotation {
	return predicate.LoadBalancerAnnotation(sql.FieldLT(FieldCreatedAt, v))
}

// CreatedAtLTE applies the LTE predicate on the "created_at" field.
func CreatedAtLTE(v time.Time) predicate.LoadBalancerAnnotation {
	return predicate.LoadBalancerAnnotation(sql.FieldLTE(FieldCreatedAt, v))
}

// UpdatedAtEQ applies the EQ predicate on the "updated_at" field.
func UpdatedAtEQ(v time.Time) predicate.LoadBalancerAnnotation {
	return predicate.LoadBalancerAnnotation(sql.FieldEQ(FieldUpdatedAt, v))
}

// UpdatedAtNEQ applies the NEQ predicate on the "updated_at" field.
func UpdatedAtNEQ(v time.Time) predicate.LoadBalancerAnnotation {
	return predicate.LoadBalancerAnnotation(sql.FieldNEQ(FieldUpdatedAt, v))
}

// UpdatedAtIn applies the In predicate on the "updated_at" field.
func UpdatedAtIn(vs ...time.Time) predicate.LoadBalancerAnnotation {
	return predicate.LoadBalancerAnnotation(sql.FieldIn(FieldUpdatedAt, vs...))
}

// UpdatedAtNotIn applies the NotIn predicate on the "updated_at" field.
func UpdatedAtNotIn(vs ...time.Time) predicate.LoadBalancerAnnotation {
	return predicate.LoadBalancerAnnotation(sql.FieldNotIn(FieldUpdatedAt, vs...))
}

// UpdatedAtGT applies the GT predicate on the "updated_at" field.
func UpdatedAtGT(v time.Time) predicate.LoadBalancerAnnotation {
	return predicate.LoadBalancerAnnotation(sql.FieldGT(FieldUpdatedAt, v))
}

// UpdatedAtGTE applies the GTE predicate on the "updated_at" field.
func UpdatedAtGTE(v time.Time) predicate.LoadBalancerAnnotation {
	return predicate.LoadBalancerAnnotation(sql.FieldGTE(FieldUpdatedAt, v))
}

// UpdatedAtLT applies the LT predicate on the "updated_at" field.
func UpdatedAtLT(v time.Time) predicate.LoadBalancerAnnotation {
	return predicate.LoadBalancerAnnotation(sql.FieldLT(FieldUpdatedAt, v))
}

// UpdatedAtLTE applies the LTE predicate on the "updated_at" field.
func UpdatedAtLTE(v time.Time) predicate.LoadBalancerAnnotation {
	return predicate.LoadBalancerAnnotation(sql.FieldLTE(FieldUpdatedAt, v))
}

// LoadBalancerIDEQ applies the EQ predicate on the "load_balancer_id" field.
func LoadBalancerIDEQ(v gidx.PrefixedID) predicate.LoadBalancerAnnotation {
	return predicate.LoadBalancerAnnotation(sql.FieldEQ(FieldLoadBalancerID, v))
}

// LoadBalancerIDNEQ applies the NEQ predicate on the "load_balancer_id" field.
func LoadBalancerIDNEQ(v gidx.PrefixedID) predicate.LoadBalancerAnnotation {
	return predicate.LoadBalancerAnnotation(sql.FieldNEQ(FieldLoadBalancerID, v))
}

// LoadBalancerIDIn applies the In predicate on the "load_balancer_id" field.
func LoadBalancerIDIn(vs ...gidx.PrefixedID) predicate.LoadBalancerAnnotation {
	return predicate.LoadBalancerAnnotation(sql.FieldIn(FieldLoadBalancerID, vs...))
}

// LoadBalancerIDNotIn applies the NotIn predicate on the "load_balancer_id" field.
func LoadBalancerIDNotIn(vs ...gidx.PrefixedID) predicate.LoadBalancerAnnotation {
	return predicate.LoadBalancerAnnotation(sql.FieldNotIn(FieldLoadBalancerID, vs...))
}

// LoadBalancerIDGT applies the GT predicate on the "load_balancer_id" field.
func LoadBalancerIDGT(v gidx.PrefixedID) predicate.LoadBalancerAnnotation {
	return predicate.LoadBalancerAnnotation(sql.FieldGT(FieldLoadBalancerID, v))
}

// LoadBalancerIDGTE applies the GTE predicate on the "load_balancer_id" field.
func LoadBalancerIDGTE(v gidx.PrefixedID) predicate.LoadBalancerAnnotation {
	return predicate.LoadBalancerAnnotation(sql.FieldGTE(FieldLoadBalancerID, v))
}

// LoadBalancerIDLT applies the LT predicate on the "load_balancer_id" field.
func LoadBalancerIDLT(v gidx.PrefixedID) predicate.LoadBalancerAnnotation {
	return predicate.LoadBalancerAnnotation(sql.FieldLT(FieldLoadBalancerID, v))
}

// LoadBalancerIDLTE applies the LTE predicate on the "load_balancer_id" field.
func LoadBalancerIDLTE(v gidx.PrefixedID) predicate.LoadBalancerAnnotation {
	return predicate.LoadBalancerAnnotation(sql.FieldLTE(FieldLoadBalancerID, v))
}

// LoadBalancerIDContains applies the Contains predicate on the "load_balancer_id" field.
func LoadBalancerIDContains(v gidx.PrefixedID) predicate.LoadBalancerAnnotation {
	vc := string(v)
	return predicate.LoadBalancerAnnotation(sql.FieldContains(FieldLoadBalancerID, vc))
}

// LoadBalancerIDHasPrefix applies the HasPrefix predicate on the "load_balancer_id" field.
func LoadBalancerIDHasPrefix(v gidx.PrefixedID) predicate.LoadBalancerAnnotation {
	vc := string(v)
	return predicate.LoadBalancerAnnotation(sql.FieldHasPrefix(FieldLoadBalancerID, vc))
}

// LoadBalancerIDHasSuffix applies the HasSuffix predicate on the "load_balancer_id" field.
func LoadBalancerIDHasSuffix(v gidx.PrefixedID) predicate.LoadBalancerAnnotation {
	vc := string(v)
	return predicate.LoadBalancerAnnotation(sql.FieldHasSuffix(FieldLoadBalancerID, vc))
}

// LoadBalancerIDEqualFold applies the EqualFold predicate on the "load_balancer_id" field.
func LoadBalancerIDEqualFold(v gidx.PrefixedID) predicate.LoadBalancerAnnotation {
	vc := string(v)
	return predicate.LoadBalancerAnnotation(sql.FieldEqualFold(FieldLoadBalancerID, vc))
}

// LoadBalancerIDContainsFold applies the ContainsFold predicate on the "load_balancer_id" field.
func LoadBalancerIDContainsFold(v gidx.PrefixedID) predicate.LoadBalancerAnnotation {
	vc := string(v)
	return predicate.LoadBalancerAnnotation(sql.FieldContainsFold(FieldLoadBalancerID, vc))
}

// HasLoadBalancer applies the HasEdge predicate on the "load_balancer" edge.
func HasLoadBalancer() predicate.LoadBalancerAnnotation {
	return predicate.LoadBalancerAnnotation(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, LoadBalancerTable, LoadBalancerColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasLoadBalancerWith applies the HasEdge predicate on the "load_balancer" edge with a given conditions (other predicates).
func HasLoadBalancerWith(preds ...predicate.LoadBalancer) predicate.LoadBalancerAnnotation {
	return predicate.LoadBalancerAnnotation(func(s *sql.Selector) {
		step := newLoadBalancerStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.LoadBalancerAnnotation) predicate.LoadBalancerAnnotation {
	return predicate.LoadBalancerAnnotation(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for _, p := range predicates {
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.LoadBalancerAnnotation) predicate.LoadBalancerAnnotation {
	return predicate.LoadBalancerAnnotation(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for i, p := range predicates {
			if i > 0 {
				s1.Or()
			}
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Not applies the not operator on the given predicate.
func Not(p predicate.LoadBalancerAnnotation) predicate.LoadBalancerAnnotation {
	return predicate.LoadBalancerAnnotation(func(s *sql.Selector) {
		p(s.Not())
	})
}

// DataHasKey checks if Data contains given value
func DataHasKey(v string) predicate.LoadBalancerAnnotation {
	return predicate.LoadBalancerAnnotation(func(s *sql.Selector) {
		s.Where(sqljson.HasKey(s.C(FieldData), sqljson.DotPath(v)))
	})
}