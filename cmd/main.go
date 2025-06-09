package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"transfer-system/adapters/controllers"
	"transfer-system/adapters/repositories"
	"transfer-system/adapters/utils"
	"transfer-system/adapters/web"
	"transfer-system/domain/services"
	"transfer-system/infrastructure/datastore"
	"transfer-system/pkg/logger"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Failed load keys")
	}

	baseLogger := logger.NewLogger()

	db, err := datastore.NewDatabase(os.Getenv("DB_URL"), baseLogger)
	if err != nil {
		baseLogger.Fatal("Failed to connect to database: ", err)
	}
	defer db.Close()

	ctxTimeout := time.Duration(60) * time.Second

	accountRepository := &repositories.AccountRepositoryPostgre{
		DB: db,
	}
	accountService := &services.AccountServiceImpl{
		DB:                db,
		AccountRepository: accountRepository,
		CtxTimeout:        ctxTimeout,
	}
	accountController := &controllers.AccountController{
		AccountService: accountService,
	}

	e := echo.New()

	web.AccountRouter(accountController, e)

	e.Use(logger.LogTrafficMiddleware)

	// Run server in a goroutine
	go func() {
		log.Printf("Server is running on port %s", os.Getenv("APP_PORT"))
		if err := e.Start(fmt.Sprintf(":%s", os.Getenv("APP_PORT"))); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatalf("HTTP server error: %v", err)
		}
	}()

	// graceful shutdown
	wait := utils.GracefullShutdown(context.Background(), 5*time.Second, map[string]utils.Operation{
		"database": func(ctx context.Context) error {
			return db.Close()
		},
		"http-server": func(ctx context.Context) error {
			return e.Shutdown(ctx)
		},
	})

	<-wait
}
