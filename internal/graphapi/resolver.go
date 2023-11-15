package graphapi

import (
	"fmt"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/labstack/echo/v4"
	"github.com/wundergraph/graphql-go-tools/pkg/playground"
	"go.infratographer.com/x/gqlgenx/oteltracing"
	"go.uber.org/zap"

	ent "go.infratographer.com/load-balancer-api/internal/ent/generated"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

const (
	graphPath      = "query"
	playgroundPath = "playground"
)

var graphFullPath = fmt.Sprintf("/%s", graphPath)

// Resolver provides a graph response resolver
type Resolver struct {
	client   *ent.Client
	logger   *zap.SugaredLogger
	metadata Metadata
}

// Option is a function that modifies a resolver
type Option func(*Resolver)

// NewResolver returns a resolver configured with the given ent client
func NewResolver(client *ent.Client, logger *zap.SugaredLogger, opts ...Option) *Resolver {
	r := &Resolver{
		client: client,
		logger: logger,
	}

	for _, opt := range opts {
		opt(r)
	}

	return r
}

// TODO - @rizzza - This should be an all-purpose supergraph client

// WithMetadataClient sets the metadata client on the resolver
func WithMetadataClient(m Metadata) func(*Resolver) {
	return func(r *Resolver) {
		r.metadata = m
	}
}

// Handler is an http handler wrapping a Resolver
type Handler struct {
	r              *Resolver
	graphqlHandler http.Handler
	playground     *playground.Playground
	middleware     []echo.MiddlewareFunc
}

// Handler returns an http handler for a graph resolver
func (r *Resolver) Handler(withPlayground bool, middleware ...echo.MiddlewareFunc) *Handler {
	srv := handler.NewDefaultServer(
		NewExecutableSchema(
			Config{
				Resolvers: r,
			},
		),
	)

	srv.Use(oteltracing.Tracer{})

	h := &Handler{
		r:              r,
		middleware:     middleware,
		graphqlHandler: srv,
	}

	if withPlayground {
		h.playground = playground.New(playground.Config{
			PathPrefix:          "/",
			PlaygroundPath:      playgroundPath,
			GraphqlEndpointPath: graphFullPath,
		})
	}

	return h
}

// Handler returns the http.HandlerFunc for the GraphAPI
func (h *Handler) Handler() http.HandlerFunc {
	return h.graphqlHandler.ServeHTTP
}

// Routes ...
func (h *Handler) Routes(e *echo.Group) {
	e.Use(h.middleware...)

	e.POST("/"+graphPath, func(c echo.Context) error {
		h.graphqlHandler.ServeHTTP(c.Response(), c.Request())
		return nil
	})

	if h.playground != nil {
		handlers, err := h.playground.Handlers()
		if err != nil {
			h.r.logger.Fatal("error configuring playground handlers", "error", err)
			return
		}

		for i := range handlers {
			// with the function we need to dereference the handler so that it remains
			// the same in the function below
			hCopy := handlers[i].Handler

			e.GET(handlers[i].Path, func(c echo.Context) error {
				hCopy.ServeHTTP(c.Response(), c.Request())
				return nil
			})
		}
	}
}
