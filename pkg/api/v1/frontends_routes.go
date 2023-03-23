package api

import "github.com/gin-gonic/gin"

// addFrontendRoutes adds the frontend routes to the router
func (r *Router) addFrontendRoutes(rg *gin.RouterGroup) {
	rg.GET("/frontends/:frontend_id", r.frontendGet)
	rg.GET("/loadbalancers/:load_balancer_id/frontends", r.frontendList)

	rg.POST("/loadbalancers/:load_balancer_id/frontends", r.frontendCreate)

	rg.PUT("/frontends/:frontend_id", r.frontendUpdate)
	rg.PUT("/loadbalancers/:load_balancer_id/frontends", r.frontendUpdate)

	rg.DELETE("/frontends/:frontend_id", r.frontendDelete)
	rg.DELETE("/loadbalancers/:load_balancer_id/frontends", r.frontendDelete)
}
