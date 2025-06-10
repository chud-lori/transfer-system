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

type TestableAccountRepository struct {
	*repositories.AccountRepositoryPostgre
	mockQueryRow func(ctx context.Context, query string, args ...interface{}) (int64, error)
	mockFindRow  func(ctx context.Context, query string, args ...interface{}) (*entities.Account, error)
}

func (r *TestableAccountRepository) Save(ctx context.Context, tx ports.Transaction, account *entities.Account) (*entities.Account, error) {
	if r.mockQueryRow != nil {
		id, err := r.mockQueryRow(ctx, "", account.AccountID, account.Balance)
		if err != nil {
			return nil, err
		}
		account.AccountID = id
		return account, nil
	}
	return r.AccountRepositoryPostgre.Save(ctx, tx, account)
}

func (r *TestableAccountRepository) FindById(ctx context.Context, tx ports.Transaction, id int64) (*entities.Account, error) {
	if r.mockFindRow != nil {
		return r.mockFindRow(ctx, "", id)
	}
	return r.AccountRepositoryPostgre.FindById(ctx, tx, id)
}

func TestAccountRepositoryPostgre_Save(t *testing.T) {
	// Setup logger context
	baseLogger := logrus.New()
	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, baseLogger)

	tests := []struct {
		name          string
		account       *entities.Account
		mockReturn    int64
		mockError     error
		expectedError bool
		expectedID    int64
	}{
		{
			name: "successful save",
			account: &entities.Account{
				AccountID: 1,
				Balance:   decimal.NewFromFloat(100.50),
			},
			mockReturn:    123,
			mockError:     nil,
			expectedError: false,
			expectedID:    123,
		},
		{
			name: "database error",
			account: &entities.Account{
				AccountID: 1,
				Balance:   decimal.NewFromFloat(100.50),
			},
			mockReturn:    0,
			mockError:     errors.New("database error"),
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup testable repository
			repo := &TestableAccountRepository{
				AccountRepositoryPostgre: &repositories.AccountRepositoryPostgre{
					DB: &mocks.MockDatabase{},
				},
				mockQueryRow: func(ctx context.Context, query string, args ...interface{}) (int64, error) {
					return tt.mockReturn, tt.mockError
				},
			}

			// Execute
			result, err := repo.Save(ctx, &mocks.MockTransaction{}, tt.account)

			// Assert
			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectedID, result.AccountID)
			}
		})
	}
}

func TestAccountRepositoryPostgre_FindById(t *testing.T) {
	// Setup logger context
	baseLogger := logrus.New()
	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, baseLogger)

	tests := []struct {
		name            string
		accountID       int64
		mockAccount     *entities.Account
		mockError       error
		expectedError   bool
		expectedAccount *entities.Account
	}{
		{
			name:      "successful find",
			accountID: 1,
			mockAccount: &entities.Account{
				AccountID: 1,
				Balance:   decimal.NewFromFloat(100.50),
			},
			mockError:     nil,
			expectedError: false,
			expectedAccount: &entities.Account{
				AccountID: 1,
				Balance:   decimal.NewFromFloat(100.50),
			},
		},
		{
			name:            "account not found",
			accountID:       999,
			mockAccount:     nil,
			mockError:       sql.ErrNoRows,
			expectedError:   true,
			expectedAccount: nil,
		},
		{
			name:            "database error",
			accountID:       1,
			mockAccount:     nil,
			mockError:       errors.New("database error"),
			expectedError:   true,
			expectedAccount: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup testable repository
			repo := &TestableAccountRepository{
				AccountRepositoryPostgre: &repositories.AccountRepositoryPostgre{
					DB: &mocks.MockDatabase{},
				},
				mockFindRow: func(ctx context.Context, query string, args ...interface{}) (*entities.Account, error) {
					return tt.mockAccount, tt.mockError
				},
			}

			// Execute
			result, err := repo.FindById(ctx, &mocks.MockTransaction{}, tt.accountID)

			// Assert
			if tt.expectedError {
				assert.Error(t, err)
				if tt.mockError == sql.ErrNoRows {
					assert.Equal(t, sql.ErrNoRows, err)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectedAccount.AccountID, result.AccountID)
				assert.True(t, tt.expectedAccount.Balance.Equal(result.Balance))
			}
		})
	}
}

// Benchmark tests
func BenchmarkAccountRepositoryPostgre_Save(b *testing.B) {
	// Setup
	baseLogger := logrus.New()
	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, baseLogger)

	repo := &TestableAccountRepository{
		AccountRepositoryPostgre: &repositories.AccountRepositoryPostgre{},
		mockQueryRow: func(ctx context.Context, query string, args ...interface{}) (int64, error) {
			return 1, nil
		},
	}

	account := &entities.Account{
		AccountID: 1,
		Balance:   decimal.NewFromFloat(100.0),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = repo.Save(ctx, &mocks.MockTransaction{}, account)
	}
}
