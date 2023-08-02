package client

// OriginNode represents a single origin
type OriginNode struct {
	ID         string
	Name       string
	Target     string
	PortNumber int64
	Active     bool
}

// OriginEdges represents an edge of an origin
type OriginEdges struct {
	Node OriginNode
}

// Origins represents a list of origin edges
type Origins struct {
	Edges []OriginEdges
}

// PoolNode represents a single pool
type PoolNode struct {
	ID       string
	Name     string
	Protocol string
	Origins  Origins
}

// PoolEdges represents an edge of a pool
type PoolEdges struct {
	Node PoolNode
}

// Pools represents a list of pool edges
type Pools struct {
	Edges []PoolEdges
}

// PortNode represents a single port
type PortNode struct {
	ID     string
	Name   string
	Number int64
	Pools  Pools
}

// PortEdges represents an edge of a port
type PortEdges struct {
	Node PortNode
}

// Ports represents a list of port edges
type Ports struct {
	Edges []PortEdges
}

// OwnerNode represents a single owner
type OwnerNode struct {
	ID string
}

// LocationNode represents a single location
type LocationNode struct {
	ID string
}

// LoadBalancer represents a single load balancer
type LoadBalancer struct {
	ID          string
	Name        string
	Owner       OwnerNode
	Location    LocationNode
	IPAddresses []IPAddress `graphql:"IPAddresses" json:"IPAddresses"`
	Ports       Ports
}

// GetLoadBalancer represents the graphql query for a load balancer
type GetLoadBalancer struct {
	LoadBalancer LoadBalancer `graphql:"loadBalancer(id: $id)"`
}

// IPAddress represents a single ip address
type IPAddress struct {
	ID       string
	IP       string
	Reserved bool
}

// Readable version of the above:
// type GetLoadBalancer struct {
// 	LoadBalancer struct {
// 		ID    string
//      Owner string
// 		Name  string
//		IPAddressableFragment
// 		Ports struct {
// 			Edges []struct {
// 				Node struct {
// 					Name   string
// 					Number int64
// 					Pools  struct {
//						Edges []struct {
//							Node struct {
// 								Name     string
// 								Protocol string
// 								Origins  struct {
// 									Edges []struct {
// 										Node struct {
// 											Name       string
// 											Target     string
// 											PortNumber int64
// 											Active     bool
//										}
//									}
// 								}
// 							}
// 						}
// 					}
// 				}
// 			}
// 		}
// 	} `graphql:"loadBalancer(id: $id)"`
// }
