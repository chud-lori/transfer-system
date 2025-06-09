package ports

import (
	"github.com/labstack/echo/v4"
)

type AccountController interface {
	Create(ctx echo.Context) error
	FindById(ctx echo.Context) error
}
