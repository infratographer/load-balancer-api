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

// Pool holds the schema definition for the Pool entity.
type Pool struct {
	ent.Schema
}

// Mixin to use for Pool type
func (Pool) Mixin() []ent.Mixin {
	return []ent.Mixin{
		entx.NewTimestampMixin(),
	}
}

// Fields of the Instance.
func (Pool) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			GoType(gidx.PrefixedID("")).
			DefaultFunc(func() gidx.PrefixedID { return gidx.MustNewID(PoolPrefix) }).
			Unique().
			Immutable(),
		field.String("name").
			NotEmpty().
			Annotations(
				entgql.OrderField("name"),
			),
		field.Enum("protocol").
			Values("tcp", "udp").
			Annotations(
				entgql.OrderField("protocol"),
			),
		field.String("owner_id").
			GoType(gidx.PrefixedID("")).
			Immutable().
			NotEmpty().
			Annotations(
				entgql.Type("ID"),
				entgql.Skip(entgql.SkipWhereInput, entgql.SkipMutationUpdateInput),
				pubsubinfo.EventsHookAdditionalSubject("owner"),
			),
	}
}

// Edges of the Instance.
func (Pool) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("ports", Port.Type),
		edge.From("origins", Origin.Type).
			Ref("pool").
			Annotations(
				entgql.RelayConnection(),
			),
	}
}

// Indexes of the Pool
func (Pool) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("owner_id"),
	}
}

// Annotations for the Pool
func (Pool) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entx.GraphKeyDirective("id"),
		pubsubinfo.EventsHookSubjectName("load-balancer-pool"),
		entgql.Type("LoadBalancerPool"),
		prefixIDDirective(PoolPrefix),
		entgql.RelayConnection(),
		entgql.QueryField(),
		entgql.Mutations(entgql.MutationCreate(), entgql.MutationUpdate()),
	}
}
