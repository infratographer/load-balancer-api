package api

import (
	"github.com/gin-gonic/gin"
)

// parsePoolID parses and validates a pool ID from the request path if the path param is found
func (r *Router) parsePoolID(c *gin.Context) (string, error) {
	pool := struct {
		id string `uri:"pool_id" binding:"required"`
	}{}

	if err := r.parseUUIDFromURI(c, &pool); err != nil {
		return "", err
	}

	if err := r.parseUUID(pool.id); err != nil {
		return "", err
	}

	return pool.id, nil
}
