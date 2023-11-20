package graphapi

import "errors"

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
)
