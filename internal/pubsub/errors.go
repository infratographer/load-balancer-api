package pubsub

import "errors"

var (
	ErrInvalidActorURN      = errors.New("invalid actor urn")
	ErrInvalidTenantURN     = errors.New("invalid tenant urn")
	ErrInvalidAssignmentURN = errors.New("invalid assignment urn")
	ErrInvalidURN           = errors.New("invalid urn")
)
