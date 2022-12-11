// Package dbtools helps test code that interacts with the database
package dbtools

import (
	"os"
	"testing"

	// import the crdbpgx for automatic retries of errors for crdb that support retry
	_ "github.com/cockroachdb/cockroach-go/v2/crdb/crdbpgx"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Register the Postgres driver.
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestDBURI is the URI for the test database
var TestDBURI = os.Getenv("LOADBALANCERAPI_CRDB_URI")
var testDB *sqlx.DB

func testDatastore(t *testing.T) error {
	// don't setup the datastore if we already have one
	if testDB != nil {
		return nil
	}

	// Uncomment when you are having database issues with your tests and need to see the db logs
	// Hidden by default because it can be noisy and make it harder to read normal failures.
	// You can also enable at the beginning of your test and then disable it again at the end
	// boil.DebugMode = true

	db, err := sqlx.Open("postgres", TestDBURI)
	if err != nil {
		return err
	}

	testDB = db

	return nil
}

// DatabaseTest allows you to run tests that interact with the database
func DatabaseTest(t *testing.T) *sqlx.DB {
	if testing.Short() {
		t.Skip("skipping database test in short mode")
	}

	err := testDatastore(t)
	require.NoError(t, err, "Unexpected error getting connection to test datastore")

	return testDB
}

// CleanUpTables deletes all rows from the specified tables for the specified tenant, this should be run per test that uses DatabaseTest
func CleanUpTables(t *testing.T, tenantID uuid.UUID, tables ...string) {
	for _, table := range tables {
		res, err := testDB.Exec("DELETE FROM "+table+" WHERE tenant_id = $1", tenantID.String())
		assert.Nil(t, err)
		total, err := res.RowsAffected()
		assert.Nil(t, err)
		assert.Equal(t, int64(1), total)
	}
}
