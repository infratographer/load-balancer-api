package metadata

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
