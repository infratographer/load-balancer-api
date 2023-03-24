package pubsub

import (
	"go.infratographer.com/x/pubsubx"
)

// NewPortMessage creates a new port message
func NewPortMessage(actorURN string, tenantURN string, portURN string, additionalSubjectURNs ...string) (*pubsubx.Message, error) {
	return newMessage(actorURN, portURN, additionalSubjectURNs...), nil
}
