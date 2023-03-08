package api

import (
	"testing"

	natssrv "github.com/nats-io/nats-server/v2/server"

	"go.infratographer.com/load-balancer-api/internal/pubsub"
)

var natsSrv *natssrv.Server

func TestMain(m *testing.M) {
	srv, err := pubsub.StartNatsServer()
	if err != nil {
		panic(err)
	}

	natsSrv = srv

	defer natsSrv.Shutdown()

	m.Run()
}
