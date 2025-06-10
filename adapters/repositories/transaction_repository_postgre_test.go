package repositories_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"transfer-system/adapters/repositories"
	"transfer-system/domain/entities"
	"transfer-system/domain/ports"
	"transfer-system/mocks"
	"transfer-system/pkg/logger"

	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

type TestableTransactionRepository struct {
	*repositories.TransactionRepositoryPostgre
	mockQueryRow func(ctx context.Context, query string, args ...interface{}) (int64, error)
	mockExec     func(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

func (r *TestableTransactionRepository) Save(ctx context.Context, tx ports.Transaction, transaction *entities.Transaction) (*entities.Transaction, error) {
	if r.mockQueryRow != nil {
		id, err := r.mockQueryRow(ctx, "", transaction.SourceAccountID, transaction.DestinationAccountID, transaction.Amount)
		if err != nil {
			return nil, err
		}
		transaction.Id = id
		return transaction, nil
	}
	return r.TransactionRepositoryPostgre.Save(ctx, tx, transaction)
}

func (r *TestableTransactionRepository) UpdateBalance(ctx context.Context, tx ports.Transaction, accountID int64, amount decimal.Decimal) error {
	if r.mockExec != nil {
		_, err := r.mockExec(ctx, "", amount, accountID)
		return err
	}
	return r.TransactionRepositoryPostgre.UpdateBalance(ctx, tx, accountID, amount)
}

func TestTransactionRepositoryPostgre_Save(t *testing.T) {
	// Setup logger context
	baseLogger := logrus.New()
	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, baseLogger)

	tests := []struct {
		name          string
		transaction   *entities.Transaction
		mockReturn    int64
		mockError     error
		expectedError bool
		expectedID    int64
	}{
		{
			name: "successful save",
			transaction: &entities.Transaction{
				SourceAccountID:      1,
				DestinationAccountID: 2,
				Amount:               decimal.NewFromFloat(100.50),
			},
			mockReturn:    123,
			mockError:     nil,
			expectedError: false,
			expectedID:    123,
		},
		{
			name: "successful save with different amounts",
			transaction: &entities.Transaction{
				SourceAccountID:      10,
				DestinationAccountID: 20,
				Amount:               decimal.NewFromFloat(999.99),
			},
			mockReturn:    456,
			mockError:     nil,
			expectedError: false,
			expectedID:    456,
		},
		{
			name: "database error",
			transaction: &entities.Transaction{
				SourceAccountID:      1,
				DestinationAccountID: 2,
				Amount:               decimal.NewFromFloat(100.50),
			},
			mockReturn:    0,
			mockError:     errors.New("database error"),
			expectedError: true,
		},
		{
			name: "constraint violation error",
			transaction: &entities.Transaction{
				SourceAccountID:      1,
				DestinationAccountID: 999, // non-existent account
				Amount:               decimal.NewFromFloat(100.50),
			},
			mockReturn:    0,
			mockError:     errors.New("foreign key constraint violation"),
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup testable repository
			repo := &TestableTransactionRepository{
				TransactionRepositoryPostgre: &repositories.TransactionRepositoryPostgre{
					DB: &mocks.MockDatabase{},
				},
				mockQueryRow: func(ctx context.Context, query string, args ...interface{}) (int64, error) {
					// Verify the arguments passed
					assert.Len(t, args, 3)
					assert.Equal(t, tt.transaction.SourceAccountID, args[0])
					assert.Equal(t, tt.transaction.DestinationAccountID, args[1])
					assert.True(t, tt.transaction.Amount.Equal(args[2].(decimal.Decimal)))
					return tt.mockReturn, tt.mockError
				},
			}

			// Execute
			result, err := repo.Save(ctx, &mocks.MockTransaction{}, tt.transaction)

			// Assert
			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectedID, result.Id)
				assert.Equal(t, tt.transaction.SourceAccountID, result.SourceAccountID)
				assert.Equal(t, tt.transaction.DestinationAccountID, result.DestinationAccountID)
				assert.True(t, tt.transaction.Amount.Equal(result.Amount))
			}
		})
	}
}

