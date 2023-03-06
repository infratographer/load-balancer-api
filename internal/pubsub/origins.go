package pubsub

import (
	"go.infratographer.com/x/pubsubx"
)

// NewOriginMessage creates a new origin message
func NewOriginMessage(actorURN string, tenantURN string, originURN string, additionalSubjectURNs ...string) (*pubsubx.Message, error) {
	return newMessage(actorURN, originURN, additionalSubjectURNs...), nil
}
