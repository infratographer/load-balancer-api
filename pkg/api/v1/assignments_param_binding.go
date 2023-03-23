package api

import (
	"github.com/gin-gonic/gin"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"go.infratographer.com/load-balancer-api/internal/models"
)

func (r *Router) assignmentParamsBinding(c *gin.Context) ([]qm.QueryMod, error) {
	mods := []qm.QueryMod{}

	tenantID, err := r.parseTenantID(c)
	if err != nil {
		return nil, err
	}

	mods = append(mods, models.AssignmentWhere.TenantID.EQ(tenantID))

	// queryParams := []string{"frontend_id", "pool_id"}

	query := struct {
		FrontendID string `form:"frontend_id"`
		PoolID     string `form:"pool_id"`
	}{}

	if err := c.ShouldBindQuery(&query); err != nil {
		return nil, err
	}

	if query.FrontendID != "" {
		mods = append(mods, models.AssignmentWhere.FrontendID.EQ(query.FrontendID))
	}

	if query.PoolID != "" {
		mods = append(mods, models.AssignmentWhere.PoolID.EQ(query.PoolID))
	}

	return mods, nil
}
