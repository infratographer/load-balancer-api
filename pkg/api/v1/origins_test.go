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
	nats "github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"go.infratographer.com/load-balancer-api/internal/httptools"
	"go.infratographer.com/load-balancer-api/internal/pubsub"
	"go.infratographer.com/x/pubsubx"
)

const (
	loadBalancerOriginSubjectCreate = "com.infratographer.events.load-balancer-origin.create.global"
	loadBalancerOriginSubjectDelete = "com.infratographer.events.load-balancer-origin.delete.global"
	loadBalancerOriginBaseUrn       = "urn:infratographer:load-balancer-origin:"
)

func createOrigin(t *testing.T, srv *httptest.Server, name string, poolID string) (*origin, func(*testing.T)) {
	t.Helper()

	body := fmt.Sprintf(`{"disabled": true,"name": "%s", "target": "1.1.1.1", "port": 80}`, name)

	doHTTPTest(t, &httpTest{
		name:   "create origin",
		body:   body,
		status: http.StatusOK,
		path:   srv.URL + "/v1/pools/" + poolID + "/origins",
		method: http.MethodPost,
	})

	// Get the origin
	req, err := http.NewRequest(http.MethodGet, srv.URL+"/v1/pools/"+poolID+"/origins?slug="+slug.Make(name), nil) //nolint
	assert.NoError(t, err)

	res, err := http.DefaultClient.Do(req) //nolint
	assert.NoError(t, err)

	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)

	origin := response{}
	err = json.NewDecoder(res.Body).Decode(&origin)
	assert.NoError(t, err)

	return (*origin.Origins)[0], func(t *testing.T) {
		t.Helper()

		// Delete the origin
		req, err := http.NewRequest(http.MethodDelete, srv.URL+"/v1/pools/"+poolID+"/origins?slug="+slug.Make(name), nil) //nolint
		assert.NoError(t, err)

		res, err := http.DefaultClient.Do(req) //nolint
		assert.NoError(t, err)

		defer res.Body.Close()
	}
}

func TestCreateOrigins(t *testing.T) {
	nsrv := newNatsTestServer(t, "load-balancer-api-test", "com.infratographer.events.>")
	defer nsrv.Shutdown()

	srv := newTestServer(t, nsrv.ClientURL())
	defer srv.Close()

	// create a pool for testing
	testPool, cleanupPool := createPool(t, srv, "testPool01", uuid.New().String())
	defer cleanupPool(t)

	// create a pubsub client for subscribing to NATS events
	subscriber := newPubSubClient(t, nsrv.ClientURL())
	msgChan := make(chan *nats.Msg, 10)

	// create a new nats subscription on the server created above
	subscription, err := subscriber.ChanSubscribe(
		context.TODO(),
		"com.infratographer.events.load-balancer-origin.>",
		msgChan,
		"load-balancer-api-test",
	)

	assert.NoError(t, err)

	defer func() {
		if err := subscription.Unsubscribe(); err != nil {
			t.Error(err)
		}
	}()

	baseURL := srv.URL + "/v1/pools/" + testPool.ID + "/origins"

	tests := []struct {
		name       string
		fakeBody   string
		tenant     string
		wantStatus int
	}{
		{
			name: "happy path",
			fakeBody: `{
					"disabled": false,
					"name": "testorigin01",
					"target": "1.1.1.1",
					"port": 31337
				}`,
			wantStatus: http.StatusOK,
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

			testOrigin := struct {
				Version  string `json:"version"`
				Message  string `json:"message"`
				Status   int    `json:"status"`
				OriginID string `json:"origin_id"`
			}{}

			err = json.NewDecoder(createResp.Body).Decode(&testOrigin)

			assert.NoError(t, err)

			select {
			case msg := <-msgChan:
				pMsg := &pubsubx.Message{}
				err = json.Unmarshal(msg.Data, pMsg)
				assert.NoError(t, err)

				assert.Equal(t, loadBalancerOriginSubjectCreate, msg.Subject)
				assert.Equal(t, someTestJWTURN, pMsg.ActorURN)
				assert.Equal(t, pubsub.CreateEventType, pMsg.EventType)
			case <-time.After(natsMsgSubTimeout):
				t.Error("failed to receive nats message for create")
			}

			deleteRequest, err := http.NewRequestWithContext(
				context.TODO(),
				http.MethodDelete,
				fmt.Sprintf("%s?origin_id=%s", baseURL, testOrigin.OriginID),
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

				assert.Equal(t, loadBalancerOriginSubjectDelete, msg.Subject)
				assert.Equal(t, someTestJWTURN, pMsg.ActorURN)
				assert.Equal(t, pubsub.DeleteEventType, pMsg.EventType)
			case <-time.After(natsMsgSubTimeout):
				t.Error("failed to receive nats message for delete")
			}
		})
	}
}

