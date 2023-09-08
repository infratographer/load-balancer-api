package validations

import "errors"

// ErrInvalidIPAddress is returned when the given string is not a valid IP address
var ErrInvalidIPAddress = errors.New("invalid ip address")

// ErrRestrictedPort is returned when the given port is restricted
var ErrRestrictedPort = errors.New("port number restricted")
