package schematype

import (
	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"go.infratographer.com/x/entx"
	"go.infratographer.com/x/gidx"
)

// LoadBalancerStatus holds the schema definition for the LoadBalancerStatus entity.
type LoadBalancerStatus struct {
	ent.Schema
}

// Mixin of the LoadBalancerStatus
func (LoadBalancerStatus) Mixin() []ent.Mixin {
	return []ent.Mixin{
		entx.NamespacedDataMixin{},
		entx.NewTimestampMixin(),
	}
}

// Fields of the LoadBalancerStatus.
func (LoadBalancerStatus) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			GoType(gidx.PrefixedID("")).
			DefaultFunc(func() gidx.PrefixedID { return gidx.MustNewID(LoadBalancerStatusPrefix) }).
			Unique().
			Immutable(),
		field.String("load_balancer_id").
			GoType(gidx.PrefixedID("")).
			Immutable().
			Annotations(
				entgql.Type("ID"),
				entgql.Skip(entgql.SkipWhereInput, entgql.SkipMutationUpdateInput),
			),
		field.String("source").
			Immutable().
			NotEmpty().
			Annotations(
				entgql.Skip(entgql.SkipMutationUpdateInput),
			),
	}
}

// Indexes of the LoadBalancerStatus
func (LoadBalancerStatus) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("load_balancer_id"),
		index.Fields("load_balancer_id", "namespace", "source"),
		index.Fields("namespace", "data").Annotations(
			entsql.IndexTypes(map[string]string{
				dialect.Postgres: "GIN",
			}),
		),
	}
}

// Edges of the LoadBalancerStatus
func (LoadBalancerStatus) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("load_balancer", LoadBalancer.Type).
			Unique().
			Required().
			Immutable().
			Field("load_balancer_id").
			Annotations(),
	}
}

// Annotations for the LoadBalancerStatus
func (LoadBalancerStatus) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entx.GraphKeyDirective("id"),
		entgql.RelayConnection(),
	}
}

// Hooks configures actions to take before and after mutations.
func (LoadBalancerStatus) Hooks() []ent.Hook {
	return []ent.Hook{}
}
