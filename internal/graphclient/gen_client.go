// Code generated by github.com/Yamashou/gqlgenc, DO NOT EDIT.

package graphclient

import (
	"context"
	"net/http"
	"time"

	"github.com/Yamashou/gqlgenc/client"
	"go.infratographer.com/x/gidx"
)

type GraphClient interface {
	GetLoadBalancer(ctx context.Context, id gidx.PrefixedID, httpRequestOptions ...client.HTTPRequestOption) (*GetLoadBalancer, error)
	GetLoadBalancerPool(ctx context.Context, id gidx.PrefixedID, httpRequestOptions ...client.HTTPRequestOption) (*GetLoadBalancerPool, error)
	GetLoadBalancerPoolOrigin(ctx context.Context, id gidx.PrefixedID, originid gidx.PrefixedID, httpRequestOptions ...client.HTTPRequestOption) (*GetLoadBalancerPoolOrigin, error)
	GetLoadBalancerPort(ctx context.Context, id gidx.PrefixedID, portid gidx.PrefixedID, httpRequestOptions ...client.HTTPRequestOption) (*GetLoadBalancerPort, error)
	GetLoadBalancerProvider(ctx context.Context, id gidx.PrefixedID, httpRequestOptions ...client.HTTPRequestOption) (*GetLoadBalancerProvider, error)
	GetOwnerLoadBalancers(ctx context.Context, id gidx.PrefixedID, orderBy *LoadBalancerOrder, httpRequestOptions ...client.HTTPRequestOption) (*GetOwnerLoadBalancers, error)
	LoadBalancerCreate(ctx context.Context, input CreateLoadBalancerInput, httpRequestOptions ...client.HTTPRequestOption) (*LoadBalancerCreate, error)
	LoadBalancerDelete(ctx context.Context, id gidx.PrefixedID, httpRequestOptions ...client.HTTPRequestOption) (*LoadBalancerDelete, error)
	LoadBalancerOriginCreate(ctx context.Context, input CreateLoadBalancerOriginInput, httpRequestOptions ...client.HTTPRequestOption) (*LoadBalancerOriginCreate, error)
	LoadBalancerOriginDelete(ctx context.Context, id gidx.PrefixedID, httpRequestOptions ...client.HTTPRequestOption) (*LoadBalancerOriginDelete, error)
	LoadBalancerOriginUpdate(ctx context.Context, id gidx.PrefixedID, input UpdateLoadBalancerOriginInput, httpRequestOptions ...client.HTTPRequestOption) (*LoadBalancerOriginUpdate, error)
	LoadBalancerPoolCreate(ctx context.Context, input CreateLoadBalancerPoolInput, httpRequestOptions ...client.HTTPRequestOption) (*LoadBalancerPoolCreate, error)
	LoadBalancerPoolDelete(ctx context.Context, id gidx.PrefixedID, httpRequestOptions ...client.HTTPRequestOption) (*LoadBalancerPoolDelete, error)
	LoadBalancerPoolUpdate(ctx context.Context, id gidx.PrefixedID, input UpdateLoadBalancerPoolInput, httpRequestOptions ...client.HTTPRequestOption) (*LoadBalancerPoolUpdate, error)
	LoadBalancerPortCreate(ctx context.Context, input CreateLoadBalancerPortInput, httpRequestOptions ...client.HTTPRequestOption) (*LoadBalancerPortCreate, error)
	LoadBalancerPortDelete(ctx context.Context, id gidx.PrefixedID, httpRequestOptions ...client.HTTPRequestOption) (*LoadBalancerPortDelete, error)
	LoadBalancerPortUpdate(ctx context.Context, id gidx.PrefixedID, input UpdateLoadBalancerPortInput, httpRequestOptions ...client.HTTPRequestOption) (*LoadBalancerPortUpdate, error)
	LoadBalancerProviderCreate(ctx context.Context, input CreateLoadBalancerProviderInput, httpRequestOptions ...client.HTTPRequestOption) (*LoadBalancerProviderCreate, error)
	LoadBalancerProviderDelete(ctx context.Context, id gidx.PrefixedID, httpRequestOptions ...client.HTTPRequestOption) (*LoadBalancerProviderDelete, error)
	LoadBalancerProviderUpdate(ctx context.Context, id gidx.PrefixedID, input UpdateLoadBalancerProviderInput, httpRequestOptions ...client.HTTPRequestOption) (*LoadBalancerProviderUpdate, error)
	LoadBalancerUpdate(ctx context.Context, id gidx.PrefixedID, input UpdateLoadBalancerInput, httpRequestOptions ...client.HTTPRequestOption) (*LoadBalancerUpdate, error)
}

type Client struct {
	Client *client.Client
}

func NewClient(cli *http.Client, baseURL string, options ...client.HTTPRequestOption) GraphClient {
	return &Client{Client: client.NewClient(cli, baseURL, options...)}
}

