package pubsub

import "errors"

var (
	// ErrMissingEventSubjectURN is returned when the event subject urn is missing
	ErrMissingEventSubjectURN = errors.New("missing event subject urn")

	// ErrMissingEventActorURN is returned when the event actor urn is missing
	ErrMissingEventActorURN = errors.New("missing event actor urn")

	// ErrMissingEventSource is returned when the event source is missing
	ErrMissingEventSource = errors.New("missing event source")
)
