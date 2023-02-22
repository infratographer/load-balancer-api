package api

import "github.com/labstack/echo/v4"

// addFrontendRoutes adds the frontend routes to the router
func (r *Router) addFrontendRoutes(rg *echo.Group) {
	rg.GET("/frontends/:frontend_id", r.frontendGet)
	rg.GET("/loadbalancers/:load_balancer_id/frontends", r.frontendList)

	rg.POST("/loadbalancers/:load_balancer_id/frontends", r.frontendCreate)

	rg.PUT("/frontends/:frontend_id", r.frontendUpdate)

	rg.DELETE("/frontends/:frontend_id", r.frontendDelete)
	rg.DELETE("/loadbalancers/:load_balancer_id/frontends", r.frontendDelete)
}
