package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"transfer-system/domain/entities"
	"transfer-system/domain/ports"

	"github.com/sirupsen/logrus"
)

type AccountRepositoryPostgre struct {
	DB     ports.Database
	logger *logrus.Entry
}

func (repository *AccountRepositoryPostgre) Save(ctx context.Context, tx ports.Transaction, account *entities.Account) (*entities.Account, error) {
	var id int64
	var createdAt time.Time
	query := `
            INSERT INTO accounts (balance)
            VALUES ($1)
            RETURNING id, balance`
	err := tx.QueryRowContext(ctx, query, account.Balance).Scan(&id, &createdAt)
	if err != nil {
		return nil, err
	}

	account.AccountId = id
	account.Balance = account.Balance

	return account, nil
}

func (r *AccountRepositoryPostgre) FindById(ctx context.Context, tx ports.Transaction, id int64) (*entities.Account, error) {
	account := &entities.Account{}
	query := "SELECT id, balance FROM accounts WHERE id = $1"
	err := tx.QueryRowContext(ctx, query, id).Scan(&account.AccountId, &account.Balance)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("account not found")
		}
		return nil, err
	}

	return account, nil
}
