package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"go.infratographer.com/load-balancer-api/internal/models"
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

func v1MetadataCreatedResponse(c echo.Context, id string) error {
	return c.JSON(http.StatusOK, struct {
		Version string `json:"version"`
		Message string `json:"message"`
		Status  int    `json:"status"`
		MetaID  string `json:"metadata_id"`
	}{
		Message: "resource created",
		Version: apiVersion,
		Status:  http.StatusOK,
		MetaID:  id,
	})
}

func v1AssignmentsCreatedResponse(c echo.Context, id string) error {
	return c.JSON(http.StatusOK, struct {
		Version      string `json:"version"`
		Message      string `json:"message"`
		Status       int    `json:"status"`
		AssignmentID string `json:"assignment_id"`
	}{
		Message:      "resource created",
		Version:      apiVersion,
		Status:       http.StatusOK,
		AssignmentID: id,
	})
}

func v1LoadBalancerCreatedResponse(c echo.Context, id string) error {
	return c.JSON(http.StatusOK, struct {
		Version        string `json:"version"`
		Message        string `json:"message"`
		Status         int    `json:"status"`
		LoadBalancerID string `json:"load_balancer_id"`
	}{
		Message:        "resource created",
		Version:        apiVersion,
		Status:         http.StatusOK,
		LoadBalancerID: id,
	})
}

func v1PortCreatedResponse(c echo.Context, id string) error {
	return c.JSON(http.StatusOK, struct {
		Version string `json:"version"`
		Message string `json:"message"`
		Status  int    `json:"status"`
		PortID  string `json:"port_id"`
	}{
		Message: "resource created",
		Version: apiVersion,
		Status:  http.StatusOK,
		PortID:  id,
	})
}

func v1OriginCreatedResponse(c echo.Context, id string) error {
	return c.JSON(http.StatusOK, struct {
		Version  string `json:"version"`
		Message  string `json:"message"`
		Status   int    `json:"status"`
		OriginID string `json:"origin_id"`
	}{
		Message:  "resource created",
		Version:  apiVersion,
		Status:   http.StatusOK,
		OriginID: id,
	})
}

func v1UpdateOriginResponse(c echo.Context, id string) error {
	return c.JSON(http.StatusAccepted, struct {
		Version  string `json:"version"`
		Message  string `json:"message"`
		Status   int    `json:"status"`
		OriginID string `json:"origin_id"`
	}{
		Message:  "resource updated",
		Version:  apiVersion,
		Status:   http.StatusAccepted,
		OriginID: id,
	})
}

func v1UpdateMetadataResponse(c echo.Context, id string) error {
	return c.JSON(http.StatusAccepted, struct {
		Version    string `json:"version"`
		Message    string `json:"message"`
		Status     int    `json:"status"`
		MetadataID string `json:"metadata_id"`
	}{
		Message:    "resource updated",
		Version:    apiVersion,
		Status:     http.StatusAccepted,
		MetadataID: id,
	})
}

func v1PoolCreatedResponse(c echo.Context, id string) error {
	return c.JSON(http.StatusOK, struct {
		Version string `json:"version"`
		Message string `json:"message"`
		Status  int    `json:"status"`
		PoolID  string `json:"pool_id"`
	}{
		Message: "resource created",
		Version: apiVersion,
		Status:  http.StatusOK,
		PoolID:  id,
	})
}

func v1DeleteMetadataResponse(c echo.Context) error {
	return c.JSON(http.StatusOK, struct {
		Version string `json:"version"`
		Message string `json:"message"`
		Status  int    `json:"status"`
	}{
		Version: apiVersion,
		Message: "resource deleted",
		Status:  http.StatusOK,
	})
}

func v1UpdatePoolResponse(c echo.Context, id string) error {
	return c.JSON(http.StatusAccepted, struct {
		Version string `json:"version"`
		Message string `json:"message"`
		Status  int    `json:"status"`
		PoolID  string `json:"pool_id"`
	}{
		Message: "resource updated",
		Version: apiVersion,
		Status:  http.StatusAccepted,
		PoolID:  id,
	})
}

func v1UpdatePortResponse(c echo.Context, id string) error {
	return c.JSON(http.StatusAccepted, struct {
		Version string `json:"version"`
		Message string `json:"message"`
		Status  int    `json:"status"`
		PortID  string `json:"port_id"`
	}{
		Message: "resource updated",
		Version: apiVersion,
		Status:  http.StatusAccepted,
		PortID:  id,
	})
}

