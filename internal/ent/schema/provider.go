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

	"go.infratographer.com/load-balancer-api/x/pubsubinfo"
)

// Provider holds the schema definition for the LoadBalancerProvider entity.
type Provider struct {
	ent.Schema
}

// Mixin of the Provider
func (Provider) Mixin() []ent.Mixin {
	return []ent.Mixin{
		entx.NewTimestampMixin(),
	}
}

// Fields of the Provider.
func (Provider) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			GoType(gidx.PrefixedID("")).
			DefaultFunc(func() gidx.PrefixedID { return gidx.MustNewID(LoadBalancerProviderPrefix) }).
			Unique().
			Immutable().
			Comment("The ID for the load balancer provider.").
			Annotations(
				entgql.OrderField("ID"),
			),
		field.String("name").
			NotEmpty().
			Comment("The name of the load balancer provider.").
			Annotations(
				entgql.OrderField("NAME"),
			),
		field.String("owner_id").
			GoType(gidx.PrefixedID("")).
			Immutable().
			NotEmpty().
			Comment("The ID for the owner for this load balancer.").
			Annotations(
				entgql.QueryField(),
				entgql.Type("ID"),
				entgql.Skip(entgql.SkipWhereInput, entgql.SkipMutationUpdateInput, entgql.SkipType),
				entgql.OrderField("OWNER"),
				pubsubinfo.EventsHookAdditionalSubject("owner"),
			),
	}
}

// Indexes of the Provider
func (Provider) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("owner_id"),
	}
}

// Edges of the Provider
func (Provider) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("load_balancers", LoadBalancer.Type).
			Ref("provider").
			Annotations(
				entgql.RelayConnection(),
				entgql.Skip(entgql.SkipMutationCreateInput, entgql.SkipMutationUpdateInput),
			),
	}
}

// Annotations for the Provider
func (Provider) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entx.GraphKeyDirective("id"),
		schema.Comment("Representation of a load balancer provider. Load balancer providers are responsible for provisioning and managing load balancers"),
		entgql.Type("LoadBalancerProvider"),
		prefixIDDirective(LoadBalancerProviderPrefix),
		entgql.RelayConnection(),
		entgql.Mutations(
			entgql.MutationCreate().Description("Input information to create a load balancer provider."),
			entgql.MutationUpdate().Description("Input information to update a load balancer provider."),
		),
	}
}
