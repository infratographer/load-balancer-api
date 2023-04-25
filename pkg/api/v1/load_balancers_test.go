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
	loadBalancerSubjectCreate = "com.infratographer.events.load-balancer.create.global"
	loadBalancerSubjectDelete = "com.infratographer.events.load-balancer.delete.global"
	loadBalancerBaseUrn       = "urn:infratographer:load-balancer:"
)

func TestCreateLoadBalancer(t *testing.T) {
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
		"com.infratographer.events.load-balancer.>",
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
	locationID := uuid.New().String()
	ipID := uuid.New().String()
	sizeName := "small"
	typeName := "layer-3"

	baseURL := srv.URL + "/v1/tenant/" + tenantID + "/loadbalancers"

	// create a pool for testing
	pool, cleanupPool := createPool(t, srv, "testPool01", tenantID)
	defer cleanupPool(t)

	tests := []struct {
		name       string
		fakeBody   string
		tenant     string
		lbType     string
		wantStatus int
	}{
		{
			name: "happy path no ports",
			fakeBody: fmt.Sprintf(`{
					"name": "testlb01",
					"location_id": "%s",
					"ip_address_id": "%s",
					"load_balancer_size": "%s",
					"load_balancer_type": "%s"
				}`, locationID, ipID, sizeName, typeName),
			tenant:     tenantID,
			wantStatus: http.StatusOK,
		},
		{
			name: "happy path ports no pools",
			fakeBody: fmt.Sprintf(`{
					"name": "testlb02",
					"location_id": "%s",
					"ip_address_id": "%s",
					"load_balancer_size": "%s",
					"load_balancer_type": "%s",
					"ports": [
						{
							"name": "testlb02-https",
							"port": 443
						}
					]
				}`, locationID, ipID, sizeName, typeName),
			tenant:     tenantID,
			wantStatus: http.StatusOK,
		},
		{
			name: "happy path ports with pool",
			fakeBody: fmt.Sprintf(`{
					"name": "testlb03",
					"location_id": "%s",
					"ip_address_id": "%s",
					"load_balancer_size": "%s",
					"load_balancer_type": "%s",
					"ports": [
						{
							"name": "testlb03-https",
							"port": 443,
							"pools": ["%s"]
						}
					]
				}`, locationID, ipID, sizeName, typeName, pool.ID),
			tenant:     tenantID,
			wantStatus: http.StatusOK,
		},
		{
			name: "sad path ports with nonexistent pool",
			fakeBody: fmt.Sprintf(`{
					"name": "testlb04",
					"location_id": "%s",
					"ip_address_id": "%s",
					"load_balancer_size": "%s",
					"load_balancer_type": "%s",
					"ports": [
						{
							"name": "testlb04-https",
							"port": 443,
							"pools": ["%s"]
						}
					]
				}`, locationID, ipID, sizeName, typeName, uuid.New().String()),
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

			testLoadBalancer := struct {
				Version        string `json:"version"`
				Message        string `json:"message"`
				Status         int    `json:"status"`
				LoadBalancerID string `json:"load_balancer_id"`
			}{}

			err = json.NewDecoder(createResp.Body).Decode(&testLoadBalancer)

			assert.NoError(t, err)

			select {
			case msg := <-msgChan:
				pMsg := &pubsubx.Message{}
				err = json.Unmarshal(msg.Data, pMsg)
				assert.NoError(t, err)

				assert.Equal(t, loadBalancerSubjectCreate, msg.Subject)
				assert.Equal(t, someTestJWTURN, pMsg.ActorURN)
				assert.Equal(t, pubsub.CreateEventType, pMsg.EventType)
			case <-time.After(natsMsgSubTimeout):
				t.Error("failed to receive nats message for delete")
			}

			deleteRequest, err := http.NewRequestWithContext(
				context.TODO(),
				http.MethodDelete,
				fmt.Sprintf("%s?load_balancer_id=%s", baseURL, testLoadBalancer.LoadBalancerID),
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

				assert.Equal(t, loadBalancerSubjectDelete, msg.Subject)
				assert.Equal(t, someTestJWTURN, pMsg.ActorURN)
				assert.Equal(t, pubsub.DeleteEventType, pMsg.EventType)
			case <-time.After(natsMsgSubTimeout):
				t.Error("failed to receive nats message for delete")
			}
		})
	}
}

