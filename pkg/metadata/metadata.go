package metadata

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
)

// LoadBalancerStatus is the status of a load balancer
type LoadBalancerStatus struct {
	State LoadBalancerState `json:"state"`
}
