// Package pubsub wraps nats calls
package pubsub

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"

	"go.infratographer.com/x/pubsubx"
)

// May be a config option later
var prefix = "com.infratographer.events"

func newMessage(actorURN string, subjectURN string, additionalSubjectURNs ...string) *pubsubx.Message {
	return &pubsubx.Message{
		SubjectURN:            subjectURN,
		ActorURN:              actorURN, // comes from the jwt eventually
		Timestamp:             time.Now().UTC(),
		Source:                "lbapi",
		AdditionalSubjectURNs: additionalSubjectURNs,
	}
}

// PublishCreate publishes a create event
func PublishCreate(ctx context.Context, js nats.JetStreamContext, actor, location string, data *pubsubx.Message) error {
	data.EventType = "create"

	return publish(ctx, js, "create", actor, location, data)
}

// PublishUpdate publishes an update event
func PublishUpdate(ctx context.Context, js nats.JetStreamContext, actor, location string, data *pubsubx.Message) error {
	data.EventType = "update"

	return publish(ctx, js, "update", actor, location, data)
}

// PublishDelete publishes a delete event
func PublishDelete(ctx context.Context, js nats.JetStreamContext, actor, location string, data *pubsubx.Message) error {
	data.EventType = "delete"

	return publish(ctx, js, "delete", actor, location, data)
}

// publish publishes an event
func publish(ctx context.Context, js nats.JetStreamContext, action, actor, location string, data interface{}) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	subject := fmt.Sprintf("%s.%s.%s.%s", prefix, actor, action, location)
	if _, err := js.Publish(subject, b); err != nil {
		return err
	}

	return nil
}
