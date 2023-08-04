// Package validations contains validation functions for ent fields
package validations

import "strings"

// PortName validates a port name
func PortName(s string) error {
	// MUST be at least 1 character and no more than 15 characters long
	if len(s) > 15 || len(s) < 1 {
		return ErrPortNameLength
	}

	// MUST contain only US-ASCII [ANSI.X3.4-1986] letters 'A' - 'Z' and 'a' - 'z', digits '0' - '9', and hyphens ('-', ASCII 0x2D or decimal 45)
	if strings.ContainsAny(s, "~`!@#$%^&*()_+={[}]|\\:;\"'<,>.?/") {
		return ErrPortNameInvalidChars
	}

	// MUST contain at least one letter ('A' - 'Z' or 'a' - 'z')
	if !strings.ContainsAny(s, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ") {
		return ErrPortNameOneLetter
	}

	// hyphens MUST NOT be adjacent to other hyphens
	if strings.Contains(s, "--") {
		return ErrPortNameAdjacentHyphens
	}

	// MUST NOT begin or end with a hyphen
	if strings.HasPrefix(s, "-") || strings.HasSuffix(s, "-") {
		return ErrPortNameHyphens
	}

	return nil
}
