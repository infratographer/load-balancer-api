package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func createOrigin(t *testing.T, srv *httptest.Server, poolID, tenantID string) (*origin, func(*testing.T)) {
	t.Helper()

	body := fmt.Sprintf(`[{"disabled": true,"display_name": "The Butt", "target": "1.1.1.1", "port": 80, "pool_id": "%s"}]`, poolID)

	doHTTPTest(t, &httpTest{
		name:   "create origin",
		body:   body,
		status: http.StatusOK,
		path:   srv.URL + "/v1/origins",
		method: http.MethodPost,
		tenant: tenantID,
	})

	// Get the origin
	req, err := http.NewRequest(http.MethodGet, srv.URL+"/v1/origins?slug=the-butt", nil) //nolint
	assert.NoError(t, err)

	req.Header.Set(tenantHeader, tenantID)

	res, err := http.DefaultClient.Do(req) //nolint
	assert.NoError(t, err)

	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)

	origin := response{}
	err = json.NewDecoder(res.Body).Decode(&origin)
	assert.NoError(t, err)

	return (*origin.Origins)[0], func(t *testing.T) {
		t.Helper()

		// Delete the origin
		req, err := http.NewRequest(http.MethodDelete, srv.URL+"/v1/origins?slug=create-origin", nil) //nolint
		assert.NoError(t, err)

		req.Header.Set(tenantHeader, tenantID)

		res, err := http.DefaultClient.Do(req) //nolint
		assert.NoError(t, err)

		defer res.Body.Close()
	}
}

func TestOriginRoutes(t *testing.T) {
	srv := newTestServer(t)
	defer srv.Close()

	tenantID := uuid.New().String()
	baseURL := srv.URL + "/v1/origins"
	pool, remove := createPool(t, srv, tenantID)

	// doHTTPTest is a helper function that makes a request to the server and
	// checks the response.
	//
	// To ensure test output has meaningful line references the function is
	// called individually for each test case
	// POST
	doHTTPTest(t, &httpTest{
		name:   "happy path",
		body:   `[{"disabled": true,"display_name": "The Butt", "target": "9.9.9.9", "port": 80, "pool_id": "` + pool.ID + `"}]`,
		status: http.StatusOK,
		path:   baseURL,
		method: http.MethodPost,
		tenant: tenantID,
	})

	doHTTPTest(t, &httpTest{
		name:   "happy path",
		body:   `[{"disabled": true,"display_name": "Fish are friends", "target": "9.9.8.8", "port": 80, "pool_id": "` + pool.ID + `"}]`,
		status: http.StatusOK,
		path:   baseURL,
		method: http.MethodPost,
		tenant: tenantID,
	})

	doHTTPTest(t, &httpTest{
		name:   "no tenant",
		body:   `[]`,
		status: http.StatusBadRequest,
		path:   baseURL,
		method: http.MethodPost,
		tenant: "",
	})

	doHTTPTest(t, &httpTest{
		name:   "bad body",
		body:   `bad body`,
		status: http.StatusBadRequest,
		path:   baseURL,
		method: http.MethodPost,
		tenant: tenantID,
	})

	doHTTPTest(t, &httpTest{
		name:   "missing pool_id",
		body:   `[{"disabled": true,"display_name": "the-butt", "target": "2.0.0.1", "port": 80}]`,
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
		name:   "touch the butt",
		status: http.StatusOK,
		path:   baseURL + "?slug=the-butt",
		method: http.MethodGet,
		tenant: tenantID,
	})

	// DELETE
	doHTTPTest(t, &httpTest{
		name:   "ambigous delete",
		status: http.StatusBadRequest,
		path:   baseURL + "?port=80",
		method: http.MethodDelete,
		tenant: tenantID,
	})

	doHTTPTest(t, &httpTest{
		name:   "happy path",
		status: http.StatusOK,
		path:   baseURL + "?slug=the-butt",
		method: http.MethodDelete,
		tenant: tenantID,
	})

	doHTTPTest(t, &httpTest{
		name:   "happy path 2",
		status: http.StatusOK,
		path:   baseURL + "?slug=fish-are-friends",
		method: http.MethodDelete,
		tenant: tenantID,
	})

	doHTTPTest(t, &httpTest{
		name:   "404",
		status: http.StatusNotFound,
		path:   baseURL + "?slug=fish-are-friends",
		method: http.MethodDelete,
		tenant: tenantID,
	})

	remove(t)
}
