package schema

import (
	"time"

	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"go.infratographer.com/x/entx"
	"go.infratographer.com/x/gidx"
)

// LoadBalancer holds the schema definition for the LoadBalancer entity.
type LoadBalancer struct {
	ent.Schema
}

// Mixin to use for LoadBalancer type
func (LoadBalancer) Mixin() []ent.Mixin {
	return []ent.Mixin{
		// entx.TimestampsMixin{},
		// softdelete.Mixin{},
	}
}

// Fields of the Instance.
func (LoadBalancer) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			GoType(gidx.PrefixedID("")).
			DefaultFunc(func() gidx.PrefixedID { return gidx.MustNewID(LoadBalancerPrefix) }).
			Unique().
			Immutable().
			Annotations(
				entgql.OrderField("ID"),
			),
		field.Text("name").
			NotEmpty().
			Annotations(
				entgql.OrderField("NAME"),
			),
		field.String("tenant_id").
			GoType(gidx.PrefixedID("")).
			Immutable().
			Annotations(
				entgql.QueryField(),
				entgql.Type("ID"),
				entgql.Skip(entgql.SkipWhereInput, entgql.SkipMutationUpdateInput),
			),
		field.String("location_id").
			GoType(gidx.PrefixedID("")).
			Immutable().
			NotEmpty().
			Annotations(
				entgql.Type("ID"),
				entgql.Skip(entgql.SkipWhereInput, entgql.SkipMutationUpdateInput),
				entgql.OrderField("LOCATION"),
			),
		field.String("provider_id").
			GoType(gidx.PrefixedID("")).
			Immutable().
			NotEmpty().
			Annotations(
				entgql.Type("ID"),
				entgql.Skip(^entgql.SkipMutationUpdateInput),
				entgql.OrderField("PROVIDER"),
			),
		field.Time("created_at").
			Default(time.Now).
			Immutable().
			Annotations(
				entgql.Skip(entgql.SkipMutationCreateInput, entgql.SkipMutationUpdateInput),
				entgql.OrderField("CREATED_AT"),
			),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now).
			Immutable().
			Annotations(
				entgql.Skip(entgql.SkipMutationCreateInput, entgql.SkipMutationUpdateInput),
				entgql.OrderField("UPDATED_AT"),
			),
	}
}

// Edges of the Instance.
func (LoadBalancer) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("annotations", LoadBalancerAnnotation.Type).
			Ref("load_balancer").
			Annotations(
				entgql.RelayConnection(),
			),
		edge.From("statuses", LoadBalancerStatus.Type).
			Ref("load_balancer").
			Annotations(
				entgql.RelayConnection(),
			),
		edge.To("provider", Provider.Type).
			Unique().
			Required().
			Immutable().
			Field("provider_id").
			Annotations(
				entgql.MapsTo("loadBalancerProvider"),
			),
	}
}

// Indexes of the LoadBalancer
func (LoadBalancer) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("provider_id"),
		index.Fields("location_id"),
		index.Fields("tenant_id"),
		index.Fields("created_at"),
		index.Fields("updated_at"),
	}
}

// Annotations for the LoadBalancer
func (LoadBalancer) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entx.GraphKeyDirective("id"),
		entgql.RelayConnection(),
		entgql.Mutations(entgql.MutationCreate(), entgql.MutationUpdate()),
		entgql.Implements("IPv4Addressable"),
	}
}
