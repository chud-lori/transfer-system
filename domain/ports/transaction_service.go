package ports

import (
	"context"
	"transfer-system/adapters/web/dto"
)

type TransactionService interface {
	Save(ctx context.Context, request *dto.TransactionRequest) (*dto.WebResponse, error)
}
