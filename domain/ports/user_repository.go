package ports

import (
	"context"

	"transfer-system/domain/entities"
)

type UserRepository interface {
	Save(ctx context.Context, tx Transaction, user *entities.User) (*entities.User, error)
	FindById(ctx context.Context, tx Transaction, id string) (*entities.User, error)
}
