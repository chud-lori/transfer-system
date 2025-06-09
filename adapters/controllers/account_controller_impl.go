package controllers

import (
	"net/http"
	"strconv"

	"transfer-system/adapters/web"
	"transfer-system/adapters/web/dto"
	"transfer-system/domain/ports"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type AccountController struct {
	AccountService ports.AccountService
}

func (controller *AccountController) Create(ctx echo.Context) error {
	logger, _ := ctx.Request().Context().Value("logger").(*logrus.Entry)
	accountRequest := dto.CreateAccountRequest{}

	if err := web.GetPayload(ctx, &accountRequest); err != nil {
		return ctx.JSON(http.StatusBadRequest, dto.WebResponse{
			Message: "Invalid request Payload",
			Status:  0,
			Data:    nil,
		})
	}

	accountResponse, err := controller.AccountService.Save(ctx.Request().Context(), &accountRequest)

	if err != nil {
		logger.Info("Error create controller")
		panic(err)
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
			Message: "Failed get account id",
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
