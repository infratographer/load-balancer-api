package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.infratographer.com/loadbalancerapi/internal/httptools"
)

func TestFrondendRoutes(t *testing.T) {
	srv := newTestServer(t)
	defer srv.Close()

	tenantID := uuid.New().String()
	baseURL := srv.URL + "/v1/frontends"

	nemo, cleanupLoadBalancers := createNemoLB(t, srv)
	defer cleanupLoadBalancers(t)

	loadBalancerID := (*nemo.LoadBalancers)[0].ID
	// locationID := nemo.LoadBalancer.LocationID

	req, err := http.NewRequest(http.MethodPost, baseURL, httptools.FakeBody(fmt.Sprintf(`[{"display_name": "Ears", "load_balancer_id": "%s", "port": 25},{"display_name": "Eyes", "port": 465, "load_balancer_id" : "%s"}]`, loadBalancerID, loadBalancerID))) //nolint
	assert.NoError(t, err)
	req.Header.Set(tenantHeader, tenantID)
	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	resp.Body.Close()

	req, err = http.NewRequest(http.MethodGet, srv.URL+"/v1/loadbalancers/"+loadBalancerID+"/frontends", nil) //nolint
	assert.NoError(t, err)
	req.Header.Set(tenantHeader, tenantID)
	resp, err = http.DefaultClient.Do(req)
	assert.NoError(t, err)

	feResp := response{}

	_ = json.NewDecoder(resp.Body).Decode(&feResp)

	resp.Body.Close()

	earsID := ""

	for _, fe := range *feResp.Frontends {
		if fe.Name == "Ears" {
			earsID = fe.ID
		}
	}

	var tests = []httpTest{
		{
			name:   "happy path",
			body:   fmt.Sprintf(`[{"display_name": "Mouth", "load_balancer_id": "%s", "port": 80}]`, loadBalancerID),
			status: http.StatusCreated,
			method: http.MethodPost,
			path:   baseURL,
			tenant: tenantID,
		},
		{
			name:   "duplicate",
			body:   fmt.Sprintf(`[{"display_name": "Mouth", "load_balancer_id": "%s", "port": 80}]`, loadBalancerID),
			status: http.StatusInternalServerError,
			method: http.MethodPost,
			path:   baseURL,
			tenant: tenantID,
		},
		{
			name:   "443",
			body:   fmt.Sprintf(`[{"display_name": "Mouth", "load_balancer_id": "%s", "port": 443}]`, loadBalancerID),
			status: http.StatusCreated,
			method: http.MethodPost,
			path:   baseURL,
			tenant: tenantID,
		},
		{
			name:   "negative port",
			body:   fmt.Sprintf(`[{"display_name": "Mouth", "load_balancer_id": "%s", "port": -1}]`, loadBalancerID),
			status: http.StatusInternalServerError,
			method: http.MethodPost,
			path:   baseURL,
			tenant: tenantID,
		},
		{
			name:   "zero port",
			body:   fmt.Sprintf(`[{"display_name": "Mouth", "load_balancer_id": "%s", "port": 0}]`, loadBalancerID),
			status: http.StatusInternalServerError,
			method: http.MethodPost,
			path:   baseURL,
			tenant: tenantID,
		},
		{
			name:   "port too high",
			body:   fmt.Sprintf(`[{"display_name": "Mouth", "load_balancer_id": "%s", "port": 65536}]`, loadBalancerID),
			status: http.StatusInternalServerError,
			method: http.MethodPost,
			path:   baseURL,
			tenant: tenantID,
		},
		{
			name:   "missing port",
			body:   fmt.Sprintf(`[{"display_name": "Mouth", "load_balancer_id": "%s"}]`, loadBalancerID),
			status: http.StatusInternalServerError,
			method: http.MethodPost,
			path:   baseURL,
			tenant: tenantID,
		},
		{
			name:   "missing display name",
			body:   fmt.Sprintf(`[{"load_balancer_id": "%s", "port": 80}]`, loadBalancerID),
			status: http.StatusInternalServerError,
			method: http.MethodPost,
			path:   baseURL,
			tenant: tenantID,
		},
		{
			name:   "missing load balancer id",
			body:   `[{"display_name": "Mouth", "port": 80}]`,
			status: http.StatusInternalServerError,
			method: http.MethodPost,
			path:   baseURL,
			tenant: tenantID,
		},
		{
			name:   "missing body",
			status: http.StatusUnprocessableEntity,
			method: http.MethodPost,
			path:   baseURL,
			tenant: tenantID,
		},
		{
			name:   "bad body",
			body:   `bad body`,
			status: http.StatusBadRequest,
			method: http.MethodPost,
			path:   baseURL,
			tenant: tenantID,
		},
		{
			name:   "bad tenant id",
			status: http.StatusBadRequest,
			method: http.MethodPost,
			path:   baseURL,
			tenant: "bad tenant id",
		},
		// Get
		{
			name:   "happy path",
			path:   baseURL,
			status: http.StatusOK,
			method: http.MethodGet,
			tenant: tenantID,
		},
		{
			name:   "happy path with id",
			path:   srv.URL + "/v1/loadbalancers/" + loadBalancerID + "/frontends",
			status: http.StatusOK,
			method: http.MethodGet,
			tenant: tenantID,
		},
		// Delete
		{
			name:   "ambiguous delete",
			path:   baseURL + "?display_name=Mouth",
			status: http.StatusBadRequest,
			method: http.MethodDelete,
			tenant: tenantID,
		},
		{
			name:   "delete fronted with port 80",
			path:   baseURL + "?display_name=Mouth&port=80",
			status: http.StatusNoContent,
			method: http.MethodDelete,
			tenant: tenantID,
		},
		{
			name:   "delete fronted with port 443",
			path:   baseURL + "?display_name=Mouth&port=443",
			status: http.StatusNoContent,
			method: http.MethodDelete,
			tenant: tenantID,
		},
		{
			name:   "delete fronted Ears by id",
			path:   baseURL + "/" + earsID,
			status: http.StatusNoContent,
			method: http.MethodDelete,
			tenant: tenantID,
		},
		{
			name:   "delete fronted Eyes by port ",
			path:   srv.URL + "/v1/loadbalancers/" + loadBalancerID + "/frontends?port=465&display_name=Eyes",
			status: http.StatusNoContent,
			method: http.MethodDelete,
			tenant: tenantID,
		},
	}

	doHTTPTests(t, tests)
}
