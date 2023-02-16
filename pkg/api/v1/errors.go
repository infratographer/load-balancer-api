package api

import "errors"

var (
	// ErrMissingPoolID is returned when a pool ID is missing
	ErrMissingPoolID = errors.New("pool ID is missing")

	// ErrMissingOriginTarget is returned when an origin target is missing
	ErrMissingOriginTarget = errors.New("origin target is missing")

	// ErrLoadBalancerIPMissing is returned when a load balancer IP is missing
	ErrLoadBalancerIPMissing = errors.New("load balancer IP is missing")

	// ErrLoadBalancerIPInvalid is returned when a load balancer IP is invalid
	ErrLoadBalancerIPInvalid = errors.New("load balancer IP is invalid")

	// ErrPoolProtocolInvalid is returned when a protocol is invalid
	ErrPoolProtocolInvalid = errors.New("protocol must be tcp")

	// ErrLoadBalancerIDMissing is returned when a load balancer ID is missing
	ErrLoadBalancerIDMissing = errors.New("load balancer ID is missing")

	// ErrPortOutOfRange is returned when a port is out of 1-65535 range
	ErrPortOutOfRange = errors.New("port is out of range")

	// ErrEmptyPayload is returned when a payload is empty
	ErrEmptyPayload = errors.New("payload is empty")

	// ErrAlreadyDeleted is returned when a load balancer is already deleted
	ErrAlreadyDeleted = errors.New("load balancer already deleted")

	// ErrAmbiguous is returned when a request is ambiguous
	ErrAmbiguous = errors.New("request is ambiguous")

	// ErrInvalidLoadBalancer is a generic invalid response
	ErrInvalidLoadBalancer = errors.New("invalid loadbalancer")

	// ErrInvalidUUID is returned when a UUID is invalid
	ErrInvalidUUID = errors.New("invalid UUID")

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

	// ErrIDRequired is returned when a location ID is not provided
	ErrIDRequired = errors.New("ID is required")

	// ErrNullUUID is returned when a UUID is null
	ErrNullUUID = errors.New("UUID is null")

	// ErrUUIDNotFound is returned when a UUID is not found in the path
	ErrUUIDNotFound = errors.New("UUID not found in path")

	// ErrUnauthenticatedRequest is returned when a request is not authenticated
	ErrUnauthenticatedRequest = errors.New("unauthenticated request")

	// ErrWrite is returned when a write operation fails
	ErrWrite = errors.New("failed to write location")

	// ErrDisplayNameMissing is returned when a display name is missing
	ErrDisplayNameMissing = errors.New("display name is missing")
)
