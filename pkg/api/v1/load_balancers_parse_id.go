package api

import (
	"github.com/gin-gonic/gin"
)

// parseLoadBalancerID parses and validates a load balancer ID from the request path if the path param is found
func (r *Router) parseLoadBalancerID(c *gin.Context) (string, error) {
	loadBalancer := struct {
		id string `uri:"load_balancer_id" binding:"required"`
	}{}

	if err := r.parseUUIDFromURI(c, &loadBalancer); err != nil {
		return "", err
	}

	if err := r.parseUUID(loadBalancer.id); err != nil {
		return "", err
	}

	return loadBalancer.id, nil
}
