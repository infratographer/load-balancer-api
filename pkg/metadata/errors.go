package metadata

import "errors"

var (
	// ErrStatusNotFound is returned when a status is not found in the payload
	ErrStatusNotFound = errors.New("status not found")

	// ErrInvalidStatusData is returned when the status json data is invalid
	ErrInvalidStatusData = errors.New("invalid status json data")
)
