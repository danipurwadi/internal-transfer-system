package tests

import (
	"fmt"
	"os"
	"testing"

	"github.com/danipurwadi/internal-transfer-system/app/api/apptest"
	"github.com/danipurwadi/internal-transfer-system/app/api/middleware"
	"github.com/danipurwadi/internal-transfer-system/app/transferapp"
	"github.com/danipurwadi/internal-transfer-system/business/api/dbtest"
	"github.com/danipurwadi/internal-transfer-system/business/transferbus"
	"github.com/danipurwadi/internal-transfer-system/business/transferbus/stores/transferdb"
	"github.com/danipurwadi/internal-transfer-system/foundation/docker"
	"github.com/danipurwadi/internal-transfer-system/foundation/web"
)

var c *docker.Container

func TestMain(m *testing.M) {
	code, err := run(m)
	if err != nil {
		fmt.Println(err)
	}

	os.Exit(code)
}

func run(m *testing.M) (int, error) {
	var err error

	c, err = dbtest.StartDB()
	if err != nil {
		return 1, err
	}
	defer dbtest.StopDB(c)

	return m.Run(), nil
}

func startTest(t *testing.T, testName string) *apptest.Test {
	db := dbtest.NewDatabase(t, c, testName)

	dbClient := transferdb.NewTxQueries(db.DB)
	// -------------------------------------------------------------------------
	// initialise business layer
	transferBus := transferbus.New(dbClient)

	// initialise app layer
	transferApp := transferapp.NewApp(transferBus)
	webClient := web.NewClient(middleware.Logger(db.Log), middleware.Errors(db.Log))
	transferApp.Routes(webClient)
	return apptest.New(db, webClient)
}
