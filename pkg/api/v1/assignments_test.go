package api

import (
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_Assignments(t *testing.T) {
	srv := newTestServer(t)
	defer srv.Close()

	assert.NotNil(t, srv)

	baseURL := srv.URL + "/v1/assignments"

	// Create a load balancer
	loadBalancer, cleanupLB := createLoadBalancer(t, srv, uuid.NewString())
	defer cleanupLB(t)

	// Create a frontend
	fe, cleanupFE := createFrontend(t, srv, loadBalancer.ID, loadBalancer.TenantID)
	defer cleanupFE(t)

	// Create a pool
	pool, cleanupPool := createPool(t, srv, loadBalancer.TenantID)
	defer cleanupPool(t)

	// Create an origin in the pool
	origin, cleanupOrigin := createOrigin(t, srv, pool.ID, loadBalancer.TenantID)
	defer cleanupOrigin(t)

	// create an assignment
	doHTTPTest(t, &httpTest{
		name:   "create assignment",
		method: http.MethodPost,
		path:   baseURL,
		body:   `[{"frontend_id": "` + fe.ID + `", "origin_id": "` + origin.ID + `", "load_balancer_id": "` + loadBalancer.ID + `"}]`,
		status: http.StatusOK,
		tenant: loadBalancer.TenantID,
	})

	// Get the assignments
	doHTTPTest(t, &httpTest{
		name:   "get assignments",
		method: http.MethodGet,
		path:   baseURL + "?frontend_id=" + fe.ID + "&pool_id=" + pool.ID,
		status: http.StatusOK,
		tenant: loadBalancer.TenantID,
	})
}
