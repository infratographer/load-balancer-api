package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.infratographer.com/x/crdbx"
	"go.uber.org/zap"

	"go.infratographer.com/loadbalancerapi/internal/config"
	"go.infratographer.com/loadbalancerapi/internal/httptools"
	"go.infratographer.com/loadbalancerapi/internal/x/echox"
)

func newTestServer(t *testing.T) *httptest.Server {
	db, err := crdbx.NewDB(config.AppConfig.CRDB, false)
	assert.NoError(t, err)

	e := echox.NewServer()
	r := NewRouter(db, zap.NewNop().Sugar())

	r.Routes(e)

	return httptest.NewServer(e)
}

func TestLocationRoutes(t *testing.T) {
	srv := newTestServer(t)
	defer srv.Close()

	tenantID := uuid.New()
	baseURL := srv.URL + "/v1/tenant/" + tenantID.String() + "/locations"

	var payloadTests = []struct {
		name        string
		body        string
		status      int
		tenant      string
		displayName string
	}{
		{
			name:        "make a reef",
			body:        `{"display_name": "reef", "tenant_id": "` + tenantID.String() + `"}`,
			status:      http.StatusCreated,
			displayName: "reef",
		},
		{
			name:        "make a shell",
			body:        `{"display_name": "shell", "tenant_id": "` + tenantID.String() + `"}`,
			status:      http.StatusCreated,
			displayName: "shell",
		},
		{
			name:        "make a shell, again",
			body:        `{"display_name": "shell", "tenant_id": "` + tenantID.String() + `"}`,
			status:      http.StatusInternalServerError,
			displayName: "shell",
		},
		{
			name:   "missing display_name",
			body:   `{"tenant_id": "` + tenantID.String() + `"}`,
			status: 500,
		},

		{
			name:   "bad body",
			body:   "bad body",
			status: 400,
		},
	}

	for _, tt := range payloadTests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := srv.Client().Post(baseURL, "application/json", httptools.FakeBody(tt.body)) //nolint:noctx

			assert.NoError(t, err)
			assert.Equal(t, tt.status, res.StatusCode)
			res.Body.Close()
		})
	}

	var pathTests = []struct {
		name   string
		path   string
		status int
	}{
		{
			name:   "get reef",
			path:   baseURL + "/reef",
			status: http.StatusOK,
		},
		{
			name:   "get shell",
			path:   baseURL + "/shell",
			status: http.StatusOK,
		},
		{
			name:   "get all locations",
			path:   baseURL,
			status: http.StatusOK,
		},
		{
			name:   "get inavlid tenant",
			path:   srv.URL + "/v1/tenant/invalid/locations",
			status: http.StatusBadRequest,
		},
		{
			name:   "get invalid location",
			path:   baseURL + "/invalid",
			status: http.StatusNotFound,
		},
	}

	for _, tt := range pathTests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := srv.Client().Get(tt.path) //nolint:noctx
			assert.NoError(t, err)
			assert.Equal(t, tt.status, res.StatusCode)
			res.Body.Close()
		})
	}

	var deleteTests = []struct {
		name   string
		path   string
		status int
	}{
		{
			name:   "delete reef",
			path:   baseURL + "/reef",
			status: http.StatusOK,
		},
		{
			name:   "delete shell",
			path:   baseURL + "/shell",
			status: http.StatusOK,
		},
		// {
		// 	name:   "delete invalid location",
		// 	path:   baseURL + "/invalid",
		// 	status: http.StatusNotFound,
		// },
	}

	for _, tt := range deleteTests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("DELETE", tt.path, nil) //nolint:noctx
			assert.NoError(t, err)

			req.Header.Set("Content-Type", "application/json")
			res, err := http.DefaultClient.Do(req)
			assert.NoError(t, err)
			assert.Equal(t, tt.status, res.StatusCode)
			res.Body.Close()
		})
	}
}

func createAnenmoes(t *testing.T, srv *httptest.Server) (*response, func(t *testing.T)) {
	tenantID := uuid.New().String()
	baseURL := srv.URL + "/v1/tenant/" + tenantID + "/locations"
	happyPath := `{"display_name": "anemones", "tenant_id": "` + tenantID + `"}`

	t.Run("POST anemones", func(t *testing.T) {
		res, err := srv.Client().Post(baseURL, "application/json", httptools.FakeBody(happyPath)) //nolint:noctx
		assert.NoError(t, err)
		assert.Equal(t, 201, res.StatusCode)
		res.Body.Close()
	})

	loc := response{}

	t.Run("GET anemones", func(t *testing.T) {
		res, err := srv.Client().Get(baseURL + "/anemones") //nolint:noctx
		assert.NoError(t, err)
		assert.Equal(t, 200, res.StatusCode)

		err = json.NewDecoder(res.Body).Decode(&loc)
		assert.NoError(t, err)
		res.Body.Close()
	})

	return &loc, func(t *testing.T) {
		t.Run("DELETE anemones", func(t *testing.T) {
			req, err := http.NewRequest("DELETE", baseURL+"/anemones", nil) //nolint:noctx
			assert.NoError(t, err)

			req.Header.Set("Content-Type", "application/json")
			res, err := http.DefaultClient.Do(req)
			assert.NoError(t, err)
			assert.Equal(t, 200, res.StatusCode)
			res.Body.Close()
		})
	}
}
