package db

import (
	"alwis.dev/selectify/internal/logger"
	"context"
	"github.com/jmoiron/sqlx"
)

type Transaction struct {
	Transaction *sqlx.Tx
	CanCommit   bool
}

func BeginTransaction(ctx context.Context, db *sqlx.DB) (*Transaction, error) {
	tx, err := db.Beginx()
	if err != nil {
		return nil, logger.Errorf(ctx, err, "No connection to database")
	}

	return &Transaction{
		Transaction: tx,
		CanCommit:   false,
	}, nil
}

func (tx *Transaction) End() {
	if tx.Transaction == nil {
		return
	}

	if tx.CanCommit {
		_ = tx.Transaction.Commit()
	} else {
		_ = tx.Transaction.Rollback()
	}
	tx.Transaction = nil
}
