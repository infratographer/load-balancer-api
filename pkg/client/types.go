package client

import "encoding/json"

// OriginNode is a struct that represents the OriginNode GraphQL type
type OriginNode struct {
	ID         string
	Name       string
	Target     string
	PortNumber int64
	Weight     int64
	Active     bool
}

// OriginEdges is a struct that represents the OriginEdges GraphQL type
type OriginEdges struct {
	Node OriginNode
}

// Origins is a struct that represents the Origins GraphQL type
type Origins struct {
	Edges []OriginEdges
}

// Pool is a struct that represents the Pool GraphQL type
type Pool struct {
	ID       string
	Name     string
	Protocol string
	Origins  Origins
}

// PortNode is a struct that represents the PortNode GraphQL type
type PortNode struct {
	ID     string
	Name   string
	Number int64
	Pools  []Pool
}

// PortEdges is a struct that represents the PortEdges GraphQL type
type PortEdges struct {
	Node PortNode
}

// Ports is a struct that represents the Ports GraphQL type
type Ports struct {
	Edges []PortEdges
}

// OwnerNode is a struct that represents the OwnerNode GraphQL type
type OwnerNode struct {
	ID string
}

// LocationNode is a struct that represents the LocationNode GraphQL type
type LocationNode struct {
	ID string
}

// LoadBalancer is a struct that represents the LoadBalancer GraphQL type
type LoadBalancer struct {
	ID          string
	Name        string
	Owner       OwnerNode
	Location    LocationNode
	IPAddresses []IPAddress `graphql:"IPAddresses" json:"IPAddresses"`
	Metadata    Metadata    `graphql:"metadata" json:"metadata"`
	Ports       Ports
}

// GetLoadBalancer is a struct that represents the GetLoadBalancer GraphQL query
type GetLoadBalancer struct {
	LoadBalancer LoadBalancer `graphql:"loadBalancer(id: $id)"`
}

// IPAddress is a struct that represents the IPAddress GraphQL type
type IPAddress struct {
	ID       string
	IP       string
	Reserved bool
}

// MetadataStatusNode is a struct that represents the Metadata status node GraphQL type
type MetadataStatusNode struct {
	ID                string          `graphql:"id"`
	Data              json.RawMessage `graphql:"data"`
	Source            string          `graphql:"source"`
	StatusNamespaceID string          `graphql:"statusNamespaceID"`
}

// MetadataStatusEdges is a struct that represents the Metadata status edges GraphQL type
type MetadataStatusEdges struct {
	Node MetadataStatusNode
}

// MetadataStatuses is a struct that represents the Metadata statuses GraphQL type
type MetadataStatuses struct {
	Edges []MetadataStatusEdges
}

// Metadata is a struct that represents the metadata GraphQL type
type Metadata struct {
	ID       string           `graphql:"id"`
	NodeID   string           `graphql:"nodeID"`
	Statuses MetadataStatuses `graphql:"statuses" json:"statuses"`
}

// Readable version of the above:
// type GetLoadBalancer struct {
// 	LoadBalancer struct {
// 		ID    string
//      Owner string
// 		Name  string
//      IPAddresses {
// 	      id string
// 	      ip
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
