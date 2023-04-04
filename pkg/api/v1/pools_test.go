package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"go.infratographer.com/load-balancer-api/internal/httptools"
	"go.infratographer.com/load-balancer-api/internal/pubsub"
	"go.infratographer.com/x/pubsubx"
)

const (
	natsMsgSubTimeout             = 2 * time.Second
	loadBalancerPoolSubjectCreate = "com.infratographer.events.load-balancer-pool.create.global"
	loadBalancerPoolSubjectDelete = "com.infratographer.events.load-balancer-pool.delete.global"
	loadBalancerPoolBaseUrn       = "urn:infratographer:load-balancer-pool:"
)

func TestCreatePool(t *testing.T) {
	nsrv := newNatsTestServer(t, "load-balancer-api-test", "com.infratographer.events.>")
	defer nsrv.Shutdown()

	srv := newTestServer(t, nsrv.ClientURL())
	defer srv.Close()

	// create a pubsub client for subscribing to NATS events
	subscriber := newPubSubClient(t, nsrv.ClientURL())
	msgChan := make(chan *nats.Msg, 10)

	// create a new nats subscription on the server created above
	subscription, err := subscriber.ChanSubscribe(
		context.TODO(),
		"com.infratographer.events.load-balancer-pool.>",
		msgChan,
		"load-balancer-api-test",
	)

	assert.NoError(t, err)

	defer func() {
		if err := subscription.Unsubscribe(); err != nil {
			t.Error(err)
		}
	}()

	tenantID := uuid.New().String()

	baseURL := srv.URL + "/v1/tenant/" + tenantID + "/pools"

	tests := []struct {
		name       string
		fakeBody   string
		tenant     string
		wantStatus int
	}{
		{
			name: "happy path no origins",
			fakeBody: `{
				"name": "testOrigin01",
				"protocol": "tcp"
			}`,
			tenant:     tenantID,
			wantStatus: http.StatusOK,
		},
		{
			name: "happy path with origins",
			fakeBody: `{
				"name": "testOrigin01",
				"protocol": "tcp",
				"origins": [
					{
						"disabled": false,
						"name": "testOrigin01",
						"target": "1.2.3.4",
						"port": 8443
					},
					{
						"disabled": false,
						"name": "testOrigin02",
						"target": "1.2.3.5",
						"port": 8443
					}
				]
			}`,
			tenant:     tenantID,
			wantStatus: http.StatusOK,
		},
		{
			name: "sad path conflicting origin targets",
			fakeBody: `{
				"name": "testOrigin01",
				"protocol": "tcp",
				"origins": [
					{
						"disabled": false,
						"name": "testOrigin01",
						"target": "1.2.3.4",
						"port": 8443
					},
					{
						"disabled": false,
						"name": "testOrigin02",
						"target": "1.2.3.4",
						"port": 8443
					}
				]
			}`,
			tenant:     tenantID,
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			createReq, err := http.NewRequestWithContext(
				context.TODO(),
				http.MethodPost,
				baseURL,
				httptools.FakeBody(tt.fakeBody),
			)
			assert.NoError(t, err)

			createResp, err := http.DefaultClient.Do(createReq)
			assert.NoError(t, err)
			assert.Equal(t, tt.wantStatus, createResp.StatusCode)
			defer createResp.Body.Close()

			// if we're testing a non-2xx resp, stop testing here
			if tt.wantStatus > 299 {
				return
			}

			testPool := struct {
				Version string `json:"version"`
				Message string `json:"message"`
				Status  int    `json:"status"`
				PoolID  string `json:"load_balancer_id"`
			}{}

			err = json.NewDecoder(createResp.Body).Decode(&testPool)

			assert.NoError(t, err)

			select {
			case msg := <-msgChan:
				pMsg := &pubsubx.Message{}
				err = json.Unmarshal(msg.Data, pMsg)
				assert.NoError(t, err)

				assert.Equal(t, loadBalancerPoolSubjectCreate, msg.Subject)
				assert.Equal(t, someTestJWTURN, pMsg.ActorURN)
				assert.Equal(t, pubsub.CreateEventType, pMsg.EventType)
			case <-time.After(natsMsgSubTimeout):
				t.Error("failed to receive nats message for delete")
			}

			deleteRequest, err := http.NewRequestWithContext(
				context.TODO(),
				http.MethodDelete,
				fmt.Sprintf("%s?pool_id=%s", baseURL, testPool.PoolID),
				nil,
			)

			assert.NoError(t, err)

			deleteResp, err := http.DefaultClient.Do(deleteRequest)
			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, deleteResp.StatusCode)
			defer deleteResp.Body.Close()

			select {
			case msg := <-msgChan:
				pMsg := &pubsubx.Message{}
				err = json.Unmarshal(msg.Data, pMsg)
				assert.NoError(t, err)

				assert.Equal(t, loadBalancerPoolSubjectDelete, msg.Subject)
				assert.Equal(t, someTestJWTURN, pMsg.ActorURN)
				assert.Equal(t, pubsub.DeleteEventType, pMsg.EventType)
			case <-time.After(natsMsgSubTimeout):
				t.Error("failed to receive nats message for delete")
			}
		})
	}
}

