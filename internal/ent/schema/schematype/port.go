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

	gen "go.infratographer.com/load-balancer-api/internal/ent/generated"
	"go.infratographer.com/load-balancer-api/internal/ent/generated/hook"
	"go.infratographer.com/load-balancer-api/internal/pubsub"
	"go.infratographer.com/load-balancer-api/internal/pubsub/types"

	"go.infratographer.com/x/entx"
	"go.infratographer.com/x/gidx"
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
			Annotations(
				entgql.OrderField("number"),
			),
		field.String("name").
			NotEmpty().
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
		entgql.RelayConnection(),
		entgql.Mutations(entgql.MutationCreate(), entgql.MutationUpdate()),
	}
}

// Hooks configures actions to take before and after mutations.
func (Port) Hooks() []ent.Hook {
	return []ent.Hook{
		// First hook.
		hook.On(
			func(next ent.Mutator) ent.Mutator {
				return hook.PortFunc(func(ctx context.Context, m *gen.PortMutation) (ent.Value, error) {
					var tenant, location, provider gidx.PrefixedID

					val, err := next.Mutate(ctx, m)
					if err != nil {
						return nil, err
					}

					id, ok := m.ID()
					if !ok {
						err = errors.Wrap(err, "unable to retrieve port id")
						return val, err
					}

					lb, ok := m.LoadBalancerID()
					if !ok {
						err = errors.Wrap(err, "unable to retrieve load balancer id")
						return val, err
					}

					lbLookup, err := m.Client().LoadBalancer.Get(ctx, lb)
					if err != nil {
						err = errors.Wrap(err, "unable to lookup load balancer")
						return val, err
					}

					if lbLookup.TenantID != "" {
						tenant = lbLookup.TenantID
					} else {
						err = errors.Wrap(err, "unable to lookup tenant id")
						return val, err
					}

					if lbLookup.LocationID != "" {
						location = lbLookup.LocationID
					} else {
						err = errors.Wrap(err, "unable to lookup location id")
						return val, err
					}

					if lbLookup.ProviderID != "" {
						provider = lbLookup.ProviderID
					} else {
						err = errors.Wrap(err, "unable to lookup provider id")
						return val, err
					}

					// TODO: Add actor to context once JWT integration is complete
					//  actorStr := ctx.Value("actor").(string)

					msg, err := pubsub.NewMessage(
						tenant.String(),
						pubsub.WithEventType(types.OpToAction[m.Op()]),
						pubsub.WithSource("load-balancer-api"),
						pubsub.WithSubjectID(id.String()),
						pubsub.WithAdditionalSubjectIDs(location.String(), provider.String(), lb.String()),
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
			ent.OpCreate|ent.OpUpdate|ent.OpUpdateOne|ent.OpDelete|ent.OpDeleteOne,
		),
	}
}