type Query struct {
	LoadBalancerPools    LoadBalancerPoolConnection "json:\"loadBalancerPools\" graphql:\"loadBalancerPools\""
	LoadBalancer         LoadBalancer               "json:\"loadBalancer\" graphql:\"loadBalancer\""
	LoadBalancerHistory  LoadBalancer               "json:\"loadBalancerHistory\" graphql:\"loadBalancerHistory\""
	LoadBalancerPool     LoadBalancerPool           "json:\"loadBalancerPool\" graphql:\"loadBalancerPool\""
	LoadBalancerProvider LoadBalancerProvider       "json:\"loadBalancerProvider\" graphql:\"loadBalancerProvider\""
	Entities             []Entity                   "json:\"_entities\" graphql:\"_entities\""
	Service              Service                    "json:\"_service\" graphql:\"_service\""
}
type Mutation struct {
	LoadBalancerOriginCreate   LoadBalancerOriginCreatePayload   "json:\"loadBalancerOriginCreate\" graphql:\"loadBalancerOriginCreate\""
	LoadBalancerOriginUpdate   LoadBalancerOriginUpdatePayload   "json:\"loadBalancerOriginUpdate\" graphql:\"loadBalancerOriginUpdate\""
	LoadBalancerOriginDelete   LoadBalancerOriginDeletePayload   "json:\"loadBalancerOriginDelete\" graphql:\"loadBalancerOriginDelete\""
	LoadBalancerCreate         LoadBalancerCreatePayload         "json:\"loadBalancerCreate\" graphql:\"loadBalancerCreate\""
	LoadBalancerUpdate         LoadBalancerUpdatePayload         "json:\"loadBalancerUpdate\" graphql:\"loadBalancerUpdate\""
	LoadBalancerDelete         LoadBalancerDeletePayload         "json:\"loadBalancerDelete\" graphql:\"loadBalancerDelete\""
	LoadBalancerPoolCreate     LoadBalancerPoolCreatePayload     "json:\"loadBalancerPoolCreate\" graphql:\"loadBalancerPoolCreate\""
	LoadBalancerPoolUpdate     LoadBalancerPoolUpdatePayload     "json:\"loadBalancerPoolUpdate\" graphql:\"loadBalancerPoolUpdate\""
	LoadBalancerPoolDelete     LoadBalancerPoolDeletePayload     "json:\"loadBalancerPoolDelete\" graphql:\"loadBalancerPoolDelete\""
	LoadBalancerPortCreate     LoadBalancerPortCreatePayload     "json:\"loadBalancerPortCreate\" graphql:\"loadBalancerPortCreate\""
	LoadBalancerPortUpdate     LoadBalancerPortUpdatePayload     "json:\"loadBalancerPortUpdate\" graphql:\"loadBalancerPortUpdate\""
	LoadBalancerPortDelete     LoadBalancerPortDeletePayload     "json:\"loadBalancerPortDelete\" graphql:\"loadBalancerPortDelete\""
	LoadBalancerProviderCreate LoadBalancerProviderCreatePayload "json:\"loadBalancerProviderCreate\" graphql:\"loadBalancerProviderCreate\""
	LoadBalancerProviderUpdate LoadBalancerProviderUpdatePayload "json:\"loadBalancerProviderUpdate\" graphql:\"loadBalancerProviderUpdate\""
	LoadBalancerProviderDelete LoadBalancerProviderDeletePayload "json:\"loadBalancerProviderDelete\" graphql:\"loadBalancerProviderDelete\""
}
type GetLoadBalancer struct {
	LoadBalancer struct {
		ID       gidx.PrefixedID "json:\"id\" graphql:\"id\""
		Name     string          "json:\"name\" graphql:\"name\""
		Location struct {
			ID gidx.PrefixedID "json:\"id\" graphql:\"id\""
		} "json:\"location\" graphql:\"location\""
		LoadBalancerProvider struct {
			ID gidx.PrefixedID "json:\"id\" graphql:\"id\""
		} "json:\"loadBalancerProvider\" graphql:\"loadBalancerProvider\""
		Owner struct {
			ID gidx.PrefixedID "json:\"id\" graphql:\"id\""
		} "json:\"owner\" graphql:\"owner\""
		CreatedAt time.Time "json:\"createdAt\" graphql:\"createdAt\""
		UpdatedAt time.Time "json:\"updatedAt\" graphql:\"updatedAt\""
	} "json:\"loadBalancer\" graphql:\"loadBalancer\""
}
type GetLoadBalancerPool struct {
	LoadBalancerPool struct {
		ID        gidx.PrefixedID          "json:\"id\" graphql:\"id\""
		Name      string                   "json:\"name\" graphql:\"name\""
		Protocol  LoadBalancerPoolProtocol "json:\"protocol\" graphql:\"protocol\""
		OwnerID   gidx.PrefixedID          "json:\"ownerID\" graphql:\"ownerID\""
		CreatedAt time.Time                "json:\"createdAt\" graphql:\"createdAt\""
		UpdatedAt time.Time                "json:\"updatedAt\" graphql:\"updatedAt\""
	} "json:\"loadBalancerPool\" graphql:\"loadBalancerPool\""
}
type GetLoadBalancerPoolOrigin struct {
	LoadBalancerPool struct {
		Origins struct {
			Edges []*struct {
				Node *struct {
					ID         gidx.PrefixedID "json:\"id\" graphql:\"id\""
					Name       string          "json:\"name\" graphql:\"name\""
					Target     string          "json:\"target\" graphql:\"target\""
					PortNumber int64           "json:\"portNumber\" graphql:\"portNumber\""
					Active     bool            "json:\"active\" graphql:\"active\""
					Weight     int64           "json:\"weight\" graphql:\"weight\""
					PoolID     gidx.PrefixedID "json:\"poolID\" graphql:\"poolID\""
					CreatedAt  time.Time       "json:\"createdAt\" graphql:\"createdAt\""
					UpdatedAt  time.Time       "json:\"updatedAt\" graphql:\"updatedAt\""
				} "json:\"node\" graphql:\"node\""
			} "json:\"edges\" graphql:\"edges\""
		} "json:\"origins\" graphql:\"origins\""
	} "json:\"loadBalancerPool\" graphql:\"loadBalancerPool\""
}
type GetLoadBalancerPort struct {
	LoadBalancer struct {
		Ports struct {
			Edges []*struct {
				Node *struct {
					ID           gidx.PrefixedID "json:\"id\" graphql:\"id\""
					Number       int64           "json:\"number\" graphql:\"number\""
					LoadBalancer struct {
						ID gidx.PrefixedID "json:\"id\" graphql:\"id\""
					} "json:\"loadBalancer\" graphql:\"loadBalancer\""
					CreatedAt time.Time "json:\"createdAt\" graphql:\"createdAt\""
					UpdatedAt time.Time "json:\"updatedAt\" graphql:\"updatedAt\""
				} "json:\"node\" graphql:\"node\""
			} "json:\"edges\" graphql:\"edges\""
		} "json:\"ports\" graphql:\"ports\""
	} "json:\"loadBalancer\" graphql:\"loadBalancer\""
}
type GetLoadBalancerProvider struct {
	LoadBalancerProvider struct {
		ID    gidx.PrefixedID "json:\"id\" graphql:\"id\""
		Name  string          "json:\"name\" graphql:\"name\""
		Owner struct {
			ID gidx.PrefixedID "json:\"id\" graphql:\"id\""
		} "json:\"owner\" graphql:\"owner\""
		CreatedAt time.Time "json:\"createdAt\" graphql:\"createdAt\""
		UpdatedAt time.Time "json:\"updatedAt\" graphql:\"updatedAt\""
	} "json:\"loadBalancerProvider\" graphql:\"loadBalancerProvider\""
}
type GetOwnerLoadBalancers struct {
	Entities []*struct {
		LoadBalancers struct {
			Edges []*struct {
				Node *struct {
					ID   gidx.PrefixedID "json:\"id\" graphql:\"id\""
					Name string          "json:\"name\" graphql:\"name\""
				} "json:\"node\" graphql:\"node\""
			} "json:\"edges\" graphql:\"edges\""
		} "json:\"loadBalancers\" graphql:\"loadBalancers\""
	} "json:\"_entities\" graphql:\"_entities\""
}
type LoadBalancerCreate struct {
	LoadBalancerCreate struct {
		LoadBalancer struct {
			ID                   gidx.PrefixedID "json:\"id\" graphql:\"id\""
			Name                 string          "json:\"name\" graphql:\"name\""
			LoadBalancerProvider struct {
				ID gidx.PrefixedID "json:\"id\" graphql:\"id\""
			} "json:\"loadBalancerProvider\" graphql:\"loadBalancerProvider\""
			Owner struct {
				ID gidx.PrefixedID "json:\"id\" graphql:\"id\""
			} "json:\"owner\" graphql:\"owner\""
			Location struct {
				ID gidx.PrefixedID "json:\"id\" graphql:\"id\""
			} "json:\"location\" graphql:\"location\""
			CreatedAt time.Time "json:\"createdAt\" graphql:\"createdAt\""
			UpdatedAt time.Time "json:\"updatedAt\" graphql:\"updatedAt\""
		} "json:\"loadBalancer\" graphql:\"loadBalancer\""
	} "json:\"loadBalancerCreate\" graphql:\"loadBalancerCreate\""
}
type LoadBalancerDelete struct {
	LoadBalancerDelete struct {
		DeletedID gidx.PrefixedID "json:\"deletedID\" graphql:\"deletedID\""
	} "json:\"loadBalancerDelete\" graphql:\"loadBalancerDelete\""
}
type LoadBalancerOriginCreate struct {
	LoadBalancerOriginCreate struct {
		LoadBalancerOrigin struct {
			ID         gidx.PrefixedID "json:\"id\" graphql:\"id\""
			Active     bool            "json:\"active\" graphql:\"active\""
			Name       string          "json:\"name\" graphql:\"name\""
			PortNumber int64           "json:\"portNumber\" graphql:\"portNumber\""
			Target     string          "json:\"target\" graphql:\"target\""
			Weight     int64           "json:\"weight\" graphql:\"weight\""
			PoolID     gidx.PrefixedID "json:\"poolID\" graphql:\"poolID\""
			CreatedAt  time.Time       "json:\"createdAt\" graphql:\"createdAt\""
			UpdatedAt  time.Time       "json:\"updatedAt\" graphql:\"updatedAt\""
		} "json:\"loadBalancerOrigin\" graphql:\"loadBalancerOrigin\""
	} "json:\"loadBalancerOriginCreate\" graphql:\"loadBalancerOriginCreate\""
}
type LoadBalancerOriginDelete struct {
	LoadBalancerOriginDelete struct {
		DeletedID gidx.PrefixedID "json:\"deletedID\" graphql:\"deletedID\""
	} "json:\"loadBalancerOriginDelete\" graphql:\"loadBalancerOriginDelete\""
}
type LoadBalancerOriginUpdate struct {
	LoadBalancerOriginUpdate struct {
		LoadBalancerOrigin struct {
			ID         gidx.PrefixedID "json:\"id\" graphql:\"id\""
			Active     bool            "json:\"active\" graphql:\"active\""
			Name       string          "json:\"name\" graphql:\"name\""
			PortNumber int64           "json:\"portNumber\" graphql:\"portNumber\""
			Target     string          "json:\"target\" graphql:\"target\""
			Weight     int64           "json:\"weight\" graphql:\"weight\""
			PoolID     gidx.PrefixedID "json:\"poolID\" graphql:\"poolID\""
			CreatedAt  time.Time       "json:\"createdAt\" graphql:\"createdAt\""
			UpdatedAt  time.Time       "json:\"updatedAt\" graphql:\"updatedAt\""
		} "json:\"loadBalancerOrigin\" graphql:\"loadBalancerOrigin\""
	} "json:\"loadBalancerOriginUpdate\" graphql:\"loadBalancerOriginUpdate\""
}
type LoadBalancerPoolCreate struct {
	LoadBalancerPoolCreate struct {
		LoadBalancerPool struct {
			ID        gidx.PrefixedID          "json:\"id\" graphql:\"id\""
			Name      string                   "json:\"name\" graphql:\"name\""
			Protocol  LoadBalancerPoolProtocol "json:\"protocol\" graphql:\"protocol\""
			OwnerID   gidx.PrefixedID          "json:\"ownerID\" graphql:\"ownerID\""
			CreatedAt time.Time                "json:\"createdAt\" graphql:\"createdAt\""
			UpdatedAt time.Time                "json:\"updatedAt\" graphql:\"updatedAt\""
		} "json:\"loadBalancerPool\" graphql:\"loadBalancerPool\""
	} "json:\"loadBalancerPoolCreate\" graphql:\"loadBalancerPoolCreate\""
}
type LoadBalancerPoolDelete struct {
	LoadBalancerPoolDelete struct {
		DeletedID *gidx.PrefixedID "json:\"deletedID\" graphql:\"deletedID\""
	} "json:\"loadBalancerPoolDelete\" graphql:\"loadBalancerPoolDelete\""
}
type LoadBalancerPoolUpdate struct {
	LoadBalancerPoolUpdate struct {
		LoadBalancerPool struct {
			ID        gidx.PrefixedID          "json:\"id\" graphql:\"id\""
			Name      string                   "json:\"name\" graphql:\"name\""
			Protocol  LoadBalancerPoolProtocol "json:\"protocol\" graphql:\"protocol\""
			OwnerID   gidx.PrefixedID          "json:\"ownerID\" graphql:\"ownerID\""
			CreatedAt time.Time                "json:\"createdAt\" graphql:\"createdAt\""
			UpdatedAt time.Time                "json:\"updatedAt\" graphql:\"updatedAt\""
		} "json:\"loadBalancerPool\" graphql:\"loadBalancerPool\""
	} "json:\"loadBalancerPoolUpdate\" graphql:\"loadBalancerPoolUpdate\""
}
type LoadBalancerPortCreate struct {
	LoadBalancerPortCreate struct {
		LoadBalancerPort struct {
			ID           gidx.PrefixedID "json:\"id\" graphql:\"id\""
			Name         *string         "json:\"name\" graphql:\"name\""
			Number       int64           "json:\"number\" graphql:\"number\""
			LoadBalancer struct {
				ID gidx.PrefixedID "json:\"id\" graphql:\"id\""
			} "json:\"loadBalancer\" graphql:\"loadBalancer\""
			CreatedAt time.Time "json:\"createdAt\" graphql:\"createdAt\""
			UpdatedAt time.Time "json:\"updatedAt\" graphql:\"updatedAt\""
		} "json:\"loadBalancerPort\" graphql:\"loadBalancerPort\""
	} "json:\"loadBalancerPortCreate\" graphql:\"loadBalancerPortCreate\""
}
type LoadBalancerPortDelete struct {
	LoadBalancerPortDelete struct {
		DeletedID gidx.PrefixedID "json:\"deletedID\" graphql:\"deletedID\""
	} "json:\"loadBalancerPortDelete\" graphql:\"loadBalancerPortDelete\""
}
type LoadBalancerPortUpdate struct {
	LoadBalancerPortUpdate struct {
		LoadBalancerPort struct {
			ID        gidx.PrefixedID "json:\"id\" graphql:\"id\""
			Name      *string         "json:\"name\" graphql:\"name\""
			Number    int64           "json:\"number\" graphql:\"number\""
			CreatedAt time.Time       "json:\"createdAt\" graphql:\"createdAt\""
			UpdatedAt time.Time       "json:\"updatedAt\" graphql:\"updatedAt\""
		} "json:\"loadBalancerPort\" graphql:\"loadBalancerPort\""
	} "json:\"loadBalancerPortUpdate\" graphql:\"loadBalancerPortUpdate\""
}
type LoadBalancerProviderCreate struct {
	LoadBalancerProviderCreate struct {
		LoadBalancerProvider struct {
			ID    gidx.PrefixedID "json:\"id\" graphql:\"id\""
			Name  string          "json:\"name\" graphql:\"name\""
			Owner struct {
				ID gidx.PrefixedID "json:\"id\" graphql:\"id\""
			} "json:\"owner\" graphql:\"owner\""
			CreatedAt time.Time "json:\"createdAt\" graphql:\"createdAt\""
			UpdatedAt time.Time "json:\"updatedAt\" graphql:\"updatedAt\""
		} "json:\"loadBalancerProvider\" graphql:\"loadBalancerProvider\""
	} "json:\"loadBalancerProviderCreate\" graphql:\"loadBalancerProviderCreate\""
}
type LoadBalancerProviderDelete struct {
	LoadBalancerProviderDelete struct {
		DeletedID gidx.PrefixedID "json:\"deletedID\" graphql:\"deletedID\""
	} "json:\"loadBalancerProviderDelete\" graphql:\"loadBalancerProviderDelete\""
}
type LoadBalancerProviderUpdate struct {
	LoadBalancerProviderUpdate struct {
		LoadBalancerProvider struct {
			ID        gidx.PrefixedID "json:\"id\" graphql:\"id\""
			Name      string          "json:\"name\" graphql:\"name\""
			CreatedAt time.Time       "json:\"createdAt\" graphql:\"createdAt\""
			UpdatedAt time.Time       "json:\"updatedAt\" graphql:\"updatedAt\""
		} "json:\"loadBalancerProvider\" graphql:\"loadBalancerProvider\""
	} "json:\"loadBalancerProviderUpdate\" graphql:\"loadBalancerProviderUpdate\""
}
type LoadBalancerUpdate struct {
	LoadBalancerUpdate struct {
		LoadBalancer struct {
			ID        gidx.PrefixedID "json:\"id\" graphql:\"id\""
			Name      string          "json:\"name\" graphql:\"name\""
			CreatedAt time.Time       "json:\"createdAt\" graphql:\"createdAt\""
			UpdatedAt time.Time       "json:\"updatedAt\" graphql:\"updatedAt\""
		} "json:\"loadBalancer\" graphql:\"loadBalancer\""
	} "json:\"loadBalancerUpdate\" graphql:\"loadBalancerUpdate\""
}

