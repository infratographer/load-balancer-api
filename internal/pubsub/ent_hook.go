package pubsub

import (
	"context"

	"entgo.io/ent"
	"github.com/pkg/errors"

	"go.infratographer.com/x/gidx"

	"go.infratographer.com/load-balancer-api/internal/ent/generated"
)

var (
	typeToSubject = make(map[string]string)
	opToAction    = make(map[ent.Op]string)
)

func init() {
	opToAction[ent.OpCreate] = CreateEventType
	opToAction[ent.OpDelete] = DeleteEventType
	opToAction[ent.OpDeleteOne] = DeleteEventType
	opToAction[ent.OpUpdate] = UpdateEventType
	opToAction[ent.OpUpdateOne] = UpdateEventType

	typeToSubject[generated.TypeLoadBalancer] = "load-balancer"
	typeToSubject[generated.TypeOrigin] = "load-balancer-origin"
	typeToSubject[generated.TypePool] = "load-balancer-pool"
	typeToSubject[generated.TypePort] = "load-balancer-port"
}

// Hooks is an ent client hook
func (c *Client) Hooks(next ent.Mutator) ent.Mutator {
	return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
		c.logger.Sugar().Debugw("clienthook", "Op", m.Op(), "Type", m.Type())

		// TODO: Add actor to context once JWT integration is complete
		//  actorStr := ctx.Value("actor").(string)

		fieldMap := make(map[string]ent.Value)

		for _, field := range m.Fields() {
			v, _ := m.Field(field)
			fieldMap[field] = v
		}
		// Do mutation, then publish events
		val, err := next.Mutate(ctx, m)
		if err != nil {
			return val, err
		}

		if fieldMap["tenant_id"] == nil && opToAction[m.Op()] != DeleteEventType {
			if m.Type() == generated.TypeProvider {
				p, _ := c.ec.Provider.Get(ctx, fieldMap["id"].(gidx.PrefixedID))
				fieldMap["tenant_id"] = p.TenantID
			}
		}

		msg, err := NewMessage(fieldMap["tenant_id"].(gidx.PrefixedID).String(),
			WithEventType(opToAction[m.Op()]),
			WithSource("load-balancer-api"),
		)

		c.logger.Sugar().Infow("!!!!!!!!!!mutation", "fields", fieldMap)

		location, ok := fieldMap["location_id"]
		if !ok {
			location = "api"
		}

		if e := c.publish(
			ctx,
			opToAction[m.Op()],
			typeToSubject[generated.TypeLoadBalancer],
			location.(string),
			msg,
		); e != nil {
			err = errors.Wrap(e, "failed to publish event")

			return val, err
		}

		return val, nil
	})
}
