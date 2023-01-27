package pubsub

import (
	"go.infratographer.com/x/pubsubx"
)

// NewLoadBalancerURN creates a new loadbalancer URN
func NewLoadBalancerURN(loadBalancerID string) string {
	return newURN("load-balancer", loadBalancerID)
}

// NewLoadBalancerMessage creates a new assignment message
func NewLoadBalancerMessage(actorURN string, tenantURN string, loadBalancerURN string, additionalSubjectURNs ...string) (*pubsubx.Message, error) {
	return newMessage(actorURN, loadBalancerURN, additionalSubjectURNs...), nil
}
