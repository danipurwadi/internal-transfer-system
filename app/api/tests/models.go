package apptest

import "github.com/danipurwadi/internal-transfer-system/business/api/dbtest"

// User extends the dbtest user for api test support.
type User struct {
	dbtest.Account
}

// SeedData represents users for api tests.
type SeedData struct {
	Users []User
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
