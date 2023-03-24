package api

import "github.com/labstack/echo/v4"

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
