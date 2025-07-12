package dbtest

import (
	"bytes"
	"context"
	"fmt"

	"math/rand"
	"testing"
	"time"

	"github.com/danipurwadi/internal-transfer-system/business/api/db"
	"github.com/danipurwadi/internal-transfer-system/business/transferbus"
	"github.com/danipurwadi/internal-transfer-system/business/transferbus/stores/transferdb"
	"github.com/danipurwadi/internal-transfer-system/foundation/docker"
	"github.com/danipurwadi/internal-transfer-system/foundation/logger"
	"github.com/danipurwadi/internal-transfer-system/foundation/web"
	"github.com/jackc/pgx/v5/pgxpool"
)

// StartDB starts a database instance.
func StartDB() (*docker.Container, error) {
	image := "postgres:14.8"
	port := "5432"
	dockerArgs := []string{"-e", "POSTGRES_PASSWORD=postgres"}
	appArgs := []string{"-c", "log_statement=all"}

	c, err := docker.StartContainer(image, port, dockerArgs, appArgs)
	if err != nil {
		return nil, fmt.Errorf("starting container: %w", err)
	}

	fmt.Printf("Image:       %s\n", image)
	fmt.Printf("ContainerID: %s\n", c.ID)
	fmt.Printf("HostPort:    %s\n", c.HostPort)

	return c, nil
}

// StopDB stops a running database instance.
func StopDB(c *docker.Container) {
	docker.StopContainer(c.ID)
	fmt.Println("Stopped:", c.ID)
}

// =============================================================================

// BusDomain represents all the business domain apis needed for testing.
type BusDomain struct {
	TransferBus *transferbus.Bus
}

func newBusDomains(db *pgxpool.Pool) BusDomain {
	dbClient := transferdb.NewTxQueries(db)
	transferBus := transferbus.New(dbClient)

	return BusDomain{
		TransferBus: transferBus,
	}
}

// =============================================================================

// Database owns state for running and shutting down tests.
type Database struct {
	DB        *pgxpool.Pool
	Log       *logger.Logger
	BusDomain BusDomain
	Teardown  func()
}

// NewDatabase creates a test database inside a Docker container. It creates the
// required table structure but the database is otherwise empty. It returns
// the database to use as well as a function to call at the end of the test.
func NewDatabase(t *testing.T, c *docker.Container, testName string) *Database {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dbMConfig := db.Config{
		User:       "postgres",
		Password:   "postgres",
		HostPort:   c.HostPort,
		Database:   "postgres",
		DisableTLS: true,
	}
	dbM := db.New(dbMConfig)
	// -------------------------------------------------------------------------

	const letterBytes = "abcdefghijklmnopqrstuvwxyz"
	b := make([]byte, 4)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	dbName := string(b)

	err := db.InitDatabase(ctx, dbMConfig, dbName)
	if err != nil {
		t.Fatalf("creating database %s: %v", dbName, err)
	}
	dbM.Close()

	// -------------------------------------------------------------------------

	dbConfig := db.Config{
		User:       "postgres",
		Password:   "postgres",
		HostPort:   c.HostPort,
		Database:   dbName,
		DisableTLS: true,
	}

	testDb := db.New(dbConfig)
	db.Migrate(dbConfig)

	// -------------------------------------------------------------------------

	var buf bytes.Buffer
	log := logger.New(&buf, logger.LevelInfo, "TEST", func(context.Context) string { return web.GetTraceID(ctx) })

	// -------------------------------------------------------------------------

	// teardown is the function that should be invoked when the caller is done
	// with the database.
	teardown := func() {
		t.Helper()

		testDb.Close()

		fmt.Printf("******************** LOGS (%s) ********************\n", testName)
		fmt.Print(buf.String())
		fmt.Printf("******************** LOGS (%s) ********************\n", testName)
	}

	return &Database{
		DB:        testDb,
		Log:       log,
		BusDomain: newBusDomains(testDb),
		Teardown:  teardown,
	}
}

// =============================================================================

// StringPointer is a helper to get a *string from a string. It is in the tests
// package because we normally don't want to deal with pointers to basic types
// but it's useful in some tests.
func StringPointer(s string) *string {
	return &s
}

// IntPointer is a helper to get a *int from a int. It is in the tests package
// because we normally don't want to deal with pointers to basic types but it's
// useful in some tests.
func IntPointer(i int) *int {
	return &i
}

// FloatPointer is a helper to get a *float64 from a float64. It is in the tests
// package because we normally don't want to deal with pointers to basic types
// but it's useful in some tests.
func FloatPointer(f float64) *float64 {
	return &f
}

// BoolPointer is a helper to get a *bool from a bool. It is in the tests package
// because we normally don't want to deal with pointers to basic types but it's
// useful in some tests.
func BoolPointer(b bool) *bool {
	return &b
}
