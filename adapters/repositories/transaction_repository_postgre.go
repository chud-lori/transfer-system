package repositories

import (
	"context"

	"transfer-system/domain/entities"
	"transfer-system/domain/ports"

	"github.com/sirupsen/logrus"
)

type TransactionRepositoryPostgre struct {
	DB     ports.Database
	logger *logrus.Entry
}

func (repository *TransactionRepositoryPostgre) Save(ctx context.Context, tx ports.Transaction, transaction *entities.Transaction) (*entities.Transaction, error) {

	var transactionId int64
	query := `
            INSERT INTO transactions (from_id, to_id, amount)
            VALUES ($1, $2, $3)
			RETURNING id`
	err := tx.QueryRowContext(ctx, query, transaction.SourceAccountId, transaction.DestinationAccountID, transaction.Amount).Scan(&transactionId)
	if err != nil {
		return nil, err
	}

	transaction.Id = transactionId

	return transaction, nil
}
