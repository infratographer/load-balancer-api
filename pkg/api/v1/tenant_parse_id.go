package api

import (
	"github.com/gin-gonic/gin"
)

// parseTenantID parses and validates a tenant ID from the request path if the path param is found
func (r *Router) parseTenantID(c *gin.Context) (string, error) {
	tenant := struct {
		id string `uri:"tenant_id" binding:"required"`
	}{}

	if err := r.parseUUIDFromURI(c, &tenant); err != nil {
		return "", err
	}

	if err := r.parseUUID(tenant.id); err != nil {
		return "", err
	}

	return tenant.id, nil
}
