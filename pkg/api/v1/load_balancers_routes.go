package api

import "github.com/labstack/echo/v4"

func (r *Router) addLoadBalancerRoutes(g *echo.Group) {
	g.GET("/tenant/:tenant_id/loadbalancers", r.loadBalancerList)
	g.GET("/loadbalancers/:load_balancer_id", r.loadBalancerGet)

	g.POST("/tenant/:tenant_id/loadbalancers", r.loadBalancerCreate)

	g.DELETE("/tenant/:tenant_id/loadbalancers", r.loadBalancerDelete)
	g.DELETE("/loadbalancers/:load_balancer_id", r.loadBalancerDelete)
}
