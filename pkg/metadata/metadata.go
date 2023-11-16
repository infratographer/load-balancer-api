package metadata

import (
	"encoding/json"
	"errors"
	"fmt"

	"go.infratographer.com/x/gidx"

	"go.infratographer.com/load-balancer-api/pkg/client"
)

// LoadBalancerState is the state of a load balancer
type LoadBalancerState string

// load balancer states
const (
	LoadBalancerStateCreating    LoadBalancerState = "creating"
	LoadBalancerStateTerminating LoadBalancerState = "terminating"
	LoadBalancerStateActive      LoadBalancerState = "active"
	LoadBalancerStateDeleted     LoadBalancerState = "deleted"
	LoadBalancerStateUpdating    LoadBalancerState = "updating"
)

var (
	// ErrStatusNotFound is returned when a status is not found in the payload
	ErrStatusNotFound = errors.New("status not found")

	// ErrInvalidStatusData is returned when the status json data is invalid
	ErrInvalidStatusData = errors.New("invalid status json data")
)

// ErrUnknownLoadBalancerState is returned when the load balancer state is unknown
type ErrUnknownLoadBalancerState struct {
	State string
}

func (e ErrUnknownLoadBalancerState) Error() string {
	return "unknown load balancer state: " + e.State
}

// GetLoadbalancerState returns the status of a load balancer
func GetLoadbalancerState(metadataStatuses client.MetadataStatuses, statusNamespaceID gidx.PrefixedID) (LoadBalancerState, error) {
	for _, status := range metadataStatuses.Edges {
		if status.Node.StatusNamespaceID == statusNamespaceID.String() {
			// we've found the loadbalancer status namespace we are looking fah
			data := map[string]string{}
			if err := json.Unmarshal(status.Node.Data, &data); err != nil {
				// bad data stored in metadata-api, drop the message and move on
				return "", fmt.Errorf("%w: %s", ErrInvalidStatusData, err)
			}

			if state, ok := data["state"]; ok {
				switch state {
				case string(LoadBalancerStateCreating):
					fallthrough
				case string(LoadBalancerStateTerminating):
					fallthrough
				case string(LoadBalancerStateActive):
					fallthrough
				case string(LoadBalancerStateDeleted):
					fallthrough
				case string(LoadBalancerStateUpdating):
					return LoadBalancerState(state), nil
				default:
					return "", ErrUnknownLoadBalancerState{State: state}
				}
			}
		}
	}

	return "", ErrStatusNotFound
}
