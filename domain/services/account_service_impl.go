package services

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"transfer-system/adapters/web/dto"
	"transfer-system/domain/entities"
	"transfer-system/domain/ports"
)

type AccountServiceImpl struct {
	DB                ports.Database
	AccountRepository ports.AccountRepository
	CtxTimeout        time.Duration
}

func (s *AccountServiceImpl) Save(c context.Context, request *dto.CreateAccountRequest) (*dto.WebResponse, error) {
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
			return nil, err
		}
	} else {
		return nil, errors.New("Id already exists")
	}

	account := entities.Account{
		AccountId: request.AccountID,
		Balance:   request.InitialBalance,
	}
	_, err = s.AccountRepository.Save(ctx, tx, &account)

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

func (s *AccountServiceImpl) FindById(c context.Context, id int64) (*dto.AccountResponse, error) {
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
