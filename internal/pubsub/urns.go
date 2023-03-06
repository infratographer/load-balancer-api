package pubsub

import (
	"fmt"
)

func newURN(kind, id string) string {
	return fmt.Sprintf("urn:infratographer:%s:%s", kind, id)
}

// NewTenantURN creates a new tenant URN
func NewTenantURN(id string) string {
	return newURN("tenant", id)
}

// NewLoadBalancerURN creates a new load balancer URN
func NewLoadBalancerURN(id string) string {
	return newURN("load-balancer", id)
}

// NewFrontendURN creates a new frontend URN
func NewFrontendURN(id string) string {
	return newURN("load-balancer-frontend", id)
}

// NewOriginURN creates a new origin URN
func NewOriginURN(id string) string {
	return newURN("load-balancer-origin", id)
}

// NewPoolURN creates a new assignment URN
func NewPoolURN(id string) string {
	return newURN("load-balancer-pool", id)
}

// NewAssignmentURN creates a new assignment URN
func NewAssignmentURN(id string) string {
	return newURN("load-balancer-assignment", id)
}
