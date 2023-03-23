package api

import "github.com/gin-gonic/gin"

// addAssignRoutes adds the assignment routes to the router
func (r *Router) addAssignRoutes(g *gin.RouterGroup) {
	g.GET("/tenant/:tenant_id/assignments", r.assignmentsGet)
	g.POST("/tenant/:tenant_id/assignments", r.assignmentsCreate)
	g.DELETE("/tenant/:tenant_id/assignments", r.assignmentsDelete)
}
