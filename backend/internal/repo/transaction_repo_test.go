package repo_test

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransactionRepo_BeginTransaction(t *testing.T) {
	tx, err := transactionRepo.BeginTransaction(ctx)

	require.NoError(t, err)
	require.NotNil(t, tx)
	require.NotNil(t, tx.Transaction)
	require.False(t, tx.CanCommit)

	tx.End()
}

func TestTransactionRepo_BeginTransaction_CanCommit(t *testing.T) {
	tx, err := transactionRepo.BeginTransaction(ctx)
	require.NoError(t, err)
	require.NotNil(t, tx)

	tx.CanCommit = true
	require.True(t, tx.CanCommit)
	tx.End()
}
