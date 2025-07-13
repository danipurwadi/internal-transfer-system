package apptest

import "github.com/danipurwadi/internal-transfer-system/business/api/dbtest"

// Account represents a test account.
type Account struct {
	dbtest.Account
}

// SeedData represents users for api tests.
type SeedData struct {
	Accounts []Account
}

// Table represent fields needed for running an api test.
type Table struct {
	Name       string
	URL        string
	Method     string
	StatusCode int
	Input      any
	GotResp    any
	ExpResp    any
	CmpFunc    func(got any, exp any) string
}
