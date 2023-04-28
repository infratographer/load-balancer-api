package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/sqlboiler/v4/types"

	"go.infratographer.com/load-balancer-api/internal/httptools"
)

func testMetadataToString(t *testing.T, ns, data string) string {
	t.Helper()

	md := struct {
		Namespace string     `json:"namespace"`
		Data      types.JSON `json:"data"`
	}{
		Namespace: ns,
		Data:      types.JSON(data),
	}
	bytes, err := json.Marshal(md)
	assert.NoError(t, err)

	return string(bytes)
}

func TestMetadataRoutes(t *testing.T) {
	nsrv := newNatsTestServer(t, "load-balancer-api-test", "com.infratographer.events.>")
	defer nsrv.Shutdown()

	srv := newTestServer(t, nsrv.ClientURL())
	defer srv.Close()

	lb, cleanupLoadBalancers := createLoadBalancer(t, srv, uuid.New().String())
	defer cleanupLoadBalancers(t)

	lb2, cleanupLoadBalancers2 := createLoadBalancer(t, srv, uuid.New().String())
	defer cleanupLoadBalancers2(t)

	baseURL := srv.URL + "/v1/metadata"
	baseURLLoadBalancer := srv.URL + "/v1/loadbalancers/" + lb.ID + "/metadata"

	testMetadata := struct {
		Version string `json:"version"`
		Message string `json:"message"`
		Status  int    `json:"status"`
		ID      string `json:"metadata_id"`
	}{}

	{
		// create metadata
		req1, err := http.NewRequestWithContext(
			context.TODO(),
			http.MethodPost,
			baseURLLoadBalancer,
			httptools.FakeBody(testMetadataToString(t, "rhyme", `{"owner":"mary","animal":"lamb","fleece":"white as snow"}`)),
		)
		assert.NoError(t, err)
		resp1, err := http.DefaultClient.Do(req1)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp1.StatusCode)

		assert.NoError(t, json.NewDecoder(resp1.Body).Decode(&testMetadata))

		t.Logf("metadata: %+v", testMetadata)

		fmt.Println("foo")

		resp1.Body.Close()
	}

	// Clean up the rhyme last
	defer doHTTPTest(t, &httpTest{
		name:   "delete metadata get by id",
		path:   baseURL + "/" + testMetadata.ID,
		status: http.StatusOK,
		method: http.MethodDelete,
	})

	metadataList := struct {
		Version string `json:"version"`
		Message string `json:"message"`
		Status  int    `json:"status"`
		Data    []struct {
			ID        string `json:"metadata_id"`
			Namespace string `json:"namespace"`
			Data      string `json:"data"`
		} `json:"data"`
	}{}

	{
		req2, err := http.NewRequestWithContext(
			context.TODO(),
			http.MethodGet,
			baseURLLoadBalancer,
			nil,
		)
		assert.NoError(t, err)

		resp2, err := http.DefaultClient.Do(req2)
		assert.NoError(t, err)

		assert.NoError(t, json.NewDecoder(resp2.Body).Decode(&metadataList))

		assert.Equal(t, http.StatusOK, resp2.StatusCode)

		resp2.Body.Close()
	}
	// POST tests
	doHTTPTest(t, &httpTest{
		name:   "bad lb id",
		body:   testMetadataToString(t, "rhyme", `{"owner":"mary","animal":"lamb","fleece":"white as snow"}`),
		status: http.StatusBadRequest,
		method: http.MethodPost,
		path:   srv.URL + "/v1/loadbalancers/1234/metadata",
	})

	doHTTPTest(t, &httpTest{
		name:   "wrong lb id",
		body:   testMetadataToString(t, "rhyme", `{"owner":"mary","animal":"lamb","fleece":"white as snow"}`),
		status: http.StatusNotFound,
		method: http.MethodPost,
		path:   srv.URL + "/v1/loadbalancers/" + uuid.New().String() + "/metadata",
	})

	doHTTPTest(t, &httpTest{
		name:   "happy path",
		body:   testMetadataToString(t, "joke", `{"type":"knock knock","setup":"who's there?","punchline":"banana"}`),
		status: http.StatusOK,
		method: http.MethodPost,
		path:   baseURLLoadBalancer,
	})

	// Clean up the joke after other test cases
	defer doHTTPTest(t, &httpTest{
		name:   "delete metadata by namespace",
		path:   baseURLLoadBalancer + "?namespace=joke",
		status: http.StatusOK,
		method: http.MethodDelete,
	})

	doHTTPTest(t, &httpTest{
		name:   "duplicate namespace",
		body:   testMetadataToString(t, "rhyme", `{"has_wool":true,"color":"black","greeting":"baa baa","quanity":3,"animal":"sheep"}`),
		status: http.StatusBadRequest,
		method: http.MethodPost,
		path:   baseURLLoadBalancer,
	})

	doHTTPTest(t, &httpTest{
		name:   "duplicate",
		body:   testMetadataToString(t, "rhyme", `{"owner":"mary","animal":"lamb","fleece":"white as snow"}`),
		status: http.StatusBadRequest,
		method: http.MethodPost,
		path:   baseURLLoadBalancer,
	})

	doHTTPTest(t, &httpTest{
		name:   "missing namespace",
		body:   `{"data": {"owner":"mary","animal":"lamb","fleece":"white as snow"}}`,
		status: http.StatusBadRequest,
		method: http.MethodPost,
		path:   baseURLLoadBalancer,
	})

	doHTTPTest(t, &httpTest{
		name:   "missing data",
		body:   `{"namespace": "rhyme"}`,
		status: http.StatusBadRequest,
		method: http.MethodPost,
		path:   baseURLLoadBalancer,
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

	// Get Tests
	doHTTPTest(t, &httpTest{
		name:   "get metadata get by id",
		path:   baseURL + "/" + testMetadata.ID,
		status: http.StatusOK,
		method: http.MethodGet,
	})

	doHTTPTest(t, &httpTest{
		name:   "get all metadata by lb id",
		path:   baseURLLoadBalancer,
		status: http.StatusOK,
		method: http.MethodGet,
	})

	doHTTPTest(t, &httpTest{
		name:   "get metadata by lb id and namespace",
		path:   baseURLLoadBalancer + "?namespace=rhyme",
		status: http.StatusOK,
		method: http.MethodGet,
	})

	doHTTPTest(t, &httpTest{
		name:   "get all metadata by lb id bad id",
		path:   srv.URL + "/v1/loadbalancers/1234/metadata",
		status: http.StatusBadRequest,
		method: http.MethodGet,
	})

	doHTTPTest(t, &httpTest{
		name:   "get all metadata by lb id wrong id",
		path:   srv.URL + "/v1/loadbalancers/" + lb2.ID + "/metadata",
		status: http.StatusNotFound,
		method: http.MethodGet,
	})

	// Patch tests
	doHTTPTest(t, &httpTest{
		name:   "patch rhyming metadata",
		body:   `{ "data": {"roses":"red","violets":"blue"}}`,
		path:   baseURLLoadBalancer + "?namespace=rhyme",
		status: http.StatusAccepted,
		method: http.MethodPatch,
	})

	doHTTPTest(t, &httpTest{
		name:   "ambiguous patch",
		body:   `{ "data": {"roses":"red","violets":"blue"}}`,
		path:   baseURLLoadBalancer,
		status: http.StatusBadRequest,
		method: http.MethodPatch,
	})

	doHTTPTest(t, &httpTest{
		name:   "patch metadata by bad id",
		body:   `{ "data": {"roses":"red","violets":"blue"}}`,
		path:   baseURL + "/1234/metadata?namespace=rhyme",
		status: http.StatusNotFound,
		method: http.MethodPatch,
	})

	doHTTPTest(t, &httpTest{
		name:   "patch metadata with bad body",
		body:   `bad body`,
		path:   baseURLLoadBalancer + "?namespace=rhyme",
		status: http.StatusBadRequest,
		method: http.MethodPatch,
	})

	// PUT tests

	doHTTPTest(t, &httpTest{
		name:   "put rhyming metadata",
		body:   testMetadataToString(t, "rhyme", `{"has_wool":true,"color":"black","greeting":"baa baa","quanity":2,"animal":"sheep"}`),
		path:   baseURLLoadBalancer + "?metadata_id=" + testMetadata.ID,
		status: http.StatusAccepted,
		method: http.MethodPut,
	})

	doHTTPTest(t, &httpTest{
		name:   "ambiguous put",
		body:   testMetadataToString(t, "rhyme", `{"has_wool":true,"color":"black","greeting":"baa baa","quanity":2,"animal":"sheep"}`),
		path:   baseURLLoadBalancer,
		status: http.StatusBadRequest,
		method: http.MethodPut,
	})

	doHTTPTest(t, &httpTest{
		name:   "put metadata by bad id",
		body:   testMetadataToString(t, "rhyme", `{"has_wool":true,"color":"black","greeting":"baa baa","quanity":2,"animal":"sheep"}`),
		path:   baseURL + "/1234/metadata?namespace=rhyme",
		status: http.StatusNotFound,
		method: http.MethodPut,
	})

	doHTTPTest(t, &httpTest{
		name:   "put metadata with bad body",
		body:   `bad body`,
		path:   baseURLLoadBalancer + "?namespace=rhyme",
		status: http.StatusBadRequest,
		method: http.MethodPut,
	})

	doHTTPTest(t, &httpTest{
		name:   "put metadata with bad param",
		body:   testMetadataToString(t, "rhyme", `{"has_wool":true,"color":"black","greeting":"baa baa","quanity":2,"animal":"sheep"}`),
		path:   baseURLLoadBalancer + "?namespace=1234",
		status: http.StatusNotFound,
		method: http.MethodPut,
	})

	// DELETE tests
	doHTTPTest(t, &httpTest{
		name:   "try ambiguous delete",
		path:   baseURLLoadBalancer,
		method: http.MethodDelete,
		status: http.StatusBadRequest,
	})

	doHTTPTest(t, &httpTest{
		name:   "delete metadata get by bad id",
		path:   baseURL + "/1234",
		status: http.StatusBadRequest,
		method: http.MethodDelete,
	})

	doHTTPTest(t, &httpTest{
		name:   "bad delete metadata get by id",
		path:   baseURL + "/" + uuid.New().String(),
		status: http.StatusNotFound,
		method: http.MethodDelete,
	})
}
