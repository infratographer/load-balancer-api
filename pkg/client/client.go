package client

import (
	"context"
	"net/http"
	"strings"

	graphql "github.com/hasura/go-graphql-client"
	"go.infratographer.com/x/gidx"
)

// GQLClient is an interface for a graphql client
type GQLClient interface {
	Query(tx context.Context, q interface{}, variables map[string]interface{}, options ...graphql.Option) error
}

// Client creates a new lb api client against a specific endpoint
type Client struct {
	gqlCli     GQLClient
	httpClient *http.Client
}

// Option is a function that modifies a client
type Option func(*Client)

// NewClient creates a new lb api client
func NewClient(url string, opts ...Option) *Client {
	c := &Client{
		httpClient: http.DefaultClient,
	}

	for _, opt := range opts {
		opt(c)
	}

	c.gqlCli = graphql.NewClient(url, c.httpClient)

	return c
}

// WithHTTPClient functional option to set the http client
func WithHTTPClient(cli *http.Client) Option {
	return func(c *Client) {
		c.httpClient = cli
	}
}

// GetLoadBalancer returns a load balancer by id
func (c Client) GetLoadBalancer(ctx context.Context, id string) (*LoadBalancer, error) {
	_, err := gidx.Parse(id)
	if err != nil {
		return nil, err
	}

	vars := map[string]interface{}{
		"id": graphql.ID(id),
	}

	var q GetLoadBalancer
	if err := c.gqlCli.Query(ctx, &q, vars); err != nil {
		return nil, translateGQLErr(err)
	}

	return &q.LoadBalancer, nil
}

// NodeMetadata returns the metadata-api subgraph node for a load balancer
// Once a load balancer is deleted, it is fully removed from load-balancer-api.
// There are no soft-deletes due to permissions complications.
// However, it's metadata remains to query via the node-resolver subgraph.
func (c Client) NodeMetadata(ctx context.Context, id string) (*Metadata, error) {
	//	query Test {
	//	  node(id:"loadbal-RelHb0FU59uEjlhtT7MAl") {
	//	    ... on MetadataNode {
	//	      metadata {
	//	        statuses {
	//	          totalCount
	//	          edges {
	//	            node {
	//	              data
	//	            }
	//	          }
	//	        }
	//	      }
	//	    }
	//	  }
	//	}

	_, err := gidx.Parse(id)
	if err != nil {
		return nil, err
	}

	vars := map[string]interface{}{
		"id": graphql.ID(id),
	}

	var q GetMetadataNode
	if err := c.gqlCli.Query(ctx, &q, vars); err != nil {
		return nil, translateGQLErr(err)
	}

	if q.MetadataNode.Metadata.Statuses.TotalCount == 0 {
		return nil, ErrMetadataStatusNotFound
	}

	return &q.MetadataNode.Metadata, nil
}

func translateGQLErr(err error) error {
	switch {
	case strings.Contains(err.Error(), "load_balancer not found"):
		return ErrLBNotfound
	case strings.Contains(err.Error(), "invalid or expired jwt"):
		return ErrUnauthorized
	case strings.Contains(err.Error(), "subject doesn't have access"):
		return ErrPermissionDenied
	case strings.Contains(err.Error(), "internal server error"):
		return ErrInternalServerError
	}

	return err
}
