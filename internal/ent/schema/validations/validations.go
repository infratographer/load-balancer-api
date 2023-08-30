// Package validations contains validation functions for ent fields
package validations

import (
	"net"
)

// IPAddress validates if the given string is a valid IP address
func IPAddress(ip string) error {
	if net.ParseIP(ip) == nil {
		return ErrInvalidIPAddress
	}

	return nil
}
