package controllers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
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
)

func TestAccountController_Create_Success(t *testing.T) {
	e := echo.New()
	mockService := new(mocks.MockAccountService)
	controller := &controllers.AccountController{AccountService: mockService}

	reqBody := dto.AccountRequest{
		AccountID: 12345,
		Balance:   "100.23344",
	}
	bodyBytes, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/accounts", bytes.NewReader(bodyBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	testutils.InjectLoggerToContext(c)

	acc := &entities.Account{
		AccountID: 12345,
		Balance:   decimal.NewFromFloat(100.23344),
	}

	mockService.On("Save", mock.Anything, acc).Return(nil)

	err := controller.Create(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)

	var response dto.WebResponse
	json.Unmarshal(rec.Body.Bytes(), &response)
	assert.Equal(t, "success create account", response.Message)
	assert.Equal(t, int(1), response.Status)
}

func TestAccountController_Create_InvalidBalanceFormat(t *testing.T) {
	e := echo.New()
	controller := &controllers.AccountController{}
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

func TestAccountController_FindById_Success(t *testing.T) {
	e := echo.New()
	mockService := new(mocks.MockAccountService)
	controller := &controllers.AccountController{AccountService: mockService}

	accId := int64(12345)
	expectedResp := &entities.Account{
		AccountID: accId,
		Balance:   decimal.NewFromFloat(100.23344),
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

func TestAccountController_FindById_NotFound(t *testing.T) {
	e := echo.New()
	mockService := new(mocks.MockAccountService)
	controller := &controllers.AccountController{AccountService: mockService}

	accId := int64(99999)
	mockService.On("FindById", mock.Anything, accId).Return(&entities.Account{}, errors.New("Account not found"))

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