func TestLoadBalancerRoutes(t *testing.T) {
	nsrv := newNatsTestServer(t, "load-balancer-api-test", "com.infratographer.events.>")
	defer nsrv.Shutdown()

	srv := newTestServer(t, nsrv.ClientURL())
	defer srv.Close()

	tenantID := uuid.New().String()
	baseURL := srv.URL + "/v1/loadbalancers"
	baseURLTenant := srv.URL + "/v1/tenant/" + tenantID + "/loadbalancers"
	locationID := uuid.New().String()
	missingUUID := uuid.New().String()
	testIPaddressUUIDBruce := "61b3625b-3c31-4c70-a42c-239bf2212ff1"

	// create a test load balancer named Bruce
	req1, err := http.NewRequestWithContext(
		context.TODO(),
		http.MethodPost,
		baseURLTenant,
		httptools.FakeBody(
			fmt.Sprintf(`{"name": "Bruce", "location_id": "%s", "ip_address_id": "%s","load_balancer_size": "small","load_balancer_type": "layer-3"}`, locationID, testIPaddressUUIDBruce),
		),
	)
	assert.NoError(t, err)
	resp1, err := http.DefaultClient.Do(req1)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp1.StatusCode)

	testLoadBalancer := struct {
		Version        string `json:"version"`
		Message        string `json:"message"`
		Status         int    `json:"status"`
		LoadBalancerID string `json:"load_balancer_id"`
	}{}

	_ = json.NewDecoder(resp1.Body).Decode(&testLoadBalancer)
	resp1.Body.Close()

	// cleanup test load balancer
	defer func(id string) {
		rq, err := http.NewRequestWithContext(context.TODO(), http.MethodDelete, baseURL+"/"+id, nil)
		assert.NoError(t, err)
		rs, err := http.DefaultClient.Do(rq)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rs.StatusCode)
		rs.Body.Close()
	}(testLoadBalancer.LoadBalancerID)

	doHTTPTest(t, &httpTest{
		name:   "list lbs before created",
		path:   baseURLTenant,
		status: http.StatusOK,
		method: http.MethodGet,
	})

	testIPaddressUUIDNemo := "5ff95301-07b1-4f7c-a4df-14b2003017ea"

	// POST tests
	doHTTPTest(t, &httpTest{
		name:   "happy path",
		body:   fmt.Sprintf(`{"name": "Nemo", "location_id": "%s", "ip_address_id": "%s","load_balancer_size": "small","load_balancer_type": "layer-3"}`, locationID, testIPaddressUUIDNemo),
		status: http.StatusOK,
		path:   baseURLTenant,
		method: http.MethodPost,
	})

	testIPaddressUUIDDory := "5ff95301-07b1-4f7c-a4df-14b2003017ea"

	doHTTPTest(t, &httpTest{
		name:   "happy path 2",
		body:   fmt.Sprintf(`{"name": "Dori", "location_id": "%s", "ip_address_id": "%s","load_balancer_size": "small","load_balancer_type": "layer-3"}`, locationID, testIPaddressUUIDDory),
		status: http.StatusOK,
		path:   baseURLTenant,
		method: http.MethodPost,
	})

	doHTTPTest(t, &httpTest{
		name:   "Duplicate",
		body:   fmt.Sprintf(`{"name": "Nemo", "location_id": "%s", "ip_address_id": "%s","load_balancer_size": "small","load_balancer_type": "layer-3"}`, locationID, uuid.NewString()),
		status: http.StatusInternalServerError,
		path:   baseURLTenant,
		method: http.MethodPost,
	})

	doHTTPTest(t, &httpTest{
		name:   "missing tenantID",
		body:   fmt.Sprintf(`{"name": "Nemo", "location_id": "%s", "ip_address_id": %s,"load_balancer_size": "small","load_balancer_type": "layer-3"}`, locationID, uuid.NewString()),
		status: http.StatusNotFound,
		path:   baseURL,
		method: http.MethodPost,
	})

	doHTTPTest(t, &httpTest{
		name:   "missing name",
		body:   fmt.Sprintf(`{"location_id": "%s", "ip_address_id": "%s","load_balancer_size": "small","load_balancer_type": "layer-3"}`, locationID, uuid.NewString()),
		status: http.StatusBadRequest,
		path:   baseURLTenant,
		method: http.MethodPost,
	})

	doHTTPTest(t, &httpTest{
		name:   "missing location id",
		body:   fmt.Sprintf(`{"name": "Nemo", "ip_address_id": "%s","load_balancer_size": "small","load_balancer_type": "layer-3"}`, uuid.NewString()),
		status: http.StatusBadRequest,
		path:   baseURLTenant,
		method: http.MethodPost,
	})

	doHTTPTest(t, &httpTest{
		name:   "missing ip address",
		body:   fmt.Sprintf(`{"name": "Anchor", "location_id": "%s", "load_balancer_size": "small","load_balancer_type": "layer-3"}`, locationID),
		status: http.StatusOK,
		path:   baseURLTenant,
		method: http.MethodPost,
	})

	doHTTPTest(t, &httpTest{
		name:   "missing size",
		body:   fmt.Sprintf(`{"name": "Nemo", "location_id": "%s", "ip_address_id": testIPaddressUUID,"load_balancer_type": "layer-3"}`, locationID),
		status: http.StatusBadRequest,
		path:   baseURLTenant,
		method: http.MethodPost,
	})

	doHTTPTest(t, &httpTest{
		name:   "missing type",
		body:   fmt.Sprintf(`{"name": "Nemo", "location_id": "%s", "ip_address_id": "%s","load_balancer_size": "small"}`, locationID, uuid.NewString()),
		status: http.StatusBadRequest,
		path:   baseURLTenant,
		method: http.MethodPost,
	})

	doHTTPTest(t, &httpTest{
		name:   "invalid type",
		body:   fmt.Sprintf(`{"name": "Nemo", "location_id": "%s","ip_address_id": testIPaddressUUID,"load_balancer_size": "small","load_balancer_type": "layer-12"}`, locationID),
		status: http.StatusBadRequest,
		path:   baseURLTenant,
		method: http.MethodPost,
	})

	doHTTPTest(t, &httpTest{
		name:   "bad ip address uuid",
		body:   fmt.Sprintf(`{"name": "Nemo", "location_id": "%s", "ip_address_id": "Dori","load_balancer_size": "small","load_balancer_type": "layer-3"}`, locationID),
		status: http.StatusBadRequest,
		path:   baseURLTenant,
		method: http.MethodPost,
	})

	doHTTPTest(t, &httpTest{
		name:   "empty body",
		body:   ``,
		status: http.StatusBadRequest,
		path:   baseURLTenant,
		method: http.MethodPost,
	})

	doHTTPTest(t, &httpTest{
		name:   "bad body",
		body:   `bad body`,
		status: http.StatusBadRequest,
		path:   baseURLTenant,
		method: http.MethodPost,
	})

	// PUT tests
	doHTTPTest(t, &httpTest{
		name:   "happy path update load balancer",
		body:   fmt.Sprintf(`{"name": "Bruce", "load_balancer_size": "x-large","load_balancer_type": "layer-3","location_id": "%s","ip_address_id": "%s"}`, locationID, uuid.NewString()),
		status: http.StatusAccepted,
		method: http.MethodPut,
		path:   baseURL + "/" + testLoadBalancer.LoadBalancerID,
		tenant: tenantID,
	})

	doHTTPTest(t, &httpTest{
		name:   "update load balancer missing name",
		body:   `{"load_balancer_size": "x-large","load_balancer_type": "layer-3"}`,
		status: http.StatusBadRequest,
		method: http.MethodPut,
		path:   baseURL + "/" + testLoadBalancer.LoadBalancerID,
		tenant: tenantID,
	})

	doHTTPTest(t, &httpTest{
		name:   "update load balancer, missing size",
		body:   `{"name": "Bruce","load_balancer_type": "layer-3"}`,
		status: http.StatusBadRequest,
		method: http.MethodPut,
		path:   baseURL + "/" + testLoadBalancer.LoadBalancerID,
		tenant: tenantID,
	})

	doHTTPTest(t, &httpTest{
		name:   "update load balancer, missing type",
		body:   `{"name": "Bruce", "load_balancer_size": "x-large"}`,
		status: http.StatusBadRequest,
		method: http.MethodPut,
		path:   baseURL + "/" + testLoadBalancer.LoadBalancerID,
		tenant: tenantID,
	})

	doHTTPTest(t, &httpTest{
		name:   "update load balancer, missing load balancer id",
		body:   `{"name": "Bruce", "load_balancer_size": "x-large","load_balancer_type": "layer-3"}`,
		status: http.StatusNotFound,
		method: http.MethodPut,
		path:   baseURL,
		tenant: tenantID,
	})

	// PATCH endpoints
	doHTTPTest(t, &httpTest{
		name:   "happy path patch update load balancer name",
		body:   `{"name": "Brucey"}`,
		status: http.StatusAccepted,
		method: http.MethodPatch,
		path:   baseURL + "/" + testLoadBalancer.LoadBalancerID,
		tenant: tenantID,
	})

	doHTTPTest(t, &httpTest{
		name:   "sad path patch update load balancer name",
		body:   `{"name": ""}`,
		status: http.StatusBadRequest,
		method: http.MethodPatch,
		path:   baseURL + "/" + testLoadBalancer.LoadBalancerID,
		tenant: tenantID,
	})

	doHTTPTest(t, &httpTest{
		name:   "happy path patch update load balancer size",
		body:   `{"load_balancer_size": "x-x-large"}`,
		status: http.StatusAccepted,
		method: http.MethodPatch,
		path:   baseURL + "/" + testLoadBalancer.LoadBalancerID,
		tenant: tenantID,
	})

	doHTTPTest(t, &httpTest{
		name:   "sad path patch update load balancer size",
		body:   `{"load_balancer_size": ""}`,
		status: http.StatusBadRequest,
		method: http.MethodPatch,
		path:   baseURL + "/" + testLoadBalancer.LoadBalancerID,
		tenant: tenantID,
	})

	doHTTPTest(t, &httpTest{
		name:   "happy path patch update load balancer size",
		body:   `{"load_balancer_type": "layer-3"}`, // only allowed type is layer-3
		status: http.StatusAccepted,
		method: http.MethodPatch,
		path:   baseURL + "/" + testLoadBalancer.LoadBalancerID,
		tenant: tenantID,
	})

	doHTTPTest(t, &httpTest{
		name:   "sad path patch update load balancer size",
		body:   `{"load_balancer_type": ""}`,
		status: http.StatusBadRequest,
		method: http.MethodPatch,
		path:   baseURL + "/" + testLoadBalancer.LoadBalancerID,
		tenant: tenantID,
	})

	doHTTPTest(t, &httpTest{
		name:   "happy path patch update load balancer size",
		body:   fmt.Sprintf(`{"location_id": "%s"}`, uuid.NewString()),
		status: http.StatusAccepted,
		method: http.MethodPatch,
		path:   baseURL + "/" + testLoadBalancer.LoadBalancerID,
		tenant: tenantID,
	})

	doHTTPTest(t, &httpTest{
		name:   "sad path patch update load balancer size",
		body:   `{"location_id": ""}`,
		status: http.StatusBadRequest,
		method: http.MethodPatch,
		path:   baseURL + "/" + testLoadBalancer.LoadBalancerID,
		tenant: tenantID,
	})

	doHTTPTest(t, &httpTest{
		name:   "happy path patch update load balancer size",
		body:   fmt.Sprintf(`{"ip_address_id": "%s"}`, uuid.NewString()),
		status: http.StatusAccepted,
		method: http.MethodPatch,
		path:   baseURL + "/" + testLoadBalancer.LoadBalancerID,
		tenant: tenantID,
	})

	doHTTPTest(t, &httpTest{
		name:   "sad path patch update load balancer size",
		body:   `{"ip_address_id": ""}`,
		status: http.StatusBadRequest,
		method: http.MethodPatch,
		path:   baseURL + "/" + testLoadBalancer.LoadBalancerID,
		tenant: tenantID,
	})

	// GET tests
	doHTTPTest(t, &httpTest{
		name:   "happy path",
		path:   baseURLTenant,
		status: http.StatusOK,
		method: http.MethodGet,
	})

	doHTTPTest(t, &httpTest{
		name:   "happy path nemo by name",
		path:   baseURLTenant + "?name=Nemo",
		status: http.StatusOK,
		method: http.MethodGet,
	})

	doHTTPTest(t, &httpTest{
		name:   "happy path nemo by ip",
		path:   baseURLTenant + "?ip_address_id=" + testIPaddressUUIDNemo,
		status: http.StatusOK,
		method: http.MethodGet,
	})

	doHTTPTest(t, &httpTest{
		name:   "sad path lb doesnt exist",
		path:   baseURL + "/" + missingUUID,
		status: http.StatusNotFound,
		method: http.MethodGet,
	})

	doHTTPTest(t, &httpTest{
		name:   "sad path bad uuid",
		path:   baseURL + "/123456",
		status: http.StatusBadRequest,
		method: http.MethodGet,
	})

	doHTTPTest(t, &httpTest{
		name:   "list lbs with invalid tenant id",
		path:   srv.URL + "/v1/tenant/123456/loadbalancers",
		status: http.StatusBadRequest,
		method: http.MethodGet,
	})

	doHTTPTest(t, &httpTest{
		name:   "list lbs with unknown tenant id",
		path:   srv.URL + "/v1/tenant/" + missingUUID + "/loadbalancers",
		status: http.StatusOK,
		method: http.MethodGet,
	})

	// DELETE tests
	doHTTPTest(t, &httpTest{
		name:   "delete invalid id",
		path:   baseURL + "/invalid",
		status: http.StatusBadRequest,
		method: http.MethodDelete,
	})

	doHTTPTest(t, &httpTest{
		name:   "delete lb that doesnt exist",
		path:   baseURL + "/ce94616e-3798-454d-91f3-9e3cec32bff6",
		status: http.StatusNotFound,
		method: http.MethodDelete,
	})

	doHTTPTest(t, &httpTest{
		name:   "delete lb without id",
		path:   baseURL,
		status: http.StatusNotFound,
		method: http.MethodDelete,
	})

	doHTTPTest(t, &httpTest{
		name:   "delete small load balancers",
		path:   baseURLTenant + "?load_balancer_size=small",
		status: http.StatusUnprocessableEntity,
		method: http.MethodDelete,
	})

	doHTTPTest(t, &httpTest{
		name:   "delete nemo by Name",
		path:   baseURLTenant + "?slug=nemo",
		status: http.StatusOK,
		method: http.MethodDelete,
	})

	doHTTPTest(t, &httpTest{
		name:   "delete Dori by name",
		path:   baseURLTenant + "?slug=dori",
		status: http.StatusOK,
		method: http.MethodDelete,
	})

	doHTTPTest(t, &httpTest{
		name:   "delete Dori again",
		path:   baseURLTenant + "?slug=dori",
		status: http.StatusNotFound,
		method: http.MethodDelete,
	})
}

