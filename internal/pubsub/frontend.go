package pubsub

import (
	"go.infratographer.com/x/pubsubx"
)

// NewFrontendMessage creates a new frontend message
func NewFrontendMessage(actorURN string, tenantURN string, frontendURN string, additionalSubjectURNs ...string) (*pubsubx.Message, error) {
	return newMessage(actorURN, frontendURN, additionalSubjectURNs...), nil
}
