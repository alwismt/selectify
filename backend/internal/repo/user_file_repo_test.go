package repo_test

import (
	"testing"

	"alwis.dev/selectify/internal/model"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestUserFileRepo_AddUserFileWithTx(t *testing.T) {
	t.Cleanup(func() {
		_, _ = dbConn.Exec("DELETE FROM user_file WHERE user_id = $1", testUser.ID)
	})

	tx, err := transactionRepo.BeginTransaction(ctx)
	require.NoError(t, err)
	require.NotNil(t, tx)

	firstFile := &model.UserFile{
		ID:          uuid.New(),
		UserID:      testUser.ID,
		ContentType: "image/png",
	}

	err = userFileRepo.AddUserFileWithTx(ctx, tx.Transaction, firstFile)
	require.NoError(t, err)

	tx.CanCommit = true
	tx.End()

	var storedFile model.UserFile
	err = dbConn.QueryRowx("SELECT * FROM user_file WHERE user_id = $1", testUser.ID).StructScan(&storedFile)
	require.NoError(t, err)
	require.Equal(t, firstFile.ID, storedFile.ID)
	require.Equal(t, firstFile.ContentType, storedFile.ContentType)

	tx, err = transactionRepo.BeginTransaction(ctx)
	require.NoError(t, err)
	require.NotNil(t, tx)

	secondFile := &model.UserFile{
		ID:          uuid.New(),
		UserID:      testUser.ID,
		ContentType: "image/webp",
	}

	err = userFileRepo.AddUserFileWithTx(ctx, tx.Transaction, secondFile)
	require.NoError(t, err)

	tx.CanCommit = true
	tx.End()

	err = dbConn.QueryRowx("SELECT * FROM user_file WHERE user_id = $1", testUser.ID).StructScan(&storedFile)
	require.NoError(t, err)
	require.Equal(t, secondFile.ID, storedFile.ID)
	require.Equal(t, secondFile.ContentType, storedFile.ContentType)
}

func TestUserFileRepo_GetByUserID(t *testing.T) {
	t.Cleanup(func() {
		_, _ = dbConn.Exec("DELETE FROM user_file WHERE user_id = $1", testUser.ID)
	})

	tx, err := transactionRepo.BeginTransaction(ctx)
	require.NoError(t, err)
	require.NotNil(t, tx)

	file := &model.UserFile{
		ID:          uuid.New(),
		UserID:      testUser.ID,
		ContentType: "image/png",
	}

	err = userFileRepo.AddUserFileWithTx(ctx, tx.Transaction, file)
	require.NoError(t, err)
	tx.CanCommit = true
	tx.End()

	storedFile, err := userFileRepo.GetByUserID(ctx, testUser.ID)
	require.NoError(t, err)
	require.NotNil(t, storedFile)
	require.Equal(t, file.ID, storedFile.ID)
	require.Equal(t, file.ContentType, storedFile.ContentType)

	missingFile, err := userFileRepo.GetByUserID(ctx, 99999999)
	require.NoError(t, err)
	require.Nil(t, missingFile)
}

func TestUserFileRepo_DeleteByUserIDWithTx(t *testing.T) {
	t.Cleanup(func() {
		_, _ = dbConn.Exec("DELETE FROM user_file WHERE user_id = $1", testUser.ID)
	})

	tx, err := transactionRepo.BeginTransaction(ctx)
	require.NoError(t, err)
	require.NotNil(t, tx)

	file := &model.UserFile{
		ID:          uuid.New(),
		UserID:      testUser.ID,
		ContentType: "image/png",
	}

	err = userFileRepo.AddUserFileWithTx(ctx, tx.Transaction, file)
	require.NoError(t, err)
	tx.CanCommit = true
	tx.End()

	tx, err = transactionRepo.BeginTransaction(ctx)
	require.NoError(t, err)
	require.NotNil(t, tx)

	err = userFileRepo.DeleteByUserIDWithTx(ctx, tx.Transaction, testUser.ID)
	require.NoError(t, err)
	tx.CanCommit = true
	tx.End()

	storedFile, err := userFileRepo.GetByUserID(ctx, testUser.ID)
	require.NoError(t, err)
	require.Nil(t, storedFile)
}
