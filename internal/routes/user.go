package routes

import (
	"github.com/julienschmidt/httprouter"
	"tigerhall_kittens/internal/handler"
	"tigerhall_kittens/internal/handler/middleware"
)

func RegisterUserRoutes(router *httprouter.Router) {
	userHandler := handler.NewUserHandler()
	router.POST("/api/v1/users", middleware.ServeV1Endpoint(middleware.AuthMiddlewareTwo, userHandler.CreateUser))
}
