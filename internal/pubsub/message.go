package pubsub

import (
	"time"

	"go.infratographer.com/x/gidx"
	"go.infratographer.com/x/pubsubx"
)

const (
	// DefaultMessageSource is the default source for messages
	DefaultMessageSource = "load-balancer-api"
)

// NewMessage functionally generates a new pubsub message and appends the tenantURN
// to the list of additional subject urns
func NewMessage(tenantID string, opts ...EventOption) (*pubsubx.ChangeMessage, error) {
	msg := pubsubx.ChangeMessage{
		Timestamp: time.Now().UTC(),
		Source:    DefaultMessageSource,
	}

	for _, opt := range opts {
		opt(&msg)
	}

	tenantGID, err := gidx.Parse(tenantID)
	if err != nil {
		return nil, err
	}

	msg.AdditionalSubjectIDs = append(msg.AdditionalSubjectIDs, tenantGID)

	if msg.SubjectFields == nil {
		msg.SubjectFields = make(map[string]string)
	}

	msg.SubjectFields["tenant_id"] = tenantID

	if err := validatePubsubMessage(&msg); err != nil {
		return nil, err
	}

	return &msg, nil
}

// EventOption is a functional argument for NewMessage
type EventOption func(m *pubsubx.ChangeMessage)

// WithEventType sets the event type of the message
func WithEventType(e string) EventOption {
	return func(m *pubsubx.ChangeMessage) {
		m.EventType = e
	}
}

// WithSource sets the source of the message
func WithSource(s string) EventOption {
	return func(m *pubsubx.ChangeMessage) {
		m.Source = s
	}
}

// WithActorID sets the actor urn of the message
func WithActorID(u string) EventOption {
	return func(m *pubsubx.ChangeMessage) {
		gid, _ := gidx.Parse(u)
		m.ActorID = gid
	}
}

// WithSubjectID sets the subject urn of the message
func WithSubjectID(s string) EventOption {
	return func(m *pubsubx.ChangeMessage) {
		gid, _ := gidx.Parse(s)
		m.SubjectID = gid
	}
}

// WithAdditionalSubjectIDs sets the additional subject urns of the message
func WithAdditionalSubjectIDs(a ...string) EventOption {
	return func(m *pubsubx.ChangeMessage) {
		for _, s := range a {
			gid, _ := gidx.Parse(s)
			m.AdditionalSubjectIDs = append(m.AdditionalSubjectIDs, gid)
		}
	}
}

// WithSubjectFields sets the subject fields of the message
func WithSubjectFields(f map[string]string) EventOption {
	return func(m *pubsubx.ChangeMessage) {
		m.SubjectFields = f
	}
}

// WithAdditionalData sets the additional data of the message
func WithAdditionalData(d map[string]interface{}) EventOption {
	return func(m *pubsubx.ChangeMessage) {
		m.AdditionalData = d
	}
}

// validatePubsubMessage validates a pubsub message for required fields
func validatePubsubMessage(msg *pubsubx.ChangeMessage) error {
	if msg.SubjectID.String() == "" {
		return ErrMissingEventSubjectURN
	}

	if msg.ActorID.String() == "" {
		return ErrMissingEventActorURN
	}

	if msg.Source == "" {
		return ErrMissingEventSource
	}

	return nil
}
