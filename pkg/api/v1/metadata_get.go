package api

import (
	"database/sql"
	"errors"

	"github.com/labstack/echo/v4"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"go.uber.org/zap"

	"go.infratographer.com/load-balancer-api/internal/models"
)

// metadataGet is the handler for the GET /metadatas/:metadata_id route
func (r *Router) metadataGet(c echo.Context) error {
	ctx := c.Request().Context()

	metadataID, err := r.parseUUID(c, "metadata_id")
	if err != nil {
		return v1BadRequestResponse(c, err)
	}

	mods := []qm.QueryMod{models.LoadBalancerMetadatumWhere.MetadataID.EQ(metadataID)}

	metadata, err := models.LoadBalancerMetadata(mods...).One(ctx, r.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return v1InternalServerErrorResponse(c, err)
		}
	}

	return v1MetadataResponse(c, metadata)
}

// metadataList is the handler for the GET /loadbalancers/:load_balancer_id/metadata route
func (r *Router) metadataList(c echo.Context) error {
	ctx := c.Request().Context()

	mods, err := r.metadataParamsBinding(c)
	if err != nil {
		r.logger.Error("error parsing query params", zap.Error(err))
		return v1BadRequestResponse(c, err)
	}

	mds, err := models.LoadBalancerMetadata(mods...).All(ctx, r.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return v1NotFoundResponse(c)
		}

		r.logger.Error("error loading metadata", zap.Error(err))

		return v1InternalServerErrorResponse(c, err)
	}

	return v1MetadatasResponse(c, mds)
}
