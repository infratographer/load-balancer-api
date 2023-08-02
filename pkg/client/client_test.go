package client

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shurcooL/graphql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetLoadBalancer(t *testing.T) {
	cli := Client{}

	t.Run("bad prefix", func(t *testing.T) {
		lb, err := cli.GetLoadBalancer(context.Background(), "badprefix-test")
		require.Error(t, err)
		require.Nil(t, lb)
		assert.ErrorContains(t, err, "invalid id")
	})

	t.Run("successful query", func(t *testing.T) {
		respJSON := `{
	"data": {
		"loadBalancer": {
			"id": "loadbal-randovalue",
			"name": "some lb",
			"IPAddresses": [
				{
					"id": "ipamipa-randovalue",
					"ip": "192.168.1.42",
					"reserved": false
				},
				{
					"id": "ipamipa-randovalue2",
					"ip": "192.168.1.1",
					"reserved": true
				}
			],
			"ports": {
				"edges": [
					{
						"node": {
							"name": "porty",
							"id": "loadprt-randovalue",
							"number": 80,
							"pools": {
								"edges": [
									{
										"node": {
											"id": "loadpol-pooly",
											"name": "pooly",
											"protocol": "tcp",
											"origins": {
												"edges": [
													{
														"node": {
															"id": "loadori-origin",
															"name": "origin",
															"target": "1.2.3.4",
															"portNumber": 80
														}
													}
												]
											}
										}
									}
								]
							}
						}
					}
				]
			}
		}
	}
}`

		cli.gqlCli = mustNewGQLTestClient(respJSON, http.StatusOK)
		lb, err := cli.GetLoadBalancer(context.Background(), "loadbal-randovalue")
		require.NoError(t, err)
		require.NotNil(t, lb)

		assert.Equal(t, "loadbal-randovalue", lb.ID)
		assert.Equal(t, "some lb", lb.Name)
		assert.Equal(t, "porty", lb.Ports.Edges[0].Node.Name)
		assert.Equal(t, int64(80), lb.Ports.Edges[0].Node.Number)

		require.Len(t, lb.Ports.Edges[0].Node.Pools.Edges, 1)
		assert.Equal(t, "loadpol-pooly", lb.Ports.Edges[0].Node.Pools.Edges[0].Node.ID)
		assert.Equal(t, "pooly", lb.Ports.Edges[0].Node.Pools.Edges[0].Node.Name)
		assert.Equal(t, "tcp", lb.Ports.Edges[0].Node.Pools.Edges[0].Node.Protocol)

		require.Len(t, lb.Ports.Edges[0].Node.Pools.Edges[0].Node.Origins.Edges, 1)
		assert.Equal(t, "loadori-origin", lb.Ports.Edges[0].Node.Pools.Edges[0].Node.Origins.Edges[0].Node.ID)
		assert.Equal(t, "origin", lb.Ports.Edges[0].Node.Pools.Edges[0].Node.Origins.Edges[0].Node.Name)
		assert.Equal(t, "1.2.3.4", lb.Ports.Edges[0].Node.Pools.Edges[0].Node.Origins.Edges[0].Node.Target)
		assert.Equal(t, int64(80), lb.Ports.Edges[0].Node.Pools.Edges[0].Node.Origins.Edges[0].Node.PortNumber)

		require.Len(t, lb.IPAddresses, 2)
		assert.Equal(t, "ipamipa-randovalue", lb.IPAddresses[0].ID)
		assert.Equal(t, "192.168.1.42", lb.IPAddresses[0].IP)
		assert.False(t, lb.IPAddresses[0].Reserved)

		assert.Equal(t, "ipamipa-randovalue2", lb.IPAddresses[1].ID)
		assert.Equal(t, "192.168.1.1", lb.IPAddresses[1].IP)
		assert.True(t, lb.IPAddresses[1].Reserved)
	})

	t.Run("unauthorized", func(t *testing.T) {
		respJSON := `{"message":"invalid or expired jwt"}`

		cli.gqlCli = mustNewGQLTestClient(respJSON, http.StatusUnauthorized)

		lb, err := cli.GetLoadBalancer(context.Background(), "loadbal-randovalue")
		require.Nil(t, lb)
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrUnauthorized)
	})

	t.Run("does not have permissions", func(t *testing.T) {
		respJSON := `{"message":"subject doesn't have access"}`

		cli.gqlCli = mustNewGQLTestClient(respJSON, http.StatusForbidden)

		lb, err := cli.GetLoadBalancer(context.Background(), "loadbal-randovalue")
		require.Nil(t, lb)
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrPermissionDenied)
	})

	t.Run("not found", func(t *testing.T) {
		respJSON := `{
			"data": null
			"errors": [
				{
					"message": "load_balancer not found"
				}
			]
		}`

		cli.gqlCli = mustNewGQLTestClient(respJSON, http.StatusUnauthorized)

		lb, err := cli.GetLoadBalancer(context.Background(), "loadbal-randovalue")
		require.Nil(t, lb)
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrLBNotfound)
	})

	t.Run("gql error", func(t *testing.T) {
		respJSON := `{
			"data": null
			"errors": [
				{
					"message": "failed to find or parse something"
				}
			]
		}`

		cli.gqlCli = mustNewGQLTestClient(respJSON, http.StatusOK)

		lb, err := cli.GetLoadBalancer(context.Background(), "loadbal-randovalue")
		require.Nil(t, lb)
		require.Error(t, err)
	})
}

func mustNewGQLTestClient(respJSON string, respCode int) *graphql.Client {
	mux := http.NewServeMux()
	mux.HandleFunc("/query", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(respCode)
		w.Header().Set("Content-Type", "application/json")
		_, err := io.WriteString(w, respJSON)
		if err != nil {
			panic(err)
		}
	})

	return graphql.NewClient("/query", &http.Client{Transport: localRoundTripper{handler: mux}})
}

type localRoundTripper struct {
	handler http.Handler
}

func (l localRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	w := httptest.NewRecorder()
	l.handler.ServeHTTP(w, req)

	return w.Result(), nil
}
