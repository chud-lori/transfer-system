package testutils

import (
	"context"
	"transfer-system/pkg/logger"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func InjectLoggerToContext(ctx echo.Context) {
	baseLogger := logrus.NewEntry(logrus.New())
	newCtx := context.WithValue(ctx.Request().Context(), logger.LoggerContextKey, baseLogger)
	ctx.SetRequest(ctx.Request().WithContext(newCtx))
}
