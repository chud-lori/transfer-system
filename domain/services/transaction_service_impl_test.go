package services

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"transfer-system/domain/entities"
	"transfer-system/domain/ports"
	"transfer-system/internal/testutils"
	appErrors "transfer-system/pkg/errors"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// -- Mock Definitions --

type MockTransactionRepository struct {
	mock.Mock
}

func (m *MockTransactionRepository) Save(ctx context.Context, tx ports.Transaction, trans *entities.Transaction) (*entities.Transaction, error) {
	args := m.Called(ctx, tx, trans)
	transaction, _ := args.Get(0).(*entities.Transaction)
	return transaction, args.Error(1)
}

func (m *MockTransactionRepository) UpdateBalance(ctx context.Context, tx ports.Transaction, accId int64, amount decimal.Decimal) error {
	args := m.Called(ctx, tx, accId, amount)
	return args.Error(0)
}

func TestTransactionService_Save_Success(t *testing.T) {
	ctx := testutils.InjectLoggerIntoContext(context.Background())

	mockDB := new(MockDatabase)
	mockAccRepo := new(MockAccountRepository)
	mockRepo := new(MockTransactionRepository)
	mockTx := new(MockTransaction)

	service := &TransactionServiceImpl{
		DB:                    mockDB,
		TransactionRepository: mockRepo,
		AccountRepository:     mockAccRepo,
		CtxTimeout:            2 * time.Second,
	}

	sourceAccount := &entities.Account{
		AccountID: 123,
		Balance:   decimal.NewFromFloat(100.23344),
	}
	destinationAccount := &entities.Account{
		AccountID: 456,
		Balance:   decimal.NewFromFloat(100.23344),
	}

	transaction := &entities.Transaction{
		SourceAccountID:      123,
		DestinationAccountID: 456,
		Amount:               decimal.NewFromFloat(100.23344),
	}

	mockDB.On("BeginTx", mock.Anything).Return(mockTx, nil)
	mockAccRepo.On("FindById", mock.Anything, mockTx, transaction.SourceAccountID).Return(sourceAccount, nil)
	mockAccRepo.On("FindById", mock.Anything, mockTx, transaction.DestinationAccountID).Return(destinationAccount, nil)
	mockRepo.On("Save", mock.Anything, mock.Anything, transaction).Return(transaction, nil).Once()
	mockRepo.On("UpdateBalance", mock.Anything, mockTx, transaction.SourceAccountID, transaction.Amount.Neg()).Return(nil)
	mockRepo.On("UpdateBalance", mock.Anything, mockTx, transaction.DestinationAccountID, transaction.Amount).Return(nil)
	mockTx.On("Commit").Return(nil)

	err := service.Save(ctx, transaction)
	assert.NoError(t, err)

	mockDB.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
	mockAccRepo.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}

func TestTransactionService_Save_AccountNotFound(t *testing.T) {
	ctx := testutils.InjectLoggerIntoContext(context.Background())

	mockDB := new(MockDatabase)
	mockAccRepo := new(MockAccountRepository)
	mockRepo := new(MockTransactionRepository)
	mockTx := new(MockTransaction)

	service := &TransactionServiceImpl{
		DB:                    mockDB,
		TransactionRepository: mockRepo,
		AccountRepository:     mockAccRepo,
		CtxTimeout:            2 * time.Second,
	}

	transaction := &entities.Transaction{
		SourceAccountID:      123,
		DestinationAccountID: 456,
		Amount:               decimal.NewFromFloat(100.23344),
	}

	mockDB.On("BeginTx", mock.Anything).Return(mockTx, nil)
	mockAccRepo.On("FindById", mock.Anything, mockTx, transaction.SourceAccountID).Return(nil, sql.ErrNoRows)
	// mockTx.On("Commit").Return(nil).Times(0)
	mockTx.On("Rollback").Return(nil)

	err := service.Save(ctx, transaction)

	assert.Error(t, err)

	appErr, ok := err.(*appErrors.AppError)
	assert.True(t, ok)
	assert.Equal(t, "Account Not Found", appErr.Message)

	mockDB.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
	mockAccRepo.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}

func TestTransactionService_Save_InsufficientBalance(t *testing.T) {
	ctx := testutils.InjectLoggerIntoContext(context.Background())

	mockDB := new(MockDatabase)
	mockAccRepo := new(MockAccountRepository)
	mockRepo := new(MockTransactionRepository)
	mockTx := new(MockTransaction)

	service := &TransactionServiceImpl{
		DB:                    mockDB,
		TransactionRepository: mockRepo,
		AccountRepository:     mockAccRepo,
		CtxTimeout:            2 * time.Second,
	}

	sourceAccount := &entities.Account{
		AccountID: 123,
		Balance:   decimal.NewFromFloat(10.23344),
	}
	destinationAccount := &entities.Account{
		AccountID: 456,
		Balance:   decimal.NewFromFloat(100.23344),
	}

	transaction := &entities.Transaction{
		SourceAccountID:      123,
		DestinationAccountID: 456,
		Amount:               decimal.NewFromFloat(100.23344),
	}

	mockDB.On("BeginTx", mock.Anything).Return(mockTx, nil)
	mockAccRepo.On("FindById", mock.Anything, mockTx, transaction.SourceAccountID).Return(sourceAccount, nil)
	mockAccRepo.On("FindById", mock.Anything, mockTx, transaction.DestinationAccountID).Return(destinationAccount, nil)
	mockTx.On("Rollback").Return(nil).Once()

	err := service.Save(ctx, transaction)
	assert.Error(t, err)

	appErr, ok := err.(*appErrors.AppError)
	assert.True(t, ok)
	assert.Equal(t, "Insufficient balance", appErr.Message)

	mockDB.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
	mockAccRepo.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}
