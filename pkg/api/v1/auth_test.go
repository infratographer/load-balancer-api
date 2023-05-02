package api

import (
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.infratographer.com/x/echojwtx"
)

func TestLoadbalancerGETWithAuth(t *testing.T) {
	oauthCLI, issuer, oAuthClose := echojwtx.TestOAuthClient("urn:test:loadbalancer", "")
	defer oAuthClose()

	nsrv := newNatsTestServer(t, "load-balancer-api-test", "com.infratographer.events.>")
	defer nsrv.Shutdown()

	srv, err := newTestServer(t, nsrv.ClientURL(), &echojwtx.AuthConfig{
		Issuer: issuer,
	})

	require.NoError(t, err)
	require.NotNil(t, srv)

	defer srv.Close()

	tenantID := uuid.New().String()

	doHTTPTest(t, &httpTest{
		name:   "default client authorization failure",
		method: http.MethodGet,
		path:   srv.URL + "/v1/tenant/" + tenantID + "/loadbalancers",
		status: http.StatusUnauthorized,
	})

	doHTTPTest(t, &httpTest{
		client: oauthCLI,
		name:   "OAuth client authorization success",
		method: http.MethodGet,
		path:   srv.URL + "/v1/tenant/" + tenantID + "/loadbalancers",
		status: http.StatusOK,
	})
}
