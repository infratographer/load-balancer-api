package graphclient

import (
	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
	"go.infratographer.com/x/gidx"

	ent "go.infratographer.com/load-balancer-api/internal/ent/generated"
	"go.infratographer.com/load-balancer-api/internal/graphapi"
)

// Client ...
type Client struct {
	graphClient *client.Client
}

// OrderBy ...
type OrderBy struct {
	Field     string `json:"field"`
	Direction string `json:"direction"`
}

// New returns a new client
func New(r *graphapi.Resolver) *Client {
	c := client.New(handler.NewDefaultServer(
		graphapi.NewExecutableSchema(
			graphapi.Config{Resolvers: r},
		)))

	return &Client{
		graphClient: c,
	}
}

// MustGetTenantLoadBalancers will return the loadbalancers for a tenant id. If an orderBy option is provided the returned
// load balancers will be sorted
func (c *Client) MustGetTenantLoadBalancers(id gidx.PrefixedID, orderBy *OrderBy) []*ent.LoadBalancer {
	q := `
	query ($_representations: [_Any!]!, $orderBy: LoadBalancerOrder) {
		_entities(representations: $_representations) {
			... on Tenant {
				loadBalancers(orderBy: $orderBy) {
					edges {
						node {
							id
							name
						}
					}
				}
			}
		}
	}`

	var resp struct {
		Entities []struct {
			LoadBalancers ent.LoadBalancerConnection `json:"loadBalancers"`
		} `json:"_entities"`
	}

	variables := []client.Option{
		client.Var("_representations", map[string]string{"__typename": "Tenant", "id": id.String()}),
	}

	if orderBy != nil {
		variables = append(variables, client.Var("orderBy", orderBy))
	}

	c.graphClient.MustPost(q, &resp, variables...)

	lbs := []*ent.LoadBalancer{}

	for _, edge := range resp.Entities[0].LoadBalancers.Edges {
		lbs = append(lbs, edge.Node)
	}

	return lbs
}
