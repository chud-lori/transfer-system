package services_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"transfer-system/domain/entities"
	"transfer-system/domain/services"
	"transfer-system/mocks"
	appErrors "transfer-system/pkg/errors"
	"transfer-system/pkg/logger"

	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTransactionService_Save_Success(t *testing.T) {
	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))

	mockDB := new(mocks.MockDatabase)
	mockAccRepo := new(mocks.MockAccountRepository)
	mockRepo := new(mocks.MockTransactionRepository)
	mockTx := new(mocks.MockTransaction)

	service := &services.TransactionServiceImpl{
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
	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))

	mockDB := new(mocks.MockDatabase)
	mockAccRepo := new(mocks.MockAccountRepository)
	mockRepo := new(mocks.MockTransactionRepository)
	mockTx := new(mocks.MockTransaction)

	service := &services.TransactionServiceImpl{
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
	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))

	mockDB := new(mocks.MockDatabase)
	mockAccRepo := new(mocks.MockAccountRepository)
	mockRepo := new(mocks.MockTransactionRepository)
	mockTx := new(mocks.MockTransaction)

	service := &services.TransactionServiceImpl{
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
