package locations

import "errors"

var (
	// ErrAlreadyDeleted is returned when a location is already deleted
	ErrAlreadyDeleted = errors.New("location already deleted")

	// ErrInvalid is a generic invalid response
	ErrInvalid = errors.New("invalid location")

	// ErrNameRequired is returned when a location name is not provided
	ErrNameRequired = errors.New("name is required")

	// ErrNotFound is returned when a location is not found
	ErrNotFound = errors.New("location not found")

	// ErrTenantIDRequired is returned when a tenant ID is not provided
	ErrTenantIDRequired = errors.New("tenant ID is required")

	// ErrWrite is returned when a write operation fails
	ErrWrite = errors.New("failed to write location")
)
