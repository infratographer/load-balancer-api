package schema

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

// LoadBalancerAnnotation holds the schema definition for the LoadBalancerAnnotation entity.
type LoadBalancerAnnotation struct {
	ent.Schema
}

// Mixin of the LoadBalancerAnnotation
func (LoadBalancerAnnotation) Mixin() []ent.Mixin {
	return []ent.Mixin{
		entx.NamespacedDataMixin{},
		entx.NewTimestampMixin(),
	}
}

// Fields of the LoadBalancerAnnotation.
func (LoadBalancerAnnotation) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			GoType(gidx.PrefixedID("")).
			DefaultFunc(func() gidx.PrefixedID { return gidx.MustNewID(LoadBalancerAnnotationPrefix) }).
			Unique().
			Immutable(),
		field.String("load_balancer_id").
			GoType(gidx.PrefixedID("")).
			Immutable().
			Annotations(
				entgql.Type("ID"),
				entgql.Skip(entgql.SkipWhereInput, entgql.SkipMutationUpdateInput),
			),
	}
}

// Indexes of the LoadBalancerAnnotation
func (LoadBalancerAnnotation) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("load_balancer_id"),
		index.Fields("load_balancer_id", "namespace"),
		index.Fields("namespace", "data").Annotations(
			entsql.IndexTypes(map[string]string{
				dialect.Postgres: "GIN",
			}),
		),
	}
}

// Edges of the LoadBalancerAnnotation
func (LoadBalancerAnnotation) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("load_balancer", LoadBalancer.Type).
			Unique().
			Required().
			Immutable().
			Field("load_balancer_id").
			Annotations(),
	}
}

// Annotations for the LoadBalancerAnnotation
func (LoadBalancerAnnotation) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entx.GraphKeyDirective("id"),
		entgql.RelayConnection(),
	}
}
