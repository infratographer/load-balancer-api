package pubsub

import (
	"go.infratographer.com/x/pubsubx"
)

// NewAssignmentMessage creates a new assignment event message
func NewAssignmentMessage(actorURN string, tenantURN string, assignmentURN string, additionalSubjectURNs ...string) (*pubsubx.Message, error) {
	return newMessage(actorURN, assignmentURN, additionalSubjectURNs...), nil
}
