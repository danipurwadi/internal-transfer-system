package transferbus

import (
	"time"

	transferdbgen "github.com/danipurwadi/internal-transfer-system/business/transferbus/stores/transferdb/gen"
	"github.com/shopspring/decimal"
)

type AccountBalance struct {
	AccountID int64
	Balance   decimal.Decimal
}

type NewAccount struct {
	AccountID      int64
	InitialBalance decimal.Decimal
}

type Account struct {
	AccountID        int64
	Balance          decimal.Decimal
	CreatedDate      time.Time
	LastModifiedDate time.Time
}

func fromDBAccount(dbAccount transferdbgen.Account) Account {
	return Account{
		AccountID:        dbAccount.AccountID,
		Balance:          dbAccount.Balance,
		CreatedDate:      dbAccount.CreatedDate,
		LastModifiedDate: dbAccount.LastModifiedDate,
	}
}

type Transaction struct {
	SourceAccountID      int64
	DestinationAccountID int64
	Amount               decimal.Decimal
}
