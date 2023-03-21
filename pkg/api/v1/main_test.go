package api

import (
	"testing"

	natssrv "github.com/nats-io/nats-server/v2/server"
	nats "github.com/nats-io/nats.go"

	"go.infratographer.com/load-balancer-api/internal/pubsub"
)

// newNatsTestServer creates a new nats server for testing and generates a new
// stream. The returned server should be Shutdown() when testing is done.
func newNatsTestServer(t *testing.T, stream string, subs ...string) *natssrv.Server {
	srv, err := pubsub.StartNatsServer()
	if err != nil {
		t.Error(err)
	}

	nc, err := nats.Connect(srv.ClientURL())
	if err != nil {
		t.Error(err)
	}

	defer nc.Close()

	js, err := nc.JetStream()
	if err != nil {
		t.Error(err)
	}

	if _, err = js.AddStream(&nats.StreamConfig{
		Name:     stream,
		Subjects: subs,
		Storage:  nats.MemoryStorage,
	}); err != nil {
		t.Error(err)
	}

	return srv
}
