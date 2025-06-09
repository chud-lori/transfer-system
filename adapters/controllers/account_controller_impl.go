package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"transfer-system/adapters/web"
	"transfer-system/adapters/web/dto"
	"transfer-system/domain/entities"
	"transfer-system/domain/ports"
	appErrors "transfer-system/pkg/errors"
	"transfer-system/pkg/logger"
	"transfer-system/pkg/validator"

	"github.com/labstack/echo/v4"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

type AccountController struct {
	AccountService ports.AccountService
}

func (controller *AccountController) Create(ctx echo.Context) error {
	logger, _ := ctx.Request().Context().Value(logger.LoggerContextKey).(*logrus.Entry)
	accountRequest := dto.AccountRequest{}

	if err := web.GetPayload(ctx, &accountRequest); err != nil {
		return ctx.JSON(http.StatusBadRequest, dto.WebResponse{
			Message: "Invalid request Payload",
			Status:  0,
			Data:    nil,
		})
	}

	if !validator.ValidateDecimalFormat(accountRequest.Balance) {
		logger.Errorf("Invalid initial balance format: %s", accountRequest.Balance)
		return ctx.JSON(http.StatusBadRequest, dto.WebResponse{
			Message: "Invalid initial balance format",
			Status:  0,
			Data:    nil,
		})
	}

	initialBalanceDecimal, err := decimal.NewFromString(accountRequest.Balance)
	if err != nil {
		logger.WithError(err).Error("Failed to parse initial balance")
		return ctx.JSON(http.StatusInternalServerError, dto.WebResponse{
			Message: "An internal error occurred while processing balance.",
			Status:  0,
			Data:    nil,
		})
	}

	internalServiceRequest := &entities.Account{
		AccountID: accountRequest.AccountID,
		Balance:   initialBalanceDecimal,
	}

	accountResponse, err := controller.AccountService.Save(ctx.Request().Context(), internalServiceRequest)

	if err != nil {
		var appErr *appErrors.AppError
		if errors.As(err, &appErr) {
			return ctx.JSON(appErr.StatusCode, dto.WebResponse{
				Message: appErr.Message,
				Status:  0,
				Data:    nil,
			})
		} else {
			return ctx.JSON(http.StatusInternalServerError, dto.WebResponse{
				Message: "An unexpected error occurred",
				Status:  0,
				Data:    nil,
			})
		}
	}

	response := dto.WebResponse{
		Message: "success create account",
		Status:  1,
		Data:    accountResponse,
	}

	return ctx.JSON(http.StatusCreated, response)
}

func (c *AccountController) FindById(ctx echo.Context) error {
	logger, _ := ctx.Request().Context().Value("logger").(*logrus.Entry)
	accountIdStr := ctx.Param("accountId")

	accountId, err := strconv.ParseInt(accountIdStr, 10, 64)
	if err != nil {
		logger.WithError(err).Errorf("Invalid accountId parameter: %s", accountIdStr)
		return ctx.JSON(http.StatusBadRequest, dto.WebResponse{
			Message: "Invalid accountId format. Please provide a valid number.",
			Status:  0,
			Data:    nil,
		})
	}

	account, err := c.AccountService.FindById(ctx.Request().Context(), accountId)

	if err != nil {
		logger.Error("Error find by id controller: ", err)

		return ctx.JSON(http.StatusNotFound, dto.WebResponse{
			Message: "Account id not found",
			Status:  0,
			Data:    nil,
		})
	}

	response := dto.WebResponse{
		Message: "success get account by id",
		Status:  1,
		Data:    &account,
	}

	return ctx.JSON(http.StatusOK, response)
}
