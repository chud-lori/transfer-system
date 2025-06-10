package ports

import (
	"context"

	"transfer-system/domain/entities"
)

type AccountService interface {
	Save(ctx context.Context, request *entities.Account) error
	FindById(ctx context.Context, id int64) (*entities.Account, error)
}
