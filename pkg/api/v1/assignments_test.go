package api

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_Assignments(t *testing.T) {
	nsrv := newNatsTestServer(t, "load-balancer-api-test", "com.infratographer.events.>")
	defer nsrv.Shutdown()

	srv := newTestServer(t, nsrv.ClientURL())
	defer srv.Close()

	assert.NotNil(t, srv)

	tenantID := uuid.NewString()
	baseURL := srv.URL + "/v1/tenant/" + tenantID + "/assignments"

	// Create a load balancer
	loadBalancer, cleanupLB := createLoadBalancer(t, srv, tenantID)
	defer cleanupLB(t)

	// Create a frontend
	fe, cleanupFE := createFrontend(t, srv, loadBalancer.ID)
	defer cleanupFE(t)

	// Create a pool
	pool, cleanupPool := createPool(t, srv, "marlin", loadBalancer.TenantID)
	defer cleanupPool(t)

	// Create an origin in the pool
	_, cleanupOrigin := createOrigin(t, srv, "bruce", pool.ID)
	defer cleanupOrigin(t)

	// poll2
	pool2, cleanupPool2 := createPool(t, srv, "dory", loadBalancer.TenantID)
	defer cleanupPool2(t)

	// origin2
	_, cleanupOrigin2 := createOrigin(t, srv, "chum", pool2.ID)
	defer cleanupOrigin2(t)

	// create an assignment
	doHTTPTest(t, &httpTest{
		name:   "create assignment",
		method: http.MethodPost,
		path:   baseURL,
		body:   fmt.Sprintf(`{"frontend_id": "%s", "pool_id": "%s"}`, fe.ID, pool.ID),
		status: http.StatusOK,
	})

	doHTTPTest(t, &httpTest{
		name:   "Duplicate assignment",
		method: http.MethodPost,
		path:   baseURL,
		body:   fmt.Sprintf(`{"frontend_id": "%s", "pool_id": "%s"}`, fe.ID, pool.ID),
		status: http.StatusInternalServerError,
	})

	doHTTPTest(t, &httpTest{
		name:   "create assignment2",
		method: http.MethodPost,
		path:   baseURL,
		body:   fmt.Sprintf(`{"frontend_id": "%s", "pool_id": "%s"}`, fe.ID, pool2.ID),
		status: http.StatusOK,
	})

	doHTTPTest(t, &httpTest{
		name:   "Empty body",
		method: http.MethodPost,
		path:   baseURL,
		body:   "",
		status: http.StatusBadRequest,
	})

	doHTTPTest(t, &httpTest{
		name:   "Invalid body",
		method: http.MethodPost,
		path:   baseURL,
		body:   "invalid",
		status: http.StatusBadRequest,
	})

	doHTTPTest(t, &httpTest{
		name:   "Invalid frontend",
		method: http.MethodPost,
		path:   baseURL,
		body:   fmt.Sprintf(`{"frontend_id": "%s", "pool_id": "%s"}`, "invalid", pool.ID),
		status: http.StatusBadRequest,
	})

	doHTTPTest(t, &httpTest{
		name:   "Invalid pool",
		method: http.MethodPost,
		path:   baseURL,
		body:   fmt.Sprintf(`{"frontend_id": "%s", "pool_id": "%s"}`, fe.ID, "invalid"),
		status: http.StatusBadRequest,
	})

	doHTTPTest(t, &httpTest{
		name:   "Invalid tenant",
		method: http.MethodPost,
		path:   srv.URL + "/v1/tenant/invalid/assignments",
		body:   fmt.Sprintf(`{"frontend_id": "%s", "pool_id": "%s"}`, fe.ID, pool.ID),
		status: http.StatusInternalServerError,
	})
	// Get the assignments
	doHTTPTest(t, &httpTest{
		name:   "get assignments",
		method: http.MethodGet,
		path:   baseURL + "?frontend_id=" + fe.ID + "&pool_id=" + pool.ID,
		status: http.StatusOK,
	})

	// Delete the assignment
	doHTTPTest(t, &httpTest{
		name:   "delete ambiguous",
		method: http.MethodDelete,
		path:   baseURL + "?load_balancer_id=" + loadBalancer.ID,
		status: http.StatusBadRequest,
	})

	doHTTPTest(t, &httpTest{
		name:   "delete assignment",
		method: http.MethodDelete,
		path:   baseURL + "?frontend_id=" + fe.ID + "&pool_id=" + pool.ID,
		status: http.StatusOK,
	})

	doHTTPTest(t, &httpTest{
		name:   "duplicate delete",
		method: http.MethodDelete,
		path:   baseURL + "?frontend_id=" + fe.ID + "&pool_id=" + pool.ID,
		status: http.StatusNotFound,
	})

	doHTTPTest(t, &httpTest{
		name:   "delete assignment2",
		method: http.MethodDelete,
		path:   baseURL + "?frontend_id=" + fe.ID + "&pool_id=" + pool2.ID,
		status: http.StatusOK,
	})
}
