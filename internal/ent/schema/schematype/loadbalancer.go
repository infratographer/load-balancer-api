package schematype

import (
	"context"

	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/pkg/errors"

	"go.infratographer.com/x/entx"
	"go.infratographer.com/x/gidx"

	gen "go.infratographer.com/load-balancer-api/internal/ent/generated"
	"go.infratographer.com/load-balancer-api/internal/ent/generated/hook"
	"go.infratographer.com/load-balancer-api/internal/pubsub"
	"go.infratographer.com/load-balancer-api/internal/pubsub/types"
)

// LoadBalancer holds the schema definition for the LoadBalancer entity.
type LoadBalancer struct {
	ent.Schema
}

// Mixin to use for LoadBalancer type
func (LoadBalancer) Mixin() []ent.Mixin {
	return []ent.Mixin{
		entx.NewTimestampMixin(),
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
			Comment("The ID for the load balancer.").
			Annotations(
				entgql.OrderField("ID"),
			),
		field.Text("name").
			NotEmpty().
			Comment("The name of the load balancer.").
			Annotations(
				entgql.OrderField("NAME"),
			),
		field.String("tenant_id").
			GoType(gidx.PrefixedID("")).
			Immutable().
			Comment("The ID for the tenant for this load balancer.").
			Annotations(
				entgql.QueryField(),
				entgql.Type("ID"),
				entgql.Skip(entgql.SkipWhereInput, entgql.SkipMutationUpdateInput, entgql.SkipType),
				entgql.OrderField("TENANT"),
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
		edge.From("annotations", LoadBalancerAnnotation.Type).
			Ref("load_balancer").
			Comment("Annotations for the load balancer.").
			Annotations(
				entgql.RelayConnection(),
				entgql.Skip(entgql.SkipMutationCreateInput, entgql.SkipMutationUpdateInput),
			),
		edge.From("statuses", LoadBalancerStatus.Type).
			Ref("load_balancer").
			Comment("Statuses for the load balancer.").
			Annotations(
				entgql.RelayConnection(),
				entgql.Skip(entgql.SkipMutationCreateInput, entgql.SkipMutationUpdateInput),
			),
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
		index.Fields("tenant_id"),
	}
}

// Annotations for the LoadBalancer
func (LoadBalancer) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entx.GraphKeyDirective("id"),
		schema.Comment("Representation of a load balancer."),
		entgql.Implements("IPv4Addressable"),
		entgql.RelayConnection(),
		entgql.Mutations(
			entgql.MutationCreate().Description("Input information to create a load balancer."),
			entgql.MutationUpdate().Description("Input information to update a load balancer."),
		),
	}
}

// Hooks configures actions to take before and after mutations.
func (LoadBalancer) Hooks() []ent.Hook {
	return []ent.Hook{
		hook.On(
			func(next ent.Mutator) ent.Mutator {
				return hook.LoadBalancerFunc(func(ctx context.Context, m *gen.LoadBalancerMutation) (ent.Value, error) {
					val, err := next.Mutate(ctx, m)
					if err != nil {
						return nil, err
					}

					id, ok := m.ID()
					if !ok {
						err = errors.Wrap(err, "unable to retrieve load balancer id")
						return val, err
					}

					tenant, ok := m.TenantID()
					if !ok {
						err = errors.Wrap(err, "unable to retrieve tenant id")
						return val, err
					}

					location, ok := m.LocationID()
					if !ok {
						location = "api"
					}

					provider, ok := m.ProviderID()
					if !ok {
						err = errors.Wrap(err, "unable to retrieve provider id")
						return val, err
					}

					// TODO: Add actor to context once JWT integration is complete
					//  actorStr := ctx.Value("actor").(string)

					msg, err := pubsub.NewMessage(
						tenant.String(),
						pubsub.WithEventType(types.OpToAction[m.Op()]),
						pubsub.WithSource("load-balancer-api"),
						pubsub.WithSubjectID(id.String()),
						pubsub.WithAdditionalSubjectIDs(location.String(), provider.String()),
						pubsub.WithActorID("testact-lkjasdflkjas"),
						// pubsub.WithActorID(actorStr)
					)

					if err != nil {
						err = errors.Wrap(err, "failed to create message")
						return val, err
					}

					if err := m.PubsubClient.PublishChange(ctx, types.OpToAction[m.Op()], types.TypeToSubject[gen.TypeLoadBalancer], location.String(), msg); err != nil {
						err = errors.Wrap(err, "failed to publish event")
						return val, err
					}

					return val, nil
				})
			},
			// Limit the hook only for these operations.
			ent.OpCreate|ent.OpUpdate|ent.OpDelete,
		),
	}
}
