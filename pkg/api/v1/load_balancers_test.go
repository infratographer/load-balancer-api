package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.infratographer.com/load-balancer-api/internal/httptools"
)

func TestLoadBalancerRoutes(t *testing.T) {
	srv := newTestServer(t)
	defer srv.Close()

	tenantID := uuid.New().String()
	baseURL := srv.URL + "/v1/loadbalancers"
	baseURLTenant := srv.URL + "/v1/tenant/" + tenantID + "/loadbalancers"
	locationID := uuid.New().String()
	missingUUID := uuid.New().String()

	// create a test load balancer named Bruce
	req1, err := http.NewRequestWithContext(
		context.TODO(),
		http.MethodPost,
		baseURLTenant,
		httptools.FakeBody(
			fmt.Sprintf(`{"display_name": "Bruce", "location_id": "%s", "ip_addr": "2.2.2.2","load_balancer_size": "small","load_balancer_type": "layer-3"}`, locationID),
		),
	)
	assert.NoError(t, err)
	resp1, err := http.DefaultClient.Do(req1)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp1.StatusCode)

	testLoadBalancer := struct {
		Version        string `json:"version"`
		Message        string `json:"message"`
		Status         int    `json:"status"`
		LoadBalancerID string `json:"load_balancer_id"`
	}{}

	_ = json.NewDecoder(resp1.Body).Decode(&testLoadBalancer)
	resp1.Body.Close()

	// cleanup test load balancer
	defer func(id string) {
		rq, err := http.NewRequestWithContext(context.TODO(), http.MethodDelete, baseURL+"/"+id, nil)
		assert.NoError(t, err)
		rs, err := http.DefaultClient.Do(rq)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rs.StatusCode)
		rs.Body.Close()
	}(testLoadBalancer.LoadBalancerID)

	// doHTTPTest is a helper function that makes a request to the server and
	// checks the response.
	//
	// To ensure test output has meaningful line references the function is
	// called individually for each test case
	doHTTPTest(t, &httpTest{
		name:   "list lbs before created",
		path:   baseURLTenant,
		status: http.StatusOK,
		method: http.MethodGet,
	})

	// POST tests
	doHTTPTest(t, &httpTest{
		name:   "happy path",
		body:   fmt.Sprintf(`{"display_name": "Nemo", "location_id": "%s", "ip_addr": "1.1.1.1","load_balancer_size": "small","load_balancer_type": "layer-3"}`, locationID),
		status: http.StatusOK,
		path:   baseURLTenant,
		method: http.MethodPost,
	})

	doHTTPTest(t, &httpTest{
		name:   "happy path 2",
		body:   fmt.Sprintf(`{"display_name": "Dori", "location_id": "%s", "ip_addr": "1.2.1.1","load_balancer_size": "small","load_balancer_type": "layer-3"}`, locationID),
		status: http.StatusOK,
		path:   baseURLTenant,
		method: http.MethodPost,
	})

	doHTTPTest(t, &httpTest{
		name:   "Duplicate",
		body:   fmt.Sprintf(`{"display_name": "Nemo", "location_id": "%s", "ip_addr": "1.1.1.1","load_balancer_size": "small","load_balancer_type": "layer-3"}`, locationID),
		status: http.StatusInternalServerError,
		path:   baseURLTenant,
		method: http.MethodPost,
	})

	doHTTPTest(t, &httpTest{
		name:   "missing tenantID",
		body:   fmt.Sprintf(`{"display_name": "Nemo", "location_id": "%s", "ip_addr": "1.1.1.1","load_balancer_size": "small","load_balancer_type": "layer-3"}`, locationID),
		status: http.StatusNotFound,
		path:   baseURL,
		method: http.MethodPost,
	})

	doHTTPTest(t, &httpTest{
		name:   "missing display name",
		body:   fmt.Sprintf(`{"location_id": "%s", "ip_addr": "1.1.1.1","load_balancer_size": "small","load_balancer_type": "layer-3"}`, locationID),
		status: http.StatusBadRequest,
		path:   baseURLTenant,
		method: http.MethodPost,
	})

	doHTTPTest(t, &httpTest{
		name:   "missing location id",
		body:   `{"display_name": "Nemo", "ip_addr": "1.1.1.1","load_balancer_size": "small","load_balancer_type": "layer-3"}`,
		status: http.StatusBadRequest,
		path:   baseURLTenant,
		method: http.MethodPost,
	})

	doHTTPTest(t, &httpTest{
		name:   "missing ip address",
		body:   fmt.Sprintf(`{"display_name": "Nemo", "location_id": "%s", "load_balancer_size": "small","load_balancer_type": "layer-3"}`, locationID),
		status: http.StatusBadRequest,
		path:   baseURLTenant,
		method: http.MethodPost,
	})

	doHTTPTest(t, &httpTest{
		name:   "missing size",
		body:   fmt.Sprintf(`{"display_name": "Nemo", "location_id": "%s", "ip_addr": "1.1.1.1","load_balancer_type": "layer-3"}`, locationID),
		status: http.StatusBadRequest,
		path:   baseURLTenant,
		method: http.MethodPost,
	})

	doHTTPTest(t, &httpTest{
		name:   "missing type",
		body:   fmt.Sprintf(`{"display_name": "Nemo", "location_id": "%s", "ip_addr": "1.1.1.1","load_balancer_size": "small"}`, locationID),
		status: http.StatusBadRequest,
		path:   baseURLTenant,
		method: http.MethodPost,
	})

	doHTTPTest(t, &httpTest{
		name:   "invalid type",
		body:   fmt.Sprintf(`{"display_name": "Nemo", "location_id": "%s", "ip_addr": "1.1.1.1","load_balancer_size": "small","load_balancer_type": "layer-12"}`, locationID),
		status: http.StatusBadRequest,
		path:   baseURLTenant,
		method: http.MethodPost,
	})

	doHTTPTest(t, &httpTest{
		name:   "bad ip address",
		body:   fmt.Sprintf(`{"display_name": "Nemo", "location_id": "%s", "ip_addr": "Dori","load_balancer_size": "small","load_balancer_type": "layer-3"}`, locationID),
		status: http.StatusBadRequest,
		path:   baseURLTenant,
		method: http.MethodPost,
	})

	doHTTPTest(t, &httpTest{
		name:   "ipv6",
		body:   fmt.Sprintf(`{"display_name": "Nemo", "location_id": "%s", "ip_addr": "2601::","load_balancer_size": "small","load_balancer_type": "layer-3"}`, locationID),
		status: http.StatusInternalServerError,
		path:   baseURLTenant,
		method: http.MethodPost,
	})

	doHTTPTest(t, &httpTest{
		name:   "empty body",
		body:   ``,
		status: http.StatusBadRequest,
		path:   baseURLTenant,
		method: http.MethodPost,
	})

	doHTTPTest(t, &httpTest{
		name:   "bad body",
		body:   `bad body`,
		status: http.StatusBadRequest,
		path:   baseURLTenant,
		method: http.MethodPost,
	})

	// PUT tests
	doHTTPTest(t, &httpTest{
		name:   "happy path update load balancer",
		body:   `{"display_name": "Bruce", "load_balancer_size": "x-large","load_balancer_type": "layer-3"}`,
		status: http.StatusAccepted,
		method: http.MethodPut,
		path:   baseURL + "/" + testLoadBalancer.LoadBalancerID,
		tenant: tenantID,
	})

	doHTTPTest(t, &httpTest{
		name:   "update load balancer missing display name",
		body:   `{"load_balancer_size": "x-large","load_balancer_type": "layer-3"}`,
		status: http.StatusBadRequest,
		method: http.MethodPut,
		path:   baseURL + "/" + testLoadBalancer.LoadBalancerID,
		tenant: tenantID,
	})

	doHTTPTest(t, &httpTest{
		name:   "update load balancer, missing size",
		body:   `{"display_name": "Bruce","load_balancer_type": "layer-3"}`,
		status: http.StatusBadRequest,
		method: http.MethodPut,
		path:   baseURL + "/" + testLoadBalancer.LoadBalancerID,
		tenant: tenantID,
	})

	doHTTPTest(t, &httpTest{
		name:   "update load balancer, missing type",
		body:   `{"display_name": "Bruce", "load_balancer_size": "x-large"}`,
		status: http.StatusBadRequest,
		method: http.MethodPut,
		path:   baseURL + "/" + testLoadBalancer.LoadBalancerID,
		tenant: tenantID,
	})

	doHTTPTest(t, &httpTest{
		name:   "update load balancer, missing load balancer id",
		body:   `{"display_name": "Bruce", "load_balancer_size": "x-large","load_balancer_type": "layer-3"}`,
		status: http.StatusNotFound,
		method: http.MethodPut,
		path:   baseURL,
		tenant: tenantID,
	})

	// GET tests
	doHTTPTest(t, &httpTest{
		name:   "happy path",
		path:   baseURLTenant,
		status: http.StatusOK,
		method: http.MethodGet,
	})

	doHTTPTest(t, &httpTest{
		name:   "happy path nemo by name",
		path:   baseURLTenant + "?display_name=Nemo",
		status: http.StatusOK,
		method: http.MethodGet,
	})

	doHTTPTest(t, &httpTest{
		name:   "happy path nemo by ip",
		path:   baseURLTenant + "?ip_addr=1.1.1.1",
		status: http.StatusOK,
		method: http.MethodGet,
	})

	doHTTPTest(t, &httpTest{
		name:   "sad path lb doesnt exist",
		path:   baseURL + "/" + missingUUID,
		status: http.StatusNotFound,
		method: http.MethodGet,
	})

	doHTTPTest(t, &httpTest{
		name:   "sad path bad uuid",
		path:   baseURL + "/123456",
		status: http.StatusBadRequest,
		method: http.MethodGet,
	})

	doHTTPTest(t, &httpTest{
		name:   "list lbs with invalid tenant id",
		path:   srv.URL + "/v1/tenant/123456/loadbalancers",
		status: http.StatusBadRequest,
		method: http.MethodGet,
	})

	doHTTPTest(t, &httpTest{
		name:   "list lbs with unknown tenant id",
		path:   srv.URL + "/v1/tenant/" + missingUUID + "/loadbalancers",
		status: http.StatusOK,
		method: http.MethodGet,
	})

	// DELETE tests
	doHTTPTest(t, &httpTest{
		name:   "delete invalid id",
		path:   baseURL + "/invalid",
		status: http.StatusBadRequest,
		method: http.MethodDelete,
	})

	doHTTPTest(t, &httpTest{
		name:   "delete lb that doesnt exist",
		path:   baseURL + "/ce94616e-3798-454d-91f3-9e3cec32bff6",
		status: http.StatusNotFound,
		method: http.MethodDelete,
	})

	doHTTPTest(t, &httpTest{
		name:   "delete lb without id",
		path:   baseURL,
		status: http.StatusNotFound,
		method: http.MethodDelete,
	})

	doHTTPTest(t, &httpTest{
		name:   "delete small load balancers",
		path:   baseURLTenant + "?load_balancer_size=small",
		status: http.StatusUnprocessableEntity,
		method: http.MethodDelete,
	})

	doHTTPTest(t, &httpTest{
		name:   "delete nemo by Name",
		path:   baseURLTenant + "?slug=nemo",
		status: http.StatusOK,
		method: http.MethodDelete,
	})

	doHTTPTest(t, &httpTest{
		name:   "delete Dori by name",
		path:   baseURLTenant + "?slug=dori",
		status: http.StatusOK,
		method: http.MethodDelete,
	})

	doHTTPTest(t, &httpTest{
		name:   "delete Dori again",
		path:   baseURLTenant + "?slug=dori",
		status: http.StatusNotFound,
		method: http.MethodDelete,
	})
}

