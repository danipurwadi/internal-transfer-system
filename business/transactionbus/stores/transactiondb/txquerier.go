package transactiondb

import (
	"context"

	transactiondb "github.com/danipurwadi/internal-transfer-system/business/transactionbus/stores/transactiondb/gen"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TxQuerier interface {
	transactiondb.Querier
	WithTx(tx pgx.Tx) TxQuerier
	GetTx(ctx context.Context) (pgx.Tx, error)
}

var _ TxQuerier = (*TxQueries)(nil)

func NewTxQueries(pool *pgxpool.Pool) *TxQueries {
	return &TxQueries{
		Queries: transactiondb.New(pool),
		TxnPool: pool,
	}
}

type TxQueries struct {
	*transactiondb.Queries
	TxnPool *pgxpool.Pool
}

func (q *TxQueries) GetTx(ctx context.Context) (pgx.Tx, error) {
	return q.TxnPool.Begin(ctx)
}

func (q *TxQueries) WithTx(tx pgx.Tx) TxQuerier {
	return &TxQueries{
		Queries: q.Queries.WithTx(tx),
		TxnPool: q.TxnPool,
	}
}
