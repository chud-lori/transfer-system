package testutils

import (
	"context"
	"transfer-system/pkg/logger"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func InjectLoggerIntoContext(ctx context.Context) context.Context {
	baseLogger := logrus.NewEntry(logrus.New())
	baseLogger.Logger.SetLevel(logrus.DebugLevel)
	return context.WithValue(ctx, logger.LoggerContextKey, baseLogger)
}

func InjectLoggerToContext(ctx echo.Context) {
	baseLogger := logrus.NewEntry(logrus.New())
	newCtx := context.WithValue(ctx.Request().Context(), logger.LoggerContextKey, baseLogger)
	ctx.SetRequest(ctx.Request().WithContext(newCtx))
}