func TestLoadBalancerGet(t *testing.T) {
	srv := newTestServer(t)
	defer srv.Close()

	assert.NotNil(t, srv)

	tenantID := uuid.New().String()
	baseURL := srv.URL + "/v1/loadbalancers"
	missingUUID := uuid.New().String()

	// Create a load balancer
	loadBalancer, cleanupLB := createLoadBalancer(t, srv, tenantID)
	defer cleanupLB(t)

	// Get the load balancer
	doHTTPTest(t, &httpTest{
		name:   "get loadblancer by id",
		method: http.MethodGet,
		path:   baseURL + "/" + loadBalancer.ID,
		status: http.StatusOK,
	})

	// Get an unknown load balancer
	doHTTPTest(t, &httpTest{
		name:   "get missing loadblancer by id",
		method: http.MethodGet,
		path:   baseURL + "/" + missingUUID,
		status: http.StatusNotFound,
	})

	// Get an unknown tenant
	doHTTPTest(t, &httpTest{
		name:   "get missing loadblancer by id",
		method: http.MethodGet,
		path:   srv.URL + "/v1/tenant/" + missingUUID + "/loadbalancers/" + loadBalancer.ID,
		status: http.StatusNotFound,
	})

	// Get the load balancer without id
	doHTTPTest(t, &httpTest{
		name:   "get loadblancer without id",
		method: http.MethodGet,
		path:   baseURL,
		status: http.StatusNotFound,
	})
}

