package repositories

import (
	"context"
	"fmt"

	"transfer-system/domain/entities"
	"transfer-system/domain/ports"
	"transfer-system/pkg/logger"

	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

type TransactionRepositoryPostgre struct {
	DB ports.Database
}

func (repository *TransactionRepositoryPostgre) Save(ctx context.Context, tx ports.Transaction, transaction *entities.Transaction) (*entities.Transaction, error) {
	logger, _ := ctx.Value(logger.LoggerContextKey).(logrus.FieldLogger)

	var transactionId int64
	query := `
            INSERT INTO transactions (source_id, destination_id, amount)
            VALUES ($1, $2, $3)
			RETURNING id`
	err := tx.QueryRowContext(ctx, query, transaction.SourceAccountID, transaction.DestinationAccountID, transaction.Amount).Scan(&transactionId)
	if err != nil {
		logger.WithError(err).Error("Failed to insert transaction")
		return nil, err
	}

	transaction.Id = transactionId

	return transaction, nil
}

func (repository *TransactionRepositoryPostgre) UpdateBalance(ctx context.Context, tx ports.Transaction, accountID int64, amount decimal.Decimal) error {
	logger, _ := ctx.Value(logger.LoggerContextKey).(logrus.FieldLogger)

	query := `
			UPDATE accounts
			SET balance = balance + $1
			WHERE id = $2`
	res, err := tx.ExecContext(ctx, query, amount, accountID)
	if err != nil {
		logger.WithError(err).Error("Failed to update account balance")
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		logger.WithError(err).Error("Failed to get rows affected")
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no account found with id %d", accountID)
	}

	return nil
}
