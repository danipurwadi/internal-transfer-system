package tests

import (
	"net/http"

	"github.com/danipurwadi/internal-transfer-system/app/api/apptest"
	"github.com/danipurwadi/internal-transfer-system/app/transferapp"
	"github.com/danipurwadi/internal-transfer-system/business/transferbus"
	"github.com/danipurwadi/internal-transfer-system/foundation/customerror"
	"github.com/google/go-cmp/cmp"
)

func transactionSubmission201(sd apptest.SeedData) []apptest.Table {
	table := []apptest.Table{
		{
			Name:       "basic",
			URL:        "/transactions",
			Method:     http.MethodPost,
			StatusCode: http.StatusCreated,
			Input: &transferapp.TransactionRequest{
				SourceAccountID:      sd.Accounts[0].AccountID,
				DestinationAccountID: sd.Accounts[1].AccountID,
				Amount:               "10.0",
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

func transactionSubmission400(sd apptest.SeedData) []apptest.Table {
	table := []apptest.Table{
		{
			Name:       "negativeamount",
			URL:        "/transactions",
			Method:     http.MethodPost,
			StatusCode: http.StatusBadRequest,
			Input: &transferapp.TransactionRequest{
				SourceAccountID:      1,
				DestinationAccountID: 2,
				Amount:               "-10.0",
			},
			GotResp: &customerror.Error{},
			ExpResp: toErrorPtr(customerror.Newf(customerror.InvalidArgument, transferbus.ErrNegativeBalance.Error())),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
		{
			Name:       "sameaccount",
			URL:        "/transactions",
			Method:     http.MethodPost,
			StatusCode: http.StatusBadRequest,
			Input: &transferapp.TransactionRequest{
				SourceAccountID:      1,
				DestinationAccountID: 1,
				Amount:               "10.0",
			},
			GotResp: &customerror.Error{},
			ExpResp: toErrorPtr(customerror.Newf(customerror.InvalidArgument, transferbus.ErrSameAccount.Error())),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
		{
			Name:       "insufficientfunds",
			URL:        "/transactions",
			Method:     http.MethodPost,
			StatusCode: http.StatusBadRequest,
			Input: &transferapp.TransactionRequest{
				SourceAccountID:      sd.Accounts[0].AccountID,
				DestinationAccountID: sd.Accounts[1].AccountID,
				Amount:               "1000000.0",
			},
			GotResp: &customerror.Error{},
			ExpResp: toErrorPtr(customerror.Newf(customerror.FailedPrecondition, transferbus.ErrInsufficientFunds.Error())),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
	}

	return table
}

func transactionSubmission404(sd apptest.SeedData) []apptest.Table {
	table := []apptest.Table{
		{
			Name:       "sourceaccnotfound",
			URL:        "/transactions",
			Method:     http.MethodPost,
			StatusCode: http.StatusNotFound,
			Input: &transferapp.TransactionRequest{
				SourceAccountID:      1234,
				DestinationAccountID: sd.Accounts[0].AccountID,
				Amount:               "10.0",
			},
			GotResp: &customerror.Error{},
			ExpResp: toErrorPtr(customerror.Newf(customerror.NotFound, transferbus.ErrAccNotFound.Error())),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
		{
			Name:       "destaccnotfound",
			URL:        "/transactions",
			Method:     http.MethodPost,
			StatusCode: http.StatusNotFound,
			Input: &transferapp.TransactionRequest{
				SourceAccountID:      sd.Accounts[0].AccountID,
				DestinationAccountID: 1234,
				Amount:               "10.0",
			},
			GotResp: &customerror.Error{},
			ExpResp: toErrorPtr(customerror.Newf(customerror.NotFound, transferbus.ErrAccNotFound.Error())),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
	}

	return table
}
