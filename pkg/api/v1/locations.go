package api

import (
	"net/http"
	"time"

	"github.com/google/uuid"
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

	mods := []qm.QueryMod{models.LocationWhere.TenantID.EQ(c.Param("tenant_id"))}

	ls, err := models.Locations(mods...).All(ctx, r.db)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, v1LocationSlice(ls))
}

// locationGet returns a location for a tenant by ID
func (r *Router) locationGet(c echo.Context) error {
	tenantID := c.Param("tenant_id")
	name := c.Param("name")

	ctx, span := tracer.Start(c.Request().Context(), "locationGet")
	defer span.End()

	if _, err := uuid.Parse(tenantID); err != nil {
		return err
	}

	l, err := models.Locations(
		models.LocationWhere.TenantID.EQ(tenantID),
		models.LocationWhere.DisplayName.EQ(name),
	).One(ctx, r.db)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, v1Location(l))
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
		return err
	}

	return c.JSON(http.StatusOK, v1CreatedResponse())
}

// locationDelete soft deletes a location
func (r *Router) locationDelete(c echo.Context) error {
	ctx, span := tracer.Start(c.Request().Context(), "locationDelete")
	defer span.End()

	span.SetAttributes(attribute.String("route", "locationDelete"))

	l, err := models.FindLocation(ctx, r.db, c.Param("location_id"))
	if err != nil {
		return err
	}

	_, err = l.Delete(ctx, r.db, false)
	if err != nil {
		return err
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

func v1Location(l *models.Location) any {
	type loc struct {
		CreatedAt time.Time  `json:"created_at"`
		UpdatedAt time.Time  `json:"updated_at"`
		DeletedAt *time.Time `json:"deleted_at,omitempty"`
		ID        string     `json:"id"`
		TenantID  string     `json:"tenant_id"`
		Name      string     `json:"display_name"`
	}

	return struct {
		Version  string `json:"version"`
		Location loc    `json:"location"`
	}{
		Version: "v1",
		Location: loc{
			CreatedAt: l.CreatedAt,
			UpdatedAt: l.UpdatedAt,
			DeletedAt: l.DeletedAt.Ptr(),
			ID:        l.LocationID,
			TenantID:  l.TenantID,
			Name:      l.DisplayName,
		},
	}
}

func v1LocationSlice(ls models.LocationSlice) any {
	out := []any{}

	for _, l := range ls {
		out = append(out, v1Location(l))
	}

	return struct {
		Version   string `json:"version"`
		Locations []any  `json:"locations,omitempty"`
	}{
		Version:   "v1",
		Locations: out,
	}
}

func (r *Router) addLocationRoutes(g *echo.Group) {
	g.GET("/tenant/:tenant_id/locations", r.locationsList)
	g.POST("/tenant/:tenant_id/locations", r.locationCreate)
	g.GET("/tenant/:tenant_id/locations/:name", r.locationGet)
	g.DELETE("/tenant/:tenant_id/locations/:name", r.locationDelete)
}
