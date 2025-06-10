package services

import (
	"context"
	"database/sql"
	"errors"
	"time"

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

func (s *TransactionServiceImpl) Save(c context.Context, request *entities.Transaction) error {
	// TODO: improvment to add identifier id for each transaction from client side to make it idempotent
	logger, _ := c.Value(logger.LoggerContextKey).(logrus.FieldLogger)

	ctx, cancel := context.WithTimeout(c, s.CtxTimeout)
	defer cancel()

	tx, err := s.DB.BeginTx(ctx)
	if err != nil {
		return err
	}
	// handle panic gracefully
	defer func() {
		if r := recover(); r != nil || err != nil {
			logger.Errorf("Transaction rollback due to error: %v", err)
			logger.Errorf("Transaction rollback due to panic: %v", r)
			tx.Rollback()
		}
	}()

	// check source account exist
	sourceAccount, err := s.AccountRepository.FindById(ctx, tx, request.SourceAccountID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Errorf("AccountID %d not found", request.SourceAccountID)
			return appErrors.NewBadRequestError("Account Not Found", err)
		}
		logger.WithError(err).Error("Database error")
		return appErrors.NewInternalServerError("Currently we're facing an issue", err)
	}

	// check destination account exist
	_, err = s.AccountRepository.FindById(ctx, tx, request.DestinationAccountID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Errorf("AccountID %d not found", request.SourceAccountID)
			return appErrors.NewBadRequestError("Account Not Found", err)
		}
		logger.WithError(err).Error("Database error")
		return appErrors.NewInternalServerError("Currently we're facing an issue", err)
	}

	// check if source account has sufficient balance
	if sourceAccount.Balance.LessThan(request.Amount) {
		logger.Errorf("Insufficient balance in source account id %d", request.SourceAccountID)
		// to trigger rollback
		err = appErrors.NewBadRequestError("Insufficient balance", nil)
		return err
	}

	transaction := entities.Transaction{
		Id:                   0,
		SourceAccountID:      request.SourceAccountID,
		DestinationAccountID: request.DestinationAccountID,
		Amount:               request.Amount,
	}
	_, err = s.TransactionRepository.Save(ctx, tx, &transaction)

	if err != nil {
		logger.WithError(err).Error("Failed to save transaction")
		return err
	}
	// logger.Debugf("DEBUG: Calling UpdateBalance for Destination. Tx type: %T, Tx value: %#v, AccountID: %d, Amount: %s", tx, tx, request.DestinationAccountID, request.Amount.String())

	// update balance of source account
	err = s.TransactionRepository.UpdateBalance(ctx, tx, request.SourceAccountID, request.Amount.Neg())

	if err != nil {
		logger.WithError(err).Error("Failed to update source account balance")
		return err
	}

	// update balance of destination account
	err = s.TransactionRepository.UpdateBalance(ctx, tx, request.DestinationAccountID, request.Amount)

	if err != nil {
		logger.WithError(err).Error("Failed to update destination account balance")
		return err
	}

	if err := tx.Commit(); err != nil {
		logger.WithError(err).Error("Failed to commit transaction")
		return err
	}

	return nil
}
