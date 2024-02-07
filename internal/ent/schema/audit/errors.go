package audit

import (
	"errors"
	"fmt"

	"entgo.io/ent"
)

var (
	// errUnexpectedMutation is returned when an unexpected mutation is encountered
	errUnexpectedMutation = errors.New("unexpected mutatino type")
)

// newUnexpectedMutationErorr returns the UnexpectedAuditError in string format
func newUnexpectedMutationError(m ent.Mutation) error {
	return fmt.Errorf("%s: %T", errUnexpectedMutation, m)
}
