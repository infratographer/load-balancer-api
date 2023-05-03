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
)

// Provider holds the schema definition for the LoadBalancerProvider entity.
type Provider struct {
	ent.Schema
}

// Mixin of the Provider
func (Provider) Mixin() []ent.Mixin {
	return []ent.Mixin{
		entx.TimestampsMixin{},
	}
}

// Fields of the Provider.
func (Provider) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			GoType(gidx.PrefixedID("")).
			DefaultFunc(func() gidx.PrefixedID { return gidx.MustNewID(LoadBalancerProviderPrefix) }).
			Unique().
			Immutable(),
		field.String("name").
			NotEmpty(),
		field.String("tenant_id").
			GoType(gidx.PrefixedID("")).
			Immutable().
			Annotations(
				entgql.Type("ID"),
				entgql.Skip(entgql.SkipWhereInput, entgql.SkipMutationUpdateInput),
			),
	}
}

// Indexes of the Provider
func (Provider) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("tenant_id"),
	}
}

// Edges of the Provider
func (Provider) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("load_balancers", LoadBalancer.Type).
			Ref("provider").
			Annotations(
				entgql.RelayConnection(),
			),
	}
}

// Annotations for the Provider
func (Provider) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entx.GraphKeyDirective("id"),
		entgql.Type("LoadBalancerProvider"),
		entgql.RelayConnection(),
		entgql.QueryField(),
		entgql.Mutations(entgql.MutationCreate(), entgql.MutationUpdate()),
	}
}
