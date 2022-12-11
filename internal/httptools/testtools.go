// Package httptools provides a set of tools for testing HTTP handlers.
package httptools

import (
	"bytes"
	"io"
)

// FakeBody returns a fake io.ReadCloser for testing purposes.
func FakeBody(payload string) io.ReadCloser {
	return io.NopCloser(bytes.NewBufferString(payload))
}
