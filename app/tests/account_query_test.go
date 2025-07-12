package tests

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/danipurwadi/internal-transfer-system/app/api/apptest"
	"github.com/danipurwadi/internal-transfer-system/app/transferapp"
	"github.com/danipurwadi/internal-transfer-system/business/transferbus"
	"github.com/danipurwadi/internal-transfer-system/foundation/customerror"
	"github.com/google/go-cmp/cmp"
)

func accountQuery200(sd apptest.SeedData) []apptest.Table {
	table := []apptest.Table{
		{
			Name:       "basic",
			URL:        "/accounts/" + fmt.Sprint(sd.Accounts[0].AccountId),
			Method:     http.MethodGet,
			StatusCode: http.StatusOK,
			Input:      nil,
			GotResp:    &transferapp.BalanceResponse{},
			ExpResp: &transferapp.BalanceResponse{
				AccountId: strconv.FormatInt(sd.Accounts[0].AccountId, 10),
				Balance:   sd.Accounts[0].Balance.String(),
			},
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
	}

	return table
}

func accountQuery400() []apptest.Table {
	table := []apptest.Table{
		{
			Name:       "invalidaccountid",
			URL:        "/accounts/invalid",
			Method:     http.MethodGet,
			StatusCode: http.StatusBadRequest,
			Input:      nil,
			GotResp:    &customerror.Error{},
			ExpResp:    toErrorPtr(customerror.Newf(customerror.InvalidArgument, "invalid account id")),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
	}

	return table
}

func accountQuery404() []apptest.Table {
	table := []apptest.Table{
		{
			Name:       "notfound",
			URL:        "/accounts/12345",
			Method:     http.MethodGet,
			StatusCode: http.StatusNotFound,
			Input:      nil,
			GotResp:    &customerror.Error{},
			ExpResp:    toErrorPtr(customerror.Newf(customerror.NotFound, transferbus.ErrAccNotFound.Error())),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
	}

	return table
}