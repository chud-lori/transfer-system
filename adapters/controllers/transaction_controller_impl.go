package controllers

import (
	"errors"
	"net/http"

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

type TransactionController struct {
	TransactionService ports.TransactionService
}

// Save Transaction godoc
// @Summary      Create Transaction
// @Description  Transfer amount from source account to destination account
// @Tags         Transactions
// @Accept       json
// @Produce      json
// @Param        body  body      dto.TransactionRequest  true  "Transaction payload"  example({"source_account_id":1,"destination_account_id":2,"amount":"100.00"})
// @Success      201   {object}  dto.WebResponse
// @Failure      400   {object}  dto.WebResponse
// @Failure      500   {object}  dto.WebResponse
// @Router       /transactions [post]
func (c *TransactionController) Save(ctx echo.Context) error {
	logger, _ := ctx.Request().Context().Value(logger.LoggerContextKey).(*logrus.Entry)
	transactionRequest := dto.TransactionRequest{}

	if err := web.GetPayload(ctx, &transactionRequest); err != nil {
		return ctx.JSON(http.StatusBadRequest, dto.WebResponse{
			Message: "Invalid request Payload",
			Status:  0,
			Data:    nil,
		})
	}

	if !validator.ValidateDecimalFormat(transactionRequest.Amount) {
		logger.Errorf("Invalid amount format: %s", transactionRequest.Amount)
		return ctx.JSON(http.StatusBadRequest, dto.WebResponse{
			Message: "Invalid amount format",
			Status:  0,
			Data:    nil,
		})
	}

	amountDecimal, err := decimal.NewFromString(transactionRequest.Amount)
	if err != nil {
		logger.WithError(err).Error("Failed to parse amount")
		return ctx.JSON(http.StatusInternalServerError, dto.WebResponse{
			Message: "An internal error occurred while processing balance.",
			Status:  0,
			Data:    nil,
		})
	}

	internalServiceRequest := &entities.Transaction{
		SourceAccountID:      transactionRequest.SourceAccountID,
		DestinationAccountID: transactionRequest.DestinationAccountID,
		Amount:               amountDecimal,
	}

	err = c.TransactionService.Save(ctx.Request().Context(), internalServiceRequest)

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
		Message: "transaction success",
		Status:  1,
		Data:    nil,
	}

	return ctx.JSON(http.StatusCreated, response)
}
