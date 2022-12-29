package api

import (
	"fmt"
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
	pool, cleanupPool := createPool(t, srv, "marlin", loadBalancer.TenantID)
	defer cleanupPool(t)

	// Create an origin in the pool
	_, cleanupOrigin := createOrigin(t, srv, "bruce", pool.ID, loadBalancer.TenantID)
	defer cleanupOrigin(t)

	// poll2
	pool2, cleanupPool2 := createPool(t, srv, "dory", loadBalancer.TenantID)
	defer cleanupPool2(t)

	// origin2
	_, cleanupOrigin2 := createOrigin(t, srv, "chum", pool2.ID, loadBalancer.TenantID)
	defer cleanupOrigin2(t)

	// create an assignment
	doHTTPTest(t, &httpTest{
		name:   "create assignment",
		method: http.MethodPost,
		path:   baseURL,
		body:   fmt.Sprintf(`[{"frontend_id": "%s", "pool_id": "%s","load_balancer_id":"%s"}]`, fe.ID, pool.ID, loadBalancer.ID),
		status: http.StatusOK,
		tenant: loadBalancer.TenantID,
	})

	doHTTPTest(t, &httpTest{
		name:   "Duplicate assignment",
		method: http.MethodPost,
		path:   baseURL,
		body:   fmt.Sprintf(`[{"frontend_id": "%s", "pool_id": "%s","load_balancer_id":"%s"}]`, fe.ID, pool.ID, loadBalancer.ID),
		status: http.StatusInternalServerError,
		tenant: loadBalancer.TenantID,
	})

	doHTTPTest(t, &httpTest{
		name:   "create assignment2",
		method: http.MethodPost,
		path:   baseURL,
		body:   fmt.Sprintf(`[{"frontend_id": "%s", "pool_id": "%s","load_balancer_id":"%s"}]`, fe.ID, pool2.ID, loadBalancer.ID),
		status: http.StatusOK,
		tenant: loadBalancer.TenantID,
	})

	doHTTPTest(t, &httpTest{
		name:   "Empty body",
		method: http.MethodPost,
		path:   baseURL,
		body:   "",
		status: http.StatusNotFound,
		tenant: loadBalancer.TenantID,
	})

	doHTTPTest(t, &httpTest{
		name:   "Invalid body",
		method: http.MethodPost,
		path:   baseURL,
		body:   "invalid",
		status: http.StatusBadRequest,
		tenant: loadBalancer.TenantID,
	})

	doHTTPTest(t, &httpTest{
		name:   "Invalid frontend",
		method: http.MethodPost,
		path:   baseURL,
		body:   fmt.Sprintf(`[{"frontend_id": "%s", "pool_id": "%s","load_balancer_id":"%s"}]`, "invalid", pool.ID, loadBalancer.ID),
		status: http.StatusInternalServerError,
		tenant: loadBalancer.TenantID,
	})

	doHTTPTest(t, &httpTest{
		name:   "Invalid pool",
		method: http.MethodPost,
		path:   baseURL,
		body:   fmt.Sprintf(`[{"frontend_id": "%s", "pool_id": "%s","load_balancer_id":"%s"}]`, fe.ID, "invalid", loadBalancer.ID),
		status: http.StatusInternalServerError,
		tenant: loadBalancer.TenantID,
	})

	doHTTPTest(t, &httpTest{
		name:   "Invalid load balancer",
		method: http.MethodPost,
		path:   baseURL,
		body:   fmt.Sprintf(`[{"frontend_id": "%s", "pool_id": "%s","load_balancer_id":"%s"}]`, fe.ID, pool.ID, "invalid"),
		status: http.StatusInternalServerError,
		tenant: loadBalancer.TenantID,
	})

	doHTTPTest(t, &httpTest{
		name:   "Invalid tenant",
		method: http.MethodPost,
		path:   baseURL,
		body:   fmt.Sprintf(`[{"frontend_id": "%s", "pool_id": "%s","load_balancer_id":"%s"}]`, fe.ID, pool.ID, loadBalancer.ID),
		status: http.StatusInternalServerError,
		tenant: "invalid",
	})
	// Get the assignments
	doHTTPTest(t, &httpTest{
		name:   "get assignments",
		method: http.MethodGet,
		path:   baseURL + "?frontend_id=" + fe.ID + "&pool_id=" + pool.ID,
		status: http.StatusOK,
		tenant: loadBalancer.TenantID,
	})

	// Delete the assignment
	doHTTPTest(t, &httpTest{
		name:   "delete ambiguous",
		method: http.MethodDelete,
		path:   baseURL + "?load_balancer_id=" + loadBalancer.ID,
		status: http.StatusBadRequest,
		tenant: loadBalancer.TenantID,
	})

	doHTTPTest(t, &httpTest{
		name:   "delete assignment",
		method: http.MethodDelete,
		path:   baseURL + "?frontend_id=" + fe.ID + "&pool_id=" + pool.ID,
		status: http.StatusOK,
		tenant: loadBalancer.TenantID,
	})

	doHTTPTest(t, &httpTest{
		name:   "duplicate delete",
		method: http.MethodDelete,
		path:   baseURL + "?frontend_id=" + fe.ID + "&pool_id=" + pool.ID,
		status: http.StatusNotFound,
		tenant: loadBalancer.TenantID,
	})

	doHTTPTest(t, &httpTest{
		name:   "delete assignment2",
		method: http.MethodDelete,
		path:   baseURL + "?frontend_id=" + fe.ID + "&pool_id=" + pool2.ID,
		status: http.StatusOK,
		tenant: loadBalancer.TenantID,
	})
}