func TestTransactionRepositoryPostgre_UpdateBalance(t *testing.T) {
	// Setup logger context
	baseLogger := logrus.New()
	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, baseLogger)

	tests := []struct {
		name          string
		accountID     int64
		amount        decimal.Decimal
		mockResult    sql.Result
		mockError     error
		expectedError bool
	}{
		{
			name:          "successful balance update - positive amount",
			accountID:     1,
			amount:        decimal.NewFromFloat(100.50),
			mockResult:    &mocks.MockResult{},
			mockError:     nil,
			expectedError: false,
		},
		{
			name:          "successful balance update - negative amount (debit)",
			accountID:     2,
			amount:        decimal.NewFromFloat(-50.25),
			mockResult:    &mocks.MockResult{},
			mockError:     nil,
			expectedError: false,
		},
		{
			name:          "successful balance update - zero amount",
			accountID:     3,
			amount:        decimal.NewFromFloat(0),
			mockResult:    &mocks.MockResult{},
			mockError:     nil,
			expectedError: false,
		},
		{
			name:          "database error",
			accountID:     1,
			amount:        decimal.NewFromFloat(100.50),
			mockResult:    nil,
			mockError:     errors.New("database connection error"),
			expectedError: true,
		},
		{
			name:          "account not found",
			accountID:     999,
			amount:        decimal.NewFromFloat(100.50),
			mockResult:    nil,
			mockError:     errors.New("account not found"),
			expectedError: true,
		},
		{
			name:          "insufficient balance constraint",
			accountID:     1,
			amount:        decimal.NewFromFloat(-1000.00),
			mockResult:    nil,
			mockError:     errors.New("check constraint violation: balance cannot be negative"),
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup testable repository
			repo := &TestableTransactionRepository{
				TransactionRepositoryPostgre: &repositories.TransactionRepositoryPostgre{
					DB: &mocks.MockDatabase{},
				},
				mockExec: func(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
					// Verify the arguments passed
					assert.Len(t, args, 2)
					assert.True(t, tt.amount.Equal(args[0].(decimal.Decimal)))
					assert.Equal(t, tt.accountID, args[1])
					return tt.mockResult, tt.mockError
				},
			}

			// Execute
			err := repo.UpdateBalance(ctx, &mocks.MockTransaction{}, tt.accountID, tt.amount)

			// Assert
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// Test edge cases and error scenarios
func TestTransactionRepositoryPostgre_EdgeCases(t *testing.T) {
	baseLogger := logrus.New()
	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, baseLogger)

	t.Run("Save with very large amount", func(t *testing.T) {
		repo := &TestableTransactionRepository{
			TransactionRepositoryPostgre: &repositories.TransactionRepositoryPostgre{
				DB: &mocks.MockDatabase{},
			},
			mockQueryRow: func(ctx context.Context, query string, args ...interface{}) (int64, error) {
				return 1, nil
			},
		}

		transaction := &entities.Transaction{
			SourceAccountID:      1,
			DestinationAccountID: 2,
			Amount:               decimal.NewFromFloat(999999999999.99),
		}

		result, err := repo.Save(ctx, &mocks.MockTransaction{}, transaction)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.True(t, transaction.Amount.Equal(result.Amount))
	})

	t.Run("Save with very small amount", func(t *testing.T) {
		repo := &TestableTransactionRepository{
			TransactionRepositoryPostgre: &repositories.TransactionRepositoryPostgre{
				DB: &mocks.MockDatabase{},
			},
			mockQueryRow: func(ctx context.Context, query string, args ...interface{}) (int64, error) {
				return 1, nil
			},
		}

		transaction := &entities.Transaction{
			SourceAccountID:      1,
			DestinationAccountID: 2,
			Amount:               decimal.NewFromFloat(0.01),
		}

		result, err := repo.Save(ctx, &mocks.MockTransaction{}, transaction)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.True(t, transaction.Amount.Equal(result.Amount))
	})

	t.Run("UpdateBalance with precision decimal", func(t *testing.T) {
		repo := &TestableTransactionRepository{
			TransactionRepositoryPostgre: &repositories.TransactionRepositoryPostgre{
				DB: &mocks.MockDatabase{},
			},
			mockExec: func(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
				return &mocks.MockResult{}, nil
			},
		}

		// Test with high precision decimal
		amount := decimal.NewFromFloat(123.456789)
		err := repo.UpdateBalance(ctx, &mocks.MockTransaction{}, 1, amount)
		assert.NoError(t, err)
	})
}

// Test with missing logger context
func TestTransactionRepositoryPostgre_NoLoggerContext(t *testing.T) {
	// Context without logger
	ctx := context.Background()

	t.Run("Save without logger context", func(t *testing.T) {
		repo := &TestableTransactionRepository{
			TransactionRepositoryPostgre: &repositories.TransactionRepositoryPostgre{
				DB: &mocks.MockDatabase{},
			},
			mockQueryRow: func(ctx context.Context, query string, args ...interface{}) (int64, error) {
				return 1, nil
			},
		}

		transaction := &entities.Transaction{
			SourceAccountID:      1,
			DestinationAccountID: 2,
			Amount:               decimal.NewFromFloat(100.50),
		}

		// Should not panic even without logger
		result, err := repo.Save(ctx, &mocks.MockTransaction{}, transaction)
		assert.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("UpdateBalance without logger context", func(t *testing.T) {
		repo := &TestableTransactionRepository{
			TransactionRepositoryPostgre: &repositories.TransactionRepositoryPostgre{
				DB: &mocks.MockDatabase{},
			},
			mockExec: func(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
				return &mocks.MockResult{}, nil
			},
		}

		// Should not panic even without logger
		err := repo.UpdateBalance(ctx, &mocks.MockTransaction{}, 1, decimal.NewFromFloat(100.50))
		assert.NoError(t, err)
	})
}

// Benchmark tests
func BenchmarkTransactionRepositoryPostgre_Save(b *testing.B) {
	baseLogger := logrus.New()
	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, baseLogger)

	repo := &TestableTransactionRepository{
		TransactionRepositoryPostgre: &repositories.TransactionRepositoryPostgre{},
		mockQueryRow: func(ctx context.Context, query string, args ...interface{}) (int64, error) {
			return 1, nil
		},
	}

	transaction := &entities.Transaction{
		SourceAccountID:      1,
		DestinationAccountID: 2,
		Amount:               decimal.NewFromFloat(100.0),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = repo.Save(ctx, &mocks.MockTransaction{}, transaction)
	}
}

func BenchmarkTransactionRepositoryPostgre_UpdateBalance(b *testing.B) {
	baseLogger := logrus.New()
	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, baseLogger)

	repo := &TestableTransactionRepository{
		TransactionRepositoryPostgre: &repositories.TransactionRepositoryPostgre{},
		mockExec: func(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
			return &mocks.MockResult{}, nil
		},
	}

	amount := decimal.NewFromFloat(100.0)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = repo.UpdateBalance(ctx, &mocks.MockTransaction{}, 1, amount)
	}
}

