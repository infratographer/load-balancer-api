// Package validations contains validation functions for ent fields
package validations

import (
	"net"

	"go.infratographer.com/load-balancer-api/internal/config"

	"golang.org/x/exp/slices"
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
