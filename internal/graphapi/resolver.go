package graphapi

import (
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/labstack/echo/v4"

	ent "go.infratographer.com/load-balancer-api/internal/ent/generated"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

// Resolver provides a graph response resolver
type Resolver struct {
	client *ent.Client
}

// Handler is an http handler wrapping a Resolver
type Handler struct {
	r                 *Resolver
	graphqlHandler    http.Handler
	playgroundHandler http.Handler
}

// NewHandler returns an http handler for a graph resolver
func NewHandler(client *ent.Client) *Handler {
	h := &Handler{
		r: &Resolver{
			client: client,
		},
	}

	h.graphqlHandler = handler.NewDefaultServer(
		NewExecutableSchema(
			Config{
				Resolvers: h.r,
			},
		),
	)
	h.playgroundHandler = playground.Handler("GraphQL", "/query")

	return h
}

// Routes ...
func (h *Handler) Routes(e *echo.Group) {
	e.POST("/query", func(c echo.Context) error {
		h.graphqlHandler.ServeHTTP(c.Response(), c.Request())
		return nil
	})

	e.GET("/playground", func(c echo.Context) error {
		h.playgroundHandler.ServeHTTP(c.Response(), c.Request())
		return nil
	})
}
