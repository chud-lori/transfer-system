package ports

import (
	"context"

	"transfer-system/domain/entities"

	"github.com/shopspring/decimal"
)

type TransactionRepository interface {
	Save(ctx context.Context, tx Transaction, transaction *entities.Transaction) (*entities.Transaction, error)
	UpdateBalance(ctx context.Context, tx Transaction, accountID int64, amount decimal.Decimal) error
}
