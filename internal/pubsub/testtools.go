//go:build testtools
// +build testtools

package pubsub

import (
	"errors"
	"fmt"
	"testing"
	"time"

	natssrv "github.com/nats-io/nats-server/v2/server"
	natsgo "github.com/nats-io/nats.go"
)

const (
	natsTimeout = 2 * time.Second
)

// StartNatsServer creates a new Nats server in memory.
// If stream subjects are passed, a new stream will be created
// with all subjects, using the first subject as the stream name.
func StartNatsServer() (*natssrv.Server, error) {
	const maxControlLine = 2048

	s, err := natssrv.NewServer(&natssrv.Options{
		Host:           "127.0.0.1",
		Debug:          false,
		Trace:          false,
		TraceVerbose:   false,
		Port:           natssrv.RANDOM_PORT,
		NoLog:          false,
		NoSigs:         true,
		MaxControlLine: maxControlLine,
		JetStream:      true,
	})
	if err != nil {
		return nil, fmt.Errorf("building nats server: %w", err)
	}

	// uncomment to enable nats server logging
	// s.ConfigureLogger()

	if err := natssrv.Run(s); err != nil {
		return nil, err
	}

	if !s.ReadyForConnections(natsTimeout) {
		return nil, errors.New("starting nats server: timeout") //nolint:goerr113
	}

	return s, nil
}

// WaitConnected waits the timeout for a connection
func WaitConnected(t *testing.T, c *natsgo.Conn) {
	t.Helper()

	const defaultWaitTime = 25 * time.Millisecond

	timeout := time.Now().Add(natsTimeout)
	for time.Now().Before(timeout) {
		if c.IsConnected() {
			return
		}

		time.Sleep(defaultWaitTime)
	}

	t.Fatal("client connecting timeout")
}
