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
	ErrAccNotFound       = errors.New("account not found")
	ErrAccAlreadyExist   = errors.New("account already exist")
	ErrNegativeBalance   = errors.New("negative balance")
	ErrInsufficientFunds = errors.New("insufficient funds")
)

type Bus struct {
	store transferdb.TxQuerier
}

func New(store transferdb.TxQuerier) *Bus {
	return &Bus{
		store: store,
	}
}

func (b *Bus) CreateAccount(ctx context.Context, account NewAccount) (Account, error) {
	if account.InitialBalance.IsNegative() {
		return Account{}, ErrNegativeBalance
	}

	acc, err := b.store.CreateAccount(ctx, transferdbgen.CreateAccountParams{
		AccountID:        account.AccountId,
		Balance:          account.InitialBalance,
		CreatedDate:      time.Now(),
		LastModifiedDate: time.Now(),
	})

	if err != nil {
		var pgError *pgconn.PgError
		if errors.As(err, &pgError) && pgError.Code == DuplicateKeyViolatesUniqueConstraintCode {
			return Account{}, ErrAccAlreadyExist
		}
		return Account{}, fmt.Errorf("create: %w", err)
	}

	err = b.store.CreateTransaction(ctx, transferdbgen.CreateTransactionParams{
		AccountID:   account.AccountId,
		Amount:      account.InitialBalance,
		CreatedDate: time.Now(),
	})

	if err != nil {
		return Account{}, fmt.Errorf("create transaction: %w", err)
	}

	return fromDbAccount(acc), nil
}

func (b *Bus) CreateTransaction(ctx context.Context, transaction Transaction) error {
	debitResult, err := b.store.DebitAccount(ctx, transferdbgen.DebitAccountParams{
		Amount:    transaction.Amount,
		AccountID: transaction.SourceAccountId,
	})
	if err != nil {
		return fmt.Errorf("debit account: %w", err)
	}

	// if no rows is updated, balance was too low
	if debitResult.RowsAffected() == 0 {
		return ErrInsufficientFunds
	}

	// if Debit was successful, credit the destination account
	_, err = b.store.CreditAccount(ctx, transferdbgen.CreditAccountParams{
		Amount:    transaction.Amount,
		AccountID: transaction.DestinationAccountId,
	})
	if err != nil {
		return fmt.Errorf("credit account: %w", err)
	}

	// Record Credit Transaction
	err = b.store.CreateTransaction(ctx, transferdbgen.CreateTransactionParams{
		AccountID:   transaction.SourceAccountId,
		Amount:      transaction.Amount.Neg(),
		CreatedDate: time.Now(),
	})

	if err != nil {
		return fmt.Errorf("create transaction: %w", err)
	}

	// Record Debit Transaction
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
	account, err := b.store.GetAccount(ctx, accountId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return AccountBalance{}, ErrAccNotFound
		}
		return AccountBalance{}, fmt.Errorf("get account: %d: %w", accountId, err)
	}

	resp := AccountBalance{
		AccountId: accountId,
		Balance:   account.Balance,
	}
	return resp, nil
}
