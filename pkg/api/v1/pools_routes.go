package api

import "github.com/gin-gonic/gin"

// addPoolsRoutes adds the routes for the pools API
func (r *Router) addPoolsRoutes(g *gin.RouterGroup) {
	g.GET("/tenant/:tenant_id/pools", r.poolsList)
	g.GET("/pools/:pool_id", r.poolsGet)

	g.POST("/tenant/:tenant_id/pools", r.poolCreate)

	g.DELETE("/tenant/:tenant_id/pools", r.poolDelete)
	g.DELETE("/pools/:pool_id", r.poolDelete)
}
