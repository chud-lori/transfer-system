package repositories_test

import (
	"context"
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

func TestTransactionRepositoryPostgre_Save_Valid(t *testing.T) {
	db := testutils.SetupTestDB(t)
	tx := testutils.SetupTestTx(t, db)
	defer tx.Rollback()

	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.New())

	accountRepo := &repositories.AccountRepositoryPostgre{DB: db}
	sourceAcc := &entities.Account{AccountID: 1001, Balance: decimal.NewFromFloat(1000)}
	destAcc := &entities.Account{AccountID: 1002, Balance: decimal.NewFromFloat(500)}

	_, err := accountRepo.Save(ctx, tx, sourceAcc)
	require.NoError(t, err)
	_, err = accountRepo.Save(ctx, tx, destAcc)
	require.NoError(t, err)

	repo := &repositories.TransactionRepositoryPostgre{DB: db}

	transaction := &entities.Transaction{
		SourceAccountID:      sourceAcc.AccountID,
		DestinationAccountID: destAcc.AccountID,
		Amount:               decimal.NewFromFloat(100),
	}

	result, err := repo.Save(ctx, tx, transaction)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotZero(t, result.Id)
	assert.Equal(t, transaction.Amount, result.Amount)
}

func TestTransactionRepositoryPostgre_Save_Invalid(t *testing.T) {
	db := testutils.SetupTestDB(t)
	tx := testutils.SetupTestTx(t, db)
	defer tx.Rollback()

	repo := &repositories.TransactionRepositoryPostgre{DB: db}
	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.New())

	transaction := &entities.Transaction{
		SourceAccountID:      1,
		DestinationAccountID: 2,
		Amount:               decimal.Decimal{},
	}

	_, err := repo.Save(ctx, tx, transaction)
	assert.Error(t, err)
}

func TestTransactionRepositoryPostgre_UpdateBalance_Valid(t *testing.T) {
	db := testutils.SetupTestDB(t)
	tx := testutils.SetupTestTx(t, db)
	defer tx.Rollback()

	repo := &repositories.TransactionRepositoryPostgre{DB: db}
	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.New())

	accountID := int64(1)
	initialBalance := decimal.NewFromFloat(100.0)

	// seed the account
	_, err := tx.ExecContext(ctx,
		`INSERT INTO accounts (id, balance) VALUES ($1, $2)`,
		accountID, initialBalance)
	require.NoError(t, err)

	amount := decimal.NewFromFloat(20.0)
	err = repo.UpdateBalance(ctx, tx, accountID, amount)

	assert.NoError(t, err)

	var newBalance decimal.Decimal
	row := tx.QueryRowContext(ctx, `SELECT balance FROM accounts WHERE id = $1`, accountID)
	err = row.Scan(&newBalance)
	require.NoError(t, err)

	expectedBalance := initialBalance.Add(amount)
	assert.True(t, expectedBalance.Equal(newBalance))
}

func TestTransactionRepositoryPostgre_UpdateBalance_Invalid(t *testing.T) {
	db := testutils.SetupTestDB(t)
	tx := testutils.SetupTestTx(t, db)
	defer tx.Rollback()

	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.New())

	repo := &repositories.TransactionRepositoryPostgre{DB: db}

	invalidAccountID := int64(999999) // account ID not exist
	err := repo.UpdateBalance(ctx, tx, invalidAccountID, decimal.NewFromFloat(100))

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no account found", "Expected error related to missing account ID")
}
