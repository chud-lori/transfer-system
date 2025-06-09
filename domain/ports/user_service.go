package ports

import (
	"context"

	"transfer-system/adapters/transport"
)

type UserService interface {
	Save(ctx context.Context, request *transport.UserRequest) (*transport.UserResponse, error)
	Update(ctx context.Context, request *transport.UserRequest) (*transport.UserResponse, error)
	Delete(ctx context.Context, id string) error
	FindById(ctx context.Context, id string) (*transport.UserResponse, error)
	FindAll(ctx context.Context) ([]*transport.UserResponse, error)
}
