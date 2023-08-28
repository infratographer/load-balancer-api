package validations

import "errors"

// ErrInvalidIPAddress is returned when the given string is not a valid IP address
var ErrInvalidIPAddress = errors.New("invalid ip address")
