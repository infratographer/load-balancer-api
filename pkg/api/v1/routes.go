package api

import "github.com/labstack/echo/v4"

// addAssignRoutes adds the assignment routes to the router
func (r *Router) addAssignRoutes(g *echo.Group) {
	g.GET("/tenant/:tenant_id/assignments", r.assignmentsGet)
	g.POST("/tenant/:tenant_id/assignments", r.assignmentsCreate)
	g.DELETE("/tenant/:tenant_id/assignments", r.assignmentsDelete)
}

func (r *Router) addLoadBalancerRoutes(g *echo.Group) {
	g.GET("/tenant/:tenant_id/loadbalancers", r.loadBalancerList)
	g.GET("/loadbalancers/:load_balancer_id", r.loadBalancerGet)

	g.POST("/tenant/:tenant_id/loadbalancers", r.loadBalancerCreate)

	g.PUT("/loadbalancers/:load_balancer_id", r.loadBalancerUpdate)

	g.DELETE("/tenant/:tenant_id/loadbalancers", r.loadBalancerDelete)
	g.DELETE("/loadbalancers/:load_balancer_id", r.loadBalancerDelete)
}

// addOriginsRoutes adds the origins routes to the router
func (r *Router) addOriginRoutes(g *echo.Group) {
	g.GET("/pools/:pool_id/origins", r.originsList)
	g.GET("/origins/:origin_id", r.originsGet)

	g.POST("/pools/:pool_id/origins", r.originsCreate)

	g.DELETE("/pools/:pool_id/origins", r.originsDelete)
	g.DELETE("/origins/:origin_id", r.originsDelete)
}

// addPoolsRoutes adds the routes for the pools API
func (r *Router) addPoolsRoutes(g *echo.Group) {
	g.GET("/tenant/:tenant_id/pools", r.poolsList)
	g.GET("/pools/:pool_id", r.poolsGet)

	g.POST("/tenant/:tenant_id/pools", r.poolCreate)

	g.DELETE("/tenant/:tenant_id/pools", r.poolDelete)
	g.DELETE("/pools/:pool_id", r.poolDelete)
}

// addPortRoutes adds the port routes to the router
func (r *Router) addPortRoutes(rg *echo.Group) {
	rg.GET("/ports/:port_id", r.portGet)
	rg.GET("/loadbalancers/:load_balancer_id/ports", r.portList)

	rg.POST("/loadbalancers/:load_balancer_id/ports", r.portCreate)

	rg.PUT("/ports/:port_id", r.portUpdate)
	rg.PUT("/loadbalancers/:load_balancer_id/ports", r.portUpdate)

	rg.DELETE("/ports/:port_id", r.portDelete)
	rg.DELETE("/loadbalancers/:load_balancer_id/ports", r.portDelete)
}