func v1UpdateLoadBalancerResponse(c echo.Context, id string) error {
	return c.JSON(http.StatusAccepted, struct {
		Version        string `json:"version"`
		Message        string `json:"message"`
		Status         int    `json:"status"`
		LoadBalancerID string `json:"load_balancer_id"`
	}{
		Message:        "resource updated",
		Version:        apiVersion,
		Status:         http.StatusAccepted,
		LoadBalancerID: id,
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
			CreatedAt: a.CreatedAt,
			UpdatedAt: a.UpdatedAt,
			ID:        a.AssignmentID,
			PortID:    a.PortID,
			PoolID:    a.PoolID,
			TenantID:  a.TenantID,
		})
	}

	return c.JSON(http.StatusOK, &response{
		Version:     apiVersion,
		Kind:        "assignmentsList",
		Assignments: &out,
	})
}

func v1PortResponse(c echo.Context, p *models.Port) error {
	pools := make([]string, len(p.R.Assignments))
	for k, a := range p.R.Assignments {
		pools[k] = a.PoolID
	}

	return c.JSON(http.StatusOK, &response{
		Version: apiVersion,
		Kind:    "portsGet",
		Port: &port{
			CreatedAt:      p.CreatedAt,
			UpdatedAt:      p.UpdatedAt,
			DeletedAt:      p.DeletedAt.Ptr(),
			ID:             p.PortID,
			LoadBalancerID: p.LoadBalancerID,
			Port:           p.Port,
			AddressFamily:  p.AfInet,
			Name:           p.Name,
			Pools:          pools,
		},
	})
}

func v1PortsResponse(c echo.Context, ps models.PortSlice) error {
	out := make(portSlice, len(ps))

	for i, p := range ps {
		pools := make([]string, len(p.R.Assignments))
		for k, a := range p.R.Assignments {
			pools[k] = a.PoolID
		}

		out[i] = &port{
			CreatedAt:      p.CreatedAt,
			UpdatedAt:      p.UpdatedAt,
			DeletedAt:      p.DeletedAt.Ptr(),
			ID:             p.PortID,
			LoadBalancerID: p.LoadBalancerID,
			Port:           p.Port,
			AddressFamily:  p.AfInet,
			Name:           p.Name,
			Pools:          pools,
		}
	}

	return c.JSON(http.StatusOK, &response{
		Version: apiVersion,
		Kind:    "portsList",
		Ports:   &out,
	})
}

func v1LoadBalancer(c echo.Context, lb *models.LoadBalancer) error {
	pSlice := make(portSlice, len(lb.R.Ports))

	for j, p := range lb.R.Ports {
		pools := make([]string, len(p.R.Assignments))
		for k, a := range p.R.Assignments {
			pools[k] = a.PoolID
		}

		pSlice[j] = &port{
			CreatedAt:      p.CreatedAt,
			UpdatedAt:      p.UpdatedAt,
			DeletedAt:      p.DeletedAt.Ptr(),
			TenantID:       p.R.LoadBalancer.TenantID,
			LoadBalancerID: p.R.LoadBalancer.LoadBalancerID,
			ID:             p.PortID,
			Port:           p.Port,
			AddressFamily:  p.AfInet,
			Name:           p.Name,
			Pools:          pools,
		}
	}

	return c.JSON(http.StatusOK, &response{
		Version: apiVersion,
		Kind:    "loadBalancersGet",
		LoadBalancer: &loadBalancer{
			CreatedAt:   lb.CreatedAt,
			UpdatedAt:   lb.UpdatedAt,
			DeletedAt:   lb.DeletedAt.Ptr(),
			ID:          lb.LoadBalancerID,
			Name:        lb.Name,
			IPAddressID: lb.IPAddressID,
			TenantID:    lb.TenantID,
			LocationID:  lb.LocationID,
			Size:        lb.LoadBalancerSize,
			Type:        lb.LoadBalancerType,
			Ports:       pSlice,
		},
	})
}

func v1LoadBalancers(c echo.Context, lbs models.LoadBalancerSlice) error {
	out := make(loadBalancerSlice, len(lbs))

	for i, lb := range lbs {
		l := &loadBalancer{
			CreatedAt:   lb.CreatedAt,
			UpdatedAt:   lb.UpdatedAt,
			DeletedAt:   lb.DeletedAt.Ptr(),
			ID:          lb.LoadBalancerID,
			Name:        lb.Name,
			IPAddressID: lb.IPAddressID,
			TenantID:    lb.TenantID,
			LocationID:  lb.LocationID,
			Size:        lb.LoadBalancerSize,
			Type:        lb.LoadBalancerType,
		}

		pSlice := make(portSlice, len(lb.R.Ports))

		for j, p := range lb.R.Ports {
			pools := make([]string, len(p.R.Assignments))
			for k, a := range p.R.Assignments {
				pools[k] = a.PoolID
			}

			pSlice[j] = &port{
				CreatedAt:      p.CreatedAt,
				UpdatedAt:      p.UpdatedAt,
				DeletedAt:      p.DeletedAt.Ptr(),
				TenantID:       p.R.LoadBalancer.TenantID,
				LoadBalancerID: p.R.LoadBalancer.LoadBalancerID,
				ID:             p.PortID,
				Port:           p.Port,
				AddressFamily:  p.AfInet,
				Name:           p.Name,
				Pools:          pools,
			}
		}

		l.Ports = pSlice

		out[i] = l
	}

	return c.JSON(http.StatusOK, &response{
		Version:       apiVersion,
		Kind:          "loadBalancersList",
		LoadBalancers: &out,
	})
}

