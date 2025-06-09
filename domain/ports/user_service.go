package ports

import (
	"context"

	"transfer-system/adapters/transport"
)

type UserService interface {
	Save(ctx context.Context, request *transport.UserRequest) (*transport.UserResponse, error)
	FindById(ctx context.Context, id string) (*transport.UserResponse, error)
}
