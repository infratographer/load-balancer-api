package api

import "github.com/labstack/echo/v4"

// assignmentsGet handles the GET /assignments route
func (r *Router) assignmentsGet(c echo.Context) error {
	return c.JSON(200, []string{"assignments"})
}

// assignmentsPost handles the POST /assignments route
func (r *Router) assignmentsPost(c echo.Context) error {
	return v1CreatedResponse(c)
}

// assignmentsDelete handles the DELETE /assignments route
func (r *Router) assignmentsDelete(c echo.Context) error {
	return nil
}

// addAssignRoutes adds the assignment routes to the router
func (r *Router) addAssignRoutes(g *echo.Group) {
	g.GET("/assignments", r.assignmentsGet)
	g.POST("/assignments", r.assignmentsPost)
	g.DELETE("/assignments", r.assignmentsDelete)
}
