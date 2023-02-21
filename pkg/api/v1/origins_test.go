package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"github.com/stretchr/testify/assert"
)

func createOrigin(t *testing.T, srv *httptest.Server, name string, poolID string) (*origin, func(*testing.T)) {
	t.Helper()

	body := fmt.Sprintf(`{"disabled": true,"display_name": "%s", "target": "1.1.1.1", "port": 80}`, name)

	doHTTPTest(t, &httpTest{
		name:   "create origin",
		body:   body,
		status: http.StatusOK,
		path:   srv.URL + "/v1/pools/" + poolID + "/origins",
		method: http.MethodPost,
	})

	// Get the origin
	req, err := http.NewRequest(http.MethodGet, srv.URL+"/v1/pools/"+poolID+"/origins?slug="+slug.Make(name), nil) //nolint
	assert.NoError(t, err)

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
		req, err := http.NewRequest(http.MethodDelete, srv.URL+"/v1/pools/"+poolID+"/origins?slug="+slug.Make(name), nil) //nolint
		assert.NoError(t, err)

		res, err := http.DefaultClient.Do(req) //nolint
		assert.NoError(t, err)

		defer res.Body.Close()
	}
}

func TestOriginRoutes(t *testing.T) {
	srv := newTestServer(t)
	defer srv.Close()

	tenantID := uuid.New().String()
	pool, remove := createPool(t, srv, "squirt", tenantID)

	baseURL := srv.URL + "/v1/origins"
	baseURLPool := srv.URL + "/v1/pools/" + pool.ID + "/origins"
	missingUUID := uuid.New().String()

	// doHTTPTest is a helper function that makes a request to the server and
	// checks the response.
	//
	// To ensure test output has meaningful line references the function is
	// called individually for each test case
	doHTTPTest(t, &httpTest{
		name:   "list origins before created",
		status: http.StatusOK,
		path:   baseURLPool,
		method: http.MethodGet,
	})

	// POST
	doHTTPTest(t, &httpTest{
		name:   "happy path",
		body:   `{"disabled": true,"display_name": "The Butt", "target": "9.9.9.9", "port": 80}`,
		status: http.StatusOK,
		path:   baseURLPool,
		method: http.MethodPost,
	})

	doHTTPTest(t, &httpTest{
		name:   "happy path",
		body:   `{"disabled": true,"display_name": "Fish are friends", "target": "9.9.8.8", "port": 80}`,
		status: http.StatusOK,
		path:   baseURLPool,
		method: http.MethodPost,
	})

	doHTTPTest(t, &httpTest{
		name:   "list of origins",
		body:   `[{"disabled": true,"display_name": "The Butt", "target": "9.9.9.9", "port": 80},{"disabled": true,"display_name": "The Beard", "target": "9.9.9.10", "port": 80}]`,
		status: http.StatusBadRequest,
		path:   baseURLPool,
		method: http.MethodPost,
	})

	doHTTPTest(t, &httpTest{
		name:   "no pool",
		body:   `[]`,
		status: http.StatusBadRequest,
		path:   baseURLPool,
		method: http.MethodPost,
	})

	doHTTPTest(t, &httpTest{
		name:   "bad body",
		body:   `bad body`,
		status: http.StatusBadRequest,
		path:   baseURLPool,
		method: http.MethodPost,
	})

	doHTTPTest(t, &httpTest{
		name:   "missing pool_id in path",
		body:   `{"disabled": true,"display_name": "the-butt", "target": "2.0.0.1", "port": 80}`,
		status: http.StatusNotFound,
		path:   baseURL,
		method: http.MethodPost,
	})

	// GET
	doHTTPTest(t, &httpTest{
		name:   "happy path",
		status: http.StatusOK,
		path:   baseURLPool,
		method: http.MethodGet,
	})

	doHTTPTest(t, &httpTest{
		name:   "touch the butt",
		status: http.StatusOK,
		path:   baseURLPool + "?slug=the-butt",
		method: http.MethodGet,
	})

	doHTTPTest(t, &httpTest{
		name:   "missing origin uuid",
		path:   baseURL + "/" + missingUUID,
		status: http.StatusNotFound,
		method: http.MethodGet,
	})

	doHTTPTest(t, &httpTest{
		name:   "bad origin uuid",
		path:   baseURL + "/123456",
		status: http.StatusBadRequest,
		method: http.MethodGet,
	})

	doHTTPTest(t, &httpTest{
		name:   "list origins with invalid pool id",
		path:   srv.URL + "/v1/pools/123456/origins",
		status: http.StatusBadRequest,
		method: http.MethodGet,
	})

	doHTTPTest(t, &httpTest{
		name:   "list origins with unknown pool id",
		path:   srv.URL + "/v1/pools/" + uuid.New().String() + "/origins",
		status: http.StatusOK,
		method: http.MethodGet,
	})

	// DELETE
	doHTTPTest(t, &httpTest{
		name:   "ambigous delete",
		status: http.StatusBadRequest,
		path:   baseURLPool + "?port=80",
		method: http.MethodDelete,
	})

	doHTTPTest(t, &httpTest{
		name:   "happy path",
		status: http.StatusOK,
		path:   baseURLPool + "?slug=the-butt",
		method: http.MethodDelete,
	})

	doHTTPTest(t, &httpTest{
		name:   "happy path 2",
		status: http.StatusOK,
		path:   baseURLPool + "?slug=fish-are-friends",
		method: http.MethodDelete,
	})

	doHTTPTest(t, &httpTest{
		name:   "404",
		status: http.StatusNotFound,
		path:   baseURL + "?slug=fish-are-friends",
		method: http.MethodDelete,
	})

	remove(t)
}

func TestOriginsBalancerGet(t *testing.T) {
	srv := newTestServer(t)
	defer srv.Close()

	assert.NotNil(t, srv)

	baseURL := srv.URL + "/v1/origins"

	// Create a load balancer
	loadBalancer, cleanupLB := createLoadBalancer(t, srv, uuid.NewString())
	defer cleanupLB(t)

	// Create a pool
	pool, cleanupPool := createPool(t, srv, "marlin", loadBalancer.TenantID)
	defer cleanupPool(t)

	// Create an origin in the pool
	origin, cleanupOrigin := createOrigin(t, srv, "bruce", pool.ID)
	defer cleanupOrigin(t)

	// Get the origin
	doHTTPTest(t, &httpTest{
		name:   "get origin by id",
		method: http.MethodGet,
		path:   baseURL + "/" + origin.ID,
		status: http.StatusOK,
	})

	// Get an unknown origin
	doHTTPTest(t, &httpTest{
		name:   "get missing origin by id",
		method: http.MethodGet,
		path:   baseURL + "/bfad65a9-abe3-44af-82ce-64331c84b2ad",
		status: http.StatusNotFound,
	})
}
