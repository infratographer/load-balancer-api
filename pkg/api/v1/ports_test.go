package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.infratographer.com/load-balancer-api/internal/httptools"
)

func createPort(t *testing.T, srv *httptest.Server, loadBalancerID string) (*port, func(*testing.T)) {
	baseURL := srv.URL + "/v1/loadbalancers/" + loadBalancerID + "/ports"

	t.Run("create port:[POST]_"+baseURL, func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, baseURL, httptools.FakeBody(fmt.Sprintf(`{"name": "Ears", "port": 25}`))) //nolint
		assert.NoError(t, err)
		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		resp.Body.Close()
	})

	ret := &port{}

	t.Run("get port:[GET]_"+baseURL+"?slug=ears", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, baseURL, nil) //nolint
		assert.NoError(t, err)
		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)

		feResp := response{}

		_ = json.NewDecoder(resp.Body).Decode(&feResp)

		resp.Body.Close()

		for _, fe := range *feResp.Ports {
			if fe.Name == "Ears" {
				ret = fe
			}
		}
	})

	return ret, func(t *testing.T) {
		t.Run("delete port:[DELETE]_"+srv.URL+"/v1/ports/"+ret.ID, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodDelete, srv.URL+"/v1/ports/"+ret.ID, nil) //nolint
			assert.NoError(t, err)
			resp, err := http.DefaultClient.Do(req)
			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			resp.Body.Close()
		})
	}
}

