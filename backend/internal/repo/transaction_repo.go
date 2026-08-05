package repo

import (
	"alwis.dev/selectify/internal/db"
	"alwis.dev/selectify/internal/logger"
	"context"
	"github.com/jmoiron/sqlx"
)

type transactionRepoy struct {
	rwDb *sqlx.DB
}

type TransactionRepo interface {
	BeginTransaction(ctx context.Context) (*db.Transaction, error)
}

func NewTransactionRepository(db *db.DatabaseConnection) TransactionRepo {
	return &transactionRepoy{
		rwDb: db.RwDb,
	}
}

func (r *transactionRepoy) BeginTransaction(ctx context.Context) (*db.Transaction, error) {
	tx, err := db.BeginTransaction(ctx, r.rwDb)
	if err != nil {
		return nil, logger.Error(ctx, err, "Failed to start transaction")
	}

	return tx, nil
}
