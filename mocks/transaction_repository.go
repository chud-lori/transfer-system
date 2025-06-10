package mocks

import (
	"context"
	"transfer-system/domain/entities"
	"transfer-system/domain/ports"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/mock"
)

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
