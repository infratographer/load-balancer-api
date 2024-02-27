package audit

import (
	"errors"
)

var (
	// errUnexpectedMutation is returned when an unexpected mutation is encountered
	errUnexpectedMutation = errors.New("unexpected mutation type")
)
