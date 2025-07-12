package transferbus

import "github.com/shopspring/decimal"

type AccountBalance struct {
	AccountId int64
	Balance   decimal.Decimal
}

type AccountCreation struct {
	AccountId      int64
	InitialBalance decimal.Decimal
}

type Transaction struct {
	SourceAccountId      int64
	DestinationAccountId int64
	Amount               decimal.Decimal
}