func v1MetadataResponse(c echo.Context, m *models.LoadBalancerMetadatum) error {
	return c.JSON(http.StatusOK, &response{
		Version: apiVersion,
		Kind:    "metadataGet",
		Metadata: &metadata{
			CreatedAt: m.CreatedAt,
			UpdatedAt: m.UpdatedAt,
			ID:        m.MetadataID,
			Namespace: m.Namespace,
			Data:      m.Data,
		},
	})
}

func v1MetadatasResponse(c echo.Context, ms models.LoadBalancerMetadatumSlice) error {
	out := make(metadataSlice, len(ms))

	for i, m := range ms {
		out[i] = &metadata{
			CreatedAt: m.CreatedAt,
			UpdatedAt: m.UpdatedAt,
			ID:        m.MetadataID,
			Namespace: m.Namespace,
			Data:      m.Data,
		}
	}

	return c.JSON(http.StatusOK, &response{
		Version:   apiVersion,
		Kind:      "metadataList",
		Metadatas: &out,
	})
}

func v1OriginResponse(c echo.Context, o *models.Origin) error {
	return c.JSON(http.StatusOK, &response{
		Version: apiVersion,
		Kind:    "originsGet",
		Origin: &origin{
			CreatedAt:      o.CreatedAt,
			UpdatedAt:      o.UpdatedAt,
			DeletedAt:      o.DeletedAt.Ptr(),
			ID:             o.OriginID,
			Name:           o.Name,
			OriginDisabled: o.OriginUserSettingDisabled,
			OriginTarget:   o.OriginTarget,
			Port:           o.Port,
		},
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
			Name:           o.Name,
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

func v1PoolResponse(c echo.Context, p *models.Pool) error {
	originSlice := make(originSlice, len(p.R.Origins))
	for j, o := range p.R.Origins {
		originSlice[j] = &origin{
			CreatedAt:      o.CreatedAt,
			UpdatedAt:      o.UpdatedAt,
			DeletedAt:      o.DeletedAt.Ptr(),
			ID:             o.OriginID,
			Name:           o.Name,
			Port:           o.Port,
			OriginTarget:   o.OriginTarget,
			OriginDisabled: o.OriginUserSettingDisabled,
		}
	}

	return c.JSON(http.StatusOK, &response{
		Version: apiVersion,
		Kind:    "poolsGet",
		Pool: &pool{
			CreatedAt: p.CreatedAt,
			UpdatedAt: p.UpdatedAt,
			DeletedAt: p.DeletedAt.Ptr(),
			ID:        p.PoolID,
			Name:      p.Name,
			Protocol:  p.Protocol,
			TenantID:  p.TenantID,
			Origins:   originSlice,
		},
	})
}

func v1PoolsResponse(c echo.Context, ps models.PoolSlice) error {
	out := make(poolSlice, len(ps))

	for i, p := range ps {
		originSlice := make(originSlice, len(p.R.Origins))
		for j, o := range p.R.Origins {
			originSlice[j] = &origin{
				CreatedAt:      o.CreatedAt,
				UpdatedAt:      o.UpdatedAt,
				DeletedAt:      o.DeletedAt.Ptr(),
				ID:             o.OriginID,
				Name:           o.Name,
				Port:           o.Port,
				OriginTarget:   o.OriginTarget,
				OriginDisabled: o.OriginUserSettingDisabled,
			}
		}

		out[i] = &pool{
			CreatedAt: p.CreatedAt,
			UpdatedAt: p.UpdatedAt,
			DeletedAt: p.DeletedAt.Ptr(),
			ID:        p.PoolID,
			Name:      p.Name,
			Protocol:  p.Protocol,
			TenantID:  p.TenantID,
			Origins:   originSlice,
		}
	}

	return c.JSON(http.StatusOK, &response{
		Version: apiVersion,
		Kind:    "poolsList",
		Pools:   &out,
	})
}
