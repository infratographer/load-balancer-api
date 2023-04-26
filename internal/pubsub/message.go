package pubsub

import (
	"time"

	"go.infratographer.com/x/pubsubx"
)

const (
	// DefaultMessageSource is the default source for messages
	DefaultMessageSource = "load-balancer-api"
)

// NewMessage functionally generates a new pubsub message and appends the tenantURN
// to the list of additional subject urns
func NewMessage(tenantURN string, opts ...MsgOption) (*pubsubx.Message, error) {
	msg := pubsubx.Message{
		Timestamp: time.Now().UTC(),
		Source:    DefaultMessageSource,
	}

	for _, opt := range opts {
		opt(&msg)
	}

	msg.AdditionalSubjectURNs = append(msg.AdditionalSubjectURNs, tenantURN)

	if msg.SubjectFields == nil {
		msg.SubjectFields = make(map[string]string)
	}

	msg.SubjectFields["tenant_urn"] = tenantURN

	if err := validatePubsubMessage(&msg); err != nil {
		return nil, err
	}

	return &msg, nil
}

// MsgOption is a functional argument for NewMessage
type MsgOption func(m *pubsubx.Message)

// WithEventType sets the event type of the message
func WithEventType(e string) MsgOption {
	return func(m *pubsubx.Message) {
		m.EventType = e
	}
}

// WithSource sets the source of the message
func WithSource(s string) MsgOption {
	return func(m *pubsubx.Message) {
		m.Source = s
	}
}

// WithActorURN sets the actor urn of the message
func WithActorURN(u string) MsgOption {
	return func(m *pubsubx.Message) {
		m.ActorURN = u
	}
}

// WithSubjectURN sets the subject urn of the message
func WithSubjectURN(s string) MsgOption {
	return func(m *pubsubx.Message) {
		m.SubjectURN = s
	}
}

// WithAdditionalSubjectURNs sets the additional subject urns of the message
func WithAdditionalSubjectURNs(a ...string) MsgOption {
	return func(m *pubsubx.Message) {
		m.AdditionalSubjectURNs = a
	}
}

// WithSubjectFields sets the subject fields of the message
func WithSubjectFields(f map[string]string) MsgOption {
	return func(m *pubsubx.Message) {
		m.SubjectFields = f
	}
}

// WithAdditionalData sets the additional data of the message
func WithAdditionalData(d map[string]interface{}) MsgOption {
	return func(m *pubsubx.Message) {
		m.AdditionalData = d
	}
}

// validatePubsubMessage validates a pubsub message for required fields
func validatePubsubMessage(msg *pubsubx.Message) error {
	if msg.SubjectURN == "" {
		return ErrMissingEventSubjectURN
	}

	if msg.ActorURN == "" {
		return ErrMissingEventActorURN
	}

	if msg.Source == "" {
		return ErrMissingEventSource
	}

	return nil
}
