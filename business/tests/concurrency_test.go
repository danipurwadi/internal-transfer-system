
//go:build go1.18
// +build go1.18

package tests

import (
	"context"
	"sync"
	"testing"

	"github.com/danipurwadi/internal-transfer-system/business/api/dbtest"
	"github.com/danipurwadi/internal-transfer-system/business/transferbus"
	"github.com/google/go-cmp/cmp"
	"github.com/shopspring/decimal"
)

// Test_Concurrent_Transfers simulates high concurrency on a single account to check for race conditions.
func Test_Concurrent_Transfers(t *testing.T) {
	t.Parallel()

	db := dbtest.NewDatabase(t, c, "Test_Concurrent_Transfers")
	defer db.Teardown()

	busDomain := db.BusDomain
	ctx := context.Background()

	// 1. SETUP: Create two accounts with a known starting balance.
	initialBalance := decimal.NewFromInt(10000)
	transferAmount := decimal.NewFromFloat(10.12345)

	acc1, err := busDomain.TransferBus.CreateAccount(ctx, transferbus.NewAccount{
		AccountID:      1,
		InitialBalance: initialBalance,
	})
	if err != nil {
		t.Fatalf("Failed to create account 1: %v", err)
	}

	acc2, err := busDomain.TransferBus.CreateAccount(ctx, transferbus.NewAccount{
		AccountID:      2,
		InitialBalance: decimal.Zero,
	})
	if err != nil {
		t.Fatalf("Failed to create account 2: %v", err)
	}

	// 2. EXECUTE: Launch a large number of concurrent transfers from account 1 to account 2.
	const numConcurrentTransfers = 99 
	var wg sync.WaitGroup
	wg.Add(numConcurrentTransfers)

	for i := 0; i < numConcurrentTransfers; i++ {
		go func() {
			defer wg.Done()
			err := busDomain.TransferBus.CreateTransaction(ctx, transferbus.Transaction{
				SourceAccountID:      acc1.AccountID,
				DestinationAccountID: acc2.AccountID,
				Amount:               transferAmount,
			})
			if err != nil {
				// We can log the error, but we don't want to fail the test here,
				// as some transaction failures (like temporary deadlocks) might be expected under high load.
				// The final balance check is the source of truth.
				t.Logf("Transaction failed (this might be ok): %v", err)
			}
		}()
	}

	wg.Wait()

	// 3. VERIFY: Check the final balances.
	finalAcc1, err := busDomain.TransferBus.GetBalance(ctx, acc1.AccountID)
	if err != nil {
		t.Fatalf("Failed to get final balance for account 1: %v", err)
	}

	finalAcc2, err := busDomain.TransferBus.GetBalance(ctx, acc2.AccountID)
	if err != nil {
		t.Fatalf("Failed to get final balance for account 2: %v", err)
	}

	// The total amount transferred should be numConcurrentTransfers * transferAmount.
	expectedAmountTransferred := transferAmount.Mul(decimal.NewFromInt(numConcurrentTransfers))

	expectedFinalAcc1Balance := initialBalance.Sub(expectedAmountTransferred)
	expectedFinalAcc2Balance := expectedAmountTransferred

	if !finalAcc1.Balance.Equal(expectedFinalAcc1Balance) {
		t.Errorf("Account 1 final balance is incorrect. Got %s, Expected %s", finalAcc1.Balance, expectedFinalAcc1Balance)
		t.Error(cmp.Diff(finalAcc1.Balance, expectedFinalAcc1Balance))
	}

	if !finalAcc2.Balance.Equal(expectedFinalAcc2Balance) {
		t.Errorf("Account 2 final balance is incorrect. Got %s, Expected %s", finalAcc2.Balance, expectedFinalAcc2Balance)
		t.Error(cmp.Diff(finalAcc2.Balance, expectedFinalAcc2Balance))
	}

	t.Logf("Final balance for account 1 is correct: %s", finalAcc1.Balance)
	t.Logf("Final balance for account 2 is correct: %s", finalAcc2.Balance)
}
