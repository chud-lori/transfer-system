package mocks

import (
	"context"
	"transfer-system/domain/entities"

	"github.com/stretchr/testify/mock"
)

type MockAccountService struct {
	mock.Mock
}

func (m *MockAccountService) Save(ctx context.Context, acc *entities.Account) error {
	args := m.Called(ctx, acc)
	return args.Error(0)
}

func (m *MockAccountService) FindById(ctx context.Context, accountId int64) (*entities.Account, error) {
	args := m.Called(ctx, accountId)
	return args.Get(0).(*entities.Account), args.Error(1)
}
