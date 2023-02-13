package api

import "github.com/labstack/echo/v4"

// addOriginsRoutes adds the origins routes to the router
func (r *Router) addOriginRoutes(g *echo.Group) {
	g.GET("/origins", r.originsList)
	g.GET("/origins/:origin_id", r.originsGet)

	g.POST("/origins", r.originsCreate)

	g.DELETE("/origins", r.originsDelete)
	g.DELETE("/origins/:origin_id", r.originsDelete)
}
