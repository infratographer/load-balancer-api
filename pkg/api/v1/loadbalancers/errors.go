package loadbalancers

import "errors"

var (
	// ErrAlreadyDeleted is returned when a load balancer is already deleted
	ErrAlreadyDeleted = errors.New("load balancer already deleted")

	// ErrInvalid is a generic invalid response
	ErrInvalid = errors.New("invalid loadbalancer")

	// ErrIPAddressRequired is returned when a IP Address is not provided
	ErrIPAddressRequired = errors.New("ip address is required")

	// ErrNameRequired is returned when a location name is not provided
	ErrNameRequired = errors.New("name is required")

	// ErrNotFound is returned when a ip address is not found
	ErrNotFound = errors.New("ip address not found")

	// ErrTenantIDRequired is returned when a tenant ID is not provided
	ErrTenantIDRequired = errors.New("tenant ID is required")

	// ErrTypeInvalid is returned when a type is not valid
	ErrTypeInvalid = errors.New("type is invalid")

	// ErrTypeRequired is returned when a type is not provided
	ErrTypeRequired = errors.New("type is required")

	// ErrLocationIDRequired is returned when a location ID is not provided
	ErrLocationIDRequired = errors.New("location ID is required")

	// ErrIPAddressInvalid is returned when a IP Address is not valid
	ErrIPAddressInvalid = errors.New("ip address is invalid")

	// ErrIPv4Required is returned when a IP Address is not valid
	ErrIPv4Required = errors.New("ip address is invalid")

	// ErrSizeRequired is returned when a size is not provided
	ErrSizeRequired = errors.New("size is required")

	// ErrWrite is returned when a write operation fails
	ErrWrite = errors.New("failed to write loadbalancer")
)
