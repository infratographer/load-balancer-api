package api

import "github.com/labstack/echo/v4"

// v1routes adds the v1 routes to an echo group
func (r *Router) v1Routes(g *echo.Group) {
	// assignment routes
	g.GET("/tenant/:tenant_id/assignments", r.assignmentsGet)
	g.POST("/tenant/:tenant_id/assignments", r.assignmentsCreate)
	g.DELETE("/tenant/:tenant_id/assignments", r.assignmentsDelete)

	// port routes
	g.GET("/ports/:port_id", r.portGet)
	g.GET("/loadbalancers/:load_balancer_id/ports", r.portList)
	g.POST("/loadbalancers/:load_balancer_id/ports", r.portCreate)
	g.PUT("/ports/:port_id", r.portUpdate)
	g.PUT("/loadbalancers/:load_balancer_id/ports", r.portUpdate)
	g.DELETE("/ports/:port_id", r.portDelete)
	g.DELETE("/loadbalancers/:load_balancer_id/ports", r.portDelete)

	// load balancer routes
	g.GET("/tenant/:tenant_id/loadbalancers", r.loadBalancerList)
	g.GET("/loadbalancers/:load_balancer_id", r.loadBalancerGet)
	g.POST("/tenant/:tenant_id/loadbalancers", r.loadBalancerCreate)
	g.PUT("/loadbalancers/:load_balancer_id", r.loadBalancerUpdate)
	g.DELETE("/tenant/:tenant_id/loadbalancers", r.loadBalancerDelete)
	g.DELETE("/loadbalancers/:load_balancer_id", r.loadBalancerDelete)

	// origin routes
	g.GET("/pools/:pool_id/origins", r.originsList)
	g.GET("/origins/:origin_id", r.originsGet)
	g.POST("/pools/:pool_id/origins", r.originsCreate)
	g.DELETE("/pools/:pool_id/origins", r.originsDelete)
	g.DELETE("/origins/:origin_id", r.originsDelete)

	// pool routes
	g.GET("/tenant/:tenant_id/pools", r.poolsList)
	g.GET("/pools/:pool_id", r.poolsGet)
	g.POST("/tenant/:tenant_id/pools", r.poolCreate)
	g.DELETE("/tenant/:tenant_id/pools", r.poolDelete)
	g.DELETE("/pools/:pool_id", r.poolDelete)
}
