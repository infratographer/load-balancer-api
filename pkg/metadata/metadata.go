package metadata

import (
	"encoding/json"
	"fmt"

	"go.infratographer.com/x/gidx"

	"go.infratographer.com/load-balancer-api/pkg/client"
)

// LoadBalancerState state of a load balancer
type LoadBalancerState string

// load balancer states
const (
	LoadBalancerStateCreating    LoadBalancerState = "creating"
	LoadBalancerStateTerminating LoadBalancerState = "terminating"
	LoadBalancerStateActive      LoadBalancerState = "active"
	LoadBalancerStateDeleted     LoadBalancerState = "deleted"
	LoadBalancerStateUpdating    LoadBalancerState = "updating"
)

// LoadBalancerStatus is the status of a load balancer
type LoadBalancerStatus struct {
	State LoadBalancerState `json:"state"`
}

// GetLoadbalancerStatus returns the status of a load balancer
func GetLoadbalancerStatus(metadataStatuses client.MetadataStatuses, statusNamespaceID gidx.PrefixedID) (*LoadBalancerStatus, error) {
	if metadataStatuses.TotalCount > 0 {
		for _, s := range metadataStatuses.Edges {
			if s.Node.StatusNamespaceID == statusNamespaceID.String() {
				status := &LoadBalancerStatus{}

				if err := json.Unmarshal(s.Node.Data, status); err != nil {
					return nil, fmt.Errorf("%w: %s", ErrInvalidStatusData, err)
				}

				return status, nil
			}
		}
	}

	return nil, ErrStatusNotFound
}
