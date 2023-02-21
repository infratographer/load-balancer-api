package api

import "github.com/labstack/echo/v4"

// addPoolsRoutes adds the routes for the pools API
func (r *Router) addPoolsRoutes(g *echo.Group) {
	g.GET("/tenant/:tenant_id/pools", r.poolsList)
	g.GET("/pools/:pool_id", r.poolsGet)

	g.POST("/tenant/:tenant_id/pools", r.poolCreate)

	g.DELETE("/tenant/:tenant_id/pools", r.poolDelete)
	g.DELETE("/pools/:pool_id", r.poolDelete)
}