const GetLoadBalancerDocument = `query GetLoadBalancer ($id: ID!) {
	loadBalancer(id: $id) {
		id
		name
		location {
			id
		}
		loadBalancerProvider {
			id
		}
		owner {
			id
		}
		createdAt
		updatedAt
	}
}
`

func (c *Client) GetLoadBalancer(ctx context.Context, id gidx.PrefixedID, httpRequestOptions ...client.HTTPRequestOption) (*GetLoadBalancer, error) {
	vars := map[string]interface{}{
		"id": id,
	}

	var res GetLoadBalancer
	if err := c.Client.Post(ctx, "GetLoadBalancer", GetLoadBalancerDocument, &res, vars, httpRequestOptions...); err != nil {
		return nil, err
	}

	return &res, nil
}

const GetLoadBalancerPoolDocument = `query GetLoadBalancerPool ($id: ID!) {
	loadBalancerPool(id: $id) {
		id
		name
		protocol
		ownerID
		createdAt
		updatedAt
	}
}
`

func (c *Client) GetLoadBalancerPool(ctx context.Context, id gidx.PrefixedID, httpRequestOptions ...client.HTTPRequestOption) (*GetLoadBalancerPool, error) {
	vars := map[string]interface{}{
		"id": id,
	}

	var res GetLoadBalancerPool
	if err := c.Client.Post(ctx, "GetLoadBalancerPool", GetLoadBalancerPoolDocument, &res, vars, httpRequestOptions...); err != nil {
		return nil, err
	}

	return &res, nil
}

