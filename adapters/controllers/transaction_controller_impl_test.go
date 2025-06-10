package controllers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"transfer-system/adapters/controllers"
	"transfer-system/adapters/web/dto"
	"transfer-system/domain/entities"
	"transfer-system/internal/testutils"
	"transfer-system/mocks"
	appErrors "transfer-system/pkg/errors"
)

func TestTransactionController_Save_Success(t *testing.T) {
	e := echo.New()

	mockService := new(mocks.MockTransactionService)
	controller := &controllers.TransactionController{TransactionService: mockService}

	reqBody := dto.TransactionRequest{
		SourceAccountID:      123,
		DestinationAccountID: 456,
		Amount:               "100.12345",
	}
	bodyBytes, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewReader(bodyBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	testutils.InjectLoggerToContext(c)

	expectedEntity := &entities.Transaction{
		SourceAccountID:      reqBody.SourceAccountID,
		DestinationAccountID: reqBody.DestinationAccountID,
		Amount:               decimal.RequireFromString(reqBody.Amount),
	}

	mockService.On("Save", mock.Anything, expectedEntity).Return(nil)

	err := controller.Save(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)

	var resp dto.WebResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "transaction success", resp.Message)
	assert.Equal(t, int(1), resp.Status)
	assert.Nil(t, resp.Data)
}

func TestTransactionController_Save_InsufficientBalance(t *testing.T) {
	e := echo.New()

	mockService := new(mocks.MockTransactionService)
	controller := &controllers.TransactionController{TransactionService: mockService}

	reqBody := dto.TransactionRequest{
		SourceAccountID:      123,
		DestinationAccountID: 456,
		Amount:               "100.12345",
	}
	bodyBytes, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewReader(bodyBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	testutils.InjectLoggerToContext(c)

	expectedEntity := &entities.Transaction{
		SourceAccountID:      reqBody.SourceAccountID,
		DestinationAccountID: reqBody.DestinationAccountID,
		Amount:               decimal.RequireFromString(reqBody.Amount),
	}

	appErr := appErrors.NewBadRequestError("Insufficient balance", nil)
	mockService.On("Save", mock.Anything, expectedEntity).Return(appErr)

	err := controller.Save(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var resp dto.WebResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "Insufficient balance", resp.Message)
	assert.Equal(t, 0, resp.Status)
	assert.Nil(t, resp.Data)
}

func TestTransactionController_Save_InvalidAmountFormat(t *testing.T) {
	e := echo.New()
	mockService := new(mocks.MockTransactionService)
	controller := &controllers.TransactionController{TransactionService: mockService}

	reqBody := dto.TransactionRequest{
		SourceAccountID:      123,
		DestinationAccountID: 456,
		Amount:               "100",
	}
	bodyBytes, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewReader(bodyBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	testutils.InjectLoggerToContext(c)

	err := controller.Save(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var resp dto.WebResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "Invalid amount format", resp.Message)
	assert.Equal(t, 0, resp.Status)
	assert.Nil(t, resp.Data)

	mockService.AssertNotCalled(t, "Save")
}

func TestTransactionController_Save_AccountNotFound(t *testing.T) {
	e := echo.New()

	mockService := new(mocks.MockTransactionService)
	controller := &controllers.TransactionController{TransactionService: mockService}

	reqBody := dto.TransactionRequest{
		SourceAccountID:      123,
		DestinationAccountID: 456,
		Amount:               "100.12345",
	}
	bodyBytes, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewReader(bodyBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	testutils.InjectLoggerToContext(c)

	expectedEntity := &entities.Transaction{
		SourceAccountID:      reqBody.SourceAccountID,
		DestinationAccountID: reqBody.DestinationAccountID,
		Amount:               decimal.RequireFromString(reqBody.Amount),
	}

	appErr := appErrors.NewBadRequestError("Account Not Found", nil)
	mockService.On("Save", mock.Anything, expectedEntity).Return(appErr)

	err := controller.Save(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var resp dto.WebResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "Account Not Found", resp.Message)
	assert.Equal(t, 0, resp.Status)
	assert.Nil(t, resp.Data)
}
