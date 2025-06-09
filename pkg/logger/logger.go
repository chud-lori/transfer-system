package logger

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
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

func LogTrafficMiddleware(next http.Handler, baseLogger *logrus.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		logger := baseLogger.WithField("RequestID", requestID)

		ctx := context.WithValue(r.Context(), "logger", logger)

		r = r.WithContext(ctx)

		lrw := NewLoggingTraffic(w)

		// call the next handler
		next.ServeHTTP(lrw, r)

		// TODO: if showing source in log
        // baseLogger.SetReportCaller(true)
		//_, file, line, ok := runtime.Caller(1)
		//source := "unknown"
		//if ok {
		//    source = fmt.Sprintf("%s:%d", file, line)
		//}

		duration := time.Since(start)

		logger.WithFields(logrus.Fields{
			"method":   r.Method,
			"path":     r.URL.Path,
			"duration": duration.String(),
			"status":   lrw.statusCode,
		}).Info("Processed request")

	})
}
