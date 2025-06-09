package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"transfer-system/domain/entities"
	"transfer-system/domain/ports"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type UserRepositoryPostgre struct {
	DB     ports.Database
	logger *logrus.Entry
}

func (repository *UserRepositoryPostgre) Save(ctx context.Context, tx ports.Transaction, user *entities.User) (*entities.User, error) {
	var id string
	var createdAt time.Time
	query := `
            INSERT INTO users (email, passcode)
            VALUES ($1, $2)
            RETURNING id, created_at`
	err := tx.QueryRowContext(ctx, query, user.Email, user.Passcode).Scan(&id, &createdAt)
	if err != nil {
		return nil, err
	}

	user.Id = id
	user.Created_at = createdAt

	return user, nil
}

func (repository *UserRepositoryPostgre) Update(ctx context.Context, tx ports.Transaction, user *entities.User) (*entities.User, error) {
	query := "UPDATE users SET email = $1, passcode = $2 WHERE id = $3"
	_, err := tx.ExecContext(ctx, query, user.Email, user.Passcode, user.Id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (repository *UserRepositoryPostgre) Delete(ctx context.Context, tx ports.Transaction, id string) error {
	query := "DELETE FROM users WHERE id = $1"
	result, err := tx.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

func (r *UserRepositoryPostgre) FindById(ctx context.Context, tx ports.Transaction, id string) (*entities.User, error) {
	if _, err := uuid.Parse(id); err != nil {
		//r.logger.Info("Invalid UUID Format: ", id)
		return nil, fmt.Errorf("Invalid UUID Format")
	}

	user := &entities.User{}
	query := "SELECT id, email, created_at FROM users WHERE id = $1"
	err := tx.QueryRowContext(ctx, query, id).Scan(&user.Id, &user.Email, &user.Created_at)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	return user, nil
}

func (repository *UserRepositoryPostgre) FindAll(ctx context.Context, tx ports.Transaction) ([]*entities.User, error) {
	query := "SELECT id, email, created_at FROM users"
	rows, err := tx.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*entities.User
	for rows.Next() {
		var user entities.User
		err := rows.Scan(&user.Id, &user.Email, &user.Created_at)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