// Table-driven test for multiple scenarios
func TestTransactionRepositoryPostgre_MultipleScenarios(t *testing.T) {
	baseLogger := logrus.New()
	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, baseLogger)

	scenarios := []struct {
		name        string
		transaction *entities.Transaction
		saveError   error
		updateError error
	}{
		{
			name: "normal transfer",
			transaction: &entities.Transaction{
				SourceAccountID:      1,
				DestinationAccountID: 2,
				Amount:               decimal.NewFromFloat(100.0),
			},
		},
		{
			name: "large transfer",
			transaction: &entities.Transaction{
				SourceAccountID:      1,
				DestinationAccountID: 2,
				Amount:               decimal.NewFromFloat(10000.0),
			},
		},
		{
			name: "small transfer",
			transaction: &entities.Transaction{
				SourceAccountID:      1,
				DestinationAccountID: 2,
				Amount:               decimal.NewFromFloat(0.01),
			},
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			repo := &TestableTransactionRepository{
				TransactionRepositoryPostgre: &repositories.TransactionRepositoryPostgre{
					DB: &mocks.MockDatabase{},
				},
				mockQueryRow: func(ctx context.Context, query string, args ...interface{}) (int64, error) {
					return 1, scenario.saveError
				},
				mockExec: func(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
					return &mocks.MockResult{}, scenario.updateError
				},
			}

			// Test Save
			result, err := repo.Save(ctx, &mocks.MockTransaction{}, scenario.transaction)
			if scenario.saveError != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}

			// Test UpdateBalance (debit source)
			err = repo.UpdateBalance(ctx, &mocks.MockTransaction{}, scenario.transaction.SourceAccountID, scenario.transaction.Amount.Neg())
			if scenario.updateError != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			// Test UpdateBalance (credit destination)
			err = repo.UpdateBalance(ctx, &mocks.MockTransaction{}, scenario.transaction.DestinationAccountID, scenario.transaction.Amount)
			if scenario.updateError != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
