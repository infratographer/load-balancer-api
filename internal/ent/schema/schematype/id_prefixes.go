package schematype

const (
	// ApplicationPrefix is the prefix for all application IDs owned by load balancer API
	ApplicationPrefix string = "load"
	// LoadBalancerPrefix is the prefix for all load balancer IDs
	LoadBalancerPrefix string = ApplicationPrefix + "bal"
	// LoadBalancerAnnotationPrefix is the prefix for all load balancer annotation IDs
	LoadBalancerAnnotationPrefix string = ApplicationPrefix + "ban"
	// LoadBalancerStatusPrefix is the prefix for all load balancer Status IDs
	LoadBalancerStatusPrefix string = ApplicationPrefix + "bst"
	// LoadBalancerProviderPrefix is the prefix for all load balancer provider IDs
	LoadBalancerProviderPrefix string = ApplicationPrefix + "pvd"
	// OriginPrefix is the prefix for all origin IDs
	OriginPrefix string = ApplicationPrefix + "ogn"
	// PortPrefix is the prefix for all port IDs
	PortPrefix string = ApplicationPrefix + "prt"
	// PoolPrefix is the prefix for all pool IDs
	PoolPrefix string = ApplicationPrefix + "pol"
)
