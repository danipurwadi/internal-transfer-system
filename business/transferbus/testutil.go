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
	balance := rand.Float64() * 100
	for i := range n {
		idx++

		nu := NewAccount{
			AccountId:      idx,
			InitialBalance: decimal.NewFromFloat(balance),
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
