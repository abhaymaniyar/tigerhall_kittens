package routes

import (
	"github.com/julienschmidt/httprouter"

	"tigerhall_kittens/internal/handler"
	"tigerhall_kittens/internal/handler/middleware"
)

func RegisterAuthRoutes(router *httprouter.Router) {
	authHandler := handler.NewAuthHandler()
	router.POST("/api/v1/auth/login", middleware.ServeV1Endpoint(middleware.EmptyMiddleware, authHandler.Login))
}
