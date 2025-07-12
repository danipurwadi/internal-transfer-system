package transferbus

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/danipurwadi/internal-transfer-system/business/transferbus/stores/transferdb"
	transferdbgen "github.com/danipurwadi/internal-transfer-system/business/transferbus/stores/transferdb/gen"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	DuplicateKeyViolatesUniqueConstraintCode = "23505"
)

var (
	ErrAccNotFound     = errors.New("account not found")
	ErrAccAlreadyExist = errors.New("account already exist")
	ErrNegativeBalance = errors.New("negative balance")
)

type Bus struct {
	store transferdb.TxQuerier
}

func New(store transferdb.TxQuerier) *Bus {
	return &Bus{
		store: store,
	}
}

func (b *Bus) CreateAccount(ctx context.Context, account AccountCreation) error {
	if account.InitialBalance.IsNegative() {
		return ErrNegativeBalance
	}

	err := b.store.CreateAccount(ctx, transferdbgen.CreateAccountParams{
		AccountID:        account.AccountId,
		CreatedDate:      time.Now(),
		LastModifiedDate: time.Now(),
	})

	if err != nil {
		var pgError *pgconn.PgError
		if errors.As(err, &pgError) && pgError.Code == DuplicateKeyViolatesUniqueConstraintCode {
			return ErrAccAlreadyExist
		}
		return fmt.Errorf("create: %w", err)
	}

	err = b.store.CreateTransaction(ctx, transferdbgen.CreateTransactionParams{
		AccountID:   account.AccountId,
		Amount:      account.InitialBalance,
		CreatedDate: time.Now(),
	})

	if err != nil {
		return fmt.Errorf("create transaction: %w", err)
	}

	return nil
}

func (b *Bus) CreateTransaction(ctx context.Context, transaction Transaction) error {
	err := b.store.CreateTransaction(ctx, transferdbgen.CreateTransactionParams{
		AccountID:   transaction.SourceAccountId,
		Amount:      transaction.Amount.Neg(),
		CreatedDate: time.Now(),
	})

	if err != nil {
		return fmt.Errorf("create transaction: %w", err)
	}

	err = b.store.CreateTransaction(ctx, transferdbgen.CreateTransactionParams{
		AccountID:   transaction.DestinationAccountId,
		Amount:      transaction.Amount,
		CreatedDate: time.Now(),
	})

	if err != nil {
		return fmt.Errorf("create transaction: %w", err)
	}

	return nil
}

func (b *Bus) GetBalance(ctx context.Context, accountId int64) (AccountBalance, error) {
	// check that account exist in the first place
	_, err := b.store.GetAccount(ctx, accountId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return AccountBalance{}, ErrAccNotFound
		}
		return AccountBalance{}, fmt.Errorf("get account: %d: %w", accountId, err)
	}

	balance, err := b.store.GetBalance(ctx, accountId)
	if err != nil {
		return AccountBalance{}, fmt.Errorf("get balance: %d: %w", accountId, err)
	}

	resp := AccountBalance{
		AccountId: accountId,
		Balance:   balance,
	}
	return resp, nil
}
