package audit

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"go.infratographer.com/x/echojwtx"
)

// AuditMixin provides auditing for all records where enabled. The created_at, created_by, updated_at, and updated_by records are automatically populated when this mixin is enabled.
type AuditMixin struct {
	mixin.Schema
}

// Fields of the AuditMixin
func (AuditMixin) Fields() []ent.Field {
	return []ent.Field{
		field.String("created_by").
			Immutable().
			Optional(),
		field.String("updated_by").
			Optional(),
	}
}

// Hooks of the AuditMixin
func (AuditMixin) Hooks() []ent.Hook {
	return []ent.Hook{
		AuditHook,
	}
}

// AuditHook sets and returns the created_at, updated_at, etc., fields
func AuditHook(next ent.Mutator) ent.Mutator {
	type AuditLogger interface {
		SetCreatedBy(string)
		CreatedBy() (id string, exists bool)
		SetUpdatedBy(string)
		UpdatedBy() (id string, exists bool)
	}

	return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
		ml, ok := m.(AuditLogger)
		if !ok {
			return nil, errors.New("unexpected mutation type")
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
			fmt.Println("chacha")

		case op.Is(ent.OpUpdateOne | ent.OpUpdate):
			ml.SetUpdatedBy(actor)
		}

		return next.Mutate(ctx, m)
	})
}