const GetLoadBalancerPoolOriginDocument = `query GetLoadBalancerPoolOrigin ($id: ID!, $originid: ID!) {
	loadBalancerPool(id: $id) {
		origins(where: {id:$originid}) {
			edges {
				node {
					id
					name
					target
					portNumber
					active
					weight
					poolID
					createdAt
					updatedAt
				}
			}
		}
	}
}
`

func (c *Client) GetLoadBalancerPoolOrigin(ctx context.Context, id gidx.PrefixedID, originid gidx.PrefixedID, httpRequestOptions ...client.HTTPRequestOption) (*GetLoadBalancerPoolOrigin, error) {
	vars := map[string]interface{}{
		"id":       id,
		"originid": originid,
	}

	var res GetLoadBalancerPoolOrigin
	if err := c.Client.Post(ctx, "GetLoadBalancerPoolOrigin", GetLoadBalancerPoolOriginDocument, &res, vars, httpRequestOptions...); err != nil {
		return nil, err
	}

	return &res, nil
}

const GetLoadBalancerPortDocument = `query GetLoadBalancerPort ($id: ID!, $portid: ID!) {
	loadBalancer(id: $id) {
		ports(where: {id:$portid}) {
			edges {
				node {
					id
					number
					loadBalancer {
						id
					}
					createdAt
					updatedAt
				}
			}
		}
	}
}
`

