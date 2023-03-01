package pubsub

import (
	"testing"

	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestClient_AddStream(t *testing.T) {
	nc, err := nats.Connect("nats://nats:4222", nats.UserCredentials("/tmp/user.creds"))
	if err != nil {
		// fail open on nats
		t.Log(err)
	}

	js, err := nc.JetStream()
	if err != nil {
		// fail open on nats
		t.Log(err)
	}

	c1 := NewClient(
		WithJetreamContext(js),
		WithLogger(zap.NewNop().Sugar()),
		WithStreamName("load-balancer-api-test"),
		WithSubjectPrefix("com.infratographer.tests"),
	)

	out, err := c1.AddStream()

	assert.NoError(t, err)
	assert.Equal(t, "load-balancer-api-test", out.Config.Name)
	assert.Equal(t, nats.FileStorage, out.Config.Storage)
	assert.Equal(t, nats.LimitsPolicy, out.Config.Retention)
	assert.Equal(t, nats.DiscardNew, out.Config.Discard)

	// run AddStream a second time, no error should be returned
	_, err = c1.AddStream()
	assert.NoError(t, err)

	c2 := NewClient(
		WithJetreamContext(js),
		WithLogger(zap.NewNop().Sugar()),
		WithStreamName("load-balancer-api-more-tests"),
		WithSubjectPrefix("com.infratographer.tests"),
	)

	// AddStream should error since subjects overlap
	_, err = c2.AddStream()
	assert.Error(t, err)
}
