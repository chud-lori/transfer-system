package ports

import (
	"context"

	"transfer-system/adapters/web/dto"
)

type AccountService interface {
	Save(ctx context.Context, request *dto.CreateAccountRequest) (*dto.WebResponse, error)
	FindById(ctx context.Context, id int64) (*dto.AccountResponse, error)
}
