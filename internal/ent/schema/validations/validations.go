// Package validations contains validation functions for ent fields
package validations

import (
	"fmt"
	"net"
	"regexp"
	"strings"

	"go.infratographer.com/load-balancer-api/internal/config"

	"golang.org/x/exp/slices"
)

const (
	// maxNameLength is the maximum length of a name field
	maxNameLength = 64
	// unallowedChars is a string containing the unallowed characters in a string for a name field
	unallowedChars = `<>&'"`
)

// IPAddress validates if the given string is a valid IP address
func IPAddress(ip string) error {
	if net.ParseIP(ip) == nil {
		return ErrInvalidIPAddress
	}

	return nil
}

// RestrictedPorts validates if the given port is restricted
func RestrictedPorts(port int) error {
	if slices.Contains(config.AppConfig.RestrictedPorts, port) {
		return ErrRestrictedPort
	}

	return nil
}

// NameField validates the name field
func NameField(name string) error {
	if len(strings.TrimSpace(name)) == 0 {
		return ErrFieldEmpty
	}

	if len(name) > maxNameLength {
		return fmt.Errorf("must not be longer than %d characters", maxNameLength) // nolint: goerr113
	}

	if containsUnallowedChars(name) {
		return fmt.Errorf("must not contain the following characters: %s", unallowedChars) // nolint: goerr113
	}

	return nil
}

// OptionalNameField validates the name field when optional
func OptionalNameField(name string) error {
	if len(name) > maxNameLength {
		return fmt.Errorf("must not be longer than %d characters", maxNameLength) // nolint: goerr113
	}

	if containsUnallowedChars(name) {
		return fmt.Errorf("must not contain the following characters: %s", unallowedChars) // nolint: goerr113
	}

	return nil
}

func containsUnallowedChars(s string) bool {
	// Define the regular expression pattern to match the unallowed characters
	pattern := "[" + unallowedChars + "]"

	// Compile the regular expression pattern
	regex := regexp.MustCompile(pattern)

	// Check if the pattern matches any part of the string
	return regex.MatchString(s)
}
