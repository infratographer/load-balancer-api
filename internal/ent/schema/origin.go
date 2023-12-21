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

	"go.infratographer.com/load-balancer-api/internal/ent/schema/validations"
	"go.infratographer.com/load-balancer-api/x/pubsubinfo"
)

var defaultOriginWeight int32 = 100

// Origin holds the schema definition for the Origin entity.
type Origin struct {
	ent.Schema
}

// Mixin to use for Origin type
func (Origin) Mixin() []ent.Mixin {
	return []ent.Mixin{
		entx.NewTimestampMixin(),
	}
}

// Fields of the Instance.
func (Origin) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			GoType(gidx.PrefixedID("")).
			DefaultFunc(func() gidx.PrefixedID { return gidx.MustNewID(OriginPrefix) }).
			Unique().
			Immutable(),
		field.String("name").
			NotEmpty().
			Annotations(
				entgql.OrderField("name"),
			),
		field.Int32("weight").
			Default(defaultOriginWeight).
			Annotations(
				entgql.OrderField("weight"),
			),
		field.String("target").
			NotEmpty().
			Validate(validations.IPAddress).
			// Comment("origin address").
			Annotations(
				entgql.OrderField("target"),
			),
		field.Int("port_number").
			Min(minPort).
			Max(maxPort).
			Annotations(
				entgql.OrderField("number"),
			),
		field.Bool("active").
			Default(true).
			Annotations(
				entgql.OrderField("active"),
			),
		field.String("pool_id").
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
func (Origin) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("pool", Pool.Type).
			Unique().
			Required().
			Immutable().
			Field("pool_id").
			Annotations(),
	}
}

// Indexes of the Origin
func (Origin) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("pool_id"),
	}
}

// Annotations for the Origin
func (Origin) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entx.GraphKeyDirective("id"),
		pubsubinfo.EventsHookSubjectName("load-balancer-origin"),
		entgql.Type("LoadBalancerOrigin"),
		prefixIDDirective(OriginPrefix),
		entgql.RelayConnection(),
		entgql.Mutations(entgql.MutationCreate(), entgql.MutationUpdate()),
	}
}
