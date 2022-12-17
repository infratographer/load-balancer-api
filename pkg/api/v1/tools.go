package api

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.infratographer.com/loadbalancerapi/internal/config"
	"go.infratographer.com/loadbalancerapi/internal/httptools"
	"go.infratographer.com/loadbalancerapi/internal/x/echox"
	"go.infratographer.com/x/crdbx"
	"go.uber.org/zap"
)

type httpTest struct {
	name        string
	body        string
	tenant      string
	method      string
	displayName string
	path        string
	status      int
}

func newTestServer(t *testing.T) *httptest.Server {
	db, err := crdbx.NewDB(config.AppConfig.CRDB, false)
	assert.NoError(t, err)

	e := echox.NewServer()
	r := NewRouter(db, zap.NewNop().Sugar())

	r.Routes(e)

	return httptest.NewServer(e)
}

// doHTTPReqTest is a helper function to test the HTTP request
func doHTTPTests(t *testing.T, tests []httpTest) {
	t.Helper()

	var body io.Reader

	for _, tt := range tests {
		t.Run(tt.name+":["+tt.method+"] "+tt.path, func(t *testing.T) {
			if tt.body != "" {
				body = httptools.FakeBody(tt.body)
			} else {
				body = nil
			}

			req, err := http.NewRequest(tt.method, tt.path, body) //nolint

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set(tenantHeader, tt.tenant)

			res, err := http.DefaultClient.Do(req)
			assert.NoError(t, err)
			assert.Equal(t, tt.status, res.StatusCode)
			res.Body.Close()
		})
	}
}
