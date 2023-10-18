// Package testutils provides some utilities that may be useful for testing
package testutils

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"

	"go.infratographer.com/permissions-api/pkg/permissions/mockpermissions"
)

// MockPermissions creates a context from the given context with mocks for permission-api methods
func MockPermissions(ctx context.Context) context.Context {
	// mock permissions
	perms := new(mockpermissions.MockPermissions)
	perms.On("CreateAuthRelationships", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	perms.On("DeleteAuthRelationships", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	ctx = perms.ContextWithHandler(ctx)

	return ctx
}

// IfErrPanic conditionally panics on err with msg
func IfErrPanic(msg string, err error) {
	if err != nil {
		log.Panicf("%s err: %s", msg, err.Error())
	}
}

// ChannelReceiveWithTimeout returns the next message from channel chan or panics if it timesout before
func ChannelReceiveWithTimeout[T any](t *testing.T, channel <-chan T, timeout time.Duration) (msg T) {
	select {
	case msg = <-channel:
	case <-time.After(timeout):
		t.Fatal("timed out waiting to receive from channel")
	}

	return
}
