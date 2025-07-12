package transferapp

import (
	"strconv"

	"github.com/danipurwadi/internal-transfer-system/business/transferbus"
	"github.com/danipurwadi/internal-transfer-system/foundation/customerror"
	"github.com/danipurwadi/internal-transfer-system/foundation/validate"
	"github.com/shopspring/decimal"
)

type BalanceResponse struct {
	AccountId string `json:"account_id"`
	Balance   string `json:"balance"`
}

func fromBusAccBalance(accBalance transferbus.AccountBalance) BalanceResponse {
	return BalanceResponse{
		AccountId: strconv.FormatInt(accBalance.AccountId, 10),
		Balance:   accBalance.Balance.String(),
	}
}

type AccountCreationRequest struct {
	AccountId      int64  `json:"account_id" validate:"required,min=1"`
	InitialBalance string `json:"initial_balance" validate:"required"`
}

// Validate checks if the data in the model is considered clean.
func (r AccountCreationRequest) Validate() error {
	if err := validate.Check(r); err != nil {
		return customerror.Newf(customerror.FailedPrecondition, "validate: %s", err)
	}
	return nil
}

func toBusAccCreation(req AccountCreationRequest) (transferbus.AccountCreation, error) {
	decimalBalance, err := decimal.NewFromString(req.InitialBalance)
	if err != nil {
		return transferbus.AccountCreation{}, err
	}

	return transferbus.AccountCreation{
		AccountId:      req.AccountId,
		InitialBalance: decimalBalance,
	}, nil
}

type TransactionRequest struct {
	SourceAccountId      int64  `json:"source_account_id" validate:"required,min=1"`
	DestinationAccountId int64  `json:"destination_account_id" validate:"required,min=1"`
	Amount               string `json:"amount" validate:"required"`
}

func toBusTransaction(req TransactionRequest) (transferbus.Transaction, error) {
	decimalAmount, err := decimal.NewFromString(req.Amount)
	if err != nil {
		return transferbus.Transaction{}, err
	}
	return transferbus.Transaction{
		SourceAccountId:      req.SourceAccountId,
		DestinationAccountId: req.DestinationAccountId,
		Amount:               decimalAmount,
	}, nil
}
