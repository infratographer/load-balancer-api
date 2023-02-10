package api

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"github.com/stretchr/testify/assert"
)

// createPool creates a pool with the given display name and protocol.
func createPool(t *testing.T, srv *httptest.Server, name string, tenantID string) (*pool, func(t *testing.T)) {
	t.Helper()

	body := `[{"display_name": "` + name + `", "protocol": "tcp"}]`

	baseURL := srv.URL + "/v1/pools"

	doHTTPTest(t, &httpTest{
		name:   "create pool",
		method: http.MethodPost,
		path:   baseURL,
		body:   body,
		status: http.StatusOK,
		tenant: tenantID,
	})

	pool := response{}

	req, err := http.NewRequest(http.MethodGet, baseURL+"?slug="+slug.Make(name), nil) //nolint
	assert.NoError(t, err)

	req.Header.Set(tenantHeader, tenantID)

	res, err := http.DefaultClient.Do(req) //nolint
	assert.Equal(t, http.StatusOK, res.StatusCode)

	bytes, err := io.ReadAll(res.Body)
	assert.NoError(t, err)

	err = json.Unmarshal(bytes, &pool)
	assert.NoError(t, err)

	res.Body.Close()

	return (*pool.Pools)[0], func(t *testing.T) {
		t.Helper()

		req, err := http.NewRequest(http.MethodDelete, baseURL+"/"+(*pool.Pools)[0].ID, nil) //nolint
		assert.NoError(t, err)

		req.Header.Set(tenantHeader, tenantID)

		res, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, res.StatusCode)
		res.Body.Close()
	}
}

func TestPoolRoutes(t *testing.T) {
	srv := newTestServer(t)
	defer srv.Close()

	tenantID := uuid.New().String()
	baseURL := srv.URL + "/v1/pools"
	missingUUID := uuid.New().String()

	// doHTTPTest is a helper function that makes a request to the server and
	// checks the response.
	//
	// To ensure test output has meaningful line references the function is
	// called individually for each test case
	doHTTPTest(t, &httpTest{
		name:   "get pools before create",
		status: http.StatusOK,
		path:   baseURL,
		method: http.MethodGet,
		tenant: tenantID,
	})

	// POST
	doHTTPTest(t, &httpTest{
		name:   "happy path",
		body:   `[{"display_name": "Nemo", "protocol": "tcp"}]`,
		status: http.StatusOK,
		path:   baseURL,
		method: http.MethodPost,
		tenant: tenantID,
	})

	doHTTPTest(t, &httpTest{
		name:   "duplicate",
		body:   `[{"display_name": "Nemo", "protocol": "tcp"}]`,
		status: http.StatusInternalServerError,
		path:   baseURL,
		method: http.MethodPost,
		tenant: tenantID,
	})

	doHTTPTest(t, &httpTest{
		name:   "missing display name",
		body:   `[{"protocol": "tcp"}]`,
		status: http.StatusBadRequest,
		path:   baseURL,
		method: http.MethodPost,
		tenant: tenantID,
	})

	doHTTPTest(t, &httpTest{
		name:   "missing protocol",
		body:   `[{"display_name": "Bruce"}]`,
		status: http.StatusOK,
		path:   baseURL,
		method: http.MethodPost,
		tenant: tenantID,
	})

	doHTTPTest(t, &httpTest{
		name:   "invalid protocol",
		body:   `[{"display_name": "Nemo", "protocol": "invalid"}]`,
		status: http.StatusBadRequest,
		path:   baseURL,
		method: http.MethodPost,
		tenant: tenantID,
	})

	doHTTPTest(t, &httpTest{
		name:   "invalid body",
		body:   `invalid`,
		status: http.StatusBadRequest,
		path:   baseURL,
		method: http.MethodPost,
		tenant: tenantID,
	})

	// GET
	doHTTPTest(t, &httpTest{
		name:   "happy path",
		status: http.StatusOK,
		path:   baseURL,
		method: http.MethodGet,
		tenant: tenantID,
	})

	doHTTPTest(t, &httpTest{
		name:   "happy path with query",
		status: http.StatusOK,
		path:   baseURL + "?display_name=Nemo",
		method: http.MethodGet,
		tenant: tenantID,
	})

	doHTTPTest(t, &httpTest{
		name:   "not found with query",
		status: http.StatusOK,
		path:   baseURL + "?slug=NotNemo",
		method: http.MethodGet,
		tenant: tenantID,
	})

	doHTTPTest(t, &httpTest{
		name:   "list pools with invalid tenant id",
		path:   baseURL,
		status: http.StatusBadRequest,
		method: http.MethodGet,
		tenant: "123456",
	})

	doHTTPTest(t, &httpTest{
		name:   "list pools with unknown tenant id",
		path:   baseURL,
		status: http.StatusOK,
		method: http.MethodGet,
		tenant: missingUUID,
	})

	doHTTPTest(t, &httpTest{
		name:   "not found with path param",
		status: http.StatusNotFound,
		path:   baseURL + "/" + missingUUID,
		method: http.MethodGet,
		tenant: tenantID,
	})

	doHTTPTest(t, &httpTest{
		name:   "bad UUID in path param",
		status: http.StatusBadRequest,
		path:   baseURL + "/123456",
		method: http.MethodGet,
		tenant: tenantID,
	})

	// DELETE
	doHTTPTest(t, &httpTest{
		name:   "happy path",
		status: http.StatusOK,
		path:   baseURL + "?display_name=Nemo",
		method: http.MethodDelete,
		tenant: tenantID,
	})
}

func TestPoolsGet(t *testing.T) {
	srv := newTestServer(t)
	defer srv.Close()

	assert.NotNil(t, srv)

	baseURL := srv.URL + "/v1/pools"

	// Create a load balancer
	loadBalancer, cleanupLB := createLoadBalancer(t, srv, uuid.NewString())
	defer cleanupLB(t)

	// Create a pool
	pool, cleanupPool := createPool(t, srv, "marlin", loadBalancer.TenantID)
	defer cleanupPool(t)

	// Get the pool
	doHTTPTest(t, &httpTest{
		name:   "get pool by id",
		method: http.MethodGet,
		path:   baseURL + "/" + pool.ID,
		status: http.StatusOK,
		tenant: loadBalancer.TenantID,
	})

	// Get an unknown pool
	doHTTPTest(t, &httpTest{
		name:   "pool not found",
		method: http.MethodGet,
		path:   baseURL + "/bfad65a9-abe3-44af-82ce-64331c84b2ad",
		status: http.StatusNotFound,
		tenant: loadBalancer.TenantID,
	})
}
