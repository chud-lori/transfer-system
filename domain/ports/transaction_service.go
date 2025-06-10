package ports

import (
	"context"
	"transfer-system/domain/entities"
)

type TransactionService interface {
	Save(ctx context.Context, request *entities.Transaction) error
}
