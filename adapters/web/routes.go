package web

import (
	"net/http"

	"transfer-system/domain/ports"
)

func UserRouter(controller ports.UserController, serve *http.ServeMux) {
	serve.HandleFunc("POST /api/user", controller.Create)
	serve.HandleFunc("PUT /api/user", controller.Update)
	serve.HandleFunc("DELETE /api/user/{userId}", controller.Delete)
	serve.HandleFunc("GET /api/user/{userId}", controller.FindById)
	serve.HandleFunc("GET /api/user", controller.FindAll)
}
