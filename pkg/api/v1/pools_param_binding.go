package api

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"go.uber.org/zap"

	"go.infratographer.com/load-balancer-api/internal/models"
)

// poolsParamsBinding return a set of query mods based
// on the query parameters and path parameters
func (r *Router) poolsParamsBinding(c *gin.Context, relationships ...string) ([]qm.QueryMod, error) {
	mods := []qm.QueryMod{}

	// optional tenant_id in the request path
	if tenantID, err := r.parseTenantID(c); err != nil {
		if !errors.Is(err, ErrUUIDNotFound) {
			return nil, err
		}
	} else {
		// found tenant_id in path so add to query mods
		mods = append(mods, models.PoolWhere.TenantID.EQ(tenantID))
		r.logger.Debug("path param", zap.String("tenant_id", tenantID))
	}

	poolID := c.Param("pool_id")
	if poolID != "" {
		if _, err := uuid.Parse(poolID); err != nil {
			return nil, ErrInvalidUUID
		}

		mods = append(mods, models.PoolWhere.PoolID.EQ(poolID))
		r.logger.Debug("path param", zap.String("pool_id", poolID))
	}

	// queryParams := []string{"slug", "protocol", "name"}

	query := struct {
		Slug     string `form:"slug"`
		Protocol string `form:"protocol"`
		Name     string `form:"name"`
	}{}

	if err := c.ShouldBindQuery(&query); err != nil {
		return nil, err
	}

	if query.Slug != "" {
		mods = append(mods, models.PoolWhere.Slug.EQ(query.Slug))
	}

	if query.Protocol != "" {
		mods = append(mods, models.PoolWhere.Protocol.EQ(query.Protocol))
	}

	if query.Name != "" {
		mods = append(mods, models.PoolWhere.Name.EQ(query.Name))
	}

	// append relationships to query

	for _, rel := range relationships {
		r.logger.Debug("appending relationships to query", zap.String("relationship", rel))
		mods = append(mods, qm.Load(rel))
	}

	return mods, nil
}
