package controllers

import (
	"net/http"

	"transfer-system/adapters/transport"
	"transfer-system/domain/ports"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type UserController struct {
	UserService ports.UserService
}

// func NewUserController(service ports.UserService) *UserController {
// 	return &UserController{UserService: service}
// }

func GetPayload(ctx echo.Context, result interface{}) error {
	if err := ctx.Bind(result); err != nil {
		return err
	}
	return nil
}

type WebResponse struct {
	Message string      `json:"message"`
	Status  int         `json:"status"`
	Data    interface{} `json:"data"`
}

// func (controller *UserController) Create(ctx echo.Context) error {
// 	logger := r.Context().Value("logger").(*logrus.Entry)
// 	userRequest := transport.UserRequest{}
// 	GetPayload(r, &userRequest)

// 	userResponse, err := controller.UserService.Save(r.Context(), &userRequest)
// 	if err != nil {
// 		logger.Error("Failed to create user: ", err)
// 		WriteResponse(w, WebResponse{
// 			Message: "Failed to create user",
// 			Status:  0,
// 			Data:    nil,
// 		}, http.StatusBadRequest)
// 		return
// 	}

// 	response := WebResponse{
// 		Message: "success save user",
// 		Status:  1,
// 		Data:    userResponse,
// 	}
// 	WriteResponse(w, &response, http.StatusCreated)
// }

func (controller *UserController) Create(ctx echo.Context) error {
	logger, _ := ctx.Request().Context().Value("logger").(*logrus.Entry)
	userRequest := transport.UserRequest{}

	if err := GetPayload(ctx, &userRequest); err != nil {
		return ctx.JSON(http.StatusBadRequest, WebResponse{
			Message: "Invalid request Payload",
			Status:  0,
			Data:    nil,
		})
	}

	userResponse, err := controller.UserService.Save(ctx.Request().Context(), &userRequest)

	if err != nil {
		logger.Info("Error create controller")
		panic(err)
	}

	response := WebResponse{
		Message: "success save user",
		Status:  1,
		Data:    userResponse,
	}

	return ctx.JSON(http.StatusCreated, response)
}

// func (c *UserController) FindById(w http.ResponseWriter, r *http.Request) {
// 	logger := r.Context().Value("logger").(*logrus.Entry)
// 	userId := r.PathValue("userId")

// 	user, err := c.UserService.FindById(r.Context(), userId)
// 	if err != nil {
// 		if err.Error() == "user not found" {
// 			WriteResponse(w, WebResponse{
// 				Message: "User not found",
// 				Status:  0,
// 				Data:    nil,
// 			}, http.StatusNotFound)
// 			return
// 		}
// 		if err.Error() == "Invalid UUID Format" {
// 			WriteResponse(w, WebResponse{
// 				Message: "Invalid user ID format",
// 				Status:  0,
// 				Data:    nil,
// 			}, http.StatusBadRequest)
// 			return
// 		}
// 		logger.Error("Failed to find user: ", err)
// 		WriteResponse(w, WebResponse{
// 			Message: "Internal server error",
// 			Status:  0,
// 			Data:    nil,
// 		}, http.StatusInternalServerError)
// 		return
// 	}

// 	response := WebResponse{
// 		Message: "success get user by id",
// 		Status:  1,
// 		Data:    user,
// 	}
// 	WriteResponse(w, &response, http.StatusOK)
// }

func (c *UserController) FindById(ctx echo.Context) error {
	logger, _ := ctx.Request().Context().Value("logger").(*logrus.Entry)
	userId := ctx.Param("userId")

	user, err := c.UserService.FindById(ctx.Request().Context(), userId)

	if err != nil {
		logger.Error("Error find by id controller: ", err)

		return ctx.JSON(http.StatusNotFound, WebResponse{
			Message: "Failed get user id",
			Status:  0,
			Data:    nil,
		})
	}

	response := WebResponse{
		Message: "success get user by id",
		Status:  1,
		Data:    &user,
	}

	return ctx.JSON(http.StatusOK, response)
}
