package ports

import (
	"context"

	"transfer-system/domain/entities"
)

type TransactionRepository interface {
	Save(ctx context.Context, tx Transaction, transaction *entities.Transaction) (*entities.Transaction, error)
}
