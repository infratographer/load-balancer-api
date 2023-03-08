package pubsub

import (
	"go.infratographer.com/x/pubsubx"
)

// NewLoadBalancerMessage creates a new loadbalancer event message
func NewLoadBalancerMessage(actorURN string, tenantURN string, loadBalancerURN string, additionalSubjectURNs ...string) (*pubsubx.Message, error) {
	return newMessage(actorURN, loadBalancerURN, additionalSubjectURNs...), nil
}
