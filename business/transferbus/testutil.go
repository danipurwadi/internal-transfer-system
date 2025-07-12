package transferbus

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/shopspring/decimal"
)

// TestNewAccounts is a helper method for testing.
func TestNewAccounts(n int) []NewAccount {
	newUsrs := make([]NewAccount, n)

	idx := rand.Int63n(10000)
	balance := rand.Int63n(100)
	diff := int64(10) // determines the balance difference between acc[i] and acc[i-1]
	for i := range n {
		idx++

		nu := NewAccount{
			AccountId:      idx,
			InitialBalance: decimal.NewFromInt(balance + diff*idx),
		}

		newUsrs[i] = nu
	}

	return newUsrs
}

// TestSeedAccounts is a helper method for testing.
func TestSeedAccounts(ctx context.Context, n int, api *Bus) ([]Account, error) {
	newUsrs := TestNewAccounts(n)

	accs := make([]Account, len(newUsrs))
	for i, nu := range newUsrs {
		acc, err := api.CreateAccount(ctx, nu)
		if err != nil {
			return nil, fmt.Errorf("seeding user: idx: %d : %w", i, err)
		}

		accs[i] = acc
	}

	return accs, nil
}
