//go:build testtools
// +build testtools

package api

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"go.infratographer.com/load-balancer-api/internal/config"
	"go.infratographer.com/load-balancer-api/internal/httptools"
	"go.infratographer.com/load-balancer-api/internal/pubsub"
	"go.infratographer.com/load-balancer-api/internal/x/echox"
	"go.infratographer.com/x/crdbx"
	"go.uber.org/zap"
)

type httpTest struct {
	name   string
	body   string
	tenant string
	method string
	path   string
	status int
}

func newTestServer(t *testing.T, natsURL string) *httptest.Server {
	db, err := crdbx.NewDB(config.AppConfig.CRDB, false)
	if err != nil {
		t.Fatal(err)
	}

	dbx := sqlx.NewDb(db, "postgres")
	e := echox.NewServer()

	nc, err := nats.Connect(natsURL)
	if err != nil {
		// fail open on nats
		t.Log(err)
	}

	js, err := nc.JetStream()
	if err != nil {
		// fail open on nats
		t.Log(err)
	}

	ps := pubsub.NewClient(
		pubsub.WithJetreamContext(js),
		pubsub.WithLogger(zap.NewNop().Sugar()),
		pubsub.WithStreamName("load-balancer-api"),
		pubsub.WithSubjectPrefix("com.infratographer.events"),
	)

	r := NewRouter(dbx, zap.NewNop().Sugar(), ps)

	r.Routes(e)

	return httptest.NewServer(e)
}

// doHTTPReqTest is a helper function to test the HTTP request
func doHTTPTest(t *testing.T, tt *httpTest) {
	t.Helper()

	var body io.Reader

	t.Run(tt.name+":["+tt.method+"] "+tt.path, func(t *testing.T) {
		t.Helper()
		if tt.body != "" {
			body = httptools.FakeBody(tt.body)
		} else {
			body = nil
		}

		req, err := http.NewRequest(tt.method, tt.path, body) //nolint

		req.Header.Set("Content-Type", "application/json")

		res, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, tt.status, res.StatusCode)
		res.Body.Close()
	})
}
