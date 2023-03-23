package api

import "github.com/gin-gonic/gin"

func (r *Router) addLoadBalancerRoutes(g *gin.RouterGroup) {
	g.GET("/tenant/:tenant_id/loadbalancers", r.loadBalancerList)
	g.GET("/loadbalancers/:load_balancer_id", r.loadBalancerGet)

	g.POST("/tenant/:tenant_id/loadbalancers", r.loadBalancerCreate)

	g.PUT("/loadbalancers/:load_balancer_id", r.loadBalancerUpdate)

	g.DELETE("/tenant/:tenant_id/loadbalancers", r.loadBalancerDelete)
	g.DELETE("/loadbalancers/:load_balancer_id", r.loadBalancerDelete)
}
