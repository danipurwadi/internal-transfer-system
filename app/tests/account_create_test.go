package tests

import (
	"net/http"

	"github.com/danipurwadi/internal-transfer-system/app/api/apptest"
	"github.com/danipurwadi/internal-transfer-system/app/transferapp"
	"github.com/danipurwadi/internal-transfer-system/business/transferbus"
	"github.com/danipurwadi/internal-transfer-system/foundation/customerror"
	"github.com/google/go-cmp/cmp"
)

func accountCreation200() []apptest.Table {
	table := []apptest.Table{
		{
			Name:       "basic",
			URL:        "/accounts",
			Method:     http.MethodPost,
			StatusCode: http.StatusCreated,
			Input: &transferapp.AccountCreationRequest{
				AccountID:      2,
				InitialBalance: "100.12345",
			},
			GotResp: nil,
			ExpResp: nil,
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
	}

	return table
}

func accountCreation400() []apptest.Table {
	table := []apptest.Table{
		{
			Name:       "missingaccountid",
			URL:        "/accounts",
			Method:     http.MethodPost,
			StatusCode: http.StatusBadRequest,
			Input: &transferapp.AccountCreationRequest{
				InitialBalance: "100.12345",
			},
			GotResp: &customerror.Error{},
			ExpResp: toErrorPtr(customerror.Newf(customerror.FailedPrecondition, "validate: [{\"field\":\"account_id\",\"error\":\"account_id is a required field\"}]")),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
		{
			Name:       "negativebalance",
			URL:        "/accounts",
			Method:     http.MethodPost,
			StatusCode: http.StatusBadRequest,
			Input: &transferapp.AccountCreationRequest{
				AccountID:      2,
				InitialBalance: "-100.12345",
			},
			GotResp: &customerror.Error{},
			ExpResp: toErrorPtr(customerror.Newf(customerror.InvalidArgument, transferbus.ErrNegativeBalance.Error())),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
		{
			Name:       "emptyinitialbalance",
			URL:        "/accounts",
			Method:     http.MethodPost,
			StatusCode: http.StatusBadRequest,
			Input: &transferapp.AccountCreationRequest{
				AccountID:      2,
				InitialBalance: "",
			},
			GotResp: &customerror.Error{},
			ExpResp: toErrorPtr(customerror.Newf(customerror.FailedPrecondition, "validate: [{\"field\":\"initial_balance\",\"error\":\"initial_balance is a required field\"}]")),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
	}

	return table
}

func accountCreation409(sd apptest.SeedData) []apptest.Table {
	table := []apptest.Table{
		{
			Name:       "duplicateaccountid",
			URL:        "/accounts",
			Method:     http.MethodPost,
			StatusCode: http.StatusConflict,
			Input: &transferapp.AccountCreationRequest{
				AccountID:      sd.Accounts[0].AccountID,
				InitialBalance: "100.12345",
			},
			GotResp: &customerror.Error{},
			ExpResp: toErrorPtr(customerror.Newf(customerror.AlreadyExists, transferbus.ErrAccAlreadyExist.Error())),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
	}

	return table
}
