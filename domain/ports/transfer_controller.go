package ports

import (
	"github.com/labstack/echo/v4"
)

type TransferController interface {
	Transaction(ctx echo.Context) error
}
