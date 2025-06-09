package ports

import (
	"github.com/labstack/echo/v4"
)

type TransactionController interface {
	Save(ctx echo.Context) error
}
