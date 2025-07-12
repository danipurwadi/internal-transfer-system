package transferbus

import (
	"time"

	transferdbgen "github.com/danipurwadi/internal-transfer-system/business/transferbus/stores/transferdb/gen"
	"github.com/shopspring/decimal"
)

type AccountBalance struct {
	AccountId int64
	Balance   decimal.Decimal
}

type NewAccount struct {
	AccountId      int64
	InitialBalance decimal.Decimal
}

type Account struct {
	AccountId        int64
	Balance          decimal.Decimal
	CreatedDate      time.Time
	LastModifiedDate time.Time
}

func fromDbAccount(dbAccount transferdbgen.Account) Account {
	return Account{
		AccountId:        dbAccount.AccountID,
		Balance:          dbAccount.Balance,
		CreatedDate:      dbAccount.CreatedDate,
		LastModifiedDate: dbAccount.LastModifiedDate,
	}
}

type Transaction struct {
	SourceAccountId      int64
	DestinationAccountId int64
	Amount               decimal.Decimal
}
