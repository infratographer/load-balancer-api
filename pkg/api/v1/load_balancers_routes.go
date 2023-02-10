package api

import "github.com/labstack/echo/v4"

func (r *Router) addLoadBalancerRoutes(g *echo.Group) {
	g.GET("/loadbalancers", r.loadBalancerList)
	g.GET("/loadbalancers/:load_balancer_id", r.loadBalancerGet)

	g.POST("/loadbalancers", r.loadBalancerCreate)

	g.DELETE("/loadbalancers", r.loadBalancerDelete)
	g.DELETE("/loadbalancers/:load_balancer_id", r.loadBalancerDelete)
}