func (c *Client) GetLoadBalancerPort(ctx context.Context, id gidx.PrefixedID, portid gidx.PrefixedID, httpRequestOptions ...client.HTTPRequestOption) (*GetLoadBalancerPort, error) {
	vars := map[string]interface{}{
		"id":     id,
		"portid": portid,
	}

	var res GetLoadBalancerPort
	if err := c.Client.Post(ctx, "GetLoadBalancerPort", GetLoadBalancerPortDocument, &res, vars, httpRequestOptions...); err != nil {
		return nil, err
	}

	return &res, nil
}

const GetLoadBalancerProviderDocument = `query GetLoadBalancerProvider ($id: ID!) {
	loadBalancerProvider(id: $id) {
		id
		name
		owner {
			id
		}
		createdAt
		updatedAt
	}
}
`

func (c *Client) GetLoadBalancerProvider(ctx context.Context, id gidx.PrefixedID, httpRequestOptions ...client.HTTPRequestOption) (*GetLoadBalancerProvider, error) {
	vars := map[string]interface{}{
		"id": id,
	}

	var res GetLoadBalancerProvider
	if err := c.Client.Post(ctx, "GetLoadBalancerProvider", GetLoadBalancerProviderDocument, &res, vars, httpRequestOptions...); err != nil {
		return nil, err
	}

	return &res, nil
}