func createLoadBalancer(t *testing.T, srv *httptest.Server, locationID string) (*loadBalancer, func(t *testing.T)) {
	tenantID := uuid.New().String()
	baseURL := srv.URL + "/v1/tenant/" + tenantID + "/loadbalancers"

	test := &httpTest{
		name:   "create nemo lb",
		body:   fmt.Sprintf(`{"display_name": "Nemo", "location_id": "%s", "ip_addr": "1.1.1.1","load_balancer_size": "small","load_balancer_type": "layer-3"}`, locationID),
		path:   baseURL,
		method: http.MethodPost,
		status: http.StatusOK,
	}

	doHTTPTest(t, test)

	// get loadbalancer by name
	loadbalancer := response{}

	t.Run("get nemo by name:[POST] "+baseURL+"?display_name=Nemo", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, baseURL+"?display_name=Nemo", nil) //nolint
		assert.NoError(t, err)

		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		err = json.NewDecoder(resp.Body).Decode(&loadbalancer)
		assert.NoError(t, err)
		resp.Body.Close()
	})

	return (*loadbalancer.LoadBalancers)[0], func(t *testing.T) {
		test := &httpTest{
			name:   "delete nemo",
			path:   baseURL + "?slug=nemo",
			method: http.MethodDelete,
			status: http.StatusOK,
		}

		doHTTPTest(t, test)
	}
}
