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

type AccountServiceImpl struct {
	DB                ports.Database
	AccountRepository ports.AccountRepository
	CtxTimeout        time.Duration
}

func (s *AccountServiceImpl) Save(c context.Context, request *dto.CreateAccountRequest) (*dto.WebResponse, error) {
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

	_, err = s.AccountRepository.FindById(ctx, tx, request.AccountID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {

		} else {
			logger.WithError(err).Error("Database error")
			return nil, appErrors.NewInternalServerError("Currently we're facing an issue", err)
		}
	} else {
		logger.Errorf("AccountID %d already exists", request.AccountID)
		return nil, appErrors.NewBadRequestError("AccountId already exists", err)
	}

	account := entities.Account{
		AccountId: request.AccountID,
		Balance:   request.InitialBalance,
	}
	_, err = s.AccountRepository.Save(ctx, tx, &account)

	if err != nil {
		logger.WithError(err).Error("Database error")
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

func (s *AccountServiceImpl) FindById(c context.Context, id int64) (*dto.AccountResponse, error) {
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

	accountResult, err := s.AccountRepository.FindById(ctx, tx, id)
	if err != nil {
		logger.WithError(err).Error("Database error")
		return nil, err
	}

	accountResponse := &dto.AccountResponse{
		AccountID: accountResult.AccountId,
		Balance:   accountResult.Balance,
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return accountResponse, nil
}
