package dbtest

import "github.com/danipurwadi/internal-transfer-system/business/transferbus"

// User represents an app user specified for the test.
type User struct {
	Accounts     []transferbus.AccountCreation
	Transactions []transferbus.Transaction
}

// SeedData represents data that was seeded for the test.
type SeedData struct {
	Users []User
}
