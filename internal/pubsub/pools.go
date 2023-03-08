package pubsub

import (
	"go.infratographer.com/x/pubsubx"
)

// NewPoolMessage creates a new pool event message
func NewPoolMessage(actorURN string, tenantURN string, poolURN string, additionalSubjectURNs ...string) (*pubsubx.Message, error) {
	return newMessage(actorURN, poolURN, additionalSubjectURNs...), nil
}
