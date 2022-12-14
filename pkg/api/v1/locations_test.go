package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"go.infratographer.com/loadbalancerapi/internal/httptools"
)

func TestLocationRoutes(t *testing.T) {
	srv := newTestServer(t)
	defer srv.Close()

	tenantID := uuid.New()
	baseURL := srv.URL + "/v1/locations"

	doHTTPTests(t, []httpTest{
		// POST tests
		{
			name:        "make a reef",
			body:        `{"display_name": "reef"}`,
			path:        baseURL,
			status:      http.StatusCreated,
			displayName: "reef",
			tenant:      tenantID.String(),
			method:      "POST",
		},
		{
			name:        "make a shell",
			body:        `{"display_name": "shell"}`,
			path:        baseURL,
			status:      http.StatusCreated,
			displayName: "shell",
			tenant:      tenantID.String(),
			method:      "POST",
		},
		{
			name:        "make a shell, again",
			body:        `{"display_name": "shell"}`,
			path:        baseURL,
			status:      http.StatusInternalServerError,
			displayName: "shell",
			tenant:      tenantID.String(),
			method:      "POST",
		},
		{
			name:   "bad body",
			body:   "bad body",
			path:   baseURL,
			status: http.StatusInternalServerError,
			tenant: tenantID.String(),
			method: "POST",
		},
		// Get tests
		{
			name:   "get reef",
			path:   baseURL + "/reef",
			status: http.StatusOK,
			tenant: tenantID.String(),
			method: "GET",
		},
		{
			name:   "get shell",
			path:   baseURL + "/shell",
			status: http.StatusOK,
			tenant: tenantID.String(),
			method: "GET",
		},
		{
			name:   "get all locations",
			path:   baseURL,
			status: http.StatusOK,
			tenant: tenantID.String(),
			method: "GET",
		},
		{
			name:   "get inavlid tenant",
			path:   baseURL,
			status: http.StatusBadRequest,
			tenant: "invalid",
			method: "GET",
		},
		{
			name:   "get invalid location",
			path:   baseURL + "/invalid",
			status: http.StatusNotFound,
			tenant: tenantID.String(),
			method: "GET",
		},
		// Delete tests
		{
			name:   "delete reef",
			path:   baseURL + "/reef",
			status: http.StatusNoContent,
			tenant: tenantID.String(),
			method: "DELETE",
		},
		{
			name:   "delete shell",
			path:   baseURL + "/shell",
			status: http.StatusNoContent,
			tenant: tenantID.String(),
			method: "DELETE",
		},
	})
}

func createAnemones(t *testing.T, srv *httptest.Server) (*response, func(t *testing.T)) {
	tenantID := uuid.New().String()
	baseURL := srv.URL + "/v1/locations"
	happyPath := `{"display_name": "anemones"}`

	t.Run("POST anemones", func(t *testing.T) {
		req, err := http.NewRequest("POST", baseURL, httptools.FakeBody(happyPath)) //nolint
		assert.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Infratographer-Tenant-ID", tenantID)

		res, err := srv.Client().Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, res.StatusCode)
		res.Body.Close()
	})

	loc := response{}

	t.Run("GET anemones", func(t *testing.T) {
		req, err := http.NewRequest("GET", baseURL, nil) //nolint
		assert.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Infratographer-Tenant-ID", tenantID)

		res, err := srv.Client().Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode)

		err = json.NewDecoder(res.Body).Decode(&loc)
		assert.NoError(t, err)
		res.Body.Close()
	})

	return &loc, func(t *testing.T) {
		t.Run("DELETE anemones", func(t *testing.T) {
			req, err := http.NewRequest("DELETE", baseURL+"/anemones", nil) //nolint:noctx
			assert.NoError(t, err)

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("X-Infratographer-Tenant-ID", tenantID)

			res, err := http.DefaultClient.Do(req)
			assert.NoError(t, err)
			assert.Equal(t, http.StatusNoContent, res.StatusCode)
			res.Body.Close()
		})
	}
}
