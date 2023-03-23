package api

import (
	"github.com/gin-gonic/gin"
)

// parseFrontEndID parses and validates a front end ID from the request path if the path param is found
func (r *Router) parseFrontEndID(c *gin.Context) (string, error) {
	frontEnd := struct {
		id string `uri:"frontend_id" binding:"required"`
	}{}

	if err := r.parseUUIDFromURI(c, &frontEnd); err != nil {
		return "", err
	}

	if err := r.parseUUID(frontEnd.id); err != nil {
		return "", err
	}

	return frontEnd.id, nil
}