// createPool creates a pool with the given display name and protocol.
func createPool(t *testing.T, srv *httptest.Server, name string, tenantID string) (*pool, func(t *testing.T)) {
	t.Helper()

	body := `{"name": "` + name + `", "protocol": "tcp"}`

	baseURL := srv.URL + "/v1/pools"
	baseURLTenant := srv.URL + "/v1/tenant/" + tenantID + "/pools"

	doHTTPTest(t, &httpTest{
		name:   "create pool",
		method: http.MethodPost,
		path:   baseURLTenant,
		body:   body,
		status: http.StatusOK,
	})

	pool := response{}

	req, err := http.NewRequest(http.MethodGet, baseURLTenant+"?slug="+slug.Make(name), nil) //nolint
	assert.NoError(t, err)

	res, err := http.DefaultClient.Do(req) //nolint
	assert.Equal(t, http.StatusOK, res.StatusCode)

	bytes, err := io.ReadAll(res.Body)
	assert.NoError(t, err)

	err = json.Unmarshal(bytes, &pool)
	assert.NoError(t, err)

	res.Body.Close()

	return (*pool.Pools)[0], func(t *testing.T) {
		t.Helper()

		req, err := http.NewRequest(http.MethodDelete, baseURL+"/"+(*pool.Pools)[0].ID, nil) //nolint
		assert.NoError(t, err)

		res, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, res.StatusCode)
		res.Body.Close()
	}
}

