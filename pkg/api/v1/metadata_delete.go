package api

import (
	"database/sql"
	"errors"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"go.infratographer.com/load-balancer-api/internal/models"
)

// metadataDelete is the handler for the DELETE /metadata/:metadata_id route
func (r *Router) metadataDelete(c echo.Context) error {
	ctx := c.Request().Context()

	mods, err := r.metadataParamsBinding(c)
	if err != nil {
		r.logger.Error("error parsing query params", zap.Error(err))
		return v1BadRequestResponse(c, err)
	}

	mds, err := models.LoadBalancerMetadata(mods...).All(ctx, r.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return v1InternalServerErrorResponse(c, err)
		}
	}

	if len(mds) == 0 {
		return v1NotFoundResponse(c)
	} else if len(mds) > 1 {
		return v1BadRequestResponse(c, ErrAmbiguous)
	}

	metadata := mds[0]

	_, err = metadata.Delete(ctx, r.db, false)
	if err != nil {
		r.logger.Error("error deleting metadata", zap.Error(err))
		return v1InternalServerErrorResponse(c, err)
	}

	return v1DeleteMetadataResponse(c)
}
