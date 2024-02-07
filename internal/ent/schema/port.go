package schema

import (
	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"go.infratographer.com/x/entx"
	"go.infratographer.com/x/gidx"

	"go.infratographer.com/load-balancer-api/internal/ent/schema/audit"
	"go.infratographer.com/load-balancer-api/internal/ent/schema/validations"
	"go.infratographer.com/load-balancer-api/x/pubsubinfo"
)

const (
	minPort = 1
	maxPort = 65535
)

// Port holds the schema definition for the Port entity.
type Port struct {
	ent.Schema
}

// Mixin to use for Port type
func (Port) Mixin() []ent.Mixin {
	return []ent.Mixin{
		entx.NewTimestampMixin(),
		audit.AuditMixin{},
	}
}

// Fields of the Instance.
func (Port) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			GoType(gidx.PrefixedID("")).
			DefaultFunc(func() gidx.PrefixedID { return gidx.MustNewID(PortPrefix) }).
			Unique().
			Immutable(),
		field.Int("number").
			Min(minPort).
			Max(maxPort).
			Validate(validations.RestrictedPorts).
			Annotations(
				entgql.OrderField("number"),
			),
		field.String("name").
			Annotations(
				entgql.OrderField("name"),
			),
		field.String("load_balancer_id").
			GoType(gidx.PrefixedID("")).
			Immutable().
			NotEmpty().
			Annotations(
				entgql.Type("ID"),
				entgql.Skip(entgql.SkipWhereInput, entgql.SkipMutationUpdateInput),
			),
	}
}

// Edges of the Instance.
func (Port) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("pools", Pool.Type).Ref("ports"),
		edge.To("load_balancer", LoadBalancer.Type).
			Unique().
			Required().
			Immutable().
			Field("load_balancer_id").
			Annotations(),
	}
}

// Indexes of the Port
func (Port) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("load_balancer_id"),
		index.Fields("load_balancer_id", "number").Unique(),
	}
}

// Annotations for the Port
func (Port) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entx.GraphKeyDirective("id"),
		entgql.Type("LoadBalancerPort"),
		prefixIDDirective(PortPrefix),
		entgql.RelayConnection(),
		entgql.Mutations(entgql.MutationCreate(), entgql.MutationUpdate()),
		pubsubinfo.EventsHookSubjectName("load-balancer-port"),
	}
}
