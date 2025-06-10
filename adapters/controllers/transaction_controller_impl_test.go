package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"transfer-system/adapters/web/dto"
	"transfer-system/domain/entities"
	"transfer-system/internal/testutils"
)

type MockTransactionService struct {
	mock.Mock
}

func (m *MockTransactionService) Save(ctx context.Context, acc *entities.Account) error {
	args := m.Called(ctx, acc)
	return args.Error(1)
}

func TestTransactionController_Save_Success(t *testing.T) {
	e := echo.New()
	mockService := new(MockAccountService)
	controller := &AccountController{AccountService: mockService}

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

	expectedResp := &dto.WebResponse{
		Message: "success create account",
		Status:  1,
		Data: &dto.AccountResponse{
			AccountID: 123,
			Balance:   "100.23344",
		},
	}
	mockService.On("Save", mock.Anything, mock.Anything).Return(expectedResp, nil)

	err := controller.Create(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)

	var response dto.WebResponse
	json.Unmarshal(rec.Body.Bytes(), &response)
	assert.Equal(t, "success create account", response.Message)
	assert.Equal(t, int(1), response.Status)
}

func TestTransactionController_Create_InvalidBalanceFormat(t *testing.T) {
	e := echo.New()
	controller := &AccountController{}
	reqBody := dto.AccountRequest{
		AccountID: 12345,
		Balance:   "abc", // invalid
	}
	bodyBytes, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/accounts", bytes.NewReader(bodyBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	testutils.InjectLoggerToContext(c)

	err := controller.Create(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestTransactionController_FindById_Success(t *testing.T) {
	e := echo.New()
	mockService := new(MockAccountService)
	controller := &AccountController{AccountService: mockService}

	accId := int64(12345)
	expectedResp := &dto.AccountResponse{
		AccountID: accId,
		Balance:   "100.23344",
	}
	mockService.On("FindById", mock.Anything, accId).Return(expectedResp, nil)

	req := httptest.NewRequest(http.MethodGet, "/accounts/12345", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("accountId")
	c.SetParamValues(strconv.FormatInt(accId, 10))
	testutils.InjectLoggerToContext(c)

	err := controller.FindById(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response dto.WebResponse
	json.Unmarshal(rec.Body.Bytes(), &response)
	assert.Equal(t, "success get account by id", response.Message)
	assert.Equal(t, int(1), response.Status)
}

func TestTransactionController_FindById_NotFound(t *testing.T) {
	e := echo.New()
	mockService := new(MockAccountService)
	controller := &AccountController{AccountService: mockService}

	accId := int64(99999)
	mockService.On("FindById", mock.Anything, accId).Return(&dto.AccountResponse{}, errors.New("not found"))

	req := httptest.NewRequest(http.MethodGet, "/accounts/99999", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("accountId")
	c.SetParamValues(strconv.FormatInt(accId, 10))
	testutils.InjectLoggerToContext(c)

	err := controller.FindById(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}
