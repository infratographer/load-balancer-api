package api

import (
	"github.com/labstack/echo/v4"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/types"
	"go.infratographer.com/load-balancer-api/internal/models"
	"go.uber.org/zap"
)

// metadataCreate is the handler for the POST /loadbalancers/:load_balancer_id/metadata route
func (r *Router) metadataCreate(c echo.Context) error {
	ctx := c.Request().Context()
	payload := struct {
		Namespace string     `json:"namespace"`
		Data      types.JSON `json:"data"`
	}{}

	if err := c.Bind(&payload); err != nil {
		r.logger.Error("error binding payload", zap.Error(err))
		return v1BadRequestResponse(c, err)
	}

	lbID, err := r.parseUUID(c, "load_balancer_id")
	if err != nil {
		return v1BadRequestResponse(c, err)
	}

	// validate load balancer exists
	if _, err = models.LoadBalancers(
		models.LoadBalancerWhere.LoadBalancerID.EQ(lbID),
	).One(ctx, r.db); err != nil {
		r.logger.Error("error fetching load balancer", zap.Error(err))
		return v1BadRequestResponse(c, err)
	}

	metadataModel := &models.LoadBalancerMetadatum{
		LoadBalancerID: lbID,
		Namespace:      payload.Namespace,
		Data:           payload.Data,
	}

	if err := metadataModel.Insert(ctx, r.db, boil.Infer()); err != nil {
		r.logger.Error("error inserting metadata", zap.Error(err))
		return v1BadRequestResponse(c, err)
	}

	return v1MetadataCreatedResponse(c, metadataModel.MetadataID)
}
