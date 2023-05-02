//go:build testtools
// +build testtools

package api

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"go.infratographer.com/load-balancer-api/internal/config"
	"go.infratographer.com/load-balancer-api/internal/httptools"
	"go.infratographer.com/load-balancer-api/internal/pubsub"
	"go.infratographer.com/x/crdbx"
	"go.infratographer.com/x/echojwtx"
	"go.uber.org/zap"
)

type httpTest struct {
	client *http.Client
	name   string
	body   string
	tenant string
	method string
	path   string
	status int
}

func newTestServer(t *testing.T, natsURL string, authConfig *echojwtx.AuthConfig) (*httptest.Server, error) {
	var middleware []echo.MiddlewareFunc

	db, err := crdbx.NewDB(config.AppConfig.CRDB, false)
	if err != nil {
		t.Fatal(err)
	}

	dbx := sqlx.NewDb(db, "postgres")

	e := echo.New()

	if authConfig != nil {
		auth, err := echojwtx.NewAuth(context.Background(), *authConfig)
		if err != nil {
			return nil, err
		}

		middleware = append(middleware, auth.Middleware())
	}

	r := NewRouter(
		dbx,
		newPubSubClient(t, natsURL),
		WithLogger(zap.NewNop()),
		WithMiddleware(middleware...),
	)

	r.Routes(e.Group("/"))

	return httptest.NewServer(e), nil
}

func newPubSubClient(t *testing.T, url string) *pubsub.Client {
	nc, err := nats.Connect(url)
	if err != nil {
		// fail open on nats
		t.Error(err)
	}

	js, err := nc.JetStream()
	if err != nil {
		// fail open on nats
		t.Error(err)
	}

	return pubsub.NewClient(
		pubsub.WithJetreamContext(js),
		pubsub.WithLogger(zap.NewNop().Sugar()),
		pubsub.WithStreamName("load-balancer-api-test"),
		pubsub.WithSubjectPrefix("com.infratographer.events"),
	)
}

// doHTTPReqTest is a helper function to test the HTTP request
func doHTTPTest(t *testing.T, tt *httpTest) {
	t.Helper()

	var (
		body   io.Reader
		client = http.DefaultClient
	)

	if tt.client != nil {
		client = tt.client
	}

	t.Run(tt.name+":["+tt.method+"] "+tt.path, func(t *testing.T) {
		t.Helper()
		if tt.body != "" {
			body = httptools.FakeBody(tt.body)
		} else {
			body = nil
		}

		req, err := http.NewRequestWithContext(context.Background(), tt.method, tt.path, body) //nolint

		req.Header.Set("Content-Type", "application/json")

		res, err := client.Do(req)
		assert.NoError(t, err)

		defer res.Body.Close()

		resBody, err := io.ReadAll(res.Body)
		assert.NoError(t, err)

		t.Logf("response body %s", string(resBody))

		assert.Equal(t, tt.status, res.StatusCode)
	})
}