func TestPortRoutes(t *testing.T) {
	nsrv := newNatsTestServer(t, "load-balancer-api-test", "com.infratographer.events.>")
	defer nsrv.Shutdown()

	srv := newTestServer(t, nsrv.ClientURL())
	defer srv.Close()

	lb, cleanupLoadBalancers := createLoadBalancer(t, srv, uuid.New().String())
	defer cleanupLoadBalancers(t)

	loadBalancerID := lb.ID

	baseURL := srv.URL + "/v1/ports"
	baseURLLoadBalancer := srv.URL + "/v1/loadbalancers/" + loadBalancerID + "/ports"
	missingUUID := uuid.New().String()

	req1, err := http.NewRequest(http.MethodPost, baseURLLoadBalancer, httptools.FakeBody(`{"name": "Ears", "port": 25}`)) //nolint
	assert.NoError(t, err)
	resp1, err := http.DefaultClient.Do(req1)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp1.StatusCode)
	resp1.Body.Close()

	req2, err := http.NewRequest(http.MethodPost, baseURLLoadBalancer, httptools.FakeBody(`{"name": "Eyes", "port": 465}`)) //nolint
	assert.NoError(t, err)
	resp2, err := http.DefaultClient.Do(req2)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp2.StatusCode)
	resp2.Body.Close()

	req3, err := http.NewRequest(http.MethodGet, baseURLLoadBalancer, nil) //nolint
	assert.NoError(t, err)
	resp3, err := http.DefaultClient.Do(req3)
	assert.NoError(t, err)

	feResp := response{}

	_ = json.NewDecoder(resp3.Body).Decode(&feResp)

	resp3.Body.Close()

	earsID := ""

	for _, fe := range *feResp.Ports {
		if fe.Name == "Ears" {
			earsID = fe.ID
		}
	}

	// POST tests
	doHTTPTest(t, &httpTest{
		name:   "happy path",
		body:   `{"name": "Mouth", "port": 80}`,
		status: http.StatusOK,
		method: http.MethodPost,
		path:   baseURLLoadBalancer,
	})

	doHTTPTest(t, &httpTest{
		name:   "duplicate",
		body:   `{"name": "Mouth", "port": 80}`,
		status: http.StatusInternalServerError,
		method: http.MethodPost,
		path:   baseURLLoadBalancer,
	})

	doHTTPTest(t, &httpTest{
		name:   "443",
		body:   `{"name": "TLS Mouth", "port": 443}`,
		status: http.StatusOK,
		method: http.MethodPost,
		path:   baseURLLoadBalancer,
	})

	doHTTPTest(t, &httpTest{
		name:   "list of ports",
		body:   `[{"name": "Mouth", "port": 80},{"name": "Beard", "port": 443}]`,
		status: http.StatusBadRequest,
		method: http.MethodPost,
		path:   baseURLLoadBalancer,
	})

	doHTTPTest(t, &httpTest{
		name:   "negative port",
		body:   `{"name": "Mouth", "port": -1}`,
		status: http.StatusBadRequest,
		method: http.MethodPost,
		path:   baseURLLoadBalancer,
	})

	doHTTPTest(t, &httpTest{
		name:   "zero port",
		body:   `{"name": "Mouth", "port": 0}`,
		status: http.StatusBadRequest,
		method: http.MethodPost,
		path:   baseURLLoadBalancer,
	})

	doHTTPTest(t, &httpTest{
		name:   "port too high",
		body:   `{"name": "Mouth", "port": 65536}`,
		status: http.StatusBadRequest,
		method: http.MethodPost,
		path:   baseURLLoadBalancer,
	})

	doHTTPTest(t, &httpTest{
		name:   "missing port",
		body:   `{"name": "Mouth"}`,
		status: http.StatusBadRequest,
		method: http.MethodPost,
		path:   baseURLLoadBalancer,
	})

	doHTTPTest(t, &httpTest{
		name:   "missing display name",
		body:   `{"port": 80}`,
		status: http.StatusBadRequest,
		method: http.MethodPost,
		path:   baseURLLoadBalancer,
	})

	doHTTPTest(t, &httpTest{
		name:   "invalid load balancer id",
		body:   `{"name": "Mouth", "port": 80}`,
		status: http.StatusBadRequest,
		method: http.MethodPost,
		path:   srv.URL + "/v1/loadbalancers/1234/ports",
	})

	doHTTPTest(t, &httpTest{
		name:   "missing body",
		status: http.StatusBadRequest,
		method: http.MethodPost,
		path:   baseURLLoadBalancer,
	})

	doHTTPTest(t, &httpTest{
		name:   "bad body",
		body:   `bad body`,
		status: http.StatusBadRequest,
		method: http.MethodPost,
		path:   baseURLLoadBalancer,
	})

	// PUT tests
	doHTTPTest(t, &httpTest{
		name:   "happy path update port",
		body:   `{"name": "LeftEar", "port": 8080}`,
		status: http.StatusAccepted,
		method: http.MethodPut,
		path:   baseURL + "/" + earsID,
	})

	doHTTPTest(t, &httpTest{
		name:   "update port port too low",
		body:   `{"name": "LeftEar", "port": -1}`,
		status: http.StatusBadRequest,
		method: http.MethodPut,
		path:   baseURL + "/" + earsID,
	})

	doHTTPTest(t, &httpTest{
		name:   "update port port two high",
		body:   `{"name": "LeftEar", "port": 131337}`,
		status: http.StatusBadRequest,
		method: http.MethodPut,
		path:   baseURL + "/" + earsID,
	})

	doHTTPTest(t, &httpTest{
		name:   "update port missing display name",
		body:   `{"port": 8080}`,
		status: http.StatusBadRequest,
		method: http.MethodPut,
		path:   baseURL + "/" + earsID,
	})

	doHTTPTest(t, &httpTest{
		name:   "update port missing port",
		body:   `{"name": "LeftEar"}`,
		status: http.StatusBadRequest,
		method: http.MethodPut,
		path:   baseURL + "/" + earsID,
	})

	doHTTPTest(t, &httpTest{
		name:   "missing port id",
		body:   `{"name": "LeftEar", "port": 8080}`,
		status: http.StatusNotFound,
		method: http.MethodPut,
		path:   baseURL,
	})

	doHTTPTest(t, &httpTest{
		name:   "port not found",
		body:   `{"name": "Plain Mouth", "port": 80}`,
		status: http.StatusAccepted,
		method: http.MethodPut,
		path:   baseURLLoadBalancer + "?port=80",
	})

	// Get Tests
	doHTTPTest(t, &httpTest{
		name:   "happy path",
		path:   baseURL + "/" + earsID,
		status: http.StatusOK,
		method: http.MethodGet,
	})

	doHTTPTest(t, &httpTest{
		name:   "happy path with id",
		path:   baseURLLoadBalancer,
		status: http.StatusOK,
		method: http.MethodGet,
	})

	doHTTPTest(t, &httpTest{
		name:   "not found with query",
		path:   baseURLLoadBalancer + "?slug=not_found",
		status: http.StatusOK,
		method: http.MethodGet,
	})

	doHTTPTest(t, &httpTest{
		name:   "not found with path param",
		status: http.StatusNotFound,
		path:   baseURL + "/" + missingUUID,
		method: http.MethodGet,
	})

	doHTTPTest(t, &httpTest{
		name:   "bad UUID in path param",
		status: http.StatusBadRequest,
		path:   baseURL + "/123456",
		method: http.MethodGet,
	})

	// Delete
	doHTTPTest(t, &httpTest{
		name:   "slug not found",
		path:   baseURLLoadBalancer + "?slug=not_found",
		status: http.StatusNotFound,
		method: http.MethodDelete,
	})

	doHTTPTest(t, &httpTest{
		name:   "port not found",
		path:   baseURLLoadBalancer + "?port=404",
		status: http.StatusNotFound,
		method: http.MethodDelete,
	})

	doHTTPTest(t, &httpTest{
		name:   "delete port with port 80",
		path:   baseURLLoadBalancer + "?slug=mouth&port=80",
		status: http.StatusOK,
		method: http.MethodDelete,
	})

	doHTTPTest(t, &httpTest{
		name:   "delete port with port 443",
		path:   baseURLLoadBalancer + "?slug=tls-mouth&port=443",
		status: http.StatusOK,
		method: http.MethodDelete,
	})

	doHTTPTest(t, &httpTest{
		name:   "delete port Ears by id",
		path:   baseURL + "/" + earsID,
		status: http.StatusOK,
		method: http.MethodDelete,
	})

	doHTTPTest(t, &httpTest{
		name:   "delete port Eyes by port ",
		path:   baseURLLoadBalancer + "?port=465&name=Eyes",
		status: http.StatusOK,
		method: http.MethodDelete,
	})
}
