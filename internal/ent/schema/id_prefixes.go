package schema

const (
	ApplicationPrefix            string = "load"
	LoadBalancerPrefix           string = ApplicationPrefix + "bal"
	LoadBalancerAnnotationPrefix string = ApplicationPrefix + "ban"
	LoadBalancerStatusPrefix     string = ApplicationPrefix + "bst"
	LoadBalancerProviderPrefix   string = ApplicationPrefix + "pvd"
	OriginPrefix                 string = ApplicationPrefix + "ogn"
	PortPrefix                   string = ApplicationPrefix + "prt"
	PoolPrefix                   string = ApplicationPrefix + "pol"
)
