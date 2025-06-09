package services

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"transfer-system/adapters/web/dto"
	"transfer-system/domain/entities"
	"transfer-system/domain/ports"
	appErrors "transfer-system/pkg/errors"
	"transfer-system/pkg/logger"

	"github.com/sirupsen/logrus"
)

type TransactionServiceImpl struct {
	DB                    ports.Database
	TransactionRepository ports.TransactionRepository
	AccountRepository     ports.AccountRepository
	CtxTimeout            time.Duration
}

func (s *TransactionServiceImpl) Save(c context.Context, request *dto.TransactionRequest) (*dto.WebResponse, error) {
	logger, _ := c.Value(logger.LoggerContextKey).(logrus.FieldLogger)

	ctx, cancel := context.WithTimeout(c, s.CtxTimeout)
	defer cancel()

	tx, err := s.DB.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	// handle panic gracefully
	defer func() {
		if r := recover(); r != nil || err != nil {
			tx.Rollback()
		}
	}()

	// check if balance is sufficient
	sourceAccount, err := s.AccountRepository.FindById(ctx, tx, request.SourceAccountID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Errorf("AccountID %d not found", request.SourceAccountID)
			return nil, appErrors.NewBadRequestError("AccountId already exists", err)
		}
		logger.WithError(err).Error("Database error")
		return nil, appErrors.NewInternalServerError("Currently we're facing an issue", err)
	}

	transaction := entities.Transaction{
		Id:                   0,
		SourceAccountId:      request.SourceAccountID,
		DestinationAccountID: request.DestinationAccountID,
		Amount:               request.Amount,
	}
	_, err = s.TransactionRepository.Save(ctx, tx, &transaction)

	if err != nil {
		return nil, err
	}

	accountResponse := &dto.WebResponse{
		Message: "success create account",
		Status:  1,
		Data:    nil,
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return accountResponse, nil
}
