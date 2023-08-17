package graphapi

import "errors"

// ErrPortNumberInUse is returned when a port number is already in use.
var ErrPortNumberInUse = errors.New("port number already in use")
