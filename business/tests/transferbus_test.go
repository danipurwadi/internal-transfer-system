package tests

import (
	"context"
	"fmt"
	"runtime/debug"
	"sort"
	"testing"
	"time"

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

	sd, err := accountSeedData(db)
	if err != nil {
		t.Fatalf("Seeding error: %s", err)
	}

	unittest.Run(t, accountCreation(db), "account-creation")
	unittest.Run(t, accountQuery(db, sd), "account-query")
	unittest.Run(t, transactionSubmission(db, sd), "transaction-submission")
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

func accountCreation(db *dbtest.Database) []unittest.Table {
	table := []unittest.Table{
		{
			Name: "basic",
			ExpResp: transferbus.Account{
				AccountID: 1,
				Balance:   decimal.NewFromFloat(100.12345),
			},
			ExcFunc: func(ctx context.Context) any {
				nu := transferbus.NewAccount{
					AccountID:      1,
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
			Name: "correctlyroundup",
			ExpResp: transferbus.Account{
				AccountID: 2,
				Balance:   decimal.NewFromFloat(100.12346),
			},
			ExcFunc: func(ctx context.Context) any {
				nu := transferbus.NewAccount{
					AccountID:      2,
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
			Name: "correctlyrounddown",
			ExpResp: transferbus.Account{
				AccountID: 3,
				Balance:   decimal.NewFromFloat(100.12345),
			},
			ExcFunc: func(ctx context.Context) any {
				nu := transferbus.NewAccount{
					AccountID:      3,
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
		{
			Name:    "negativebalance",
			ExpResp: transferbus.ErrNegativeBalance,
			ExcFunc: func(ctx context.Context) any {
				nu := transferbus.NewAccount{
					AccountID:      4,
					InitialBalance: decimal.NewFromFloat(-100.12345),
				}

				resp, err := db.BusDomain.TransferBus.CreateAccount(ctx, nu)
				if err != nil {
					return err
				}

				return resp
			},
			CmpFunc: func(got any, exp any) string {
				gotResp := got.(error).Error()
				expResp := exp.(error).Error()
				return cmp.Diff(gotResp, expResp)
			},
		},
	}

	return table
}

func accountQuery(db *dbtest.Database, sd dbtest.SeedData) []unittest.Table {
	accs := sd.Accounts

	sort.Slice(accs, func(i, j int) bool {
		return accs[i].AccountID <= accs[j].AccountID
	})

	table := []unittest.Table{
		{
			Name:    "basic",
			ExpResp: sd.Accounts[0].Account,
			ExcFunc: func(ctx context.Context) any {

				resp, err := db.BusDomain.TransferBus.GetBalance(ctx, sd.Accounts[0].AccountID)
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

				if gotResp.CreatedDate.Format(time.RFC3339) == expResp.CreatedDate.Format(time.RFC3339) {
					expResp.CreatedDate = gotResp.CreatedDate
				}

				if gotResp.LastModifiedDate.Format(time.RFC3339) == expResp.LastModifiedDate.Format(time.RFC3339) {
					expResp.LastModifiedDate = gotResp.LastModifiedDate
				}

				if gotResp.Balance.String() == expResp.Balance.String() {
					expResp.Balance = gotResp.Balance
				}

				return cmp.Diff(gotResp, expResp)
			},
		},
	}

	return table
}

func transactionSubmission(db *dbtest.Database, sd dbtest.SeedData) []unittest.Table {
	accs := sd.Accounts

	sort.Slice(accs, func(i, j int) bool {
		return accs[i].AccountID <= accs[j].AccountID
	})

	type accountBalances struct {
		SourceBalance      decimal.Decimal
		DestinationBalance decimal.Decimal
	}

	validAmount := decimal.NewFromFloat(12.12345)
	exceedAmount := accs[0].Balance.Add(decimal.NewFromInt(100))

	table := []unittest.Table{
		{
			Name: "basic",
			ExpResp: accountBalances{
				SourceBalance:      accs[0].Account.Balance.Sub(validAmount),
				DestinationBalance: accs[1].Account.Balance.Add(validAmount),
			},
			ExcFunc: func(ctx context.Context) any {
				r := transferbus.Transaction{
					SourceAccountID:      accs[0].AccountID,
					DestinationAccountID: accs[1].AccountID,
					Amount:               validAmount,
				}
				err := db.BusDomain.TransferBus.CreateTransaction(ctx, r)
				if err != nil {
					return err
				}

				acc1, err := db.BusDomain.TransferBus.GetBalance(ctx, accs[0].AccountID)
				if err != nil {
					return err
				}

				acc2, err := db.BusDomain.TransferBus.GetBalance(ctx, accs[1].AccountID)
				if err != nil {
					return err
				}

				return accountBalances{
					SourceBalance:      acc1.Balance,
					DestinationBalance: acc2.Balance,
				}
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, exists := got.(accountBalances)
				if !exists {
					return "error occurred"
				}

				expResp := exp.(accountBalances)
				return cmp.Diff(gotResp, expResp)
			},
		},
		{
			Name:    "insufficientbalance",
			ExpResp: transferbus.ErrInsufficientFunds,
			ExcFunc: func(ctx context.Context) any {
				r := transferbus.Transaction{
					SourceAccountID:      accs[0].AccountID,
					DestinationAccountID: accs[1].AccountID,
					Amount:               exceedAmount,
				}
				err := db.BusDomain.TransferBus.CreateTransaction(ctx, r)
				if err != nil {
					return err
				}
				return nil
			},
			CmpFunc: func(got any, exp any) string {
				gotResp := got.(error).Error()
				expResp := exp.(error).Error()
				return cmp.Diff(gotResp, expResp)
			},
		},
		{
			Name:    "invalidaccountid",
			ExpResp: transferbus.ErrAccNotFound,
			ExcFunc: func(ctx context.Context) any {
				// declare the new user id as the sum of all ids to guarantee id is not found
				invalidUserID := int64(0)
				for _, u := range sd.Accounts {
					invalidUserID += u.AccountID
				}

				r := transferbus.Transaction{
					SourceAccountID:      accs[0].AccountID,
					DestinationAccountID: invalidUserID,
					Amount:               validAmount,
				}
				err := db.BusDomain.TransferBus.CreateTransaction(ctx, r)
				if err != nil {
					return err
				}
				return nil
			},
			CmpFunc: func(got any, exp any) string {
				gotResp := got.(error).Error()
				expResp := exp.(error).Error()
				return cmp.Diff(gotResp, expResp)
			},
		},
	}
	return table
}