const GetOwnerLoadBalancersDocument = `query GetOwnerLoadBalancers ($id: ID!, $orderBy: LoadBalancerOrder) {
	_entities(representations: {__typename:"ResourceOwner",id:$id}) {
		... on ResourceOwner {
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
}
`

func (c *Client) GetOwnerLoadBalancers(ctx context.Context, id gidx.PrefixedID, orderBy *LoadBalancerOrder, httpRequestOptions ...client.HTTPRequestOption) (*GetOwnerLoadBalancers, error) {
	vars := map[string]interface{}{
		"id":      id,
		"orderBy": orderBy,
	}

	var res GetOwnerLoadBalancers
	if err := c.Client.Post(ctx, "GetOwnerLoadBalancers", GetOwnerLoadBalancersDocument, &res, vars, httpRequestOptions...); err != nil {
		return nil, err
	}

	return &res, nil
}

const LoadBalancerCreateDocument = `mutation LoadBalancerCreate ($input: CreateLoadBalancerInput!) {
	loadBalancerCreate(input: $input) {
		loadBalancer {
			id
			name
			loadBalancerProvider {
				id
			}
			owner {
				id
			}
			location {
				id
			}
			createdAt
			updatedAt
		}
	}
}
`

func (c *Client) LoadBalancerCreate(ctx context.Context, input CreateLoadBalancerInput, httpRequestOptions ...client.HTTPRequestOption) (*LoadBalancerCreate, error) {
	vars := map[string]interface{}{
		"input": input,
	}

	var res LoadBalancerCreate
	if err := c.Client.Post(ctx, "LoadBalancerCreate", LoadBalancerCreateDocument, &res, vars, httpRequestOptions...); err != nil {
		return nil, err
	}

	return &res, nil
}

const LoadBalancerDeleteDocument = `mutation LoadBalancerDelete ($id: ID!) {
	loadBalancerDelete(id: $id) {
		deletedID
	}
}
`

func (c *Client) LoadBalancerDelete(ctx context.Context, id gidx.PrefixedID, httpRequestOptions ...client.HTTPRequestOption) (*LoadBalancerDelete, error) {
	vars := map[string]interface{}{
		"id": id,
	}

	var res LoadBalancerDelete
	if err := c.Client.Post(ctx, "LoadBalancerDelete", LoadBalancerDeleteDocument, &res, vars, httpRequestOptions...); err != nil {
		return nil, err
	}

	return &res, nil
}

const LoadBalancerOriginCreateDocument = `mutation LoadBalancerOriginCreate ($input: CreateLoadBalancerOriginInput!) {
	loadBalancerOriginCreate(input: $input) {
		loadBalancerOrigin {
			id
			active
			name
			portNumber
			target
			weight
			poolID
			createdAt
			updatedAt
		}
	}
}
`

func (c *Client) LoadBalancerOriginCreate(ctx context.Context, input CreateLoadBalancerOriginInput, httpRequestOptions ...client.HTTPRequestOption) (*LoadBalancerOriginCreate, error) {
	vars := map[string]interface{}{
		"input": input,
	}

	var res LoadBalancerOriginCreate
	if err := c.Client.Post(ctx, "LoadBalancerOriginCreate", LoadBalancerOriginCreateDocument, &res, vars, httpRequestOptions...); err != nil {
		return nil, err
	}

	return &res, nil
}

