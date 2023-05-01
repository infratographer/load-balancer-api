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

// NewPortURN creates a new port URN
func NewPortURN(id string) string {
	return newURN("load-balancer-port", id)
}

// NewOriginURN creates a new origin URN
func NewOriginURN(id string) string {
	return newURN("load-balancer-origin", id)
}

// NewMetadataURN creates a new metadata URN
func NewMetadataURN(id string) string {
	return newURN("load-balancer-metadata", id)
}

// NewPoolURN creates a new pool URN
func NewPoolURN(id string) string {
	return newURN("load-balancer-pool", id)
}

// NewAssignmentURN creates a new assignment URN
func NewAssignmentURN(id string) string {
	return newURN("load-balancer-assignment", id)
}
