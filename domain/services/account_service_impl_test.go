package services

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"transfer-system/domain/entities"
	"transfer-system/domain/ports"
	"transfer-system/pkg/logger"

	appErrors "transfer-system/pkg/errors"

	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// -- Mock Definitions --

type MockTransaction struct {
	mock.Mock
}

func (m *MockTransaction) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	argsM := m.Called(ctx, query, args)
	return argsM.Get(0).(sql.Result), argsM.Error(1)
}
func (m *MockTransaction) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return &sql.Row{} // stubbed, not used
}
func (m *MockTransaction) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return nil, nil // stubbed, not used
}
func (m *MockTransaction) Commit() error {
	return m.Called().Error(0)
}
func (m *MockTransaction) Rollback() error {
	return m.Called().Error(0)
}

type MockDatabase struct {
	mock.Mock
}

func (m *MockDatabase) BeginTx(ctx context.Context) (ports.Transaction, error) {
	args := m.Called(ctx)
	return args.Get(0).(ports.Transaction), args.Error(1)
}
func (m *MockDatabase) Close() error {
	return nil
}

type MockAccountRepository struct {
	mock.Mock
}

func (m *MockAccountRepository) FindById(ctx context.Context, tx ports.Transaction, id int64) (*entities.Account, error) {
	args := m.Called(ctx, tx, id)
	account, _ := args.Get(0).(*entities.Account)
	return account, args.Error(1)
}

func (m *MockAccountRepository) Save(ctx context.Context, tx ports.Transaction, acc *entities.Account) (*entities.Account, error) {
	args := m.Called(ctx, tx, acc)
	account, _ := args.Get(0).(*entities.Account)
	return account, args.Error(1)
}

func TestAccountService_Save_Success(t *testing.T) {
	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))

	mockDB := new(MockDatabase)
	mockRepo := new(MockAccountRepository)
	mockTx := new(MockTransaction)

	service := &AccountServiceImpl{
		DB:                mockDB,
		AccountRepository: mockRepo,
		CtxTimeout:        2 * time.Second,
	}

	account := &entities.Account{
		AccountID: 12345,
		Balance:   decimal.NewFromFloat(100.23344),
	}

	mockDB.On("BeginTx", mock.Anything).Return(mockTx, nil)
	mockRepo.On("FindById", mock.Anything, mockTx, account.AccountID).Return(nil, sql.ErrNoRows)
	mockRepo.On("Save", mock.Anything, mockTx, account).Return(account, nil)
	mockTx.On("Commit").Return(nil)

	err := service.Save(ctx, account)
	assert.NoError(t, err)

	mockDB.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}

func TestAccountService_Save_AccountExists(t *testing.T) {
	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))

	mockDB := new(MockDatabase)
	mockRepo := new(MockAccountRepository)
	mockTx := new(MockTransaction)

	service := &AccountServiceImpl{
		DB:                mockDB,
		AccountRepository: mockRepo,
		CtxTimeout:        time.Second * 2,
	}

	existingAcc := &entities.Account{AccountID: 12345, Balance: decimal.NewFromFloat(100.23344)}
	mockDB.On("BeginTx", mock.Anything).Return(mockTx, nil)
	mockRepo.On("FindById", mock.Anything, mockTx, existingAcc.AccountID).Return(existingAcc, nil)
	mockTx.On("Commit").Return(nil)
	mockTx.On("Rollback").Return(nil) // Ensure it's set even if not called (defer func includes it)

	err := service.Save(ctx, existingAcc)
	assert.Error(t, err)

	appErr, ok := err.(*appErrors.AppError)
	assert.True(t, ok)
	assert.Equal(t, "AccountId already exists", appErr.Message)
	assert.Equal(t, 400, appErr.StatusCode)
}

func TestAccountService_FindById_Success(t *testing.T) {
	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))

	mockDB := new(MockDatabase)
	mockRepo := new(MockAccountRepository)
	mockTx := new(MockTransaction)

	service := &AccountServiceImpl{
		DB:                mockDB,
		AccountRepository: mockRepo,
		CtxTimeout:        time.Second * 2,
	}

	acc := &entities.Account{
		AccountID: 12345,
		Balance:   decimal.NewFromFloat(100.23344),
	}

	mockDB.On("BeginTx", mock.Anything).Return(mockTx, nil)
	mockRepo.On("FindById", mock.Anything, mockTx, acc.AccountID).Return(acc, nil)
	mockTx.On("Commit").Return(nil)
	mockTx.On("Rollback").Return(nil)

	resp, err := service.FindById(ctx, acc.AccountID)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, acc.AccountID, resp.AccountID)
	assert.Equal(t, acc.Balance, resp.Balance)
}

func TestAccountService_FindById_NotFound(t *testing.T) {
	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))

	mockDB := new(MockDatabase)
	mockRepo := new(MockAccountRepository)
	mockTx := new(MockTransaction)

	service := &AccountServiceImpl{
		DB:                mockDB,
		AccountRepository: mockRepo,
		CtxTimeout:        time.Second * 2,
	}

	accID := int64(1)

	mockDB.On("BeginTx", mock.Anything).Return(mockTx, nil)
	mockRepo.On("FindById", mock.Anything, mockTx, accID).Return(&entities.Account{}, sql.ErrNoRows)
	mockTx.On("Commit").Return(nil)
	mockTx.On("Rollback").Return(nil)

	resp, err := service.FindById(ctx, accID)
	assert.Nil(t, resp)

	appErr, ok := err.(*appErrors.AppError)
	assert.True(t, ok)
	assert.Equal(t, "Account not found", appErr.Message)
}
