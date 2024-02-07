package audit

import (
	"context"

	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"go.infratographer.com/x/echojwtx"
)

// Mixin provides auditing for all records where enabled. The created_at, created_by, updated_at, and updated_by records are automatically populated when this mixin is enabled.
type Mixin struct {
	mixin.Schema
}

// Fields of the Mixin
func (Mixin) Fields() []ent.Field {
	return []ent.Field{
		field.String("created_by").
			Immutable().
			Annotations(
				entgql.Skip(entgql.SkipMutationCreateInput, entgql.SkipMutationUpdateInput),
				entgql.OrderField("CREATED_BY"),
			).
			Optional(),
		field.String("updated_by").
			Annotations(
				entgql.Skip(entgql.SkipMutationCreateInput, entgql.SkipMutationUpdateInput),
				entgql.OrderField("UPDATED_BY"),
			).
			Optional(),
	}
}

// Hooks of the Mixin
func (Mixin) Hooks() []ent.Hook {
	return []ent.Hook{
		Hook,
	}
}

// Hook sets and returns the created_at, updated_at, etc., fields
func Hook(next ent.Mutator) ent.Mutator {
	type AuditLogger interface {
		SetCreatedBy(string)
		CreatedBy() (id string, exists bool)
		SetUpdatedBy(string)
		UpdatedBy() (id string, exists bool)
	}

	return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
		ml, ok := m.(AuditLogger)
		if !ok {
			return nil, errUnexpectedMutation
		}
		actor := "unknown-actor"
		id, ok := ctx.Value(echojwtx.ActorCtxKey).(string)
		if ok {
			actor = id
		}

		switch op := m.Op(); {
		case op.Is(ent.OpCreate):
			ml.SetCreatedBy(actor)
			ml.SetUpdatedBy(actor)

		case op.Is(ent.OpUpdateOne | ent.OpUpdate):
			ml.SetUpdatedBy(actor)
		}

		return next.Mutate(ctx, m)
	})
}