func TestLoadBalancerGet(t *testing.T) {
	nsrv := newNatsTestServer(t, "load-balancer-api-test", "com.infratographer.events.>")
	defer nsrv.Shutdown()

	srv := newTestServer(t, nsrv.ClientURL())
	defer srv.Close()

	assert.NotNil(t, srv)

	locationID := uuid.New().String()
	baseURL := srv.URL + "/v1/loadbalancers"
	missingUUID := uuid.New().String()

	// Create a load balancer to use for testing
	lb, cleanupLB := createLoadBalancer(t, srv, locationID)
	defer cleanupLB(t)

	doHTTPTest(t, &httpTest{
		name:   "get a list of loadblancer in the tenant",
		method: http.MethodGet,
		path:   srv.URL + "/v1/tenant/" + lb.TenantID + "/loadbalancers",
		status: http.StatusOK,
	})

	doHTTPTest(t, &httpTest{
		name:   "get a list of loadblancer in specified location",
		method: http.MethodGet,
		path:   baseURL + "/locations/" + locationID,
		status: http.StatusOK,
	})

	doHTTPTest(t, &httpTest{
		name:   "get a list of loadblancer in invalid location",
		method: http.MethodGet,
		path:   baseURL + "/locations/bad-uuid",
		status: http.StatusBadRequest,
	})

	doHTTPTest(t, &httpTest{
		name:   "get a list of loadblancer with no location",
		method: http.MethodGet,
		path:   baseURL + "/locations/",
		status: http.StatusNotFound,
	})

	doHTTPTest(t, &httpTest{
		name:   "get loadblancer by id",
		method: http.MethodGet,
		path:   baseURL + "/" + lb.ID,
		status: http.StatusOK,
	})

	doHTTPTest(t, &httpTest{
		name:   "get missing loadblancer by id",
		method: http.MethodGet,
		path:   baseURL + "/" + missingUUID,
		status: http.StatusNotFound,
	})

	doHTTPTest(t, &httpTest{
		name:   "get loadblancer by id on unknown tenant",
		method: http.MethodGet,
		path:   srv.URL + "/v1/tenant/" + missingUUID + "/loadbalancers/" + lb.ID,
		status: http.StatusNotFound,
	})

	doHTTPTest(t, &httpTest{
		name:   "get loadblancer without id",
		method: http.MethodGet,
		path:   baseURL,
		status: http.StatusNotFound,
	})

	// Test loadbalancer list response
	listReq, err := http.NewRequestWithContext(
		context.TODO(),
		http.MethodGet,
		srv.URL+"/v1/tenant/"+lb.TenantID+"/loadbalancers",
		nil,
	)
	assert.NoError(t, err)

	listResp, err := http.DefaultClient.Do(listReq)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, listResp.StatusCode)

	defer listResp.Body.Close()

	ca := lb.CreatedAt.Format(time.RFC3339Nano)
	ua := lb.UpdatedAt.Format(time.RFC3339Nano)
	testLBListExpected := fmt.Sprintf(`{"version":"v1","kind":"loadBalancersList","load_balancers":[{"created_at":"%s","updated_at":"%s","id":"%s","ip_address_id":"%s","tenant_id":"%s","name":"%s","location_id":"%s","load_balancer_size":"%s","load_balancer_type":"%s","ports":[]}]}`+"\n", ca, ua, lb.ID, lb.IPAddressID, lb.TenantID, lb.Name, lb.LocationID, lb.Size, lb.Type)
	testLBList, err := io.ReadAll(listResp.Body)
	assert.NoError(t, err)
	assert.Equal(t, testLBListExpected, string(testLBList))

	// Test loadbalancer get by id response
	getReq, err := http.NewRequestWithContext(
		context.TODO(),
		http.MethodGet,
		baseURL+"/"+lb.ID,
		nil,
	)
	assert.NoError(t, err)

	getResp, err := http.DefaultClient.Do(getReq)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, listResp.StatusCode)

	defer getResp.Body.Close()

	testLBGetExpected := fmt.Sprintf(`{"version":"v1","kind":"loadBalancersGet","load_balancer":{"created_at":"%s","updated_at":"%s","id":"%s","ip_address_id":"%s","tenant_id":"%s","name":"%s","location_id":"%s","load_balancer_size":"%s","load_balancer_type":"%s","ports":[]}}`+"\n", ca, ua, lb.ID, lb.IPAddressID, lb.TenantID, lb.Name, lb.LocationID, lb.Size, lb.Type)
	testLBGet, err := io.ReadAll(getResp.Body)
	assert.NoError(t, err)
	assert.Equal(t, testLBGetExpected, string(testLBGet))
}

