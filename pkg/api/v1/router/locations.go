package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"go.infratographer.sh/loadbalancerapi/pkg/api/v1/locations"
)

const (
	// LocationsBaseURI is the path prefix for all locations endpoints
	LocationsBaseURI = "/locations"

	// TenantLocationsURI is the path for a location by Tenant UUID
	TenantLocationsURI = TenantPrefix + LocationsBaseURI

	// TenantLocationsNameURI is the path for a location by UUID and name
	TenantLocationsNameURI = TenantLocationsURI + "/:name"
)

// addLocationRoutes adds the routes for this API version to a router group
func (r *Router) addLocationRoutes(rg *gin.RouterGroup) {
	rg.POST(LocationsBaseURI, r.createLocation)

	rg.GET(TenantLocationsURI, r.getLocations)

	rg.GET(TenantLocationsNameURI, r.getLocationByName)
	rg.DELETE(TenantLocationsNameURI, r.deleteLocationByName)
}

func locationSuccess(c *gin.Context, obj interface{}) {
	uri := uriWithoutQueryParams(c)
	r := newResponse("resources found")
	r.Location = obj
	r.Links = &responseLinks{
		Self: &link{Href: uri},
	}

	c.Header("Location", uri)
	c.JSON(http.StatusCreated, r)
}

func locationsSuccess(c *gin.Context, obj interface{}) {
	uri := uriWithoutQueryParams(c)
	r := newResponse("resources found")
	r.Locations = obj
	r.Links = &responseLinks{
		Self: &link{Href: uri},
	}

	c.Header("Location", uri)
	c.JSON(http.StatusCreated, r)
}

// createLocation creates a new location
func (r *Router) createLocation(c *gin.Context) {
	loc, err := locations.NewLocation(c)
	if err != nil {
		badRequestResponse(c, locations.ErrInvalid.Error(), err)
		return
	}

	if err := loc.Create(c, r.db); err != nil {
		badRequestResponse(c, locations.ErrWrite.Error(), err)
		return
	}

	createdResponse(c)
}

// getLocation returns all locations
func (r *Router) getLocations(c *gin.Context) {
	tenant, err := uuid.Parse(c.Param("tenant"))
	if err != nil {
		badRequestResponse(c, "could not parse tenant", err)
		return
	}

	locations, err := locations.GetLocations(c, r.db, tenant)
	if err != nil {
		badRequestResponse(c, "could not get locations", err)
		return
	}

	locationsSuccess(c, locations)
}

func (r *Router) getLocationByName(c *gin.Context) {
	tenant, err := uuid.Parse(c.Param("tenant"))
	if err != nil {
		badRequestResponse(c, "could not parse tenant", err)
		return
	}

	loc := locations.Location{
		Name:     c.Param("name"),
		TenantID: tenant,
	}

	err = loc.Find(c, r.db)
	if err != nil {
		notFoundResponse(c, locations.ErrNotFound.Error())
		return
	}

	locationSuccess(c, loc)
}

// deleteLocationByID deletes a location by ID
func (r *Router) deleteLocationByName(c *gin.Context) {
	tenant, err := uuid.Parse(c.Param("tenant"))
	if err != nil {
		badRequestResponse(c, "could not parse tenant", err)
		return
	}

	loc := locations.Location{
		Name:     c.Param("name"),
		TenantID: tenant,
	}

	err = loc.Find(c, r.db)
	if err != nil {
		notFoundResponse(c, locations.ErrNotFound.Error())
		return
	}

	err = loc.Delete(c, r.db)
	if err != nil {
		badRequestResponse(c, locations.ErrWrite.Error(), err)
		return
	}

	deletedResponse(c)
}
