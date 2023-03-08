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
func StartNatsServer(streamSubjects ...string) (*natssrv.Server, error) {
	const maxControlLine = 2048

	s, err := natssrv.NewServer(&natssrv.Options{
		Host:           "127.0.0.1",
		Port:           natssrv.RANDOM_PORT,
		NoLog:          true,
		NoSigs:         true,
		MaxControlLine: maxControlLine,
		JetStream:      true,
	})
	if err != nil {
		return nil, fmt.Errorf("building nats server: %w", err)
	}

	//nolint:errcheck // we don't care about the error here
	go natssrv.Run(s)

	if !s.ReadyForConnections(natsTimeout) {
		return nil, errors.New("starting nats server: timeout") //nolint:goerr113
	}

	if len(streamSubjects) != 0 {
		nc, err := natsgo.Connect(s.ClientURL())
		if err != nil {
			return nil, fmt.Errorf("stream seed failed to connect to server: %w", err)
		}

		defer nc.Close()

		js, err := nc.JetStream()
		if err != nil {
			return nil, fmt.Errorf("stream seed failed to establish JetStream: %w", err)
		}

		_, err = js.AddStream(&natsgo.StreamConfig{
			Name:     streamSubjects[0],
			Subjects: []string{streamSubjects[0] + ".>"},
			Storage:  natsgo.MemoryStorage,
		})
		if err != nil {
			return nil, fmt.Errorf("stream seed failed to create JetStream: %w", err)
		}
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
