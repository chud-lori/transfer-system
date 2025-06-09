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

// Add this new function at the top
func NewLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logrus.InfoLevel)
	logger.SetOutput(os.Stdout)
	return logger
}

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

		logger := baseLogger.WithField("RequestID", requestID)

		newCtx := context.WithValue(request.Context(), "logger", logger)
		ctx.SetRequest(request.WithContext(newCtx))

		lrw := NewLoggingTraffic(response.Writer)
		response.Writer = lrw

		// call the next handler
		err := next(ctx)
		// next.ServeHTTP(lrw, r)

		// TODO: if showing source in log
		// baseLogger.SetReportCaller(true)
		//_, file, line, ok := runtime.Caller(1)
		//source := "unknown"
		//if ok {
		//    source = fmt.Sprintf("%s:%d", file, line)
		//}

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
