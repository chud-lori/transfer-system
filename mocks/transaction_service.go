package mocks

import (
	"context"
	"transfer-system/domain/entities"

	"github.com/stretchr/testify/mock"
)

type MockTransactionService struct {
	mock.Mock
}

func (m *MockTransactionService) Save(ctx context.Context, req *entities.Transaction) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}
