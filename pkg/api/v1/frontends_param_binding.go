package api

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"go.uber.org/zap"

	"go.infratographer.com/load-balancer-api/internal/models"
)

// frontendParamsBinding binds the request path and query params to a slice of query mods
// for use with sqlboiler. It returns an error if an invalid uuid is provided
// for the load_balancer_id or frontend_id in the request path. It also iterates the
// expected query params and appends them to the slice of query mods if they are present
// in the request.
func (r *Router) frontendParamsBinding(c *gin.Context) ([]qm.QueryMod, error) {
	var (
		err            error
		loadBalancerID string
		frontendID     string
	)

	mods := []qm.QueryMod{}

	// optional load_balancer_id in the request path
	if loadBalancerID, err = r.parseLoadBalancerID(c); err != nil {
		if !errors.Is(err, ErrUUIDNotFound) {
			return nil, err
		}
	} else {
		// found load_balancer_id in path so add to query mods
		mods = append(mods, models.FrontendWhere.LoadBalancerID.EQ(loadBalancerID))
		r.logger.Debug("path param", zap.String("load_balancer_id", loadBalancerID))
	}

	// optional frontend_id in the request path
	if frontendID, err = r.parseFrontEndID(c); err != nil {
		if !errors.Is(err, ErrUUIDNotFound) {
			return nil, err
		}
	} else {
		// found frontend_id in path so add to query mods
		mods = append(mods, models.FrontendWhere.FrontendID.EQ(frontendID))
		r.logger.Debug("path param", zap.String("frontend_id", frontendID))
	}

	// query params
	// queryParams := []string{"port", "load_balancer_id", "slug", "af_inet"}

	query := struct {
		Port           int64  `form:"port"`
		LoadBalancerID string `form:"load_balancer_id"`
		Slug           string `form:"slug"`
		AFInet         string `form:"af_inet"`
	}{}

	if err := c.ShouldBindQuery(&query); err != nil {
		return nil, err
	}

	if query.Port != 0 {
		mods = append(mods, models.FrontendWhere.Port.EQ(query.Port))
	}

	if query.LoadBalancerID != "" {
		mods = append(mods, models.FrontendWhere.LoadBalancerID.EQ(query.LoadBalancerID))
	}

	if query.Slug != "" {
		mods = append(mods, models.FrontendWhere.Slug.EQ(query.Slug))
	}

	if query.AFInet != "" {
		mods = append(mods, models.FrontendWhere.AfInet.EQ(query.AFInet))
	}

	return mods, nil
}
