package graphapi

import "errors"

// ErrPortNumberInUse is returned when a port number is already in use.
var ErrPortNumberInUse = errors.New("port number already in use")

// ErrRestrictedPortNumber is returned when a port number is restricted.
var ErrRestrictedPortNumber = errors.New("port number restricted")

// ErrOwnerConflict is returned when ownerships of resources are conflicting
var ErrOwnerConflict = errors.New("conflicting ownership")
