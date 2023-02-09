package api

import (
	"github.com/labstack/echo/v4"
)

// addAssignRoutes adds the assignment routes to the router
func (r *Router) addAssignRoutes(g *echo.Group) {
	g.GET("/assignments", r.assignmentsGet)
	g.POST("/assignments", r.assignmentsCreate)
	g.DELETE("/assignments", r.assignmentsDelete)
}
