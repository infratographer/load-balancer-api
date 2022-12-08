package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type response struct {
	PageSize      int            `json:"page_size,omitempty"`
	Page          int            `json:"page,omitempty"`
	PageCount     int            `json:"page_count,omitempty"`
	TotalPages    int            `json:"total_pages,omitempty"`
	TotalCount    int64          `json:"total_count,omitempty"`
	Links         *responseLinks `json:"_links,omitempty"`
	Version       string         `json:"version,omitempty"`
	Message       string         `json:"message,omitempty"`
	Error         string         `json:"error,omitempty"`
	Slug          string         `json:"slug,omitempty"`
	Location      interface{}    `json:"location,omitempty"`
	Locations     interface{}    `json:"locations,omitempty"`
	LoadBalancer  interface{}    `json:"load_balancer,omitempty"`
	LoadBalancers interface{}    `json:"load_balancers,omitempty"`
}

// responseLinks represent links that could be returned on a page
type responseLinks struct {
	Self     *link `json:"self,omitempty"`
	First    *link `json:"first,omitempty"`
	Previous *link `json:"previous,omitempty"`
	Next     *link `json:"next,omitempty"`
	Last     *link `json:"last,omitempty"`
}

// link represents an address to a page
type link struct {
	Href string `json:"href,omitempty"`
}

// newResponse creates a new response
func newResponse(message string) *response {
	return &response{
		Message: message,
		Version: Version,
	}
}

// notFoundResponse writes a 404 response with the given message
func notFoundResponse(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, newResponse(message))
}

func badRequestResponse(c *gin.Context, message string, err error) {
	resp := newResponse(message)
	resp.Error = err.Error()
	c.JSON(http.StatusBadRequest, resp)
}

func createdResponse(c *gin.Context) {
	uri := uriWithoutQueryParams(c)
	resp := newResponse("resource created")
	resp.Links = &responseLinks{
		Self: &link{Href: uri},
	}

	c.Header("Location", uri)
	c.JSON(http.StatusCreated, resp)
}

func deletedResponse(c *gin.Context) {
	c.JSON(http.StatusOK, newResponse("resource deleted"))
}

// func updatedResponse(c *gin.Context, slug string) {
// 	r := &RecordResponse{
// 		Message: "resource updated",
// 		Slug:    slug,
// 		Links: &RecordResponseLinks{
// 			Self: &Link{Href: uriWithoutQueryParams(c)},
// 		},
// 	}

// 	c.JSON(http.StatusOK, r)
// }

// func dbErrorResponse(c *gin.Context, err error) {
// 	if errors.Is(err, sql.ErrNoRows) {
// 		c.JSON(http.StatusNotFound, &response{Message: "resource not found", Error: err.Error()})
// 	} else {
// 		c.JSON(http.StatusInternalServerError, &response{Message: "datastore error", Error: err.Error()})
// 	}
// }

func uriWithoutQueryParams(c *gin.Context) string {
	uri := c.Request.URL
	uri.RawQuery = ""

	return uri.String()
}
