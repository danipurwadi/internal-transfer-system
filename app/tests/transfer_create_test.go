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
				SourceAccountId:      sd.Accounts[0].AccountId,
				DestinationAccountId: sd.Accounts[1].AccountId,
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
				SourceAccountId:      1,
				DestinationAccountId: 2,
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
				SourceAccountId:      1,
				DestinationAccountId: 1,
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
				SourceAccountId:      sd.Accounts[0].AccountId,
				DestinationAccountId: sd.Accounts[1].AccountId,
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
				SourceAccountId:      1234,
				DestinationAccountId: sd.Accounts[0].AccountId,
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
				SourceAccountId:      sd.Accounts[0].AccountId,
				DestinationAccountId: 1234,
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
