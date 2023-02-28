package pubsub

import (
	"errors"

	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

// Client is an event bus client with some configuration
type Client struct {
	js             nats.JetStreamContext
	logger         *zap.SugaredLogger
	prefix, stream string
}

// Option is a functional configuration option for governor eventing
type Option func(c *Client)

// NewClient configures and establishes a new event bus client connection
func NewClient(opts ...Option) *Client {
	client := Client{
		logger: zap.NewNop().Sugar(),
	}

	for _, opt := range opts {
		opt(&client)
	}

	return &client
}

// WithJetreamContext sets the nats jetstream context
func WithJetreamContext(js nats.JetStreamContext) Option {
	return func(c *Client) {
		c.js = js
	}
}

// WithStreamName sets the nats stream name
func WithStreamName(s string) Option {
	return func(c *Client) {
		c.stream = s
	}
}

// WithSubjectPrefix sets the nats subject prefix
func WithSubjectPrefix(p string) Option {
	return func(c *Client) {
		c.prefix = p
	}
}

// WithLogger sets the client logger
func WithLogger(l *zap.SugaredLogger) Option {
	return func(c *Client) {
		c.logger = l
	}
}

// AddStream checks if a stream exists and attempts to create it if it doesn't. Currently we don't
// currently check that the stream is configured identically to the desired configuration.
func (c *Client) AddStream() (*nats.StreamInfo, error) {
	info, err := c.js.StreamInfo(c.stream)
	if err == nil {
		c.logger.Debug("got info for stream, assuming stream exists", zap.Any("nats.stream.info", info.Config))
		return info, nil
	} else if !errors.Is(err, nats.ErrStreamNotFound) {
		return nil, err
	}

	c.logger.Debug("nats stream not found, attempting to create it")

	return c.js.AddStream(&nats.StreamConfig{
		Name: c.stream,
		Subjects: []string{
			c.prefix + ".<",
		},
		Storage:   nats.FileStorage,
		Retention: nats.LimitsPolicy,
		Discard:   nats.DiscardNew,
	})
}
