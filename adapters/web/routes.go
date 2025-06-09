package web

import (
	"transfer-system/domain/ports"

	"github.com/labstack/echo/v4"
)

func UserRouter(controller ports.UserController, e *echo.Echo) {
	// e.POST("/api/user", controller.Create)
	// e.PUT("/api/user/:userId", controller.Update)
	// e.DELETE("/api/user/:userId", controller.Delete)
	// e.GET("/api/user/:userId", controller.FindById)
	// e.GET("/api/user", controller.FindAll)
	e.POST("/accounts", controller.Create)
	e.GET("/accounts/:accountId", controller.FindById)
}

func TransferRouter(controller ports.TransferController, e *echo.Echo) {
	e.POST("/transactions", controller.Transaction)
}
