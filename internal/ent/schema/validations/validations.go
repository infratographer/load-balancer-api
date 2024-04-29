// Package validations contains validation functions for ent fields
package validations

import (
	"fmt"
	"net"
	"strings"

	"go.infratographer.com/load-balancer-api/internal/config"

	"golang.org/x/exp/slices"
)

const (
	// maxNameLength is the maximum length of a name field
	maxNameLength = 64
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

	return nil
}
