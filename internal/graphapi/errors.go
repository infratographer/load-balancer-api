package graphapi

import "errors"

var (
	// ErrPortNumberInUse is returned when a port number is already in use.
	ErrPortNumberInUse = errors.New("port number already in use")

	// ErrRestrictedPortNumber is returned when a port number is restricted.
	ErrRestrictedPortNumber = errors.New("port number restricted")

	// ErrInternalServerError is returned when an internal error occurs.
	ErrInternalServerError = errors.New("internal server error")
)
