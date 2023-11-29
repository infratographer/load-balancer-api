package client

import "encoding/json"

// OriginNode is a struct that represents the OriginNode GraphQL type
type OriginNode struct {
	ID         string `graphql:"id" json:"id"`
	Name       string `graphql:"name" json:"name"`
	Target     string `graphql:"target" json:"target"`
	PortNumber int64  `graphql:"portNumber" json:"portNumber"`
	Weight     int64  `graphql:"weight" json:"weight"`
	Active     bool   `graphql:"active" json:"active"`
}

// OriginEdges is a struct that represents the OriginEdges GraphQL type
type OriginEdges struct {
	Node OriginNode `graphql:"node" json:"node"`
}

// Origins is a struct that represents the Origins GraphQL type
type Origins struct {
	Edges []OriginEdges `graphql:"edges" json:"edges"`
}

// Pool is a struct that represents the Pool GraphQL type
type Pool struct {
	ID       string  `graphql:"id"`
	Name     string  `graphql:"name" json:"name"`
	Protocol string  `graphql:"protocol" json:"protocol"`
	Origins  Origins `graphql:"origins" json:"origins"`
}

// PortNode is a struct that represents the PortNode GraphQL type
type PortNode struct {
	ID     string `graphql:"id" json:"id"`
	Name   string `graphql:"name" json:"name"`
	Number int64  `graphql:"number" json:"number"`
	Pools  []Pool `graphql:"pools" json:"pools"`
}

// PortEdges is a struct that represents the PortEdges GraphQL type
type PortEdges struct {
	Node PortNode `graphql:"node" json:"node"`
}

// Ports is a struct that represents the Ports GraphQL type
type Ports struct {
	Edges []PortEdges `graphql:"edges" json:"edges"`
}

// OwnerNode is a struct that represents the OwnerNode GraphQL type
type OwnerNode struct {
	ID string `graphql:"id" json:"id"`
}

// LocationNode is a struct that represents the LocationNode GraphQL type
type LocationNode struct {
	ID string `graphql:"id" json:"id"`
}

// LoadBalancer is a struct that represents the LoadBalancer GraphQL type
type LoadBalancer struct {
	ID          string       `graphql:"id" json:"id"`
	Name        string       `graphql:"name" json:"name"`
	Owner       OwnerNode    `graphql:"owner" json:"owner"`
	Location    LocationNode `graphql:"location" json:"location"`
	IPAddresses []IPAddress  `graphql:"IPAddresses" json:"IPAddresses"`
	Metadata    Metadata     `graphql:"metadata" json:"metadata"`
	Ports       Ports        `graphql:"ports" json:"ports"`
}

// GetLoadBalancer is a struct that represents the GetLoadBalancer GraphQL query
type GetLoadBalancer struct {
	LoadBalancer LoadBalancer `graphql:"loadBalancer(id: $id)"`
}

// IPAddress is a struct that represents the IPAddress GraphQL type
type IPAddress struct {
	ID       string `graphql:"id" json:"id"`
	IP       string `graphql:"ip" json:"ip"`
	Reserved bool   `graphql:"reserved" json:"reserved"`
}

// MetadataStatusNode is a struct that represents the Metadata status node GraphQL type
type MetadataStatusNode struct {
	ID                string          `graphql:"id" json:"id"`
	Data              json.RawMessage `graphql:"data"`
	Source            string          `graphql:"source" json:"source"`
	StatusNamespaceID string          `graphql:"statusNamespaceID" json:"statusNamespaceID"`
}

// MetadataStatusEdges is a struct that represents the Metadata status edges GraphQL type
type MetadataStatusEdges struct {
	Node MetadataStatusNode `graphql:"node" json:"node"`
}

// MetadataStatuses is a struct that represents the Metadata statuses GraphQL type
type MetadataStatuses struct {
	TotalCount int                   `graphql:"totalCount" json:"totalCount"`
	Edges      []MetadataStatusEdges `graphql:"edges" json:"edges"`
}

// Metadata is a struct that represents the metadata GraphQL type
type Metadata struct {
	ID       string           `graphql:"id" json:"id"`
	NodeID   string           `graphql:"nodeID" json:"nodeID"`
	Statuses MetadataStatuses `graphql:"statuses" json:"statuses"`
}

// MetadataNodeFragment is a struct that represents the MetadataNodeFragment GraphQL fragment
type MetadataNodeFragment struct {
	Metadata Metadata `graphql:"metadata" json:"metadata"`
}

// MetadataNode is a struct that represents the MetadataNode GraphQL type
type MetadataNode struct {
	MetadataNodeFragment `graphql:"... on MetadataNode"`
}

// GetMetadataNode is a struct that represents the node-resolver subgraph query
type GetMetadataNode struct {
	MetadataNode MetadataNode `graphql:"node(id: $id)"`
}

// Readable version of the above:
// type GetLoadBalancer struct {
// 	LoadBalancer struct {
// 		ID    string
//      Owner string
// 		Name  string
//      IPAddresses {
// 	      id string
// 	      ip string
//  	}
// 	    metadata struct {
// 	      id     string
// 	      nodeID string
// 	      statuses struct {
// 		    edges []struct {
// 		      node struct {
// 		        source            string
// 		        statusNamespaceID string
// 		        id                string
// 		        data              json bytes
// 		      }
// 		    }
// 	      }
// 	    }
// 		Ports struct {
// 			Edges []struct {
// 				Node struct {
// 					Name   string
// 					Number int64
// 					Pools  []struct {
// 						Name     string
// 						Protocol string
// 						Origins  struct {
// 							Edges []struct {
// 								Node struct {
// 									Name       string
// 									Target     string
// 									PortNumber int64
// 									Active     bool
// 								}
// 							}
// 						}
// 					}
// 				}
// 			}
// 		}
// 	} `graphql:"loadBalancer(id: $id)"`
// }