const LoadBalancerOriginDeleteDocument = `mutation LoadBalancerOriginDelete ($id: ID!) {
	loadBalancerOriginDelete(id: $id) {
		deletedID
	}
}
`

func (c *Client) LoadBalancerOriginDelete(ctx context.Context, id gidx.PrefixedID, httpRequestOptions ...client.HTTPRequestOption) (*LoadBalancerOriginDelete, error) {
	vars := map[string]interface{}{
		"id": id,
	}

	var res LoadBalancerOriginDelete
	if err := c.Client.Post(ctx, "LoadBalancerOriginDelete", LoadBalancerOriginDeleteDocument, &res, vars, httpRequestOptions...); err != nil {
		return nil, err
	}

	return &res, nil
}

const LoadBalancerOriginUpdateDocument = `mutation LoadBalancerOriginUpdate ($id: ID!, $input: UpdateLoadBalancerOriginInput!) {
	loadBalancerOriginUpdate(id: $id, input: $input) {
		loadBalancerOrigin {
			id
			active
			name
			portNumber
			target
			weight
			poolID
			createdAt
			updatedAt
		}
	}
}
`

func (c *Client) LoadBalancerOriginUpdate(ctx context.Context, id gidx.PrefixedID, input UpdateLoadBalancerOriginInput, httpRequestOptions ...client.HTTPRequestOption) (*LoadBalancerOriginUpdate, error) {
	vars := map[string]interface{}{
		"id":    id,
		"input": input,
	}

	var res LoadBalancerOriginUpdate
	if err := c.Client.Post(ctx, "LoadBalancerOriginUpdate", LoadBalancerOriginUpdateDocument, &res, vars, httpRequestOptions...); err != nil {
		return nil, err
	}

	return &res, nil
}

const LoadBalancerPoolCreateDocument = `mutation LoadBalancerPoolCreate ($input: CreateLoadBalancerPoolInput!) {
	loadBalancerPoolCreate(input: $input) {
		loadBalancerPool {
			id
			name
			protocol
			ownerID
			createdAt
			updatedAt
		}
	}
}
`

func (c *Client) LoadBalancerPoolCreate(ctx context.Context, input CreateLoadBalancerPoolInput, httpRequestOptions ...client.HTTPRequestOption) (*LoadBalancerPoolCreate, error) {
	vars := map[string]interface{}{
		"input": input,
	}

	var res LoadBalancerPoolCreate
	if err := c.Client.Post(ctx, "LoadBalancerPoolCreate", LoadBalancerPoolCreateDocument, &res, vars, httpRequestOptions...); err != nil {
		return nil, err
	}

	return &res, nil
}

const LoadBalancerPoolDeleteDocument = `mutation LoadBalancerPoolDelete ($id: ID!) {
	loadBalancerPoolDelete(id: $id) {
		deletedID
	}
}
`

func (c *Client) LoadBalancerPoolDelete(ctx context.Context, id gidx.PrefixedID, httpRequestOptions ...client.HTTPRequestOption) (*LoadBalancerPoolDelete, error) {
	vars := map[string]interface{}{
		"id": id,
	}

	var res LoadBalancerPoolDelete
	if err := c.Client.Post(ctx, "LoadBalancerPoolDelete", LoadBalancerPoolDeleteDocument, &res, vars, httpRequestOptions...); err != nil {
		return nil, err
	}

	return &res, nil
}

const LoadBalancerPoolUpdateDocument = `mutation LoadBalancerPoolUpdate ($id: ID!, $input: UpdateLoadBalancerPoolInput!) {
	loadBalancerPoolUpdate(id: $id, input: $input) {
		loadBalancerPool {
			id
			name
			protocol
			ownerID
			createdAt
			updatedAt
		}
	}
}
`

func (c *Client) LoadBalancerPoolUpdate(ctx context.Context, id gidx.PrefixedID, input UpdateLoadBalancerPoolInput, httpRequestOptions ...client.HTTPRequestOption) (*LoadBalancerPoolUpdate, error) {
	vars := map[string]interface{}{
		"id":    id,
		"input": input,
	}

	var res LoadBalancerPoolUpdate
	if err := c.Client.Post(ctx, "LoadBalancerPoolUpdate", LoadBalancerPoolUpdateDocument, &res, vars, httpRequestOptions...); err != nil {
		return nil, err
	}

	return &res, nil
}

const LoadBalancerPortCreateDocument = `mutation LoadBalancerPortCreate ($input: CreateLoadBalancerPortInput!) {
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
}
`

func (c *Client) LoadBalancerPortCreate(ctx context.Context, input CreateLoadBalancerPortInput, httpRequestOptions ...client.HTTPRequestOption) (*LoadBalancerPortCreate, error) {
	vars := map[string]interface{}{
		"input": input,
	}

	var res LoadBalancerPortCreate
	if err := c.Client.Post(ctx, "LoadBalancerPortCreate", LoadBalancerPortCreateDocument, &res, vars, httpRequestOptions...); err != nil {
		return nil, err
	}

	return &res, nil
}

