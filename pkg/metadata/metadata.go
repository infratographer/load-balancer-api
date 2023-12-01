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

	LoadBalancerStateIPAssigned   LoadBalancerState = "ip-address.assigned"
	LoadBalancerStateIPUnassigned LoadBalancerState = "ip-address.unassigned"

	LoadBalancerAPISource string = "load-balancer-api"
)

// LoadBalancerStatus is the status of a load balancer
type LoadBalancerStatus struct {
	State LoadBalancerState `json:"state"`
}

// GetLoadbalancerStatus searches through the list of metadata status for the requested status of a load balancer using namespace and source
func GetLoadbalancerStatus(metadataStatuses client.MetadataStatuses, statusNamespaceID gidx.PrefixedID, source string) (*LoadBalancerStatus, error) {
	if metadataStatuses.TotalCount > 0 {
		for _, s := range metadataStatuses.Edges {
			if s.Node.StatusNamespaceID == statusNamespaceID.String() && s.Node.Source == source {
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
