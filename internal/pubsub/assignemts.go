package pubsub

import (
	"fmt"

	"go.infratographer.com/x/pubsubx"
)

// NewAssignmentURN creates a new assignment URN
func NewAssignmentURN(assignmentID string) string {
	return newURN("assignment", assignmentID)
}

func newURN(kind, ID string) string {
	return fmt.Sprintf("urn:infratographer:infratographer.com:%s:%s/", kind, ID)
}

// NewAssignmentMessage creates a new assignment message
func NewAssignmentMessage(actorURN string, tenantURN string, assignmentURN string, additionalSubjectURNs ...string) (*pubsubx.Message, error) {
	return newMessage(actorURN, assignmentURN, additionalSubjectURNs...), nil
}
