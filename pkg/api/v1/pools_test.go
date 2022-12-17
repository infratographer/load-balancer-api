package api

import (
	"net/http"
	"testing"

	"github.com/google/uuid"
)

func TestPoolRoutes(t *testing.T) {
	srv := newTestServer(t)
	defer srv.Close()

	tenantID := uuid.New().String()
	baseURL := srv.URL + "/v1/pools"

	doHTTPTests(t, []httpTest{
		// POST
		{
			name:   "happy path",
			body:   `[{"display_name": "Nemo", "protocol": "tcp"}]`,
			status: http.StatusCreated,
			path:   baseURL,
			method: http.MethodPost,
			tenant: tenantID,
		},
		{
			name:   "duplicate",
			body:   `[{"display_name": "Nemo", "protocol": "tcp"}]`,
			status: http.StatusCreated,
			path:   baseURL,
			method: http.MethodPost,
			tenant: tenantID,
		},
		{
			name:   "missing display name",
			body:   `[{"protocol": "tcp"}]`,
			status: http.StatusBadRequest,
			path:   baseURL,
			method: http.MethodPost,
			tenant: tenantID,
		},
		{
			name:   "missing protocol",
			body:   `[{"display_name": "Nemo"}]`,
			status: http.StatusCreated,
			path:   baseURL,
			method: http.MethodPost,
			tenant: tenantID,
		},
		{
			name:   "invalid protocol",
			body:   `[{"display_name": "Nemo", "protocol": "invalid"}]`,
			status: http.StatusBadRequest,
			path:   baseURL,
			method: http.MethodPost,
			tenant: tenantID,
		},
		{
			name:   "invalid body",
			body:   `[{"display_name": "Nemo", "protocol": "tcp"`,
			status: http.StatusBadRequest,
			path:   baseURL,
			method: http.MethodPost,
			tenant: tenantID,
		},
		// GET
		{
			name:   "happy path",
			status: http.StatusOK,
			path:   baseURL,
			method: http.MethodGet,
			tenant: tenantID,
		},
		{
			name:   "happy path with query",
			status: http.StatusOK,
			path:   baseURL + "?display_name=Nemo",
			method: http.MethodGet,
			tenant: tenantID,
		},
		{
			name:   "not found with query",
			status: http.StatusNotFound,
			path:   baseURL + "?display_name=NotNemo",
			method: http.MethodGet,
			tenant: tenantID,
		},
		// DELETE
		{
			name:   "happy path",
			status: http.StatusNoContent,
			path:   baseURL + "?display_name=Nemo",
			method: http.MethodDelete,
			tenant: tenantID,
		},
	})
}
