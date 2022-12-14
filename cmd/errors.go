package cmd

import "errors"

var (
	// ErrAuditFilePathRequired is returned when a audit file path is missing
	ErrAuditFilePathRequired = errors.New("audit file path is required and cannot be empty")
)
