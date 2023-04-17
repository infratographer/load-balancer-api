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
	nats "github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"go.infratographer.com/load-balancer-api/internal/httptools"
	"go.infratographer.com/load-balancer-api/internal/pubsub"
	"go.infratographer.com/x/pubsubx"
)

const (
	loadBalancerPortSubjectCreate = "com.infratographer.events.load-balancer-port.create.global"
	loadBalancerPortSubjectDelete = "com.infratographer.events.load-balancer-port.delete.global"
	loadBalancerPortBaseUrn       = "urn:infratographer:load-balancer-port:"
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

func TestCreatePorts(t *testing.T) {
	nsrv := newNatsTestServer(t, "load-balancer-api-test", "com.infratographer.events.>")
	defer nsrv.Shutdown()

	srv := newTestServer(t, nsrv.ClientURL())
	defer srv.Close()

	testLb, cleanupLb := createLoadBalancer(t, srv, uuid.New().String())
	defer cleanupLb(t)

	// create a pool for testing
	testPool, cleanupPool := createPool(t, srv, "testPool01", testLb.TenantID)
	defer cleanupPool(t)

	// create a pubsub client for subscribing to NATS events
	subscriber := newPubSubClient(t, nsrv.ClientURL())
	msgChan := make(chan *nats.Msg, 10)

	// create a new nats subscription on the server created above
	subscription, err := subscriber.ChanSubscribe(
		context.TODO(),
		"com.infratographer.events.load-balancer-port.>",
		msgChan,
		"load-balancer-api-test",
	)

	assert.NoError(t, err)

	defer func() {
		if err := subscription.Unsubscribe(); err != nil {
			t.Error(err)
		}
	}()

	baseURL := srv.URL + "/v1/loadbalancers/" + testLb.ID + "/ports"

	tests := []struct {
		name       string
		fakeBody   string
		tenant     string
		wantStatus int
	}{
		{
			name: "happy path no pools",
			fakeBody: `{
					"name": "testport01",
					"port": 443
				}`,
			wantStatus: http.StatusOK,
		},
		{
			name: "happy path with pools",
			fakeBody: fmt.Sprintf(`{
					"name": "testport01",
					"port": 443,
					"pools": ["%s"]
				}`, testPool.ID),
			wantStatus: http.StatusOK,
		},
		{
			name: "sad path nonexistent pools",
			fakeBody: fmt.Sprintf(`{
					"name": "testport01",
					"port": 443,
					"pools": ["%s"]
				}`, uuid.New().String()),
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

			testPort := struct {
				Version string `json:"version"`
				Message string `json:"message"`
				Status  int    `json:"status"`
				PortID  string `json:"port_id"`
			}{}

			err = json.NewDecoder(createResp.Body).Decode(&testPort)

			assert.NoError(t, err)

			select {
			case msg := <-msgChan:
				pMsg := &pubsubx.Message{}
				err = json.Unmarshal(msg.Data, pMsg)
				assert.NoError(t, err)

				assert.Equal(t, loadBalancerPortSubjectCreate, msg.Subject)
				assert.Equal(t, someTestJWTURN, pMsg.ActorURN)
				assert.Equal(t, pubsub.CreateEventType, pMsg.EventType)
			case <-time.After(natsMsgSubTimeout):
				t.Error("failed to receive nats message for create")
			}

			deleteRequest, err := http.NewRequestWithContext(
				context.TODO(),
				http.MethodDelete,
				fmt.Sprintf("%s?port_id=%s", baseURL, testPort.PortID),
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

				assert.Equal(t, loadBalancerPortSubjectDelete, msg.Subject)
				assert.Equal(t, someTestJWTURN, pMsg.ActorURN)
				assert.Equal(t, pubsub.DeleteEventType, pMsg.EventType)
			case <-time.After(natsMsgSubTimeout):
				t.Error("failed to receive nats message for delete")
			}
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
		name:   "negative port",
		body:   `{"name": "Mouth", "port": -1}`,
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

	// PATCH tests
	doHTTPTest(t, &httpTest{
		name:   "happy path patch port name",
		body:   `{"name": "LeftEars"}`,
		status: http.StatusAccepted,
		method: http.MethodPatch,
		path:   baseURL + "/" + earsID,
	})

	doHTTPTest(t, &httpTest{
		name:   "patch port wrong type",
		body:   `{"port": "foobar"}`,
		status: http.StatusBadRequest,
		method: http.MethodPatch,
		path:   baseURL + "/" + earsID,
	})

	doHTTPTest(t, &httpTest{
		name:   "patch port port too low",
		body:   `{"port": -1}`,
		status: http.StatusBadRequest,
		method: http.MethodPatch,
		path:   baseURL + "/" + earsID,
	})

	doHTTPTest(t, &httpTest{
		name:   "patch port port two high",
		body:   `{"port": 131337}`,
		status: http.StatusBadRequest,
		method: http.MethodPatch,
		path:   baseURL + "/" + earsID,
	})

	doHTTPTest(t, &httpTest{
		name:   "patch port not found",
		body:   `{"name": "Plain Mouth", "port": 80}`,
		status: http.StatusAccepted,
		method: http.MethodPatch,
		path:   baseURLLoadBalancer + "?port=80",
	})

	// Get Tests
	doHTTPTest(t, &httpTest{
		name:   "happy path get by id",
		path:   baseURL + "/" + earsID,
		status: http.StatusOK,
		method: http.MethodGet,
	})

	doHTTPTest(t, &httpTest{
		name:   "happy path list with loadbalancer id",
		path:   baseURLLoadBalancer,
		status: http.StatusOK,
		method: http.MethodGet,
	})

	doHTTPTest(t, &httpTest{
		name:   "list ports by slug",
		path:   baseURLLoadBalancer + "?slug=ears",
		status: http.StatusOK,
		method: http.MethodGet,
	})

	doHTTPTest(t, &httpTest{
		name:   "list ports by name",
		path:   baseURLLoadBalancer + "?name=ears",
		status: http.StatusOK,
		method: http.MethodGet,
	})

	doHTTPTest(t, &httpTest{
		name:   "list ports by name",
		path:   baseURLLoadBalancer + "?port=25",
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

func TestPortssGet(t *testing.T) {
	nsrv := newNatsTestServer(t, "load-balancer-api-test", "com.infratographer.events.>")
	defer nsrv.Shutdown()

	srv := newTestServer(t, nsrv.ClientURL())
	defer srv.Close()

	assert.NotNil(t, srv)

	// Create a load balancer to use for testing
	lb, cleanupLB := createLoadBalancer(t, srv, uuid.NewString())
	defer cleanupLB(t)

	// TODO create test port with pools

	// Create a port to use for testing
	port, cleanupPort := createPort(t, srv, lb.ID)
	defer cleanupPort(t)

	baseURL := srv.URL + "/v1/loadbalancers/" + lb.ID + "/ports"

	doHTTPTest(t, &httpTest{
		name:   "get port by id",
		method: http.MethodGet,
		path:   baseURL + "?port_id=" + port.ID,
		status: http.StatusOK,
	})

	doHTTPTest(t, &httpTest{
		name:   "port not found",
		method: http.MethodGet,
		path:   baseURL + "/bfad65a9-abe3-44af-82ce-64331c84b2ad",
		status: http.StatusNotFound,
	})

	// Test port list response
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

	ca := port.CreatedAt.Format(time.RFC3339Nano)
	ua := port.UpdatedAt.Format(time.RFC3339Nano)
	testPortsListExpected := fmt.Sprintf(`{"version":"v1","kind":"portsList","ports":[{"created_at":"%s","updated_at":"%s","id":"%s","tenant_id":"","load_balancer_id":"%s","name":"%s","address_family":"%s","port":%d,"pools":[]}]}`+"\n", ca, ua, port.ID, lb.ID, port.Name, port.AddressFamily, port.Port)
	testPortsList, err := io.ReadAll(listResp.Body)
	assert.NoError(t, err)
	assert.Equal(t, testPortsListExpected, string(testPortsList))

	// Test ports get by id from list endpoint response
	getListReq, err := http.NewRequestWithContext(
		context.TODO(),
		http.MethodGet,
		baseURL+"?port_id="+port.ID,
		nil,
	)
	assert.NoError(t, err)

	getListResp, err := http.DefaultClient.Do(getListReq)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, getListResp.StatusCode)

	defer getListResp.Body.Close()

	testPortsGetListExpected := fmt.Sprintf(`{"version":"v1","kind":"portsList","ports":[{"created_at":"%s","updated_at":"%s","id":"%s","tenant_id":"","load_balancer_id":"%s","name":"%s","address_family":"%s","port":%d,"pools":[]}]}`+"\n", ca, ua, port.ID, lb.ID, port.Name, port.AddressFamily, port.Port)
	testPortsGetList, err := io.ReadAll(getListResp.Body)
	assert.NoError(t, err)
	assert.Equal(t, testPortsGetListExpected, string(testPortsGetList))

	// Test ports get by id from top level response
	getReq, err := http.NewRequestWithContext(
		context.TODO(),
		http.MethodGet,
		srv.URL+"/v1/ports/"+port.ID,
		nil,
	)
	assert.NoError(t, err)

	getResp, err := http.DefaultClient.Do(getReq)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, getResp.StatusCode)

	defer getResp.Body.Close()

	testPortsGetExpected := fmt.Sprintf(`{"version":"v1","kind":"portsGet","port":{"created_at":"%s","updated_at":"%s","id":"%s","tenant_id":"","load_balancer_id":"%s","name":"%s","address_family":"%s","port":%d,"pools":[]}}`+"\n", ca, ua, port.ID, lb.ID, port.Name, port.AddressFamily, port.Port)
	testPortsGet, err := io.ReadAll(getResp.Body)
	assert.NoError(t, err)
	assert.Equal(t, testPortsGetExpected, string(testPortsGet))
}
