package repositories_test

import (
	"context"
	"database/sql"
	"testing"
	"transfer-system/adapters/repositories"
	"transfer-system/domain/entities"
	"transfer-system/internal/testutils"
	"transfer-system/pkg/logger"

	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccountRepositoryPostgre_Save_Success(t *testing.T) {
	db := testutils.SetupTestDB(t)
	tx := testutils.SetupTestTx(t, db)
	defer tx.Rollback()

	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.New())

	repo := &repositories.AccountRepositoryPostgre{DB: db}

	account := &entities.Account{
		AccountID: 12345678,
		Balance:   decimal.NewFromFloat(250.75),
	}

	// Execute
	savedAccount, err := repo.Save(ctx, tx, account)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, savedAccount)
	assert.Equal(t, account.AccountID, savedAccount.AccountID)
	assert.True(t, account.Balance.Equal(savedAccount.Balance))
}

func TestAccountRepositoryPostgre_FindById_Success(t *testing.T) {
	// Setup
	db := testutils.SetupTestDB(t)
	tx := testutils.SetupTestTx(t, db)
	defer tx.Rollback()

	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.New())

	repo := &repositories.AccountRepositoryPostgre{DB: db}

	account := &entities.Account{
		AccountID: 987654321,
		Balance:   decimal.NewFromFloat(999.99),
	}

	// Insert first so we can retrieve
	_, err := repo.Save(ctx, tx, account)
	require.NoError(t, err)

	// Execute
	fetchedAccount, err := repo.FindById(ctx, tx, account.AccountID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, fetchedAccount)
	assert.Equal(t, account.AccountID, fetchedAccount.AccountID)
	assert.True(t, account.Balance.Equal(fetchedAccount.Balance))
}

func TestAccountRepositoryPostgre_Save_Error(t *testing.T) {
	db := testutils.SetupTestDB(t)
	tx := testutils.SetupTestTx(t, db)
	defer tx.Rollback()

	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.New())

	repo := &repositories.AccountRepositoryPostgre{DB: db}

	account := &entities.Account{
		AccountID: 10001,
		Balance:   decimal.NewFromFloat(50.0),
	}

	// Insert once
	_, err := repo.Save(ctx, tx, account)
	require.NoError(t, err)

	// Attempt to insert again with same ID
	_, err = repo.Save(ctx, tx, account)
	assert.Error(t, err, "Expected error due to duplicate primary key")
}

func TestAccountRepositoryPostgre_FindById_NotFound(t *testing.T) {
	db := testutils.SetupTestDB(t)
	tx := testutils.SetupTestTx(t, db)
	defer tx.Rollback()

	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.New())

	repo := &repositories.AccountRepositoryPostgre{DB: db}

	// Try to find an ID that doesn't exist
	nonExistentID := int64(999999)
	acc, err := repo.FindById(ctx, tx, nonExistentID)

	assert.Error(t, err)
	assert.Nil(t, acc)
	assert.Equal(t, sql.ErrNoRows, err)
}
