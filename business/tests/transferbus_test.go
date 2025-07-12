package tests

import (
	"context"
	"fmt"
	"runtime/debug"
	"testing"

	"github.com/danipurwadi/internal-transfer-system/business/api/dbtest"
	"github.com/danipurwadi/internal-transfer-system/business/transferbus"
	"github.com/danipurwadi/internal-transfer-system/foundation/unittest"
	"github.com/google/go-cmp/cmp"
	"github.com/shopspring/decimal"
)

func Test_Transfer(t *testing.T) {
	t.Parallel()
	db := dbtest.NewDatabase(t, c, "Test_Transfer")
	defer func() {
		if r := recover(); r != nil {
			t.Log(r)
			t.Error(string(debug.Stack()))
		}
		db.Teardown()
	}()

	_, err := accountSeedData(db)
	if err != nil {
		t.Fatalf("Seeding error: %s", err)
	}

	unittest.Run(t, accountCreate(db), "account-create")
}

func accountSeedData(db *dbtest.Database) (dbtest.SeedData, error) {
	ctx := context.Background()
	busDomain := db.BusDomain

	// -------------------------------------------------------------------------
	accs, err := transferbus.TestSeedAccounts(ctx, 2, busDomain.TransferBus)
	if err != nil {
		return dbtest.SeedData{}, fmt.Errorf("seeding users : %w", err)
	}

	tu1 := dbtest.Account{
		Account: accs[0],
	}

	tu2 := dbtest.Account{
		Account: accs[1],
	}

	// -------------------------------------------------------------------------

	sd := dbtest.SeedData{
		Accounts: []dbtest.Account{tu1, tu2},
	}

	return sd, nil
}

func accountCreate(db *dbtest.Database) []unittest.Table {
	table := []unittest.Table{
		{
			Name: "basic",
			ExpResp: transferbus.Account{
				AccountId: 1,
				Balance:   decimal.NewFromFloat(100.12345),
			},
			ExcFunc: func(ctx context.Context) any {
				nu := transferbus.NewAccount{
					AccountId:      1,
					InitialBalance: decimal.NewFromFloat(100.12345),
				}

				resp, err := db.BusDomain.TransferBus.CreateAccount(ctx, nu)
				if err != nil {
					return err
				}

				return resp
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, exists := got.(transferbus.Account)
				if !exists {
					return "error occurred"
				}

				expResp := exp.(transferbus.Account)

				expResp.CreatedDate = gotResp.CreatedDate
				expResp.LastModifiedDate = gotResp.LastModifiedDate

				return cmp.Diff(gotResp, expResp)
			},
		},
		{
			Name: "rounding-up",
			ExpResp: transferbus.Account{
				AccountId: 2,
				Balance:   decimal.NewFromFloat(100.12346),
			},
			ExcFunc: func(ctx context.Context) any {
				nu := transferbus.NewAccount{
					AccountId:      2,
					InitialBalance: decimal.NewFromFloat(100.1234599999),
				}

				resp, err := db.BusDomain.TransferBus.CreateAccount(ctx, nu)
				if err != nil {
					return err
				}

				return resp
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, exists := got.(transferbus.Account)
				if !exists {
					return "error occurred"
				}

				expResp := exp.(transferbus.Account)

				expResp.CreatedDate = gotResp.CreatedDate
				expResp.LastModifiedDate = gotResp.LastModifiedDate

				return cmp.Diff(gotResp, expResp)
			},
		},
		{
			Name: "rounding-down",
			ExpResp: transferbus.Account{
				AccountId: 3,
				Balance:   decimal.NewFromFloat(100.12345),
			},
			ExcFunc: func(ctx context.Context) any {
				nu := transferbus.NewAccount{
					AccountId:      3,
					InitialBalance: decimal.NewFromFloat(100.1234549999),
				}

				resp, err := db.BusDomain.TransferBus.CreateAccount(ctx, nu)
				if err != nil {
					return err
				}

				return resp
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, exists := got.(transferbus.Account)
				if !exists {
					return "error occurred"
				}

				expResp := exp.(transferbus.Account)

				expResp.CreatedDate = gotResp.CreatedDate
				expResp.LastModifiedDate = gotResp.LastModifiedDate

				return cmp.Diff(gotResp, expResp)
			},
		},
	}

	return table
}
