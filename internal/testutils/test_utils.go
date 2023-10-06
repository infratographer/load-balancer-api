// Package testutils provides some utilities that may be useful for testing
package testutils

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"entgo.io/ent/dialect"
	_ "github.com/lib/pq"           // used by the ent client using ParseDBURI return values
	_ "github.com/mattn/go-sqlite3" // used by the ent client using ParseDBURI return values
	"github.com/testcontainers/testcontainers-go/modules/postgres"

	"go.infratographer.com/load-balancer-api/x/testcontainersx"
)

var testDBURI = os.Getenv("LOADBALANCERAPI_TESTDB_URI")

// ParseDBURI parses the kind of query language from TESTDB_URI env var and initializes DBContainer as required
func ParseDBURI(ctx context.Context) (string, string, *testcontainersx.DBContainer) {
	switch {
	// if you don't pass in a database we default to an in memory sqlite
	case testDBURI == "":
		return dialect.SQLite, "file:ent?mode=memory&cache=shared&_fk=1", nil
	case strings.HasPrefix(testDBURI, "sqlite://"):
		return dialect.SQLite, strings.TrimPrefix(testDBURI, "sqlite://"), nil
	case strings.HasPrefix(testDBURI, "postgres://"), strings.HasPrefix(testDBURI, "postgresql://"):
		return dialect.Postgres, testDBURI, nil
	case strings.HasPrefix(testDBURI, "docker://"):
		dbImage := strings.TrimPrefix(testDBURI, "docker://")

		switch {
		case strings.HasPrefix(dbImage, "cockroach"), strings.HasPrefix(dbImage, "cockroachdb"), strings.HasPrefix(dbImage, "crdb"):
			cntr, err := testcontainersx.NewCockroachDB(ctx, dbImage)
			IfErrPanic("error starting db test container", err)

			return dialect.Postgres, cntr.URI, cntr
		case strings.HasPrefix(dbImage, "postgres"):
			_, b, _, _ := runtime.Caller(0)
			initScriptPath := filepath.Join(filepath.Dir(b), "testdata", "postgres_init.sh")

			cntr, err := testcontainersx.NewPostgresDB(ctx, dbImage,
				postgres.WithInitScripts(initScriptPath),
			)
			IfErrPanic("error starting db test container", err)

			return dialect.Postgres, cntr.URI, cntr
		default:
			panic("invalid testcontainer URI, uri: " + testDBURI)
		}

	default:
		panic("invalid DB URI, uri: " + testDBURI)
	}
}

// IfErrPanic conditionally panics on err with msg
func IfErrPanic(msg string, err error) {
	if err != nil {
		log.Panicf("%s err: %s", msg, err.Error())
	}
}

// ChannelReceiveWithTimeout returns the next message from channel chan or panics if it timesout before
func ChannelReceiveWithTimeout[T any](channel <-chan T, timeout time.Duration) (msg T) {
	select {
	case msg = <-channel:
	case <-time.After(timeout):
		log.Panicln("timed out waiting to receive from channel")
	}

	return
}
