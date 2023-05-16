// Package pubsub wraps nats calls
package pubsub

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/nats-io/nats.go"
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

// May be a config option later
var prefix = "com.infratographer.events"

// PublishCreate publishes a create event
func (c *Client) PublishCreate(ctx context.Context, actor, location string, data *pubsubx.Message) error {
	data.EventType = CreateEventType

	return c.publish(ctx, CreateEventType, actor, location, data)
}

// PublishUpdate publishes an update event
func (c *Client) PublishUpdate(ctx context.Context, actor, location string, data *pubsubx.Message) error {
	data.EventType = UpdateEventType

	return c.publish(ctx, UpdateEventType, actor, location, data)
}

// PublishDelete publishes a delete event
func (c *Client) PublishDelete(ctx context.Context, actor, location string, data *pubsubx.Message) error {
	data.EventType = DeleteEventType
	return c.publish(ctx, DeleteEventType, actor, location, data)
}

// publish publishes an event
func (c *Client) publish(_ context.Context, action, actor, location string, data interface{}) error {
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
func (c *Client) ChanSubscribe(_ context.Context, sub string, ch chan *nats.Msg, stream string) (*nats.Subscription, error) {
	return c.js.ChanSubscribe(sub, ch, nats.BindStream(stream))
}

// NewSubject builds a subject string based upon the provided
func (c *Client) NewSubject(subjects ...string) string {
	return fmt.Sprintf("%s.%s", c.prefix, strings.Join(subjects, "."))
}

func prepareMessage(subject string, data interface{}) (*nats.Msg, error) {
	msg := nats.NewMsg(subject)

	b, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	msg.Data = b

	return msg, nil
}

// PublishChange publishes a change message
func (c *Client) PublishChange(_ context.Context, subject string, data pubsubx.ChangeMessage) error {
	c.logger.Debug("publishing change message", zap.String("subject", subject))

	fmt.Println("hello")

	msg, err := prepareMessage(subject, data)
	if err != nil {
		return err
	}

	msg.Header.Set("X-INFRA9-MSG-TYPE", "change")

	if _, err := c.js.PublishMsg(msg); err != nil {
		return err
	}

	return nil
}
