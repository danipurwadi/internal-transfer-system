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

	idx := rand.Intn(10000)
	for i := 0; i < n; i++ {
		idx++

		nu := NewAccount{
			AccountId:      int64(idx),
			InitialBalance: decimal.NewFromInt(1000), // TODO: randomise this value
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