func createLoadBalancer(t *testing.T, srv *httptest.Server, locationID string) (*loadBalancer, func(t *testing.T)) {
	tenantID := uuid.New().String()
	baseURL := srv.URL + "/v1/tenant/" + tenantID + "/loadbalancers"

	test := &httpTest{
		name:   "create nemo lb",
		body:   fmt.Sprintf(`{"name": "Nemo", "location_id": "%s", "ip_addr": "1.1.1.1","load_balancer_size": "small","load_balancer_type": "layer-3"}`, locationID),
		path:   baseURL,
		method: http.MethodPost,
		status: http.StatusOK,
	}

	doHTTPTest(t, test)

	// get loadbalancer by name
	loadbalancer := response{}

	t.Run("get nemo by name:[POST] "+baseURL+"?name=Nemo", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, baseURL+"?name=Nemo", nil) //nolint
		assert.NoError(t, err)

		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		err = json.NewDecoder(resp.Body).Decode(&loadbalancer)
		assert.NoError(t, err)
		resp.Body.Close()
	})

	return (*loadbalancer.LoadBalancers)[0], func(t *testing.T) {
		test := &httpTest{
			name:   "delete nemo",
			path:   baseURL + "?slug=nemo",
			method: http.MethodDelete,
			status: http.StatusOK,
		}

		doHTTPTest(t, test)
	}
}
