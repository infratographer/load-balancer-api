package validations

import (
	"errors"
)

var (
	// ErrPortNameLength is returned when a port name is not between 1 and 15 characters long
	ErrPortNameLength = errors.New("port name must be between 1 and 15 characters long")

	// ErrPortNameHyphens is returned when a port name begins or ends with a hyphen
	ErrPortNameHyphens = errors.New("port name must not begin or end with a hyphen")

	// ErrPortNameAdjacentHyphens is returned when a port name contains adjacent hyphens
	ErrPortNameAdjacentHyphens = errors.New("port name must not contain adjacent hyphens")

	// ErrPortNameOneLetter is returned when a port name does not contain at least one letter
	ErrPortNameOneLetter = errors.New("port name must contain at least one letter A-Z or a-z")

	// ErrPortNameInvalidChars is returned when a port name contains invalid characters
	ErrPortNameInvalidChars = errors.New("port name must contain only letters A-Z or a-z, digits 0-9, and hyphens")
)
