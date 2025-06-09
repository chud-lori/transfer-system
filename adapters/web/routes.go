package web

import (
	"transfer-system/domain/ports"

	"github.com/labstack/echo/v4"
)

func AccountRouter(controller ports.AccountController, e *echo.Echo) {
	e.POST("/accounts", controller.Create)
	e.GET("/accounts/:accountId", controller.FindById)
}

func TransactionRouter(controller ports.TransactionController, e *echo.Echo) {
	e.POST("/transactions", controller.Save)
}
