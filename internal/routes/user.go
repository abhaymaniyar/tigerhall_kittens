package routes

import (
	"github.com/julienschmidt/httprouter"
	"tigerhall_kittens/internal/handler"
	"tigerhall_kittens/internal/handler/middleware"
)

func RegisterUserRoutes(router *httprouter.Router) {
	userHandler := handler.NewUserHandler()
	router.POST("/api/v1/users", middleware.AuthMiddleware(middleware.LoggingMiddleware(userHandler.CreateUser())))

	authHandler := handler.NewAuthHandler()
	router.POST("/api/v1/login", middleware.AuthMiddleware(middleware.LoggingMiddleware(authHandler.Login())))
}
