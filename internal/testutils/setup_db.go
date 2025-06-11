package testutils

import (
	"context"
	"os"
	"testing"
	"transfer-system/domain/ports"
	"transfer-system/infrastructure/datastore"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func SetupTestDB(t *testing.T) ports.Database {
	err := godotenv.Load("../../.env.test")
	require.NoError(t, err, "failed to load .env.test")

	dbURL := os.Getenv("DB_URL")
	require.NotEmpty(t, dbURL, "DB_URL must be set in .env.test")

	baseLogger := logrus.New()
	db, err := datastore.NewDatabase(dbURL, baseLogger)
	require.NoError(t, err, "failed to connect to test database")

	return db
}

func SetupTestTx(t *testing.T, db ports.Database) ports.Transaction {
	ctx := context.Background()
	tx, err := db.BeginTx(ctx)
	require.NoError(t, err, "failed to begin transaction")
	return tx
}
