package api

import (
	"github.com/labstack/echo/v4"
)

// addAssignRoutes adds the assignment routes to the router
func (r *Router) addAssignRoutes(g *echo.Group) {
	g.GET("/tenant/:tenant_id/assignments", r.assignmentsGet)
	g.POST("/tenant/:tenant_id/assignments", r.assignmentsCreate)
	g.DELETE("/tenant/:tenant_id/assignments", r.assignmentsDelete)
}
