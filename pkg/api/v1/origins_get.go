package api

import (
	"database/sql"
	"errors"

	"github.com/labstack/echo/v4"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"go.infratographer.com/load-balancer-api/internal/models"
	"go.uber.org/zap"
)

// originsGet returns an origin by id
func (r *Router) originsGet(c echo.Context) error {
	ctx := c.Request().Context()

	originID, err := r.parseUUID(c, "origin_id")
	if err != nil {
		return v1BadRequestResponse(c, err)
	}

	mods := []qm.QueryMod{models.OriginWhere.OriginID.EQ(originID)}

	origin, err := models.Origins(mods...).One(ctx, r.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return v1NotFoundResponse(c)
		}

		return v1InternalServerErrorResponse(c, err)
	}

	return v1OriginResponse(c, origin)
}

// originsList returns a list of origins
func (r *Router) originsList(c echo.Context) error {
	ctx := c.Request().Context()

	mods, err := r.originsParamsBinding(c)
	if err != nil {
		r.logger.Error("error parsing query params", zap.Error(err))
		return v1BadRequestResponse(c, err)
	}

	os, err := models.Origins(mods...).All(ctx, r.db)
	if err != nil {
		r.logger.Error("error getting origins", zap.Error(err))
		return v1InternalServerErrorResponse(c, err)
	}

	return v1OriginsResponse(c, os)
}
