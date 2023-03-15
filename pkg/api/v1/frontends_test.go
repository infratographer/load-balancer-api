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

func createFrontend(t *testing.T, srv *httptest.Server, loadBalancerID string) (*frontend, func(*testing.T)) {
	baseURL := srv.URL + "/v1/loadbalancers/" + loadBalancerID + "/frontends"

	t.Run("create frontend:[POST]_"+baseURL, func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, baseURL, httptools.FakeBody(fmt.Sprintf(`{"display_name": "Ears", "port": 25}`))) //nolint
		assert.NoError(t, err)
		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		resp.Body.Close()
	})

	ret := &frontend{}

	t.Run("get frontend:[GET]_"+baseURL+"?slug=ears", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, baseURL, nil) //nolint
		assert.NoError(t, err)
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
			resp, err := http.DefaultClient.Do(req)
			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			resp.Body.Close()
		})
	}
}

func TestFrondendRoutes(t *testing.T) {
	srv := newTestServer(t, natsSrv.ClientURL())
	defer srv.Close()

	lb, cleanupLoadBalancers := createLoadBalancer(t, srv, uuid.New().String())
	defer cleanupLoadBalancers(t)

	loadBalancerID := lb.ID

	baseURL := srv.URL + "/v1/frontends"
	baseURLLoadBalancer := srv.URL + "/v1/loadbalancers/" + loadBalancerID + "/frontends"
	missingUUID := uuid.New().String()

	req1, err := http.NewRequest(http.MethodPost, baseURLLoadBalancer, httptools.FakeBody(`{"display_name": "Ears", "port": 25}`)) //nolint
	assert.NoError(t, err)
	resp1, err := http.DefaultClient.Do(req1)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp1.StatusCode)
	resp1.Body.Close()

	req2, err := http.NewRequest(http.MethodPost, baseURLLoadBalancer, httptools.FakeBody(`{"display_name": "Eyes", "port": 465}`)) //nolint
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

	for _, fe := range *feResp.Frontends {
		if fe.Name == "Ears" {
			earsID = fe.ID
		}
	}

	// POST tests
	doHTTPTest(t, &httpTest{
		name:   "happy path",
		body:   `{"display_name": "Mouth", "port": 80}`,
		status: http.StatusOK,
		method: http.MethodPost,
		path:   baseURLLoadBalancer,
	})

	doHTTPTest(t, &httpTest{
		name:   "duplicate",
		body:   `{"display_name": "Mouth", "port": 80}`,
		status: http.StatusInternalServerError,
		method: http.MethodPost,
		path:   baseURLLoadBalancer,
	})

	doHTTPTest(t, &httpTest{
		name:   "443",
		body:   `{"display_name": "TLS Mouth", "port": 443}`,
		status: http.StatusOK,
		method: http.MethodPost,
		path:   baseURLLoadBalancer,
	})

	doHTTPTest(t, &httpTest{
		name:   "list of frontends",
		body:   `[{"display_name": "Mouth", "port": 80},{"display_name": "Beard", "port": 443}]`,
		status: http.StatusBadRequest,
		method: http.MethodPost,
		path:   baseURLLoadBalancer,
	})

	doHTTPTest(t, &httpTest{
		name:   "negative port",
		body:   `{"display_name": "Mouth", "port": -1}`,
		status: http.StatusBadRequest,
		method: http.MethodPost,
		path:   baseURLLoadBalancer,
	})

	doHTTPTest(t, &httpTest{
		name:   "zero port",
		body:   `{"display_name": "Mouth", "port": 0}`,
		status: http.StatusBadRequest,
		method: http.MethodPost,
		path:   baseURLLoadBalancer,
	})

	doHTTPTest(t, &httpTest{
		name:   "port too high",
		body:   `{"display_name": "Mouth", "port": 65536}`,
		status: http.StatusBadRequest,
		method: http.MethodPost,
		path:   baseURLLoadBalancer,
	})

	doHTTPTest(t, &httpTest{
		name:   "missing port",
		body:   `{"display_name": "Mouth"}`,
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
		body:   `{"display_name": "Mouth", "port": 80}`,
		status: http.StatusBadRequest,
		method: http.MethodPost,
		path:   srv.URL + "/v1/loadbalancers/1234/frontends",
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
		name:   "happy path update frontend",
		body:   `{"display_name": "LeftEar", "port": 8080}`,
		status: http.StatusAccepted,
		method: http.MethodPut,
		path:   baseURL + "/" + earsID,
	})

	doHTTPTest(t, &httpTest{
		name:   "update frontend port too low",
		body:   `{"display_name": "LeftEar", "port": -1}`,
		status: http.StatusBadRequest,
		method: http.MethodPut,
		path:   baseURL + "/" + earsID,
	})

	doHTTPTest(t, &httpTest{
		name:   "update frontend port two high",
		body:   `{"display_name": "LeftEar", "port": 131337}`,
		status: http.StatusBadRequest,
		method: http.MethodPut,
		path:   baseURL + "/" + earsID,
	})

	doHTTPTest(t, &httpTest{
		name:   "update frontend missing display name",
		body:   `{"port": 8080}`,
		status: http.StatusBadRequest,
		method: http.MethodPut,
		path:   baseURL + "/" + earsID,
	})

	doHTTPTest(t, &httpTest{
		name:   "update frontend missing port",
		body:   `{"display_name": "LeftEar"}`,
		status: http.StatusBadRequest,
		method: http.MethodPut,
		path:   baseURL + "/" + earsID,
	})

	doHTTPTest(t, &httpTest{
		name:   "missing frontend id",
		body:   `{"display_name": "LeftEar", "port": 8080}`,
		status: http.StatusNotFound,
		method: http.MethodPut,
		path:   baseURL,
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
	// doHTTPTest(t, &httpTest{
	// 	name:   "404",
	// 	path:   baseURLLoadBalancer + "?slug=not_found",
	// 	status: http.StatusNotFound,
	// 	method: http.MethodDelete,
	// })

	// doHTTPTest(t, &httpTest{
	// 	name:   "delete frontend with port 80",
	// 	path:   baseURLLoadBalancer + "?slug=mouth&port=80",
	// 	status: http.StatusOK,
	// 	method: http.MethodDelete,
	// })

	// doHTTPTest(t, &httpTest{
	// 	name:   "delete frontend with port 443",
	// 	path:   baseURLLoadBalancer + "?slug=tls-mouth&port=443",
	// 	status: http.StatusOK,
	// 	method: http.MethodDelete,
	// })

	doHTTPTest(t, &httpTest{
		name:   "delete frontend Ears by id",
		path:   baseURL + "/" + earsID,
		status: http.StatusOK,
		method: http.MethodDelete,
	})

	// doHTTPTest(t, &httpTest{
	// 	name:   "delete frontend Eyes by port ",
	// 	path:   baseURLLoadBalancer + "?port=465&display_name=Eyes",
	// 	status: http.StatusOK,
	// 	method: http.MethodDelete,
	// })
}
