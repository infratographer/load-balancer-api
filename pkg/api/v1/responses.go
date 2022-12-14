package api

import (
	"net/http"
	"time"

	"go.infratographer.com/loadbalancerapi/internal/models"
)

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
	LoadBalancer  *loadBalancer      `json:"load_balancer,omitempty"`
	LoadBalancers *loadBalancerSlice `json:"load_balancers,omitempty"`
	Location      *location          `json:"location,omitempty"`
	Locations     *locationSlice     `json:"locations,omitempty"`
}

func v1DeletedResponse() any {
	//
	return struct {
		DeletedAt time.Time `json:"deleted_at"`
		Message   string    `json:"message"`
		Status    int       `json:"status"`
		Version   string    `json:"version"`
	}{
		Version:   "v1",
		DeletedAt: time.Now(),
		Message:   "resource deleted",
		Status:    http.StatusNoContent,
	}
}

func v1CreatedResponse() any {
	return struct {
		Version   string    `json:"version"`
		CreatedAt time.Time `json:"created_at"`
		Message   string    `json:"message"`
		Status    int       `json:"status"`
	}{
		CreatedAt: time.Now(),
		Message:   "resource created",
		Version:   "v1",
		Status:    http.StatusCreated,
	}
}

func v1NotFoundResponse() any {
	return struct {
		Version string `json:"version"`
		Message string `json:"message"`
		Status  int    `json:"status"`
	}{
		Version: "v1",
		Message: "resource not found",
		Status:  http.StatusNotFound,
	}
}

func v1BadRequestResponse(err error) any {
	return struct {
		Version string `json:"version"`
		Message string `json:"message"`
		Error   string `json:"error"`
		Status  int    `json:"status"`
	}{
		Version: "v1",
		Message: "bad request",
		Error:   err.Error(),
		Status:  http.StatusBadRequest,
	}
}

func v1UnprocessableEntityResponse(err error) any {
	return struct {
		Version string `json:"version"`
		Message string `json:"message"`
		Error   string `json:"error"`
		Status  int    `json:"status"`
	}{
		Version: "v1",
		Message: "unprocessable entity",
		Error:   err.Error(),
		Status:  http.StatusUnprocessableEntity,
	}
}

func v1InternalServerErrorResponse(err error) any {
	return struct {
		Version string `json:"version"`
		Message string `json:"message"`
		Error   string `json:"error"`
		Status  int    `json:"status"`
	}{
		Version: "v1",
		Message: "internal server error",
		Error:   err.Error(),
		Status:  http.StatusInternalServerError,
	}
}

func v1LoadBalancer(lb *models.LoadBalancer) *response {
	return &response{
		Version: "v1",
		Kind:    "loadBalancer",
		LoadBalancer: &loadBalancer{
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
		},
	}
}

func v1LoadBalancerSlice(lbs models.LoadBalancerSlice) *response {
	out := loadBalancerSlice{}

	for _, lb := range lbs {
		out = append(out, v1LoadBalancer(lb).LoadBalancer)
	}

	return &response{
		Version:       "v1",
		Kind:          "loadBalancersList",
		LoadBalancers: &out,
	}
}

func v1Location(l *models.Location) *response {
	return &response{
		Version: "v1",
		Kind:    "location",
		Location: &location{
			CreatedAt: l.CreatedAt,
			UpdatedAt: l.UpdatedAt,
			DeletedAt: l.DeletedAt.Ptr(),
			ID:        l.LocationID,
			TenantID:  l.TenantID,
			Name:      l.DisplayName,
		},
	}
}

func v1LocationSlice(ls models.LocationSlice) *response {
	out := locationSlice{}

	for _, l := range ls {
		out = append(out, v1Location(l).Location)
	}

	return &response{
		Version:   "v1",
		Kind:      "locationsList",
		Locations: &out,
	}
}
