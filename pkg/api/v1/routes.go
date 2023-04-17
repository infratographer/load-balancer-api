package api

import "github.com/labstack/echo/v4"

// addAssignRoutes adds the assignment routes to the router
func (r *Router) addAssignRoutes(g *echo.Group) {
	g.GET("/tenant/:tenant_id/assignments", r.assignmentsList)
	g.POST("/tenant/:tenant_id/assignments", r.assignmentsCreate)
	g.DELETE("/tenant/:tenant_id/assignments", r.assignmentsDelete)
}

func (r *Router) addLoadBalancerRoutes(g *echo.Group) {
	g.GET("/tenant/:tenant_id/loadbalancers", r.loadBalancerList)
	g.GET("/loadbalancers/:load_balancer_id", r.loadBalancerGet)

	g.POST("/tenant/:tenant_id/loadbalancers", r.loadBalancerCreate)

	g.PUT("/loadbalancers/:load_balancer_id", r.loadBalancerUpdate)

	g.PATCH("/loadbalancers/:load_balancer_id", r.loadBalancerPatch)

	g.DELETE("/tenant/:tenant_id/loadbalancers", r.loadBalancerDelete)
	g.DELETE("/loadbalancers/:load_balancer_id", r.loadBalancerDelete)
}

// addMetadataRoutes adds the metadata routes to the router
func (r *Router) addMetadataRoutes(g *echo.Group) {
	g.GET("/loadbalancers/:load_balancer_id/metadatas", r.metadataList)
	g.GET("/metadatas/:metadata_id", r.metadataGet)

	g.POST("/loadbalancers/:load_balancer_id/metadata", r.metadataCreate)

	g.DELETE("/metadata/:metadata_id", r.metadataDelete)
}

// addOriginsRoutes adds the origins routes to the router
func (r *Router) addOriginRoutes(g *echo.Group) {
	g.GET("/pools/:pool_id/origins", r.originsList)
	g.GET("/origins/:origin_id", r.originsGet)

	g.POST("/pools/:pool_id/origins", r.originsCreate)

	g.PUT("/pools/:pools/origins", r.originUpdate)

	g.PATCH("/pools/:pools/origins", r.originPatch)

	g.DELETE("/pools/:pool_id/origins", r.originsDelete)
	g.DELETE("/origins/:origin_id", r.originsDelete)
}

// addPoolsRoutes adds the routes for the pools API
func (r *Router) addPoolsRoutes(g *echo.Group) {
	g.GET("/tenant/:tenant_id/pools", r.poolsList)
	g.GET("/pools/:pool_id", r.poolsGet)

	g.POST("/tenant/:tenant_id/pools", r.poolCreate)

	g.PUT("/pools/:pool_id", r.poolUpdate)

	g.PATCH("/pools/:pool_id", r.poolPatch)

	g.DELETE("/tenant/:tenant_id/pools", r.poolDelete)
	g.DELETE("/pools/:pool_id", r.poolDelete)
}

// addPortRoutes adds the port routes to the router
func (r *Router) addPortRoutes(g *echo.Group) {
	g.GET("/ports/:port_id", r.portGet)
	g.GET("/loadbalancers/:load_balancer_id/ports", r.portList)

	g.POST("/loadbalancers/:load_balancer_id/ports", r.portCreate)

	g.PUT("/ports/:port_id", r.portUpdate)
	g.PUT("/loadbalancers/:load_balancer_id/ports", r.portUpdate)

	g.PATCH("/ports/:port_id", r.portPatch)
	g.PATCH("/loadbalancers/:load_balancer_id/ports", r.portPatch)

	g.DELETE("/ports/:port_id", r.portDelete)
	g.DELETE("/loadbalancers/:load_balancer_id/ports", r.portDelete)
}
