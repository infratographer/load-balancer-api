package graphapi_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.infratographer.com/x/echojwtx"
)

func TestJWTEnabledLoadbalancerGETWithAuthClient(t *testing.T) {
	oauthCLI, issuer, oAuthClose := echojwtx.TestOAuthClient("urn:test:loadbalancer", "")
	defer oAuthClose()

	srv, err := newTestServer(&echojwtx.AuthConfig{
		Issuer: issuer,
	})

	require.NoError(t, err)
	require.NotNil(t, srv)

	defer srv.Close()

	ctx := context.Background()
	lb1 := (&LoadBalancerBuilder{}).MustNew(ctx)

	resp, err := newGraphTestClient(
		withGraphClientHTTPClient(oauthCLI),
		withGraphClientServerURL(srv.URL+"/query"),
	).GetLoadBalancer(ctx, lb1.ID)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, lb1.ID, resp.LoadBalancer.ID)
}

func TestJWTENabledLoadbalancerGETWithDefaultClient(t *testing.T) {
	_, issuer, oAuthClose := echojwtx.TestOAuthClient("urn:test:loadbalancer", "")
	defer oAuthClose()

	srv, err := newTestServer(&echojwtx.AuthConfig{
		Issuer: issuer,
	})

	require.NoError(t, err)
	require.NotNil(t, srv)

	defer srv.Close()

	ctx := context.Background()
	lb1 := (&LoadBalancerBuilder{}).MustNew(ctx)

	resp, err := newGraphTestClient(
		withGraphClientHTTPClient(http.DefaultClient),
		withGraphClientServerURL(srv.URL+"/query"),
	).GetLoadBalancer(ctx, lb1.ID)

	require.Error(t, err, "Expected an authorization error")
	require.Nil(t, resp)
	assert.ErrorContains(t, err, `{"networkErrors":{"code":401`)
}
