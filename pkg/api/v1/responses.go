package api

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"go.infratographer.com/loadbalancerapi/internal/models"
)

type frontend struct {
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
	ID             string     `json:"id"`
	TenantID       string     `json:"tenant_id"`
	LoadBalancerID string     `json:"load_balancer_id"`
	Name           string     `json:"display_name"`
	AddressFamily  string     `json:"address_family"`
	Port           int64      `json:"port"`
}

type frontendSlice []*frontend

type loadBalancer struct {
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty"`
	ID         string     `json:"id"`
	TenantID   string     `json:"tenant_id"`
	IPAddress  string     `json:"ip_address"`
	Name       string     `json:"display_name"`
	LocationID string     `json:"location_id"`
	Size       string     `json:"load_balancer_size"`
	Type       string     `json:"load_balancer_type"`
}

type loadBalancerSlice []*loadBalancer

type location struct {
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
	ID        string     `json:"id"`
	TenantID  string     `json:"tenant_id"`
	Name      string     `json:"display_name"`
}

type locationSlice []*location

type response struct {
	Version       string             `json:"version"`
	Kind          string             `json:"kind"`
	Frontends     *frontendSlice     `json:"frontends,omitempty"`
	LoadBalancers *loadBalancerSlice `json:"load_balancers,omitempty"`
	Locations     *locationSlice     `json:"locations,omitempty"`
}

func v1DeletedResponse(c echo.Context) error {
	return c.JSON(http.StatusNoContent, struct {
		DeletedAt time.Time `json:"deleted_at"`
		Message   string    `json:"message"`
		Status    int       `json:"status"`
		Version   string    `json:"version"`
	}{
		Version:   "v1",
		DeletedAt: time.Now(),
		Message:   "resource deleted",
		Status:    http.StatusNoContent,
	})
}

func v1CreatedResponse(c echo.Context) error {
	return c.JSON(http.StatusCreated, struct {
		Version   string    `json:"version"`
		CreatedAt time.Time `json:"created_at"`
		Message   string    `json:"message"`
		Status    int       `json:"status"`
	}{
		CreatedAt: time.Now(),
		Message:   "resource created",
		Version:   "v1",
		Status:    http.StatusCreated,
	})
}

func v1NotFoundResponse(c echo.Context) error {
	return c.JSON(http.StatusNotFound, struct {
		Version string `json:"version"`
		Message string `json:"message"`
		Status  int    `json:"status"`
	}{
		Version: "v1",
		Message: "resource not found",
		Status:  http.StatusNotFound,
	})
}

func v1BadRequestResponse(c echo.Context, err error) error {
	return c.JSON(http.StatusBadRequest, struct {
		Version string `json:"version"`
		Message string `json:"message"`
		Error   string `json:"error"`
		Status  int    `json:"status"`
	}{
		Version: "v1",
		Message: "bad request",
		Error:   err.Error(),
		Status:  http.StatusBadRequest,
	})
}

func v1UnprocessableEntityResponse(c echo.Context, err error) error {
	return c.JSON(http.StatusUnprocessableEntity, struct {
		Version string `json:"version"`
		Message string `json:"message"`
		Error   string `json:"error"`
		Status  int    `json:"status"`
	}{
		Version: "v1",
		Message: "unprocessable entity",
		Error:   err.Error(),
		Status:  http.StatusUnprocessableEntity,
	})
}

func v1InternalServerErrorResponse(c echo.Context, err error) error {
	return c.JSON(http.StatusInternalServerError, struct {
		Version string `json:"version"`
		Message string `json:"message"`
		Error   string `json:"error"`
		Status  int    `json:"status"`
	}{
		Version: "v1",
		Message: "internal server error",
		Error:   err.Error(),
		Status:  http.StatusInternalServerError,
	})
}

func v1Frontends(c echo.Context, fs models.FrontendSlice) error {
	out := frontendSlice{}
	for _, f := range fs {
		out = append(out, &frontend{
			CreatedAt:      f.CreatedAt,
			UpdatedAt:      f.UpdatedAt,
			DeletedAt:      f.DeletedAt.Ptr(),
			ID:             f.FrontendID,
			LoadBalancerID: f.LoadBalancerID,
			Port:           f.Port,
			AddressFamily:  f.AfInet,
			Name:           f.DisplayName,
			TenantID:       f.TenantID,
		})
	}

	return c.JSON(http.StatusOK, &response{
		Version:   "v1",
		Kind:      "frontendsList",
		Frontends: &out,
	})
}

func v1LoadBalancers(c echo.Context, lbs models.LoadBalancerSlice) error {
	out := loadBalancerSlice{}

	for _, lb := range lbs {
		out = append(out, &loadBalancer{
			CreatedAt:  lb.CreatedAt,
			UpdatedAt:  lb.UpdatedAt,
			DeletedAt:  lb.DeletedAt.Ptr(),
			ID:         lb.LoadBalancerID,
			Name:       lb.DisplayName,
			IPAddress:  lb.IPAddr,
			TenantID:   lb.TenantID,
			LocationID: lb.LocationID,
			Size:       lb.LoadBalancerSize,
			Type:       lb.LoadBalancerType,
		})
	}

	return c.JSON(http.StatusOK, &response{
		Version:       "v1",
		Kind:          "loadBalancersList",
		LoadBalancers: &out,
	})
}

func v1Locations(c echo.Context, ls models.LocationSlice) error {
	out := locationSlice{}

	for _, l := range ls {
		out = append(out, &location{
			CreatedAt: l.CreatedAt,
			UpdatedAt: l.UpdatedAt,
			DeletedAt: l.DeletedAt.Ptr(),
			ID:        l.LocationID,
			TenantID:  l.TenantID,
			Name:      l.DisplayName,
		})
	}

	return c.JSON(http.StatusOK, &response{
		Version:   "v1",
		Kind:      "locationsList",
		Locations: &out,
	})
}