func TestOriginRoutes(t *testing.T) {
	nsrv := newNatsTestServer(t, "load-balancer-api-test", "com.infratographer.events.>")
	defer nsrv.Shutdown()

	srv := newTestServer(t, nsrv.ClientURL())
	defer srv.Close()

	tenantID := uuid.New().String()
	pool, remove := createPool(t, srv, "squirt", tenantID)

	origin, cleanup := createOrigin(t, srv, "testorigin01", pool.ID)
	defer cleanup(t)

	baseURL := srv.URL + "/v1/origins"
	baseURLPool := srv.URL + "/v1/pools/" + pool.ID + "/origins"
	missingUUID := uuid.New().String()

	// doHTTPTest is a helper function that makes a request to the server and
	// checks the response.
	//
	// To ensure test output has meaningful line references the function is
	// called individually for each test case
	doHTTPTest(t, &httpTest{
		name:   "list origins before created",
		status: http.StatusOK,
		path:   baseURLPool,
		method: http.MethodGet,
	})

	// POST
	doHTTPTest(t, &httpTest{
		name:   "happy path",
		body:   `{"disabled": true,"name": "The Butt", "target": "9.9.9.9", "port": 80}`,
		status: http.StatusOK,
		path:   baseURLPool,
		method: http.MethodPost,
	})

	doHTTPTest(t, &httpTest{
		name:   "happy path",
		body:   `{"disabled": true,"name": "Fish are friends", "target": "9.9.8.8", "port": 80}`,
		status: http.StatusOK,
		path:   baseURLPool,
		method: http.MethodPost,
	})

	doHTTPTest(t, &httpTest{
		name:   "list of origins",
		body:   `[{"disabled": true,"name": "The Butt", "target": "9.9.9.9", "port": 80},{"disabled": true,"name": "The Beard", "target": "9.9.9.10", "port": 80}]`,
		status: http.StatusBadRequest,
		path:   baseURLPool,
		method: http.MethodPost,
	})

	doHTTPTest(t, &httpTest{
		name:   "no pool",
		body:   `[]`,
		status: http.StatusBadRequest,
		path:   baseURLPool,
		method: http.MethodPost,
	})

	doHTTPTest(t, &httpTest{
		name:   "bad body",
		body:   `bad body`,
		status: http.StatusBadRequest,
		path:   baseURLPool,
		method: http.MethodPost,
	})

	doHTTPTest(t, &httpTest{
		name:   "missing pool_id in path",
		body:   `{"disabled": true,"name": "the-butt", "target": "2.0.0.1", "port": 80}`,
		status: http.StatusNotFound,
		path:   baseURL,
		method: http.MethodPost,
	})

	// PUT
	doHTTPTest(t, &httpTest{
		name:   "happy path update origin by id",
		body:   `{"disabled": true,"name": "testorigin01", "target": "2.2.2.2", "port": 8080}`,
		status: http.StatusAccepted,
		path:   baseURLPool + "?origin_id=" + origin.ID,
		method: http.MethodPut,
	})

	doHTTPTest(t, &httpTest{
		name:   "happy path update origin by target and port",
		body:   `{"disabled": false,"name": "testorigin01", "target": "1.1.1.1", "port": 80}`,
		status: http.StatusAccepted,
		path:   baseURLPool + "?origin_target=2.2.2.2&port=8080",
		method: http.MethodPut,
	})

	doHTTPTest(t, &httpTest{
		name:   "sad path update origin missing target",
		body:   `{"disabled": true,"name": "testorigin01", "port": 80}`,
		status: http.StatusBadRequest,
		path:   baseURLPool + "?origin_id=" + origin.ID,
		method: http.MethodPut,
	})

	doHTTPTest(t, &httpTest{
		name:   "sad path update origin missing port",
		body:   `{"disabled": true,"name": "testorigin01", "target": "1.1.1.1"}`,
		status: http.StatusBadRequest,
		path:   baseURLPool + "?origin_id=" + origin.ID,
		method: http.MethodPut,
	})

	doHTTPTest(t, &httpTest{
		name:   "sad path update origin bad port",
		body:   `{"disabled": true,"name": "testorigin01", "target": "1.1.1.1", "port":-1}`,
		status: http.StatusBadRequest,
		path:   baseURLPool + "?origin_id=" + origin.ID,
		method: http.MethodPut,
	})

	// PATCH
	doHTTPTest(t, &httpTest{
		name:   "happy path patch origin name by id",
		body:   `{"name": "testorigin02"}`,
		status: http.StatusAccepted,
		path:   baseURLPool + "?origin_id=" + origin.ID,
		method: http.MethodPatch,
	})

	doHTTPTest(t, &httpTest{
		name:   "happy path patch origin target by id",
		body:   `{"target": "2.2.2.2"}`,
		status: http.StatusAccepted,
		path:   baseURLPool + "?origin_id=" + origin.ID,
		method: http.MethodPatch,
	})

	doHTTPTest(t, &httpTest{
		name:   "happy path patch origin port by id",
		body:   `{"port": 8080}`,
		status: http.StatusAccepted,
		path:   baseURLPool + "?origin_id=" + origin.ID,
		method: http.MethodPatch,
	})

	doHTTPTest(t, &httpTest{
		name:   "happy path patch origin name by target and port",
		body:   `{"name": "testorigin02"}`,
		status: http.StatusAccepted,
		path:   baseURLPool + "?origin_target=2.2.2.2&port=8080",
		method: http.MethodPatch,
	})

	doHTTPTest(t, &httpTest{
		name:   "happy path patch origin target by target and port",
		body:   `{"target": "1.1.1.1"}`,
		status: http.StatusAccepted,
		path:   baseURLPool + "?origin_target=2.2.2.2&port=8080",
		method: http.MethodPatch,
	})

	doHTTPTest(t, &httpTest{
		name:   "happy path patch origin port by target and port",
		body:   `{"port": 80}`,
		status: http.StatusAccepted,
		path:   baseURLPool + "?origin_target=1.1.1.1&port=8080",
		method: http.MethodPatch,
	})

	doHTTPTest(t, &httpTest{
		name:   "sad path update origin empty target",
		body:   `{"target": "", "port": 80}`,
		status: http.StatusBadRequest,
		path:   baseURLPool + "?origin_id=" + origin.ID,
		method: http.MethodPatch,
	})

	doHTTPTest(t, &httpTest{
		name:   "sad path update origin empty port",
		body:   `{"target": "1.1.1.1", "port": 0}`,
		status: http.StatusBadRequest,
		path:   baseURLPool + "?origin_id=" + origin.ID,
		method: http.MethodPatch,
	})

	// GET
	doHTTPTest(t, &httpTest{
		name:   "happy path",
		status: http.StatusOK,
		path:   baseURLPool,
		method: http.MethodGet,
	})

	doHTTPTest(t, &httpTest{
		name:   "touch the butt",
		status: http.StatusOK,
		path:   baseURLPool + "?slug=the-butt",
		method: http.MethodGet,
	})

	doHTTPTest(t, &httpTest{
		name:   "missing origin uuid",
		path:   baseURL + "/" + missingUUID,
		status: http.StatusNotFound,
		method: http.MethodGet,
	})

	doHTTPTest(t, &httpTest{
		name:   "bad origin uuid",
		path:   baseURL + "/123456",
		status: http.StatusBadRequest,
		method: http.MethodGet,
	})

	doHTTPTest(t, &httpTest{
		name:   "list origins with invalid pool id",
		path:   srv.URL + "/v1/pools/123456/origins",
		status: http.StatusBadRequest,
		method: http.MethodGet,
	})

	doHTTPTest(t, &httpTest{
		name:   "list origins with unknown pool id",
		path:   srv.URL + "/v1/pools/" + uuid.New().String() + "/origins",
		status: http.StatusOK,
		method: http.MethodGet,
	})

	// DELETE
	doHTTPTest(t, &httpTest{
		name:   "ambigous delete",
		status: http.StatusBadRequest,
		path:   baseURLPool + "?port=80",
		method: http.MethodDelete,
	})

	doHTTPTest(t, &httpTest{
		name:   "happy path",
		status: http.StatusOK,
		path:   baseURLPool + "?slug=the-butt",
		method: http.MethodDelete,
	})

	doHTTPTest(t, &httpTest{
		name:   "happy path 2",
		status: http.StatusOK,
		path:   baseURLPool + "?slug=fish-are-friends",
		method: http.MethodDelete,
	})

	doHTTPTest(t, &httpTest{
		name:   "404",
		status: http.StatusNotFound,
		path:   baseURL + "?slug=fish-are-friends",
		method: http.MethodDelete,
	})

	remove(t)
}

