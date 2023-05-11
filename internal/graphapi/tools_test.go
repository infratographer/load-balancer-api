//go:build testtools
// +build testtools

package graphapi_test

import (
	"bytes"
	"context"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"entgo.io/ent/dialect"
	"github.com/99designs/gqlgen/graphql/handler"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"go.uber.org/zap"

	"go.infratographer.com/x/echojwtx"
	"go.infratographer.com/x/echox"

	ent "go.infratographer.com/load-balancer-api/internal/ent/generated"
	"go.infratographer.com/load-balancer-api/internal/graphapi"
	"go.infratographer.com/load-balancer-api/internal/graphclient"
	"go.infratographer.com/load-balancer-api/x/testcontainersx"
)

const (
	tenantPrefix   = "testtnt"
	locationPrefix = "testloc"
	lbPrefix       = "loadbal"
)

var (
	TestDBURI   = os.Getenv("LOADBALANCERAPI_TESTDB_URI")
	EntClient   *ent.Client
	DBContainer *testcontainersx.DBContainer
)

func TestMain(m *testing.M) {
	// setup the database if needed
	setupDB()
	// run the tests
	code := m.Run()
	// teardown the database
	teardownDB()
	// return the test response code
	os.Exit(code)
}

func parseDBURI(ctx context.Context) (string, string, *testcontainersx.DBContainer) {
	switch {
	// if you don't pass in a database we default to an in memory sqlite
	case TestDBURI == "":
		return dialect.SQLite, "file:ent?mode=memory&cache=shared&_fk=1", nil
	case strings.HasPrefix(TestDBURI, "sqlite://"):
		return dialect.SQLite, strings.TrimPrefix(TestDBURI, "sqlite://"), nil
	case strings.HasPrefix(TestDBURI, "postgres://"), strings.HasPrefix(TestDBURI, "postgresql://"):
		return dialect.Postgres, TestDBURI, nil
	case strings.HasPrefix(TestDBURI, "docker://"):
		dbImage := strings.TrimPrefix(TestDBURI, "docker://")

		switch {
		case strings.HasPrefix(dbImage, "cockroach"), strings.HasPrefix(dbImage, "cockroachdb"), strings.HasPrefix(dbImage, "crdb"):
			cntr, err := testcontainersx.NewCockroachDB(ctx, dbImage)
			errPanic("error starting db test container", err)

			return dialect.Postgres, cntr.URI, cntr
		case strings.HasPrefix(dbImage, "postgres"):
			cntr, err := testcontainersx.NewPostgresDB(ctx, dbImage,
				postgres.WithInitScripts(filepath.Join("testdata", "postgres_init.sh")),
			)
			errPanic("error starting db test container", err)

			return dialect.Postgres, cntr.URI, cntr
		default:
			panic("invalid testcontainer URI, uri: " + TestDBURI)
		}

	default:
		panic("invalid DB URI, uri: " + TestDBURI)
	}
}

func setupDB() {
	// don't setup the datastore if we already have one
	if EntClient != nil {
		return
	}

	ctx := context.Background()

	dia, uri, cntr := parseDBURI(ctx)

	c, err := ent.Open(dia, uri, ent.Debug())
	if err != nil {
		errPanic("failed terminating test db container after failing to connect to the db", cntr.Container.Terminate(ctx))
		errPanic("failed opening connection to database:", err)
	}

	switch dia {
	case dialect.SQLite:
		// Run automatic migrations for SQLite
		errPanic("failed creating db scema", c.Schema.Create(ctx))
	case dialect.Postgres:
		log.Println("Running database migrations")

		cmd := exec.Command("atlas", "migrate", "apply",
			"--dir", "file://../../db/migrations",
			"--url", uri,
		)

		// write all output to stdout and stderr as it comes through
		var stdBuffer bytes.Buffer
		mw := io.MultiWriter(os.Stdout, &stdBuffer)

		cmd.Stdout = mw
		cmd.Stderr = mw

		// Execute the command
		errPanic("atlas returned an error running database migrations", cmd.Run())
	}

	EntClient = c
}

func teardownDB() {
	ctx := context.Background()

	if EntClient != nil {
		errPanic("teardown failed to close database connection", EntClient.Close())
	}

	if DBContainer != nil {
		errPanic("teardown failed to terminate test db container", DBContainer.Container.Terminate(ctx))
	}
}

func errPanic(msg string, err error) {
	if err != nil {
		log.Panicf("%s err: %s", msg, err.Error())
	}
}

type graphTestClient struct {
	srvURL     string
	httpClient *http.Client
}

type graphClientOptions func(*graphTestClient)

func withGraphClientServerURL(url string) graphClientOptions {
	return func(g *graphTestClient) {
		g.srvURL = url
	}
}

func withGraphClientHTTPClient(httpcli *http.Client) graphClientOptions {
	return func(g *graphTestClient) {
		g.httpClient = httpcli
	}
}

func newGraphTestClient(options ...graphClientOptions) graphclient.GraphClient {
	g := &graphTestClient{
		srvURL: "graph",
		httpClient: &http.Client{Transport: localRoundTripper{handler: handler.NewDefaultServer(
			graphapi.NewExecutableSchema(
				graphapi.Config{Resolvers: graphapi.NewResolver(EntClient, zap.NewNop().Sugar())},
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

func newTestServer(authConfig *echojwtx.AuthConfig) (*httptest.Server, error) {
	echoCfg := echox.Config{}

	if authConfig != nil {
		auth, err := echojwtx.NewAuth(context.Background(), *authConfig)
		if err != nil {
			return nil, err
		}

		echoCfg = echoCfg.WithMiddleware(auth.Middleware())
	}

	srv, err := echox.NewServer(zap.NewNop(), echoCfg, nil)
	if err != nil {
		return nil, err
	}

	r := graphapi.NewResolver(EntClient, zap.NewNop().Sugar())
	srv.AddHandler(r.Handler(false))

	return httptest.NewServer(srv.Handler()), nil
}

func newString(s string) *string {
	return &s
}

func newBool(b bool) *bool {
	return &b
}

func newInt64(i int64) *int64 {
	return &i
}
