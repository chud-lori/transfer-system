package mocks

import (
	"context"
	"transfer-system/domain/entities"
	"transfer-system/domain/ports"

	"github.com/stretchr/testify/mock"
)

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
