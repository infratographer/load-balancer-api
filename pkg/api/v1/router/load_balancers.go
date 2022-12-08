package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"go.infratographer.sh/loadbalancerapi/pkg/api/v1/loadbalancers"
)

const (
	// LoadBalancersBaseURI is the path prefix for all loadbalancers endpoints
	LoadBalancersBaseURI = "/loadbalancers"

	// TenantLoadBalancersURI is the path for a loadbalancer by Tenant UUID
	TenantLoadBalancersURI = TenantPrefix + LoadBalancersBaseURI

	// TenantLoadBalancersIPAddressURI is the path for a loadbalancer by UUID and IP
	TenantLoadBalancersIPAddressURI = TenantLoadBalancersURI + "/:ip_address"
)

// addLoadBalancerRoutes adds the routes for this API version to a router group
func (r *Router) addLoadBalancerRoutes(rg *gin.RouterGroup) {
	rg.POST(LoadBalancersBaseURI, r.createLoadBalancer)

	// rg.GET(TenantLoadBalancersURI, r.getLoadBalancers)

	rg.GET(TenantLoadBalancersIPAddressURI, r.getLoadBalancerByIPAddress)
	rg.DELETE(TenantLoadBalancersIPAddressURI, r.deleteLoadBalancerByIPAddress)
}

func loadBalancerSuccess(c *gin.Context, obj interface{}) {
	uri := uriWithoutQueryParams(c)
	r := newResponse("resources found")
	r.LoadBalancer = obj
	r.Links = &responseLinks{
		Self: &link{Href: uri},
	}

	c.Header("Location", uri)
	c.JSON(http.StatusCreated, r)
}

func loadBalancersSuccess(c *gin.Context, obj interface{}) {
	uri := uriWithoutQueryParams(c)
	r := newResponse("resources found")
	r.LoadBalancers = obj
	r.Links = &responseLinks{
		Self: &link{Href: uri},
	}

	c.Header("Location", uri)
	c.JSON(http.StatusCreated, r)
}

// createLoadBalancer creates a new loadbalancer
func (r *Router) createLoadBalancer(c *gin.Context) {
	var lb loadbalancers.LoadBalancer
	if err := c.BindJSON(&lb); err != nil {
		badRequestResponse(c, loadbalancers.ErrInvalid.Error(), err)
		return
	}

	if err := lb.Create(c, r.db); err != nil {
		badRequestResponse(c, loadbalancers.ErrWrite.Error(), err)
		return
	}

	createdResponse(c)
}

// getLoadBalancerByIPAddress gets a loadbalancer by IP address
func (r *Router) getLoadBalancerByIPAddress(c *gin.Context) {
	tenant, err := uuid.Parse(c.Param("tenant"))
	if err != nil {
		badRequestResponse(c, "could not parse tenant", err)
		return
	}

	lb := loadbalancers.LoadBalancer{
		IPAddress: c.Param("ip_address"),
		TenantID:  tenant,
	}

	err = lb.Find(c, r.db)
	if err != nil {
		r.logger.Debugw("could not find loadbalancer", "error", err)

		notFoundResponse(c, err.Error())

		return
	}

	loadBalancerSuccess(c, lb)
}

// deleteLoadBalancerByIPAddress deletes a loadbalancer by IP address
func (r *Router) deleteLoadBalancerByIPAddress(c *gin.Context) {
	tenant, err := uuid.Parse(c.Param("tenant"))
	if err != nil {
		badRequestResponse(c, "could not parse tenant", err)
		return
	}

	lb := loadbalancers.LoadBalancer{
		IPAddress: c.Param("ip_address"),
		TenantID:  tenant,
	}

	err = lb.Delete(c, r.db)
	if err != nil {
		notFoundResponse(c, err.Error())

		return
	}

	deletedResponse(c)
}
