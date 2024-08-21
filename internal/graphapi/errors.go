package graphapi

import (
	"errors"
	"fmt"
)

var (
	// ErrPortNumberInUse is returned when a port number is already in use.
	ErrPortNumberInUse = errors.New("port number already in use")

	// ErrRestrictedPortNumber is returned when a port number is restricted.
	ErrRestrictedPortNumber = errors.New("port number restricted")

	// ErrPortNotFound is returned when one or more ports are not found
	ErrPortNotFound = errors.New("one or more ports not found")

	// ErrPoolNotFound is returned when one or more pools are not found
	ErrPoolNotFound = errors.New("one or more pools not found")

	// ErrInternalServerError is returned when an internal error occurs.
	ErrInternalServerError = errors.New("internal server error")

	// ErrLoadBalancerLimitReached is returned when the load balancer limit has been reached for an owner.
	ErrLoadBalancerLimitReached = errors.New("load balancer limit reached")

	// ErrFieldEmpty is returned when a required field is empty.
	ErrFieldEmpty = errors.New("must not be empty")

	// ErrInvalidCharacters is returned when an invalid input is provided
	ErrInvalidCharacters = errors.New("valid characters are A-Z a-z 0-9 _ -")
)

// ErrInvalidField is returned when an invalid input is provided.
type ErrInvalidField struct {
	field string
	err   error
}

// Error implements the error interface.
func (e *ErrInvalidField) Error() string {
	return fmt.Sprintf("%s: %v", e.field, e.err)
}

func newInvalidFieldError(field string, err error) *ErrInvalidField {
	return &ErrInvalidField{field: field, err: err}
}
