package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.infratographer.com/loadbalancerapi/internal/httptools"
)

func TestLoadBalancerRoutes(t *testing.T) {
	srv := newTestServer(t)
	defer srv.Close()

	tenantID := uuid.New().String()
	baseURL := srv.URL + "/v1/tenant/" + tenantID + "/loadbalancers"

	locationResp, cleanupAnemones := createAnenmoes(t, srv)
	defer cleanupAnemones(t)

	locationID := locationResp.Location.ID

	var payloadTests = []struct {
		name   string
		body   string
		status int
	}{
		{
			name:   "happy path",
			body:   fmt.Sprintf(`[{"display_name": "Nemo", "location_id": "%s", "ip_addr": "1.1.1.1","load_balancer_size": "small","load_balancer_type": "layer-3"}]`, locationID),
			status: http.StatusCreated,
		},
		{
			name:   "happy path 2",
			body:   fmt.Sprintf(`[{"display_name": "Dori", "location_id": "%s", "ip_addr": "1.2.1.1","load_balancer_size": "small","load_balancer_type": "layer-3"}]`, locationID),
			status: http.StatusCreated,
		},
		{
			name:   "Duplicate",
			body:   fmt.Sprintf(`[{"display_name": "Nemo", "location_id": "%s", "ip_addr": "1.1.1.1","load_balancer_size": "small","load_balancer_type": "layer-3"}]`, locationID),
			status: http.StatusBadRequest,
		},
		{
			name:   "missing display name",
			body:   fmt.Sprintf(`[{"location_id": "%s", "ip_addr": "1.1.1.1","load_balancer_size": "small","load_balancer_type": "layer-3"}]`, locationID),
			status: http.StatusBadRequest,
		},
		{
			name:   "missing location id",
			body:   `[{"display_name": "Nemo", "ip_addr": "1.1.1.1","load_balancer_size": "small","load_balancer_type": "layer-3"}]`,
			status: http.StatusBadRequest,
		},
		{
			name:   "missing size",
			body:   fmt.Sprintf(`[{"display_name": "Nemo", "location_id": "%s", "ip_addr": "1.1.1.1","load_balancer_type": "layer-3"}]`, locationID),
			status: http.StatusBadRequest,
		},
		{
			name:   "missing type",
			body:   fmt.Sprintf(`[{"display_name": "Nemo", "location_id": "%s", "ip_addr": "1.1.1.1","load_balancer_size": "small"}]`, locationID),
			status: http.StatusBadRequest,
		},
		{
			name:   "invalid type",
			body:   fmt.Sprintf(`[{"display_name": "Nemo", "location_id": "%s", "ip_addr": "1.1.1.1","load_balancer_size": "small","load_balancer_type": "layer-12"}]`, locationID),
			status: http.StatusBadRequest,
		},
		{
			name:   "bad ip address",
			body:   fmt.Sprintf(`[{"display_name": "Nemo", "location_id": "%s", "ip_addr": "Dori","load_balancer_size": "small","load_balancer_type": "layer-3"}]`, locationID),
			status: http.StatusBadRequest,
		},
		{
			name:   "ipv6",
			body:   fmt.Sprintf(`[{"display_name": "Nemo", "location_id": "%s", "ip_addr": "2601::","load_balancer_size": "small","load_balancer_type": "layer-3"}]`, locationID),
			status: http.StatusBadRequest,
		},
		{
			name:   "empty body",
			body:   `[]`,
			status: http.StatusUnprocessableEntity,
		},
		{
			name:   "bad body",
			body:   `bad body`,
			status: http.StatusBadRequest,
		},
	}

	for _, tt := range payloadTests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := http.Post(baseURL, "application/json", httptools.FakeBody(tt.body)) //nolint
			assert.NoError(t, err)
			assert.Equal(t, tt.status, resp.StatusCode)
			resp.Body.Close()
		})
	}

	// get nemo load balancer id
	nemo := response{}
	resp, err := http.Get(baseURL + "?display_name=Nemo") //nolint
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	err = json.NewDecoder(resp.Body).Decode(&nemo)
	assert.NoError(t, err)
	assert.Equal(t, 36 /*uuid*/, len(nemo.LoadBalancer.ID))

	resp.Body.Close()

	var getTests = []struct {
		name   string
		url    string
		status int
	}{
		{
			name:   "happy path",
			url:    baseURL,
			status: http.StatusOK,
		},
		{
			name:   "happy path nemo by name",
			url:    baseURL + "?display_name=Nemo",
			status: http.StatusOK,
		},
		{
			name:   "happy path nemo by ip",
			url:    baseURL + "?ip_addr=1.1.1.1",
			status: http.StatusOK,
		},
		{
			name:   "happy path nemo by id",
			url:    baseURL + "/" + nemo.LoadBalancer.ID,
			status: http.StatusOK,
		},
	}

	for _, tt := range getTests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := http.Get(tt.url) //nolint
			assert.NoError(t, err)
			assert.Equal(t, tt.status, resp.StatusCode)
			resp.Body.Close()
		})
	}

	// delete nemo by id
	t.Run("delete nemo by id", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodDelete, baseURL+"/"+nemo.LoadBalancer.ID, nil)
		assert.NoError(t, err)

		resp, err = http.DefaultClient.Do(req) //nolint
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		resp.Body.Close()
	})

	// delete dori by ip
	t.Run("delete dori by ip", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodDelete, baseURL+"?ip_addr=1.2.1.1", nil) //nolint
		assert.NoError(t, err)

		resp, err = http.DefaultClient.Do(req) //nolint
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		resp.Body.Close()
	})
}

// nolint
func createNemoLB(t *testing.T, srv *httptest.Server) (*response, func(t *testing.T)) {
	tenantID := uuid.New().String()
	loc, cleanupLoc := createAnenmoes(t, srv)
	locationID := loc.Location.ID
	baseURL := srv.URL + "/v1/tenant/" + tenantID + "/loadbalancers"
	// post to create nemo
	t.Run("create nemo", func(t *testing.T) {

		nemoBody := fmt.Sprintf(`[{
		"tenant_id": "%s",
		"location_id": "%s",
		"load_balancer_size": "small",
		"load_balancer_type": "layer-3",
		"ip_addr": "2.2.2.2",
		"display_name": "Nemo",
	}]`, tenantID, locationID)

		resp, err := http.Post(baseURL, "application/json", httptools.FakeBody(nemoBody)) //nolint
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		resp.Body.Close()
	})

	// get nemo by name
	nemo := response{}

	t.Run("get nemo by name", func(t *testing.T) {
		resp, err := http.Get(baseURL + "?display_name=Nemo")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		err = json.NewDecoder(resp.Body).Decode(&nemo)
		assert.NoError(t, err)
		assert.Equal(t, 36 /*uuid*/, len(nemo.LoadBalancer.ID))
		resp.Body.Close()
	})

	return &nemo, func(t *testing.T) {
		// delete nemo by id
		t.Run("delete nemo by id", func(t *testing.T) {
			req, err := http.NewRequest(http.MethodDelete, baseURL+"/"+nemo.LoadBalancer.ID, nil)
			assert.NoError(t, err)

			resp, err := http.DefaultClient.Do(req)
			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)

			resp.Body.Close()

			cleanupLoc(t)
		})
	}
}
