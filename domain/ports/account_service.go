package ports

import (
	"context"

	"transfer-system/adapters/web/dto"
	"transfer-system/domain/entities"
)

type AccountService interface {
	Save(ctx context.Context, request *entities.Account) (*dto.WebResponse, error)
	FindById(ctx context.Context, id int64) (*dto.AccountResponse, error)
}
