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

	_ "transfer-system/docs"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title Transfer System API
// @version 1.0
// @description A simple API for managing accounts and transactions in a transfer system

// @host localhost:8080
// @BasePath /
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

	// Initialize repositories and services for account
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

	// Initialize repositories and services for transaction
	transactionRepository := &repositories.TransactionRepositoryPostgre{
		DB: db,
	}
	transactionService := &services.TransactionServiceImpl{
		DB:                    db,
		TransactionRepository: transactionRepository,
		AccountRepository:     accountRepository,
		CtxTimeout:            ctxTimeout,
	}
	transactionController := &controllers.TransactionController{
		TransactionService: transactionService,
	}

	e := echo.New()
	e.GET("/docs/*", echoSwagger.WrapHandler)

	web.AccountRouter(accountController, e)
	web.TransactionRouter(transactionController, e)

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
