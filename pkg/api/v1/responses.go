package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"go.infratographer.com/loadbalancerapi/internal/models"
)

func v1DeletedResponse(c echo.Context) error {
	return c.JSON(http.StatusOK, struct {
		Message string `json:"message"`
		Status  int    `json:"status"`
		Version string `json:"version"`
	}{
		Version: apiVersion,
		Message: "resource deleted",
		Status:  http.StatusOK,
	})
}

func v1CreatedResponse(c echo.Context) error {
	return c.JSON(http.StatusOK, struct {
		Version string `json:"version"`
		Message string `json:"message"`
		Status  int    `json:"status"`
	}{
		Message: "resource created",
		Version: apiVersion,
		Status:  http.StatusOK,
	})
}

func v1NotFoundResponse(c echo.Context) error {
	return c.JSON(http.StatusNotFound, struct {
		Version string `json:"version"`
		Message string `json:"message"`
		Status  int    `json:"status"`
	}{
		Version: apiVersion,
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
		Version: apiVersion,
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
		Version: apiVersion,
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
		Version: apiVersion,
		Message: "internal server error",
		Error:   err.Error(),
		Status:  http.StatusInternalServerError,
	})
}

func v1Assignments(c echo.Context, as models.AssignmentSlice) error {
	out := assignmentSlice{}
	for _, a := range as {
		out = append(out, &assignment{
			CreatedAt:      a.CreatedAt,
			UpdatedAt:      a.UpdatedAt,
			DeletedAt:      a.DeletedAt.Ptr(),
			ID:             a.AssignmentID,
			LoadBalancerID: a.LoadBalancerID,
			FrontendID:     a.FrontendID,
			PoolID:         a.PoolID,
			TenantID:       a.TenantID,
		})
	}

	return c.JSON(http.StatusOK, &response{
		Version:     apiVersion,
		Kind:        "assignmentsList",
		Assignments: &out,
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
		Version:   apiVersion,
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
		Version:       apiVersion,
		Kind:          "loadBalancersList",
		LoadBalancers: &out,
	})
}

func v1OriginsResponse(c echo.Context, os models.OriginSlice) error {
	out := originSlice{}

	for _, o := range os {
		out = append(out, &origin{
			CreatedAt:      o.CreatedAt,
			UpdatedAt:      o.UpdatedAt,
			DeletedAt:      o.DeletedAt.Ptr(),
			ID:             o.OriginID,
			Name:           o.DisplayName,
			TenantID:       o.TenantID,
			OriginDisabled: o.OriginUserSettingDisabled,
			OriginTarget:   o.OriginTarget,
			Port:           o.Port,
		})
	}

	return c.JSON(http.StatusOK, &response{
		Version: apiVersion,
		Kind:    "originsList",
		Origins: &out,
	})
}

func v1PoolsResponse(c echo.Context, ps models.PoolSlice) error {
	out := poolSlice{}

	for _, p := range ps {
		out = append(out, &pool{
			CreatedAt: p.CreatedAt,
			UpdatedAt: p.UpdatedAt,
			DeletedAt: p.DeletedAt.Ptr(),
			ID:        p.PoolID,
			Name:      p.DisplayName,
			Protocol:  p.Protocol,
			TenantID:  p.TenantID,
		})
	}

	return c.JSON(http.StatusOK, &response{
		Version: apiVersion,
		Kind:    "poolsList",
		Pools:   &out,
	})
}
