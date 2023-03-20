// Package pubsub wraps nats calls
package pubsub

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"go.infratographer.com/x/pubsubx"
	"go.uber.org/zap"
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
func (c *Client) PublishCreate(ctx context.Context, actor, location string, data *pubsubx.Message) error {
	data.EventType = "create"

	return c.publish(ctx, "create", actor, location, data)
}

// PublishUpdate publishes an update event
func (c *Client) PublishUpdate(ctx context.Context, actor, location string, data *pubsubx.Message) error {
	data.EventType = "update"

	return c.publish(ctx, "update", actor, location, data)
}

// PublishDelete publishes a delete event
func (c *Client) PublishDelete(ctx context.Context, actor, location string, data *pubsubx.Message) error {
	data.EventType = "delete"
	return c.publish(ctx, "delete", actor, location, data)
}

// publish publishes an event
func (c *Client) publish(ctx context.Context, action, actor, location string, data interface{}) error {
	subject := fmt.Sprintf("%s.%s.%s.%s", prefix, actor, action, location)
	c.logger.Debug("publishing nats message", zap.String("nats.subject", subject))

	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if _, err := c.js.Publish(subject, b); err != nil {
		return err
	}

	return nil
}

// ChanSubscribe creates a subcription and returns messages on a channel
func (c *Client) ChanSubscribe(ctx context.Context, sub string, ch chan *nats.Msg, stream string) (*nats.Subscription, error) {
	return c.js.ChanSubscribe(sub, ch, nats.BindStream(stream))
}