func TestOriginsBalancerGet(t *testing.T) {
	nsrv := newNatsTestServer(t, "load-balancer-api-test", "com.infratographer.events.>")
	defer nsrv.Shutdown()

	srv := newTestServer(t, nsrv.ClientURL())
	defer srv.Close()

	assert.NotNil(t, srv)

	// Create a load balancer to use for testing
	loadBalancer, cleanupLB := createLoadBalancer(t, srv, uuid.NewString())
	defer cleanupLB(t)

	// Create a pool in the same tenant to use for testing
	pool, cleanupPool := createPool(t, srv, "marlin", loadBalancer.TenantID)
	defer cleanupPool(t)

	// Create an origin (in the pool) to use for testing
	origin, cleanupOrigin := createOrigin(t, srv, "bruce", pool.ID)
	defer cleanupOrigin(t)

	baseURL := srv.URL + "/v1/pools/" + pool.ID + "/origins"

	doHTTPTest(t, &httpTest{
		name:   "get origin by id",
		method: http.MethodGet,
		path:   baseURL + "?origin_id=" + origin.ID,
		status: http.StatusOK,
	})

	doHTTPTest(t, &httpTest{
		name:   "get missing origin by query param id",
		method: http.MethodGet,
		path:   baseURL + "?origin_id=bfad65a9-abe3-44af-82ce-64331c84b2ad",
		status: http.StatusOK,
	})

	doHTTPTest(t, &httpTest{
		name:   "get missing origin by id",
		method: http.MethodGet,
		path:   srv.URL + "/v1/origins/bfad65a9-abe3-44af-82ce-64331c84b2ad",
		status: http.StatusNotFound,
	})

	// Test origins list response
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

	ca := origin.CreatedAt.Format(time.RFC3339Nano)
	ua := origin.UpdatedAt.Format(time.RFC3339Nano)
	testOriginsListExpected := fmt.Sprintf(`{"version":"v1","kind":"originsList","origins":[{"created_at":"%s","updated_at":"%s","id":"%s","name":"%s","port":%d,"origin_target":"%s","origin_disabled":%t}]}`+"\n", ca, ua, origin.ID, origin.Name, origin.Port, origin.OriginTarget, origin.OriginDisabled)
	testOriginsList, err := io.ReadAll(listResp.Body)
	assert.NoError(t, err)
	assert.Equal(t, testOriginsListExpected, string(testOriginsList))

	// Test origin get by id from list endpoint response
	getListReq, err := http.NewRequestWithContext(
		context.TODO(),
		http.MethodGet,
		baseURL+"?origin_id="+origin.ID,
		nil,
	)
	assert.NoError(t, err)

	getListResp, err := http.DefaultClient.Do(getListReq)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, getListResp.StatusCode)

	defer getListResp.Body.Close()

	testOriginsGetListExpected := fmt.Sprintf(`{"version":"v1","kind":"originsList","origins":[{"created_at":"%s","updated_at":"%s","id":"%s","name":"%s","port":%d,"origin_target":"%s","origin_disabled":%t}]}`+"\n", ca, ua, origin.ID, origin.Name, origin.Port, origin.OriginTarget, origin.OriginDisabled)
	testOriginsGetList, err := io.ReadAll(getListResp.Body)
	assert.NoError(t, err)
	assert.Equal(t, testOriginsGetListExpected, string(testOriginsGetList))

	// Test origin get by id from top level response
	getReq, err := http.NewRequestWithContext(
		context.TODO(),
		http.MethodGet,
		srv.URL+"/v1/origins/"+origin.ID,
		nil,
	)
	assert.NoError(t, err)

	getResp, err := http.DefaultClient.Do(getReq)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, getResp.StatusCode)

	defer getResp.Body.Close()

	testOriginsGetExpected := fmt.Sprintf(`{"version":"v1","kind":"originsGet","origin":{"created_at":"%s","updated_at":"%s","id":"%s","name":"%s","port":%d,"origin_target":"%s","origin_disabled":%t}}`+"\n", ca, ua, origin.ID, origin.Name, origin.Port, origin.OriginTarget, origin.OriginDisabled)
	testOriginsGet, err := io.ReadAll(getResp.Body)
	assert.NoError(t, err)
	assert.Equal(t, testOriginsGetExpected, string(testOriginsGet))
}
