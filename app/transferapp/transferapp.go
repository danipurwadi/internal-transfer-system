package transferapp

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/danipurwadi/internal-transfer-system/business/transferbus"
	"github.com/danipurwadi/internal-transfer-system/foundation/customerror"
	"github.com/danipurwadi/internal-transfer-system/foundation/web"
)

type App struct {
	transferbus *transferbus.Bus
}

func NewApp(bus *transferbus.Bus) *App {
	return &App{
		transferbus: bus,
	}
}

func (a *App) Routes(mux *web.Client) {
	mux.Handle(http.MethodGet, "/health", a.health)
	mux.Handle(http.MethodPost, "/accounts", a.createAccount)
	mux.Handle(http.MethodGet, "/accounts/{account_id}", a.getBalance)
	mux.Handle(http.MethodPost, "/transactions", a.createTransaction)
}

func (a *App) health(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	status := struct {
		Status bool `json:"status"`
	}{
		Status: true,
	}

	return web.Respond(ctx, w, status, http.StatusOK)
}

func (a *App) createAccount(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var req AccountCreationRequest
	if err := web.Decode(r, &req); err != nil {
		return customerror.New(customerror.FailedPrecondition, err)
	}

	account, err := toBusAccCreation(req)
	if err != nil {
		return customerror.New(customerror.FailedPrecondition, err)
	}

	err = a.transferbus.CreateAccount(ctx, account)
	if err != nil {
		return customerror.New(customerror.Internal, err)
	}

	return web.Respond(ctx, w, nil, http.StatusCreated)
}

func (a *App) getBalance(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	accountId := r.PathValue("account_id")
	if accountId == "" {
		return customerror.New(customerror.InvalidArgument, fmt.Errorf("account id is required"))
	}

	accId, err := strconv.ParseInt(accountId, 10, 0)
	if err != nil {
		return customerror.New(customerror.InvalidArgument, fmt.Errorf("invalid account id"))
	}

	balance, err := a.transferbus.GetBalance(ctx, accId)
	if err != nil {
		if errors.Is(err, transferbus.ErrAccNotFound) {
			return customerror.New(customerror.NotFound, err)
		}
		return customerror.Newf(customerror.Internal, "failed to get balance: accId[%d]: %s", accId, err)
	}

	return web.Respond(ctx, w, fromBusAccBalance(balance), http.StatusOK)
}

func (a *App) createTransaction(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var req TransactionRequest

	if err := web.Decode(r, &req); err != nil {
		return customerror.New(customerror.FailedPrecondition, err)
	}

	t, err := toBusTransaction(req)
	if err != nil {
		return customerror.New(customerror.FailedPrecondition, err)
	}

	err = a.transferbus.CreateTransaction(ctx, t)
	if err != nil {
		return customerror.New(customerror.Internal, err)
	}

	return web.Respond(ctx, w, nil, http.StatusCreated)
}
