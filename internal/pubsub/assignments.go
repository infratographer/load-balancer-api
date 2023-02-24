package pubsub

import (
	"go.infratographer.com/x/pubsubx"
)

// NewAssignmentURN creates a new assignment URN
func NewAssignmentURN(assignmentID string) string {
	return newURN("assignment", assignmentID)
}

// NewAssignmentMessage creates a new assignment message
func NewAssignmentMessage(actorURN string, tenantURN string, assignmentURN string, additionalSubjectURNs ...string) (*pubsubx.Message, error) {
	return newMessage(actorURN, assignmentURN, additionalSubjectURNs...), nil
}
