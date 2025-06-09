package ports

import (
	"context"

	"transfer-system/domain/entities"
)

type AccountRepository interface {
	Save(ctx context.Context, tx Transaction, account *entities.Account) (*entities.Account, error)
	FindById(ctx context.Context, tx Transaction, id int64) (*entities.Account, error)
}
