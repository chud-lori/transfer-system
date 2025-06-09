package logger

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type loggingTraffic struct {
	http.ResponseWriter
	statusCode int
}

func NewLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logrus.InfoLevel)
	logger.SetOutput(os.Stdout)
	return logger
}

const LoggerContextKey string = "logger"

func NewLoggingTraffic(w http.ResponseWriter) *loggingTraffic {
	return &loggingTraffic{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}
}

func (lrw *loggingTraffic) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func LogTrafficMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		start := time.Now()

		request := ctx.Request()
		response := ctx.Response()

		requestID := request.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		baseLogger := logrus.New()
		baseLogger.SetFormatter(&logrus.JSONFormatter{})
		baseLogger.SetFormatter(&logrus.JSONFormatter{})
		baseLogger.SetLevel(logrus.InfoLevel)
		baseLogger.SetOutput(os.Stdout)
		baseLogger.SetReportCaller(true)

		logger := baseLogger.WithField("RequestID", requestID)
		newCtx := context.WithValue(request.Context(), "logger", logger)
		ctx.SetRequest(request.WithContext(newCtx))

		lrw := NewLoggingTraffic(response.Writer)
		response.Writer = lrw

		err := next(ctx)

		duration := time.Since(start)

		logger.WithFields(logrus.Fields{
			"method":   request.Method,
			"path":     request.URL.Path,
			"duration": duration.String(),
			"status":   lrw.statusCode,
		}).Info("Processed request")

		return err
	}
}
