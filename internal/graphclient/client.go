package graphclient

import (
	"fmt"

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

// QueryLoadBalancer will return the results from a query loadBalancer with a given id.
func (c *Client) QueryLoadBalancer(id gidx.PrefixedID) (*LoadBalancer, error) {
	// ideal query with external references, but these aren't all setup yet
	// q := `
	// query ($id: ID!) {
	// 	loadBalancer(id: $id) {
	// 		id
	// 		name
	// 		location {
	// 			id
	// 		}
	// 		loadBalancerProvider {
	// 			id
	// 		}
	// 		tenant {
	// 			id
	// 		}
	// 		createdAt
	// 		updatedAt
	// 	}
	// }`
	q := `
	query ($id: ID!) {
		loadBalancer(id: $id) {
			id
			name
			loadBalancerProvider {
				id
			}
			createdAt
			updatedAt
		}
	}`

	var resp struct {
		LoadBalancer *LoadBalancer `json:"loadBalancer"`
	}

	variables := []client.Option{
		client.Var("id", id.String()),
	}

	err := c.graphClient.Post(q, &resp, variables...)

	return resp.LoadBalancer, err
}

// LoadBalancerCreate will return the results from a mutation loadBalancerCreate request
func (c *Client) LoadBalancerCreate(input ent.CreateLoadBalancerInput) (*LoadBalancer, error) {
	q := `
	mutation ($input: CreateLoadBalancerInput!) {
		loadBalancerCreate(input: $input) {
			loadBalancer {
				id
				name
				loadBalancerProvider {
					id
				}
				tenantID
				locationID
				createdAt
				updatedAt
			}
		}
	}`

	var resp Mutations

	variables := []client.Option{
		client.Var("input", map[string]string{
			"name":       input.Name,
			"locationID": input.LocationID.String(),
			"tenantID":   input.TenantID.String(),
			"providerID": input.ProviderID.String(),
		}),
	}

	err := c.graphClient.Post(q, &resp, variables...)

	return resp.LoadBalancerCreate.LoadBalancer, err
}

// LoadBalancerUpdate will return the results from a mutation loadBalancerUpdate request
func (c *Client) LoadBalancerUpdate(id gidx.PrefixedID, input ent.UpdateLoadBalancerInput) (*LoadBalancer, error) {
	q := `
	mutation ($id: ID!, $input: UpdateLoadBalancerInput!) {
		loadBalancerUpdate(id: $id, input: $input) {
			loadBalancer {
				id
				name
				createdAt
				updatedAt
			}
		}
	}`

	var resp Mutations

	inputAsMap := map[string]string{}

	if input.Name != nil {
		inputAsMap["name"] = *input.Name
	}

	variables := []client.Option{
		client.Var("id", id.String()),
		client.Var("input", inputAsMap),
	}

	err := c.graphClient.Post(q, &resp, variables...)

	return resp.LoadBalancerUpdate.LoadBalancer, err
}

// LoadBalancerDelete will return the results from a mutation loadBalancerDelete request
func (c *Client) LoadBalancerDelete(id gidx.PrefixedID) (gidx.PrefixedID, error) {
	q := `
	mutation ($id: ID!) {
		loadBalancerDelete(id: $id) {
			deletedID
		}
	}`

	var resp Mutations

	variables := []client.Option{client.Var("id", id.String())}
	err := c.graphClient.Post(q, &resp, variables...)

	return resp.LoadBalancerDelete.DeletedID, err
}

// LoadBalancerPortCreate will return the results from a mutation loadBalancerPortCreate request
func (c *Client) LoadBalancerPortCreate(input ent.CreateLoadBalancerPortInput) (*LoadBalancerPort, error) {
	q := `
	mutation ($input: CreateLoadBalancerPortInput!) {
		loadBalancerPortCreate(input: $input) {
			loadBalancerPort {
				id
				name
				number
				loadBalancer {
					id
				}
				createdAt
				updatedAt

			}
		}
	}`

	var resp Mutations

	variables := []client.Option{
		client.Var("input", map[string]string{
			"name":           input.Name,
			"number":         fmt.Sprint(input.Number),
			"loadBalancerID": input.LoadBalancerID.String(),
		}),
	}

	err := c.graphClient.Post(q, &resp, variables...)

	return resp.LoadBalancerPortCreate.LoadBalancerPort, err
}

// LoadBalancerPortDelete will return the results from a mutation loadBalancerDelete request
func (c *Client) LoadBalancerPortDelete(id gidx.PrefixedID) (gidx.PrefixedID, error) {
	q := `
	mutation ($id: ID!) {
		loadBalancerPortDelete(id: $id) {
			deletedID
		}
	}`

	var resp Mutations

	variables := []client.Option{client.Var("id", id.String())}
	err := c.graphClient.Post(q, &resp, variables...)

	return resp.LoadBalancerPortDelete.DeletedID, err
}

// LoadBalancerPortUpdate will return the results from a mutation loadBalancerPortUpdate request
func (c *Client) LoadBalancerPortUpdate(id gidx.PrefixedID, input ent.UpdateLoadBalancerPortInput) (*LoadBalancerPort, error) {
	q := `
	mutation ($id: ID!, $input: UpdateLoadBalancerPortInput!) {
		loadBalancerPortUpdate(id: $id, input: $input) {
			loadBalancerPort {
				id
				name
				number
				createdAt
				updatedAt
			}
		}
	}`

	var resp Mutations

	inputAsMap := map[string]interface{}{}

	if input.Name != nil {
		inputAsMap["name"] = *input.Name
	}

	if input.Number != nil {
		inputAsMap["number"] = *input.Number
	}

	variables := []client.Option{
		client.Var("id", id.String()),
		client.Var("input", inputAsMap),
	}

	err := c.graphClient.Post(q, &resp, variables...)

	return resp.LoadBalancerPortUpdate.LoadBalancerPort, err
}

// // QueryLoadBalancerPortByID will return the results from a query loadBalancer with a given id.
// func (c *Client) QueryLoadBalancerPortByID(id gidx.PrefixedID, portid gidx.PrefixedID) (*LoadBalancerPort, error) {
// 	q := `
// 	query ($id: ID!, $portid: ID!) {
// 		loadBalancer(id: $id) {
// 			ports(where: {id: $portid}) {
//         		edges{
//           			node{
//             			id
//             			number
//             			loadBalancer {
//               				id
//             			}
//             			createdAt
//             			updatedAt
//           			}
//         		}
// 			}
// 		}
// 	}`

// 	var resp struct {
// 		LoadBalancer *LoadBalancer `json:"loadBalancer"`
// 	}

// 	variables := []client.Option{
// 		client.Var("id", id.String()),
// 		client.Var("portid", portid.String()),
// 	}

// 	err := c.graphClient.Post(q, &resp, variables...)

// 	//return resp.LoadBalancer.Ports[0], err
// 	return nil, err
// 	// return resp.LoadBalancer, err
// }
