package api

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// parseUUIDFromURI parses and validates a UUID from the request path if the path param is found
func (r *Router) parseUUIDFromURI(c *gin.Context, obj any) error {
	if err := c.ShouldBindUri(&obj); err != nil {
		r.logger.Error("error binding uri", zap.Error(err))
		return err
	}

	return nil
}

// parseUUID parses and validates a UUID from the request path if the path param is found
func (r *Router) parseUUID(id string) error {
	if id != "" {
		if _, err := uuid.Parse(id); err != nil {
			r.logger.Error("error parsing uuid", zap.Error(err))
			return ErrInvalidUUID
		}

		return nil
	}

	r.logger.Error("error parsing uuid", zap.Error(ErrUUIDNotFound))

	return ErrUUIDNotFound
}
