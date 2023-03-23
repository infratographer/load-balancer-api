package api

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"go.infratographer.com/load-balancer-api/internal/models"
)

func (r *Router) originsParamsBinding(c *gin.Context) ([]qm.QueryMod, error) {
	mods := []qm.QueryMod{}

	// optional pool_id in the request path
	if poolID, err := r.parsePoolID(c); err != nil {
		if !errors.Is(err, ErrUUIDNotFound) {
			return nil, err
		}
	} else {
		// found pool_id in path so add to query mods
		mods = append(mods, models.OriginWhere.PoolID.EQ(poolID))
		r.logger.Debugw("path param", "pool_id", poolID)
	}

	originID := c.Param("origin_id")
	if originID != "" {
		if _, err := uuid.Parse(originID); err != nil {
			return nil, ErrInvalidUUID
		}

		mods = append(mods, models.OriginWhere.OriginID.EQ(originID))
		r.logger.Debugw("path param", "origin_id", originID)
	}

	// queryParams := []string{"slug", "target", "port"}

	query := struct {
		Slug   string `form:"slug"`
		Target string `form:"target"`
		Port   int64  `form:"port"`
	}{}

	if err := c.ShouldBindQuery(&query); err != nil {
		return nil, err
	}

	if query.Slug != "" {
		mods = append(mods, models.OriginWhere.Slug.EQ(query.Slug))
	}

	if query.Target != "" {
		mods = append(mods, models.OriginWhere.OriginTarget.EQ(query.Target))
	}

	if query.Port != 0 {
		mods = append(mods, models.OriginWhere.Port.EQ(query.Port))
	}

	return mods, nil
}
