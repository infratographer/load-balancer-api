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

func createFrontend(t *testing.T, srv *httptest.Server, loadBalancerID string, tenantID string) (*frontend, func(*testing.T)) {
	baseURL := srv.URL + "/v1/frontends"

	t.Run("create frontend:[POST]_"+baseURL, func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, baseURL, httptools.FakeBody(fmt.Sprintf(`{"display_name": "Ears", "load_balancer_id": "%s", "port": 25}`, loadBalancerID))) //nolint
		assert.NoError(t, err)
		req.Header.Set(tenantHeader, tenantID)
		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		resp.Body.Close()
	})

	ret := &frontend{}

	t.Run("get frontend:[GET]_"+baseURL+"?slug=ears", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, srv.URL+"/v1/loadbalancers/"+loadBalancerID+"/frontends", nil) //nolint
		assert.NoError(t, err)
		req.Header.Set(tenantHeader, tenantID)
		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)

		feResp := response{}

		_ = json.NewDecoder(resp.Body).Decode(&feResp)

		resp.Body.Close()

		for _, fe := range *feResp.Frontends {
			if fe.Name == "Ears" {
				ret = fe
			}
		}
	})

	return ret, func(t *testing.T) {
		t.Run("delete frontend:[DELETE]_"+srv.URL+"/v1/frontends/"+ret.ID, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodDelete, srv.URL+"/v1/frontends/"+ret.ID, nil) //nolint
			assert.NoError(t, err)
			req.Header.Set(tenantHeader, tenantID)
			resp, err := http.DefaultClient.Do(req)
			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			resp.Body.Close()
		})
	}
}

