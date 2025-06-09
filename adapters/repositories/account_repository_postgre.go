package repositories

import (
	"context"
	"database/sql"

	"transfer-system/domain/entities"
	"transfer-system/domain/ports"
	"transfer-system/pkg/logger"

	"github.com/sirupsen/logrus"
)

type AccountRepositoryPostgre struct {
	DB ports.Database
}

func (repository *AccountRepositoryPostgre) Save(ctx context.Context, tx ports.Transaction, account *entities.Account) (*entities.Account, error) {
	logger, _ := ctx.Value(logger.LoggerContextKey).(logrus.FieldLogger)

	var id int64
	query := `
            INSERT INTO accounts (id, balance)
            VALUES ($1, $2)
            RETURNING id`
	err := tx.QueryRowContext(ctx, query, account.AccountID, account.Balance).Scan(&id)
	if err != nil {
		logger.WithError(err).Error("Failed to insert account")
		return nil, err
	}

	account.AccountID = id
	account.Balance = account.Balance

	return account, nil
}

func (r *AccountRepositoryPostgre) FindById(ctx context.Context, tx ports.Transaction, id int64) (*entities.Account, error) {
	logger, _ := ctx.Value(logger.LoggerContextKey).(logrus.FieldLogger)
	account := &entities.Account{}
	query := "SELECT id, balance FROM accounts WHERE id = $1 FOR UPDATE"
	err := tx.QueryRowContext(ctx, query, id).Scan(&account.AccountID, &account.Balance)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		logger.WithError(err).Error("Failed to query account by ID")
		return nil, err
	}

	return account, nil
}
