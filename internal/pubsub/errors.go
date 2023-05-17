package pubsub

import "errors"

var (
	// ErrMissingEventSubjectID is returned when the event subject id is missing
	ErrMissingEventSubjectID = errors.New("missing event subject id")

	// ErrMissingEventActorID is returned when the event actor id is missing
	ErrMissingEventActorID = errors.New("missing event actor id")

	// ErrMissingEventSource is returned when the event source is missing
	ErrMissingEventSource = errors.New("missing event source")
)
