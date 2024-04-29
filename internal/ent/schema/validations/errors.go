package validations

import (
	"errors"
)

// ErrInvalidIPAddress is returned when the given string is not a valid IP address
var ErrInvalidIPAddress = errors.New("invalid ip address")

// ErrRestrictedPort is returned when the given port is restricted
var ErrRestrictedPort = errors.New("port number restricted")

// ErrFieldEmpty is returned when a field is empty
var ErrFieldEmpty = errors.New("must not be empty")
