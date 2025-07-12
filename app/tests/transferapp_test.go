package tests

import (
	"context"
	"fmt"
	"runtime/debug"
	"testing"

	"github.com/danipurwadi/internal-transfer-system/app/api/apptest"
	"github.com/danipurwadi/internal-transfer-system/business/api/dbtest"
	"github.com/danipurwadi/internal-transfer-system/business/transferbus"
)

func Test_Transfer_App(t *testing.T) {
	t.Parallel()

	// -------------------------------------------------------------------------

	apiTest := startTest(t, "Test_Transfer")
	defer func() {
		if r := recover(); r != nil {
			t.Log(r)
			t.Error(string(debug.Stack()))
		}
		apiTest.DB.Teardown()
	}()

	// -------------------------------------------------------------------------

	sd, err := userSeedData(apiTest.DB)
	if err != nil {
		t.Fatalf("Seeding error: %s", err)
	}

	// -------------------------------------------------------------------------

	apiTest.Run(t, accountCreation200(), "account-creation-200")
	apiTest.Run(t, accountCreation400(), "account-creation-400")
	apiTest.Run(t, accountCreation409(sd), "account-creation-409")

	apiTest.Run(t, accountQuery200(sd), "account-query-200")
	apiTest.Run(t, accountQuery400(), "account-query-400")
	apiTest.Run(t, accountQuery404(), "account-query-404")

	apiTest.Run(t, transactionSubmission201(sd), "transaction-submission-201")
	apiTest.Run(t, transactionSubmission400(sd), "transaction-submission-400")
	apiTest.Run(t, transactionSubmission404(sd), "transaction-submission-404")
}

func userSeedData(db *dbtest.Database) (apptest.SeedData, error) {
	ctx := context.Background()
	busDomain := db.BusDomain

	usrs, err := transferbus.TestSeedAccounts(ctx, 2, busDomain.TransferBus)
	if err != nil {
		return apptest.SeedData{}, fmt.Errorf("seeding users : %w", err)
	}

	tu1 := apptest.Account{
		Account: dbtest.Account{
			Account: usrs[0],
		},
	}

	tu2 := apptest.Account{
		Account: dbtest.Account{
			Account: usrs[1],
		},
	}
	// -------------------------------------------------------------------------

	sd := apptest.SeedData{
		Accounts: []apptest.Account{tu1, tu2},
	}

	return sd, nil
}