func TestFrondendRoutes(t *testing.T) {
	srv := newTestServer(t)
	defer srv.Close()

	lb, cleanupLoadBalancers := createLoadBalancer(t, srv, uuid.New().String())
	defer cleanupLoadBalancers(t)

	loadBalancerID := lb.ID

	tenantID := uuid.New().String()
	baseURL := srv.URL + "/v1/frontends"
	missingUUID := uuid.New().String()

	req1, err := http.NewRequest(http.MethodPost, baseURL, httptools.FakeBody(fmt.Sprintf(`{"display_name": "Ears", "load_balancer_id": "%s", "port": 25}`, loadBalancerID))) //nolint
	assert.NoError(t, err)
	req1.Header.Set(tenantHeader, tenantID)
	resp1, err := http.DefaultClient.Do(req1)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp1.StatusCode)
	resp1.Body.Close()

	req2, err := http.NewRequest(http.MethodPost, baseURL, httptools.FakeBody(fmt.Sprintf(`{"display_name": "Eyes", "port": 465, "load_balancer_id" : "%s"}`, loadBalancerID))) //nolint
	assert.NoError(t, err)
	req2.Header.Set(tenantHeader, tenantID)
	resp2, err := http.DefaultClient.Do(req2)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp2.StatusCode)
	resp2.Body.Close()

	req3, err := http.NewRequest(http.MethodGet, srv.URL+"/v1/loadbalancers/"+loadBalancerID+"/frontends", nil) //nolint
	assert.NoError(t, err)
	req3.Header.Set(tenantHeader, tenantID)
	resp3, err := http.DefaultClient.Do(req3)
	assert.NoError(t, err)

	feResp := response{}

	_ = json.NewDecoder(resp3.Body).Decode(&feResp)

	resp3.Body.Close()

	earsID := ""

	for _, fe := range *feResp.Frontends {
		if fe.Name == "Ears" {
			earsID = fe.ID
		}
	}

	// POST tests
	doHTTPTest(t, &httpTest{
		name:   "happy path",
		body:   fmt.Sprintf(`{"display_name": "Mouth", "load_balancer_id": "%s", "port": 80}`, loadBalancerID),
		status: http.StatusOK,
		method: http.MethodPost,
		path:   baseURL,
		tenant: tenantID,
	})

	doHTTPTest(t, &httpTest{
		name:   "duplicate",
		body:   fmt.Sprintf(`{"display_name": "Mouth", "load_balancer_id": "%s", "port": 80}`, loadBalancerID),
		status: http.StatusInternalServerError,
		method: http.MethodPost,
		path:   baseURL,
		tenant: tenantID,
	})

	doHTTPTest(t, &httpTest{
		name:   "443",
		body:   fmt.Sprintf(`{"display_name": "TLS Mouth", "load_balancer_id": "%s", "port": 443}`, loadBalancerID),
		status: http.StatusOK,
		method: http.MethodPost,
		path:   baseURL,
		tenant: tenantID,
	})

	doHTTPTest(t, &httpTest{
		name:   "list of frontends",
		body:   fmt.Sprintf(`[{"display_name": "Mouth", "load_balancer_id": "%s", "port": 80},{"display_name": "Beard", "load_balancer_id": "%s", "port": 443}]`, loadBalancerID, loadBalancerID),
		status: http.StatusBadRequest,
		method: http.MethodPost,
		path:   baseURL,
		tenant: tenantID,
	})

	doHTTPTest(t, &httpTest{
		name:   "negative port",
		body:   fmt.Sprintf(`{"display_name": "Mouth", "load_balancer_id": "%s", "port": -1}`, loadBalancerID),
		status: http.StatusBadRequest,
		method: http.MethodPost,
		path:   baseURL,
		tenant: tenantID,
	})

	doHTTPTest(t, &httpTest{
		name:   "zero port",
		body:   fmt.Sprintf(`{"display_name": "Mouth", "load_balancer_id": "%s", "port": 0}`, loadBalancerID),
		status: http.StatusBadRequest,
		method: http.MethodPost,
		path:   baseURL,
		tenant: tenantID,
	})

	doHTTPTest(t, &httpTest{
		name:   "port too high",
		body:   fmt.Sprintf(`{"display_name": "Mouth", "load_balancer_id": "%s", "port": 65536}`, loadBalancerID),
		status: http.StatusBadRequest,
		method: http.MethodPost,
		path:   baseURL,
		tenant: tenantID,
	})

	doHTTPTest(t, &httpTest{
		name:   "missing port",
		body:   fmt.Sprintf(`{"display_name": "Mouth", "load_balancer_id": "%s"}`, loadBalancerID),
		status: http.StatusBadRequest,
		method: http.MethodPost,
		path:   baseURL,
		tenant: tenantID,
	})

	doHTTPTest(t, &httpTest{
		name:   "missing display name",
		body:   fmt.Sprintf(`{"load_balancer_id": "%s", "port": 80}`, loadBalancerID),
		status: http.StatusBadRequest,
		method: http.MethodPost,
		path:   baseURL,
		tenant: tenantID,
	})

	doHTTPTest(t, &httpTest{
		name:   "missing load balancer id",
		body:   `{"display_name": "Mouth", "port": 80}`,
		status: http.StatusBadRequest,
		method: http.MethodPost,
		path:   baseURL,
		tenant: tenantID,
	})

	doHTTPTest(t, &httpTest{
		name:   "invalid load balancer id",
		body:   `{"display_name": "Mouth", "port": 80, "load_balancer_id": "bad id"}`,
		status: http.StatusBadRequest,
		method: http.MethodPost,
		path:   baseURL,
		tenant: tenantID,
	})

	doHTTPTest(t, &httpTest{
		name:   "missing body",
		status: http.StatusBadRequest,
		method: http.MethodPost,
		path:   baseURL,
		tenant: tenantID,
	})

	doHTTPTest(t, &httpTest{
		name:   "bad body",
		body:   `bad body`,
		status: http.StatusBadRequest,
		method: http.MethodPost,
		path:   baseURL,
		tenant: tenantID,
	})

	doHTTPTest(t, &httpTest{
		name:   "bad tenant id",
		status: http.StatusBadRequest,
		method: http.MethodPost,
		path:   baseURL,
		tenant: "bad tenant id",
	})

	// Get Tests
	doHTTPTest(t, &httpTest{
		name:   "happy path",
		path:   baseURL,
		status: http.StatusOK,
		method: http.MethodGet,
		tenant: tenantID,
	})

	doHTTPTest(t, &httpTest{
		name:   "happy pathwith id",
		path:   srv.URL + "/v1/loadbalancers/" + loadBalancerID + "/frontends",
		status: http.StatusOK,
		method: http.MethodGet,
		tenant: tenantID,
	})

	doHTTPTest(t, &httpTest{
		name:   "not found with query",
		path:   baseURL + "?slug=not_found",
		status: http.StatusOK,
		method: http.MethodGet,
		tenant: tenantID,
	})

	doHTTPTest(t, &httpTest{
		name:   "list frontends with invalid tenant id",
		path:   baseURL,
		status: http.StatusBadRequest,
		method: http.MethodGet,
		tenant: "123456",
	})

	doHTTPTest(t, &httpTest{
		name:   "list frontends with unknown tenant id",
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

	// Delete
	doHTTPTest(t, &httpTest{
		name:   "404",
		path:   baseURL + "?slug=not_found",
		status: http.StatusNotFound,
		method: http.MethodDelete,
		tenant: tenantID,
	})

	doHTTPTest(t, &httpTest{
		name:   "delete frontend with port 80",
		path:   baseURL + "?slug=mouth&port=80",
		status: http.StatusOK,
		method: http.MethodDelete,
		tenant: tenantID,
	})

	doHTTPTest(t, &httpTest{
		name:   "delete frontend with port 443",
		path:   baseURL + "?slug=tls-mouth&port=443",
		status: http.StatusOK,
		method: http.MethodDelete,
		tenant: tenantID,
	})

	doHTTPTest(t, &httpTest{
		name:   "delete frontend Ears by id",
		path:   baseURL + "/" + earsID,
		status: http.StatusOK,
		method: http.MethodDelete,
		tenant: tenantID,
	})

	doHTTPTest(t, &httpTest{
		name:   "delete frontend Eyes by port ",
		path:   srv.URL + "/v1/loadbalancers/" + loadBalancerID + "/frontends?port=465&display_name=Eyes",
		status: http.StatusOK,
		method: http.MethodDelete,
		tenant: tenantID,
	})
}
