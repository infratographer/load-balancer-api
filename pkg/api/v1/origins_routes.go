package api

import "github.com/gin-gonic/gin"

// addOriginsRoutes adds the origins routes to the router
func (r *Router) addOriginRoutes(g *gin.RouterGroup) {
	g.GET("/pools/:pool_id/origins", r.originsList)
	g.GET("/origins/:origin_id", r.originsGet)

	g.POST("/pools/:pool_id/origins", r.originsCreate)

	g.DELETE("/pools/:pool_id/origins", r.originsDelete)
	g.DELETE("/origins/:origin_id", r.originsDelete)
}
