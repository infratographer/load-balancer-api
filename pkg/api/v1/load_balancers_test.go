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

func TestLoadBalancerRoutes(t *testing.T) {
	srv := newTestServer(t)
	defer srv.Close()

	tenantID := uuid.New().String()
	baseURL := srv.URL + "/v1/loadbalancers"
	locationID := uuid.New().String()

	// doHTTPTest is a helper function that makes a request to the server and
	// checks the response.
	//
	// To ensure test output has meaningful line references the function is
	// called individually for each test case
	// POST
	doHTTPTest(t, &httpTest{
		name:   "happy path",
		body:   fmt.Sprintf(`[{"display_name": "Nemo", "location_id": "%s", "ip_addr": "1.1.1.1","load_balancer_size": "small","load_balancer_type": "layer-3"}]`, locationID),
		status: http.StatusOK,
		path:   baseURL,
		method: http.MethodPost,
		tenant: tenantID,
	})

	doHTTPTest(t, &httpTest{
		name:   "happy path 2",
		body:   fmt.Sprintf(`[{"display_name": "Dori", "location_id": "%s", "ip_addr": "1.2.1.1","load_balancer_size": "small","load_balancer_type": "layer-3"}]`, locationID),
		status: http.StatusOK,
		path:   baseURL,
		method: http.MethodPost,
		tenant: tenantID,
	})

	doHTTPTest(t, &httpTest{
		name:   "Duplicate",
		body:   fmt.Sprintf(`[{"display_name": "Nemo", "location_id": "%s", "ip_addr": "1.1.1.1","load_balancer_size": "small","load_balancer_type": "layer-3"}]`, locationID),
		status: http.StatusInternalServerError,
		path:   baseURL,
		method: http.MethodPost,
		tenant: tenantID,
	})

	doHTTPTest(t, &httpTest{
		name:   "missing display name",
		body:   fmt.Sprintf(`[{"location_id": "%s", "ip_addr": "1.1.1.1","load_balancer_size": "small","load_balancer_type": "layer-3"}]`, locationID),
		status: http.StatusUnprocessableEntity,
		path:   baseURL,
		method: http.MethodPost,
		tenant: tenantID,
	})

	doHTTPTest(t, &httpTest{
		name:   "missing location id",
		body:   `[{"display_name": "Nemo", "ip_addr": "1.1.1.1","load_balancer_size": "small","load_balancer_type": "layer-3"}]`,
		status: http.StatusUnprocessableEntity,
		path:   baseURL,
		method: http.MethodPost,
		tenant: tenantID,
	})

	doHTTPTest(t, &httpTest{
		name:   "missing ip address",
		body:   fmt.Sprintf(`[{"display_name": "Nemo", "location_id": "%s", "load_balancer_size": "small","load_balancer_type": "layer-3"}]`, locationID),
		status: http.StatusUnprocessableEntity,
		path:   baseURL,
		method: http.MethodPost,
		tenant: tenantID,
	})

	doHTTPTest(t, &httpTest{
		name:   "missing size",
		body:   fmt.Sprintf(`[{"display_name": "Nemo", "location_id": "%s", "ip_addr": "1.1.1.1","load_balancer_type": "layer-3"}]`, locationID),
		status: http.StatusUnprocessableEntity,
		path:   baseURL,
		method: http.MethodPost,
		tenant: tenantID,
	})

	doHTTPTest(t, &httpTest{
		name:   "missing type",
		body:   fmt.Sprintf(`[{"display_name": "Nemo", "location_id": "%s", "ip_addr": "1.1.1.1","load_balancer_size": "small"}]`, locationID),
		status: http.StatusUnprocessableEntity,
		path:   baseURL,
		method: http.MethodPost,
		tenant: tenantID,
	})

	doHTTPTest(t, &httpTest{
		name:   "invalid type",
		body:   fmt.Sprintf(`[{"display_name": "Nemo", "location_id": "%s", "ip_addr": "1.1.1.1","load_balancer_size": "small","load_balancer_type": "layer-12"}]`, locationID),
		status: http.StatusUnprocessableEntity,
		path:   baseURL,
		method: http.MethodPost,
		tenant: tenantID,
	})

	doHTTPTest(t, &httpTest{
		name:   "bad ip address",
		body:   fmt.Sprintf(`[{"display_name": "Nemo", "location_id": "%s", "ip_addr": "Dori","load_balancer_size": "small","load_balancer_type": "layer-3"}]`, locationID),
		status: http.StatusUnprocessableEntity,
		path:   baseURL,
		method: http.MethodPost,
		tenant: tenantID,
	})

	doHTTPTest(t, &httpTest{
		name:   "ipv6",
		body:   fmt.Sprintf(`[{"display_name": "Nemo", "location_id": "%s", "ip_addr": "2601::","load_balancer_size": "small","load_balancer_type": "layer-3"}]`, locationID),
		status: http.StatusInternalServerError,
		path:   baseURL,
		method: http.MethodPost,
		tenant: tenantID,
	})

	doHTTPTest(t, &httpTest{
		name:   "empty body",
		body:   `[]`,
		status: http.StatusUnprocessableEntity,
		path:   baseURL,
		method: http.MethodPost,
		tenant: tenantID,
	})

	doHTTPTest(t, &httpTest{
		name:   "bad body",
		body:   `bad body`,
		status: http.StatusBadRequest,
		path:   baseURL,
		method: http.MethodPost,
		tenant: tenantID,
	})

	// GET tests
	doHTTPTest(t, &httpTest{
		name:   "happy path",
		path:   baseURL,
		status: http.StatusOK,
		tenant: tenantID,
	})

	doHTTPTest(t, &httpTest{
		name:   "happy path nemo by name",
		path:   baseURL + "?display_name=Nemo",
		status: http.StatusOK,
		tenant: tenantID,
	})

	doHTTPTest(t, &httpTest{
		name:   "happy path nemo by ip",
		path:   baseURL + "?ip_addr=1.1.1.1",
		status: http.StatusOK,
		tenant: tenantID,
	})

	// DELETE tests
	doHTTPTest(t, &httpTest{
		name:   "delete invalid id",
		path:   baseURL + "/invalid",
		status: http.StatusUnprocessableEntity,
		method: http.MethodDelete,
		tenant: tenantID,
	})

	doHTTPTest(t, &httpTest{
		name:   "delete small load balancers",
		path:   baseURL + "?load_balancer_size=small",
		status: http.StatusUnprocessableEntity,
		method: http.MethodDelete,
		tenant: tenantID,
	})

	doHTTPTest(t, &httpTest{
		name:   "delete nemo by Name",
		path:   baseURL + "?slug=nemo",
		status: http.StatusOK,
		method: http.MethodDelete,
		tenant: tenantID,
	})

	doHTTPTest(t, &httpTest{
		name:   "delete Dori by name",
		path:   baseURL + "?slug=dori",
		status: http.StatusOK,
		method: http.MethodDelete,
		tenant: tenantID,
	})

	doHTTPTest(t, &httpTest{
		name:   "delete Dori again",
		path:   baseURL + "?slug=dori",
		status: http.StatusNotFound,
		method: http.MethodDelete,
		tenant: tenantID,
	})
}

func createLoadBalancer(t *testing.T, srv *httptest.Server, locationID string) (*loadBalancer, func(t *testing.T)) {
	tenantID := uuid.New().String()
	baseURL := srv.URL + "/v1/loadbalancers"

	test := &httpTest{
		name:   "create nemo lb",
		body:   fmt.Sprintf(`[{"display_name": "Nemo", "location_id": "%s", "ip_addr": "1.1.1.1","load_balancer_size": "small","load_balancer_type": "layer-3"}]`, locationID),
		path:   baseURL,
		method: http.MethodPost,
		tenant: tenantID,
		status: http.StatusOK,
	}

	doHTTPTest(t, test)

	// get loadbalancer by name
	loadbalancer := response{}

	t.Run("get nemo by name:[POST] "+baseURL+"?display_name=Nemo", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, baseURL+"?display_name=Nemo", nil) //nolint
		assert.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set(tenantHeader, tenantID)

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
			tenant: tenantID,
			path:   baseURL + "?slug=nemo",
			method: http.MethodDelete,
			status: http.StatusOK,
		}

		doHTTPTest(t, test)
	}
}