const LoadBalancerPortDeleteDocument = `mutation LoadBalancerPortDelete ($id: ID!) {
	loadBalancerPortDelete(id: $id) {
		deletedID
	}
}
`

func (c *Client) LoadBalancerPortDelete(ctx context.Context, id gidx.PrefixedID, httpRequestOptions ...client.HTTPRequestOption) (*LoadBalancerPortDelete, error) {
	vars := map[string]interface{}{
		"id": id,
	}

	var res LoadBalancerPortDelete
	if err := c.Client.Post(ctx, "LoadBalancerPortDelete", LoadBalancerPortDeleteDocument, &res, vars, httpRequestOptions...); err != nil {
		return nil, err
	}

	return &res, nil
}

const LoadBalancerPortUpdateDocument = `mutation LoadBalancerPortUpdate ($id: ID!, $input: UpdateLoadBalancerPortInput!) {
	loadBalancerPortUpdate(id: $id, input: $input) {
		loadBalancerPort {
			id
			name
			number
			createdAt
			updatedAt
		}
	}
}
`

func (c *Client) LoadBalancerPortUpdate(ctx context.Context, id gidx.PrefixedID, input UpdateLoadBalancerPortInput, httpRequestOptions ...client.HTTPRequestOption) (*LoadBalancerPortUpdate, error) {
	vars := map[string]interface{}{
		"id":    id,
		"input": input,
	}

	var res LoadBalancerPortUpdate
	if err := c.Client.Post(ctx, "LoadBalancerPortUpdate", LoadBalancerPortUpdateDocument, &res, vars, httpRequestOptions...); err != nil {
		return nil, err
	}

	return &res, nil
}

const LoadBalancerProviderCreateDocument = `mutation LoadBalancerProviderCreate ($input: CreateLoadBalancerProviderInput!) {
	loadBalancerProviderCreate(input: $input) {
		loadBalancerProvider {
			id
			name
			owner {
				id
			}
			createdAt
			updatedAt
		}
	}
}
`

func (c *Client) LoadBalancerProviderCreate(ctx context.Context, input CreateLoadBalancerProviderInput, httpRequestOptions ...client.HTTPRequestOption) (*LoadBalancerProviderCreate, error) {
	vars := map[string]interface{}{
		"input": input,
	}

	var res LoadBalancerProviderCreate
	if err := c.Client.Post(ctx, "LoadBalancerProviderCreate", LoadBalancerProviderCreateDocument, &res, vars, httpRequestOptions...); err != nil {
		return nil, err
	}

	return &res, nil
}

const LoadBalancerProviderDeleteDocument = `mutation LoadBalancerProviderDelete ($id: ID!) {
	loadBalancerProviderDelete(id: $id) {
		deletedID
	}
}
`

func (c *Client) LoadBalancerProviderDelete(ctx context.Context, id gidx.PrefixedID, httpRequestOptions ...client.HTTPRequestOption) (*LoadBalancerProviderDelete, error) {
	vars := map[string]interface{}{
		"id": id,
	}

	var res LoadBalancerProviderDelete
	if err := c.Client.Post(ctx, "LoadBalancerProviderDelete", LoadBalancerProviderDeleteDocument, &res, vars, httpRequestOptions...); err != nil {
		return nil, err
	}

	return &res, nil
}

const LoadBalancerProviderUpdateDocument = `mutation LoadBalancerProviderUpdate ($id: ID!, $input: UpdateLoadBalancerProviderInput!) {
	loadBalancerProviderUpdate(id: $id, input: $input) {
		loadBalancerProvider {
			id
			name
			createdAt
			updatedAt
		}
	}
}
`

func (c *Client) LoadBalancerProviderUpdate(ctx context.Context, id gidx.PrefixedID, input UpdateLoadBalancerProviderInput, httpRequestOptions ...client.HTTPRequestOption) (*LoadBalancerProviderUpdate, error) {
	vars := map[string]interface{}{
		"id":    id,
		"input": input,
	}

	var res LoadBalancerProviderUpdate
	if err := c.Client.Post(ctx, "LoadBalancerProviderUpdate", LoadBalancerProviderUpdateDocument, &res, vars, httpRequestOptions...); err != nil {
		return nil, err
	}

	return &res, nil
}

const LoadBalancerUpdateDocument = `mutation LoadBalancerUpdate ($id: ID!, $input: UpdateLoadBalancerInput!) {
	loadBalancerUpdate(id: $id, input: $input) {
		loadBalancer {
			id
			name
			createdAt
			updatedAt
		}
	}
}
`

func (c *Client) LoadBalancerUpdate(ctx context.Context, id gidx.PrefixedID, input UpdateLoadBalancerInput, httpRequestOptions ...client.HTTPRequestOption) (*LoadBalancerUpdate, error) {
	vars := map[string]interface{}{
		"id":    id,
		"input": input,
	}

	var res LoadBalancerUpdate
	if err := c.Client.Post(ctx, "LoadBalancerUpdate", LoadBalancerUpdateDocument, &res, vars, httpRequestOptions...); err != nil {
		return nil, err
	}

	return &res, nil
}
