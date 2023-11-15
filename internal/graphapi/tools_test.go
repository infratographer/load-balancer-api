package graphapi_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"

	metadata "go.infratographer.com/metadata-api/pkg/client"
	"go.infratographer.com/metadata-api/pkg/client/mockmetadata"
	"go.infratographer.com/permissions-api/pkg/permissions"
	"go.infratographer.com/x/echojwtx"
	"go.infratographer.com/x/echox"

	ent "go.infratographer.com/load-balancer-api/internal/ent/generated"
	"go.infratographer.com/load-balancer-api/internal/graphapi"
	"go.infratographer.com/load-balancer-api/internal/graphclient"
	"go.infratographer.com/load-balancer-api/internal/manualhooks"
	"go.infratographer.com/load-balancer-api/internal/testutils"
)

const (
	ownerPrefix    = "testown"
	locationPrefix = "testloc"
	lbPrefix       = "loadbal"
)

var EntClient *ent.Client

func TestMain(m *testing.M) {
	// setup the database
	testutils.SetupDB()

	// assign package variables
	EntClient = testutils.EntClient

	// setup the resolver hooks
	manualhooks.PubsubHooks(EntClient)

	// run the tests
	code := m.Run()

	// teardown the database
	testutils.TeardownDB()

	// return the test response code
	os.Exit(code)
}

type graphClient struct {
	srvURL     string
	httpClient *http.Client
}

type graphClientOptions func(*graphClient)

func withGraphClientServerURL(url string) graphClientOptions {
	return func(g *graphClient) {
		g.srvURL = url
	}
}

func withGraphClientHTTPClient(httpcli *http.Client) graphClientOptions {
	return func(g *graphClient) {
		g.httpClient = httpcli
	}
}

func graphTestClient(options ...graphClientOptions) graphclient.GraphClient {
	metadataMock := new(mockmetadata.MockMetadata)
	metadataMock.On("StatusUpdate", mock.Anything, mock.Anything).Return(&metadata.StatusUpdate{}, nil)

	g := &graphClient{
		srvURL: "graph",
		httpClient: &http.Client{Transport: localRoundTripper{handler: handler.NewDefaultServer(
			graphapi.NewExecutableSchema(
				graphapi.Config{Resolvers: graphapi.NewResolver(EntClient, zap.NewNop().Sugar(), graphapi.WithMetadataClient(metadataMock))},
			))}},
	}

	for _, opt := range options {
		opt(g)
	}

	return graphclient.NewClient(g.httpClient, g.srvURL)
}

// localRoundTripper is an http.RoundTripper that executes HTTP transactions
// by using handler directly, instead of going over an HTTP connection.
type localRoundTripper struct {
	handler http.Handler
}

func (l localRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	w := httptest.NewRecorder()
	l.handler.ServeHTTP(w, req)

	return w.Result(), nil
}

type testServerConfig struct {
	echoConfig        echox.Config
	handlerMiddleware []echo.MiddlewareFunc
}

type testServerOption func(*testServerConfig) error

func withAuthConfig(authConfig *echojwtx.AuthConfig) testServerOption {
	return func(tsc *testServerConfig) error {
		auth, err := echojwtx.NewAuth(context.Background(), *authConfig)
		if err != nil {
			return err
		}

		tsc.echoConfig = tsc.echoConfig.WithMiddleware(auth.Middleware())

		return nil
	}
}

func withPermissions(options ...permissions.Option) testServerOption {
	return func(tsc *testServerConfig) error {
		perms, err := permissions.New(permissions.Config{}, options...)
		if err != nil {
			return err
		}

		tsc.handlerMiddleware = append(tsc.handlerMiddleware, perms.Middleware())

		return nil
	}
}

func newTestServer(options ...testServerOption) (*httptest.Server, error) {
	tsc := new(testServerConfig)

	for _, opt := range options {
		if err := opt(tsc); err != nil {
			return nil, err
		}
	}

	srv, err := echox.NewServer(zap.NewNop(), tsc.echoConfig, nil)
	if err != nil {
		return nil, err
	}

	r := graphapi.NewResolver(EntClient, zap.NewNop().Sugar())
	srv.AddHandler(r.Handler(false, tsc.handlerMiddleware...))

	return httptest.NewServer(srv.Handler()), nil
}

func newString(s string) *string {
	return &s
}

func newBool(b bool) *bool {
	return &b
}

func newInt64(i int) *int64 {
	r := int64(i)
	return &r
}
