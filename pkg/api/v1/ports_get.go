package api

import (
	"database/sql"
	"errors"

	"github.com/labstack/echo/v4"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"go.infratographer.com/load-balancer-api/internal/models"
	"go.uber.org/zap"
)

// portList returns a list of ports
func (r *Router) portList(c echo.Context) error {
	ctx := c.Request().Context()

	mods, err := r.portParamsBinding(c)
	if err != nil {
		r.logger.Error("failed to bind port params", zap.Error(err))
		return v1BadRequestResponse(c, err)
	}

	mods = append(
		mods,
		qm.Load("Assignments"),
	)

	ports, err := models.Ports(mods...).All(ctx, r.db)
	if err != nil {
		r.logger.Error("failed to get ports", zap.Error(err))
		return v1InternalServerErrorResponse(c, err)
	}

	return v1PortsResponse(c, ports)
}

// portGet returns a port by ID
func (r *Router) portGet(c echo.Context) error {
	ctx := c.Request().Context()

	portID, err := r.parseUUID(c, "port_id")
	if err != nil {
		return v1BadRequestResponse(c, err)
	}

	mods := []qm.QueryMod{
		models.PortWhere.PortID.EQ(portID),
		qm.Load("Assignments"),
	}

	port, err := models.Ports(mods...).One(ctx, r.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return v1NotFoundResponse(c)
		}

		return v1InternalServerErrorResponse(c, err)
	}

	return v1PortResponse(c, port)
}
