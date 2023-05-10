package types

import (
	"entgo.io/ent"

	"go.infratographer.com/load-balancer-api/internal/ent/generated"
	"go.infratographer.com/load-balancer-api/internal/pubsub"
)

var (
	// TypeToSubject maps ent type to pubsub subject
	TypeToSubject = make(map[string]string)
	// OpToAction maps ent.Op to pubsub action
	OpToAction = make(map[ent.Op]string)
)

func init() {
	OpToAction[ent.OpCreate] = pubsub.CreateEventType
	OpToAction[ent.OpDelete] = pubsub.DeleteEventType
	OpToAction[ent.OpDeleteOne] = pubsub.DeleteEventType
	OpToAction[ent.OpUpdate] = pubsub.UpdateEventType
	OpToAction[ent.OpUpdateOne] = pubsub.UpdateEventType

	TypeToSubject[generated.TypeLoadBalancer] = "load-balancer"
	TypeToSubject[generated.TypeOrigin] = "load-balancer-origin"
	TypeToSubject[generated.TypePool] = "load-balancer-pool"
	TypeToSubject[generated.TypePort] = "load-balancer-port"
}