func TestPoolRoutes(t *testing.T) {
	nsrv := newNatsTestServer(t, "load-balancer-api-test", "com.infratographer.events.>")
	defer nsrv.Shutdown()

	srv := newTestServer(t, nsrv.ClientURL())
	defer srv.Close()

	// create a pubsub client for subscribing to NATS events
	subscriber := newPubSubClient(t, nsrv.ClientURL())
	msgChan := make(chan *nats.Msg, 10)

	// create a new nats subscription on the server created above
	subscription, err := subscriber.ChanSubscribe(
		context.TODO(),
		"com.infratographer.events.load-balancer-pool.>",
		msgChan,
		"load-balancer-api-test",
	)

	assert.NoError(t, err)

	defer func() {
		if err := subscription.Unsubscribe(); err != nil {
			t.Error(err)
		}
	}()

	tenantID := uuid.New().String()
	baseURL := srv.URL + "/v1/pools"
	baseURLTenant := srv.URL + "/v1/tenant/" + tenantID + "/pools"
	missingUUID := uuid.New().String()

	// doHTTPTest is a helper function that makes a request to the server and
	// checks the response.
	//
	// To ensure test output has meaningful line references the function is
	// called individually for each test case
	doHTTPTest(t, &httpTest{
		name:   "get pools before create",
		status: http.StatusOK,
		path:   baseURLTenant,
		method: http.MethodGet,
	})

	// POST
	doHTTPTest(t, &httpTest{
		name:   "happy path",
		body:   `{"name": "Nemo", "protocol": "tcp"}`,
		status: http.StatusOK,
		path:   baseURLTenant,
		method: http.MethodPost,
	})

	select {
	case msg := <-msgChan:
		pMsg := &pubsubx.Message{}
		err = json.Unmarshal(msg.Data, pMsg)
		assert.NoError(t, err)

		assert.Equal(t, loadBalancerPoolSubjectCreate, msg.Subject)
		assert.Equal(t, someTestJWTURN, pMsg.ActorURN)
		assert.Equal(t, pubsub.CreateEventType, pMsg.EventType)
	case <-time.After(natsMsgSubTimeout):
		t.Error("failed to receive nats message")
	}

	doHTTPTest(t, &httpTest{
		name:   "duplicate",
		body:   `{"name": "Nemo", "protocol": "tcp"}`,
		status: http.StatusInternalServerError,
		path:   baseURLTenant,
		method: http.MethodPost,
	})

	doHTTPTest(t, &httpTest{
		name:   "multiple pools",
		body:   `[{"name": "Nemo", "protocol": "tcp"},{"name": "Dory", "protocol": "tcp"}]`,
		status: http.StatusBadRequest,
		path:   baseURLTenant,
		method: http.MethodPost,
	})

	doHTTPTest(t, &httpTest{
		name:   "missing display name",
		body:   `{"protocol": "tcp"}`,
		status: http.StatusBadRequest,
		path:   baseURLTenant,
		method: http.MethodPost,
	})

	doHTTPTest(t, &httpTest{
		name:   "missing protocol",
		body:   `{"name": "Bruce"}`,
		status: http.StatusOK,
		path:   baseURLTenant,
		method: http.MethodPost,
	})

	select {
	case msg := <-msgChan:
		pMsg := &pubsubx.Message{}
		err = json.Unmarshal(msg.Data, pMsg)
		assert.NoError(t, err)

		assert.Equal(t, loadBalancerPoolSubjectCreate, msg.Subject)
		assert.Equal(t, someTestJWTURN, pMsg.ActorURN)
		assert.Equal(t, pubsub.CreateEventType, pMsg.EventType)
	case <-time.After(natsMsgSubTimeout):
		t.Error("failed to receive nats message")
	}

	doHTTPTest(t, &httpTest{
		name:   "invalid protocol",
		body:   `{"name": "Nemo", "protocol": "invalid"}`,
		status: http.StatusBadRequest,
		path:   baseURLTenant,
		method: http.MethodPost,
	})

	doHTTPTest(t, &httpTest{
		name:   "invalid body",
		body:   `invalid`,
		status: http.StatusBadRequest,
		path:   baseURLTenant,
		method: http.MethodPost,
	})

	// GET
	doHTTPTest(t, &httpTest{
		name:   "happy path",
		status: http.StatusOK,
		path:   baseURLTenant,
		method: http.MethodGet,
	})

	doHTTPTest(t, &httpTest{
		name:   "happy path with query",
		status: http.StatusOK,
		path:   baseURLTenant + "?name=Nemo",
		method: http.MethodGet,
		tenant: tenantID,
	})

	doHTTPTest(t, &httpTest{
		name:   "not found with query",
		status: http.StatusOK,
		path:   baseURLTenant + "?slug=NotNemo",
		method: http.MethodGet,
		tenant: tenantID,
	})

	doHTTPTest(t, &httpTest{
		name:   "get pools without tenant or pool id",
		status: http.StatusNotFound,
		path:   baseURL,
		method: http.MethodGet,
	})

	doHTTPTest(t, &httpTest{
		name:   "list pools with invalid tenant id",
		path:   srv.URL + "/v1/tenant/123456/pools",
		status: http.StatusBadRequest,
		method: http.MethodGet,
	})

	doHTTPTest(t, &httpTest{
		name:   "list pools with unknown tenant id",
		path:   srv.URL + "/v1/tenant/" + missingUUID + "/pools",
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

	// DELETE
	doHTTPTest(t, &httpTest{
		name:   "happy path",
		status: http.StatusOK,
		path:   baseURLTenant + "?name=Nemo",
		method: http.MethodDelete,
	})

	select {
	case msg := <-msgChan:
		pMsg := &pubsubx.Message{}
		err = json.Unmarshal(msg.Data, pMsg)
		assert.NoError(t, err)

		assert.Equal(t, loadBalancerPoolSubjectDelete, msg.Subject)
		assert.Equal(t, someTestJWTURN, pMsg.ActorURN)
		assert.Equal(t, pubsub.DeleteEventType, pMsg.EventType)
	case <-time.After(natsMsgSubTimeout):
		t.Error("failed to receive nats message")
	}

	// create a test pool, but don't cleanup since it will get deleted below by id
	testPool, _ := createPool(t, srv, "Anchor", tenantID)

	select {
	case msg := <-msgChan:
		pMsg := &pubsubx.Message{}
		err = json.Unmarshal(msg.Data, pMsg)
		assert.NoError(t, err)

		assert.Equal(t, loadBalancerPoolSubjectCreate, msg.Subject)
		assert.Equal(t, someTestJWTURN, pMsg.ActorURN)
		assert.Equal(t, pubsub.CreateEventType, pMsg.EventType)
		assert.Equal(t, loadBalancerPoolBaseUrn+testPool.ID, pMsg.SubjectURN)
	case <-time.After(natsMsgSubTimeout):
		t.Error("failed to receive nats message")
	}

	doHTTPTest(t, &httpTest{
		name:   "happy path test pool id",
		status: http.StatusOK,
		path:   baseURL + "/" + testPool.ID,
		method: http.MethodDelete,
	})

	select {
	case msg := <-msgChan:
		pMsg := &pubsubx.Message{}
		err = json.Unmarshal(msg.Data, pMsg)
		assert.NoError(t, err)

		assert.Equal(t, loadBalancerPoolSubjectDelete, msg.Subject)
		assert.Equal(t, someTestJWTURN, pMsg.ActorURN)
		assert.Equal(t, pubsub.DeleteEventType, pMsg.EventType)
		assert.Equal(t, loadBalancerPoolBaseUrn+testPool.ID, pMsg.SubjectURN)
	case <-time.After(natsMsgSubTimeout):
		t.Error("failed to receive nats message")
	}
}

func TestPoolsGet(t *testing.T) {
	nsrv := newNatsTestServer(t, "load-balancer-api-test", "com.infratographer.events.>")
	defer nsrv.Shutdown()

	srv := newTestServer(t, nsrv.ClientURL())
	defer srv.Close()

	assert.NotNil(t, srv)

	// Create a pool
	pool, cleanupPool := createPool(t, srv, "marlin", uuid.NewString())
	defer cleanupPool(t)

	baseURL := srv.URL + "/v1/tenant/" + pool.TenantID + "/pools"

	// Get the pool
	doHTTPTest(t, &httpTest{
		name:   "get pool by id",
		method: http.MethodGet,
		path:   srv.URL + "/v1/pools/" + pool.ID,
		status: http.StatusOK,
		tenant: pool.TenantID,
	})

	doHTTPTest(t, &httpTest{
		name:   "get pool by query param id",
		method: http.MethodGet,
		path:   baseURL + "?pool_id=" + pool.ID,
		status: http.StatusOK,
		tenant: pool.TenantID,
	})

	// Get an unknown pool
	doHTTPTest(t, &httpTest{
		name:   "pool not found",
		method: http.MethodGet,
		path:   srv.URL + "/v1/pools/bfad65a9-abe3-44af-82ce-64331c84b2ad",
		status: http.StatusNotFound,
		tenant: pool.TenantID,
	})

	// Test pool list response
	listReq, err := http.NewRequestWithContext(
		context.TODO(),
		http.MethodGet,
		baseURL,
		nil,
	)
	assert.NoError(t, err)

	listResp, err := http.DefaultClient.Do(listReq)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, listResp.StatusCode)

	defer listResp.Body.Close()

	ca := pool.CreatedAt.Format(time.RFC3339Nano)
	ua := pool.UpdatedAt.Format(time.RFC3339Nano)
	testPoolsListExpected := fmt.Sprintf(`{"version":"v1","kind":"poolsList","pools":[{"created_at":"%s","updated_at":"%s","id":"%s","tenant_id":"%s","name":"%s","protocol":"%s","origins":[]}]}`+"\n", ca, ua, pool.ID, pool.TenantID, pool.Name, pool.Protocol)
	testPoolsList, err := io.ReadAll(listResp.Body)
	assert.NoError(t, err)
	assert.Equal(t, testPoolsListExpected, string(testPoolsList))

	// Test pool get by id from list endpoint response
	getListReq, err := http.NewRequestWithContext(
		context.TODO(),
		http.MethodGet,
		baseURL+"?pool_id="+pool.ID,
		nil,
	)
	assert.NoError(t, err)

	getListResp, err := http.DefaultClient.Do(getListReq)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, getListResp.StatusCode)

	defer getListResp.Body.Close()

	testPoolsGetListExpected := fmt.Sprintf(`{"version":"v1","kind":"poolsList","pools":[{"created_at":"%s","updated_at":"%s","id":"%s","tenant_id":"%s","name":"%s","protocol":"%s","origins":[]}]}`+"\n", ca, ua, pool.ID, pool.TenantID, pool.Name, pool.Protocol)
	testPoolsGetList, err := io.ReadAll(getListResp.Body)
	assert.NoError(t, err)
	assert.Equal(t, testPoolsGetListExpected, string(testPoolsGetList))

	// Test origin get by id from top level response
	getReq, err := http.NewRequestWithContext(
		context.TODO(),
		http.MethodGet,
		srv.URL+"/v1/pools/"+pool.ID,
		nil,
	)
	assert.NoError(t, err)

	getResp, err := http.DefaultClient.Do(getReq)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, getResp.StatusCode)

	defer getResp.Body.Close()

	testPoolsGetExpected := fmt.Sprintf(`{"version":"v1","kind":"poolsList","pool":{"created_at":"%s","updated_at":"%s","id":"%s","tenant_id":"%s","name":"%s","protocol":"%s","origins":[]}}`+"\n", ca, ua, pool.ID, pool.TenantID, pool.Name, pool.Protocol)
	testPoolsGet, err := io.ReadAll(getResp.Body)
	assert.NoError(t, err)
	assert.Equal(t, testPoolsGetExpected, string(testPoolsGet))
}
