package client

import (
	"errors"
)

var (
	// ErrUnauthorized returned when the request is not authorized
	ErrUnauthorized = errors.New("client is unauthorized")

	// ErrPermissionDenied returned when the subject does not permissions to access the resource
	ErrPermissionDenied = errors.New("client does not have permissions")

	// ErrLBNotfound returned when the load balancer ID not found
	ErrLBNotfound = errors.New("loadbalancer ID not found")

	// ErrHTTPError returned when the http response is an error
	ErrHTTPError = errors.New("loadbalancer api http error")

	// ErrInternalServerError returned when the server returns an internal server error
	ErrInternalServerError = errors.New("internal server error")

	// ErrMetadataStatusNotFound returned when the status data is invalid
	ErrMetadataStatusNotFound = errors.New("metadata status not found")

	// ErrLocationNotFound returned when the status data is invalid
	ErrLocationNotFound = errors.New("location not found")
)
