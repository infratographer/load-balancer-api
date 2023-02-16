package api

import "github.com/labstack/echo/v4"

// addFrontendRoutes adds the frontend routes to the router
func (r *Router) addFrontendRoutes(rg *echo.Group) {
	rg.GET("/frontends", r.frontendList)
	rg.GET("/frontends/:frontend_id", r.frontendGet)
	rg.GET("/loadbalancers/:load_balancer_id/frontends", r.frontendGet)

	rg.POST("/frontends", r.frontendCreate)

	rg.PUT("/frontends/:frontend_id", r.frontendUpdate)

	rg.DELETE("/frontends", r.frontendDelete)
	rg.DELETE("/frontends/:frontend_id", r.frontendDelete)
	rg.DELETE("/loadbalancers/:load_balancer_id/frontends", r.frontendDelete)
}
