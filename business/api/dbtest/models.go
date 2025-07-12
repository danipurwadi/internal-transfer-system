package dbtest

import "github.com/danipurwadi/internal-transfer-system/business/transferbus"

// Account represents an app user specified for the test.
type Account struct {
	transferbus.Account
	Transactions []transferbus.Transaction
}

// SeedData represents data that was seeded for the test.
type SeedData struct {
	Accounts []Account
}
