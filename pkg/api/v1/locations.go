package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"go.infratographer.com/loadbalancerapi/internal/models"
	"go.opentelemetry.io/otel/attribute"
)

// locationsList return all the locations for a tenant
func (r *Router) locationsList(c echo.Context) error {
	ctx, span := tracer.Start(c.Request().Context(), "locationsList")
	defer span.End()

	span.SetAttributes(attribute.String("route", "locationsList"))

	tenantID, err := parseTenantID(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, v1BadRequestResponse(err))
	}

	mods := []qm.QueryMod{models.LocationWhere.TenantID.EQ(tenantID)}

	name := c.Param("name")
	if name != "" {
		mods = append(mods, models.LocationWhere.DisplayName.EQ(name))
	}

	ls, err := models.Locations(mods...).All(ctx, r.db)
	if err != nil {
		return err
	}

	switch len(ls) {
	case 0:
		return c.JSON(http.StatusNotFound, v1NotFoundResponse())
	case 1:
		return c.JSON(http.StatusOK, v1Location(ls[0]))
	default:
		return c.JSON(http.StatusOK, v1LocationSlice(ls))
	}
}

// locationCreate creates a new location for a tenant
func (r *Router) locationCreate(c echo.Context) error {
	ctx, span := tracer.Start(c.Request().Context(), "locationCreate")
	defer span.End()

	input := struct {
		Name string `json:"display_name"`
	}{}

	if err := c.Bind(&input); err != nil {
		return err
	}

	tenantID, err := parseTenantID(c)
	if err != nil {
		return err
	}

	l := &models.Location{
		TenantID:    tenantID,
		DisplayName: input.Name,
	}

	if err := valdiateLocation(l); err != nil {
		return err
	}

	if err := l.Insert(ctx, r.db, boil.Infer()); err != nil {
		return c.JSON(http.StatusInternalServerError, v1InternalServerErrorResponse(err))
	}

	return c.JSON(http.StatusCreated, v1CreatedResponse())
}

// locationDelete soft deletes a location
func (r *Router) locationDelete(c echo.Context) error {
	ctx, span := tracer.Start(c.Request().Context(), "locationDelete")
	defer span.End()

	span.SetAttributes(attribute.String("route", "locationDelete"))

	tenantID, err := parseTenantID(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, v1BadRequestResponse(err))
	}

	mods := []qm.QueryMod{
		models.LocationWhere.TenantID.EQ(tenantID),
		models.LocationWhere.DisplayName.EQ(c.Param("name")),
	}

	l, err := models.Locations(mods...).One(ctx, r.db)
	if err != nil {
		return c.JSON(http.StatusNotFound, v1NotFoundResponse())
	}

	_, err = l.Delete(ctx, r.db, false)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, v1InternalServerErrorResponse(err))
	}

	return c.JSON(http.StatusOK, v1DeletedResponse())
}

func valdiateLocation(l *models.Location) error {
	if l.DisplayName == "" {
		return ErrNameRequired
	}

	if l.TenantID == "" {
		return ErrTenantIDRequired
	}

	return nil
}

type locationResp struct {
	Version   string        `json:"version"`
	Kind      string        `json:"kind"`
	Location  *location     `json:"location,omitempty"`
	Locations locationSlice `json:"locations,omitempty"`
}

func (r *Router) addLocationRoutes(g *echo.Group) {
	g.GET("/tenant/:tenant_id/locations", r.locationsList)
	g.POST("/tenant/:tenant_id/locations", r.locationCreate)
	g.GET("/tenant/:tenant_id/locations/:name", r.locationsList)
	g.DELETE("/tenant/:tenant_id/locations/:name", r.locationDelete)
}
