// Package pubsub wraps nats calls
package pubsub

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/nats-io/nats.go"
	"github.com/spf13/viper"
	"go.infratographer.com/x/pubsubx"
	"go.uber.org/zap"
)

const (
	// CreateEventType is the create event type string
	CreateEventType = "create"
	// DeleteEventType is the delete event type string
	DeleteEventType = "delete"
	// UpdateEventType is the update event type string
	UpdateEventType = "update"
)

// PublishChange sets the action of the message and then publishes the message
func (c *Client) PublishChange(ctx context.Context, action, subject, location string, data *pubsubx.ChangeMessage) error {
	data.EventType = action
	return c.publish(ctx, action, "changes", subject, location, data)
}

// PublishEvent sets the action of the message and then publishes the message
func (c *Client) PublishEvent(ctx context.Context, action, subject, location string, data *pubsubx.EventMessage) error {
	data.EventType = action
	return c.publish(ctx, action, "events", subject, location, data)
}

// publish publishes an event
func (c *Client) publish(_ context.Context, action, eventType, subject, location string, data interface{}) error {
	prefix := viper.GetString("nats.subject-prefix")
	natsSubject := fmt.Sprintf("%s.%s.%s.%s.%s.%s", prefix, eventType, action, subject, "location", location)
	c.logger.Debug("publishing nats message", zap.String("nats.subject", natsSubject))

	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if _, err := c.js.Publish(natsSubject, b); err != nil {
		return err
	}

	return nil
}

// ChanSubscribe creates a subcription and returns messages on a channel
func (c *Client) ChanSubscribe(_ context.Context, sub string, ch chan *nats.Msg, stream string) (*nats.Subscription, error) {
	return c.js.ChanSubscribe(sub, ch, nats.BindStream(stream))
}
