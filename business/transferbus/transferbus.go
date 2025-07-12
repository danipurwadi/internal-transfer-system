package transferbus

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/danipurwadi/internal-transfer-system/business/transferbus/stores/transferdb"
	transferdbgen "github.com/danipurwadi/internal-transfer-system/business/transferbus/stores/transferdb/gen"
	"github.com/danipurwadi/internal-transfer-system/foundation/logger"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	DuplicateKeyViolatesUniqueConstraintCode = "23505"
	ViolatesForeignKeyConstraint             = "23503"
)

var (
	ErrAccNotFound       = errors.New("account not found")
	ErrAccAlreadyExist   = errors.New("account already exist")
	ErrNegativeBalance   = errors.New("negative balance")
	ErrInsufficientFunds = errors.New("insufficient funds")
	ErrSameAccount       = errors.New("source and destination account cannot be the same")
)

type Bus struct {
	log   *logger.Logger
	store transferdb.TxQuerier
}

func New(store transferdb.TxQuerier, log *logger.Logger) *Bus {
	return &Bus{
		log:   log,
		store: store,
	}
}

func (b *Bus) CreateAccount(ctx context.Context, account NewAccount) (Account, error) {
	if account.InitialBalance.IsNegative() {
		return Account{}, ErrNegativeBalance
	}
	tx, err := b.store.GetTx(ctx)
	if err != nil {
		return Account{}, fmt.Errorf("get transaction: %w", err)
	}
	dbtx := b.store.WithTx(tx)
	defer func() {
		if err := tx.Rollback(ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			b.log.Error(ctx, "rollback failed", "err", err)
		}
	}()

	acc, err := dbtx.CreateAccount(ctx, transferdbgen.CreateAccountParams{
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

	err = dbtx.CreateTransaction(ctx, transferdbgen.CreateTransactionParams{
		AccountID:   account.AccountId,
		Amount:      account.InitialBalance,
		CreatedDate: time.Now(),
	})

	if err != nil {
		return Account{}, fmt.Errorf("create transaction: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return Account{}, fmt.Errorf("commit transaction: %w", err)
	}

	return fromDbAccount(acc), nil
}

func (b *Bus) CreateTransaction(ctx context.Context, transaction Transaction) error {
	tx, err := b.store.GetTx(ctx)
	if err != nil {
		return fmt.Errorf("get transaction: %w", err)
	}
	dbtx := b.store.WithTx(tx)
	defer func() {
		if err := tx.Rollback(ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			b.log.Error(ctx, "rollback failed", "err", err)
		}
	}()
	if transaction.Amount.IsNegative() {
		return ErrNegativeBalance
	}
	if transaction.SourceAccountId == transaction.DestinationAccountId {
		return ErrSameAccount
	}

	// check that both accounts exist
	accounts, err := dbtx.GetAccounts(ctx, []int64{transaction.SourceAccountId, transaction.DestinationAccountId})
	if err != nil {
		return fmt.Errorf("get accounts: %w", err)
	}
	if len(accounts) != 2 {
		return ErrAccNotFound
	}

	debitResult, err := dbtx.DebitAccount(ctx, transferdbgen.DebitAccountParams{
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
	_, err = dbtx.CreditAccount(ctx, transferdbgen.CreditAccountParams{
		Amount:    transaction.Amount,
		AccountID: transaction.DestinationAccountId,
	})
	if err != nil {
		return fmt.Errorf("credit account: %w", err)
	}

	// Record Credit Transaction
	err = dbtx.CreateTransaction(ctx, transferdbgen.CreateTransactionParams{
		AccountID:   transaction.SourceAccountId,
		Amount:      transaction.Amount.Neg(),
		CreatedDate: time.Now(),
	})

	if err != nil {
		return fmt.Errorf("create transaction: %w", err)
	}

	// Record Debit Transaction
	err = dbtx.CreateTransaction(ctx, transferdbgen.CreateTransactionParams{
		AccountID:   transaction.DestinationAccountId,
		Amount:      transaction.Amount,
		CreatedDate: time.Now(),
	})

	if err != nil {
		return fmt.Errorf("create transaction: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}
	return nil
}

func (b *Bus) GetBalance(ctx context.Context, accountId int64) (Account, error) {
	// check that account exist in the first place
	account, err := b.store.GetAccount(ctx, accountId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Account{}, ErrAccNotFound
		}
		return Account{}, fmt.Errorf("get account: %d: %w", accountId, err)
	}

	return fromDbAccount(account), nil
}
