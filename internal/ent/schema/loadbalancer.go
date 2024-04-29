package schema

import (
	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/vektah/gqlparser/v2/ast"

	"go.infratographer.com/x/entx"
	"go.infratographer.com/x/gidx"

	"go.infratographer.com/load-balancer-api/internal/ent/schema/audit"
	"go.infratographer.com/load-balancer-api/internal/ent/schema/softdelete"
	"go.infratographer.com/load-balancer-api/internal/ent/schema/validations"
	"go.infratographer.com/load-balancer-api/x/pubsubinfo"
)

// LoadBalancer holds the schema definition for the LoadBalancer entity.
type LoadBalancer struct {
	ent.Schema
}

// Mixin to use for LoadBalancer type
func (LoadBalancer) Mixin() []ent.Mixin {
	return []ent.Mixin{
		entx.NewTimestampMixin(),
		audit.Mixin{},
		softdelete.Mixin{},
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
			Comment("The ID for the load balancer.").
			Annotations(
				entgql.OrderField("ID"),
			),
		field.Text("name").
			// NotEmpty().
			Validate(validations.NameField).
			Comment("The name of the load balancer.").
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
		field.String("location_id").
			GoType(gidx.PrefixedID("")).
			Immutable().
			NotEmpty().
			Comment("The ID for the location of this load balancer.").
			Annotations(
				entgql.Type("ID"),
				entgql.Skip(^entgql.SkipMutationCreateInput),
			),
		field.String("provider_id").
			GoType(gidx.PrefixedID("")).
			Immutable().
			NotEmpty().
			Comment("The ID for the load balancer provider for this load balancer.").
			Annotations(
				entgql.Type("ID"),
				entgql.Skip(^entgql.SkipMutationCreateInput),
			),
	}
}

// Edges of the Instance.
func (LoadBalancer) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("ports", Port.Type).
			Ref("load_balancer").
			Annotations(
				entgql.RelayConnection(),
			),
		edge.To("provider", Provider.Type).
			Unique().
			Required().
			Immutable().
			Field("provider_id").
			Comment("The load balancer provider for the load balancer.").
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
		index.Fields("owner_id"),
	}
}

// Annotations for the LoadBalancer
func (LoadBalancer) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entx.GraphKeyDirective("id"),
		pubsubinfo.EventsHookSubjectName("load-balancer"),
		schema.Comment("Representation of a load balancer."),
		prefixIDDirective(LoadBalancerPrefix),
		entgql.Implements("IPAddressable"),
		entgql.Implements("MetadataNode"),
		entgql.RelayConnection(),
		entgql.Mutations(
			entgql.MutationCreate().Description("Input information to create a load balancer."),
			entgql.MutationUpdate().Description("Input information to update a load balancer."),
		),
	}
}

func prefixIDDirective(prefix string) entgql.Annotation {
	var args []*ast.Argument
	if prefix != "" {
		args = append(args, &ast.Argument{
			Name: "prefix",
			Value: &ast.Value{
				Raw:  prefix,
				Kind: ast.StringValue,
			},
		})
	}

	return entgql.Directives(entgql.NewDirective("prefixedID", args...))
}
