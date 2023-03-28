package api

import (
	"errors"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"go.uber.org/zap"

	"go.infratographer.com/load-balancer-api/internal/models"
)

func (r *Router) originsParamsBinding(c echo.Context) ([]qm.QueryMod, error) {
	mods := []qm.QueryMod{}

	// optional pool_id in the request path
	if poolID, err := r.parseUUID(c, "pool_id"); err != nil {
		if !errors.Is(err, ErrUUIDNotFound) {
			return nil, err
		}
	} else {
		// found pool_id in path so add to query mods
		mods = append(mods, models.OriginWhere.PoolID.EQ(poolID))
		r.logger.Debug("path param", zap.String("pool.id", poolID))
	}

	originID := c.Param("origin_id")
	if originID != "" {
		if _, err := uuid.Parse(originID); err != nil {
			return nil, ErrInvalidUUID
		}

		mods = append(mods, models.OriginWhere.OriginID.EQ(originID))
		r.logger.Debug("path param", zap.String("origin.id", originID))
	}

	queryParams := []string{"slug", "target", "port"}

	qpb := echo.QueryParamsBinder(c)

	for _, qp := range queryParams {
		mods = queryParamsToQueryMods(qpb, qp, mods)

		if len(c.QueryParam(qp)) > 0 {
			r.logger.Debug("origin query parameters", zap.String("query.key", qp), zap.String("query.value", c.QueryParam(qp)))
		}
	}

	if err := qpb.BindError(); err != nil {
		return nil, err
	}

	return mods, nil
}
