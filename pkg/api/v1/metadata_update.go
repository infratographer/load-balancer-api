package api

import (
	"github.com/labstack/echo/v4"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/types"
	"go.uber.org/zap"

	"go.infratographer.com/load-balancer-api/internal/models"
)

// metadataUpdate updates an metadata about a load balancer
func (r *Router) metadataUpdate(c echo.Context) error {
	ctx := c.Request().Context()

	mods, err := r.metadataParamsBinding(c)
	if err != nil {
		r.logger.Error("failed to bind metadata params", zap.Error(err))
		return v1BadRequestResponse(c, err)
	}

	mods = append(mods, qm.Load("LoadBalancer"))

	mds, err := models.LoadBalancerMetadata(mods...).All(ctx, r.db)
	if err != nil {
		r.logger.Error("failed to get port", zap.Error(err))
		return v1InternalServerErrorResponse(c, err)
	}

	if len(mds) == 0 {
		return v1NotFoundResponse(c)
	} else if len(mds) != 1 {
		r.logger.Warn("ambiguous query ", zap.Any("metadatas", mds))
		return v1BadRequestResponse(c, ErrAmbiguous)
	}

	metadata := mds[0]

	payload := struct {
		Data types.JSON `json:"data"`
	}{}

	if err := c.Bind(&payload); err != nil {
		return v1BadRequestResponse(c, err)
	}

	// update origin
	metadata.Data = payload.Data

	return r.updateMetadata(c, metadata)
}

// metadataPatch patches an origin
func (r *Router) metadataPatch(c echo.Context) error {
	ctx := c.Request().Context()

	mods, err := r.metadataParamsBinding(c)
	if err != nil {
		r.logger.Error("failed to bind metadata params", zap.Error(err))
		return v1BadRequestResponse(c, err)
	}

	mods = append(mods, qm.Load("LoadBalancer"))

	mds, err := models.LoadBalancerMetadata(mods...).All(ctx, r.db)
	if err != nil {
		r.logger.Error("failed to get port", zap.Error(err))
		return v1InternalServerErrorResponse(c, err)
	}

	if len(mds) == 0 {
		return v1NotFoundResponse(c)
	} else if len(mds) != 1 {
		r.logger.Warn("ambiguous query ", zap.Any("metadatas", mds))
		return v1BadRequestResponse(c, ErrAmbiguous)
	}

	metadata := mds[0]

	payload := struct {
		Data *types.JSON `json:"data"`
	}{}

	if err := c.Bind(&payload); err != nil {
		r.logger.Error("failed to bind origin patch input", zap.Error(err))
		return v1BadRequestResponse(c, err)
	}

	if payload.Data != nil {
		metadata.Data = *payload.Data
	}

	return r.updateMetadata(c, metadata)
}

func (r *Router) updateMetadata(c echo.Context, metadata *models.LoadBalancerMetadatum) error {
	ctx := c.Request().Context()

	if _, err := metadata.Update(ctx, r.db, boil.Infer()); err != nil {
		r.logger.Error("failed to update port", zap.Error(err))
		return v1InternalServerErrorResponse(c, err)
	}

	return v1UpdateMetadataResponse(c, metadata.MetadataID)
}
