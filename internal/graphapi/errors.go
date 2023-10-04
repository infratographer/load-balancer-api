package graphapi

import "errors"

// ErrPortNumberInUse is returned when a port number is already in use.
var ErrPortNumberInUse = errors.New("port number already in use")

// ErrRestrictedPortNumber is returned when a port number is restricted.
var ErrRestrictedPortNumber = errors.New("port number restricted")

// ErrPortNotFound is returned when one or more ports are not found
var ErrPortNotFound = errors.New("one or more ports not found")

// ErrPoolNotFound is returned when one or more pools are not found
var ErrPoolNotFound = errors.New("one or more pools not found")
