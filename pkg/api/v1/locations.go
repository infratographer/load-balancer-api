package api

import (
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

	tenantID, err := r.parseTenantID(c)
	if err != nil {
		return v1BadRequestResponse(c, err)
	}

	mods := []qm.QueryMod{models.LocationWhere.TenantID.EQ(tenantID)}

	name := c.Param("name")
	if name != "" {
		mods = append(mods, models.LocationWhere.DisplayName.EQ(name))
	}

	ls, err := models.Locations(mods...).All(ctx, r.db)
	if err != nil {
		return v1InternalServerErrorResponse(c, err)
	}

	switch len(ls) {
	case 0:
		return v1NotFoundResponse(c)
	default:
		return v1Locations(c, ls)
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
		return v1InternalServerErrorResponse(c, err)
	}

	tenantID, err := r.parseTenantID(c)
	if err != nil {
		return v1BadRequestResponse(c, err)
	}

	l := &models.Location{
		TenantID:    tenantID,
		DisplayName: input.Name,
	}

	if err := valdiateLocation(l); err != nil {
		return v1BadRequestResponse(c, err)
	}

	if err := l.Insert(ctx, r.db, boil.Infer()); err != nil {
		return v1InternalServerErrorResponse(c, err)
	}

	return v1CreatedResponse(c)
}

// locationDelete soft deletes a location
func (r *Router) locationDelete(c echo.Context) error {
	ctx, span := tracer.Start(c.Request().Context(), "locationDelete")
	defer span.End()

	span.SetAttributes(attribute.String("route", "locationDelete"))

	tenantID, err := r.parseTenantID(c)
	if err != nil {
		return v1BadRequestResponse(c, err)
	}

	mods := []qm.QueryMod{
		models.LocationWhere.TenantID.EQ(tenantID),
		models.LocationWhere.DisplayName.EQ(c.Param("name")),
	}

	l, err := models.Locations(mods...).One(ctx, r.db)
	if err != nil {
		return v1NotFoundResponse(c)
	}

	if _, err = l.Delete(ctx, r.db, false); err != nil {
		return v1InternalServerErrorResponse(c, err)
	}

	return v1DeletedResponse(c)
}

func valdiateLocation(l *models.Location) error {
	if l.DisplayName == "" {
		return ErrDisplayNameMissing
	}

	return nil
}

func (r *Router) addLocationRoutes(g *echo.Group) {
	g.GET("/locations", r.locationsList)
	g.GET("/locations/:name", r.locationsList)

	g.POST("/locations", r.locationCreate)

	g.DELETE("/locations/:name", r.locationDelete)
}
