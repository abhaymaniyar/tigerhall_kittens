package routes

import (
	"github.com/julienschmidt/httprouter"

	"tigerhall_kittens/internal/handler"
	"tigerhall_kittens/internal/handler/middleware"
)

func RegisterTigerRoutes(router *httprouter.Router) {
	tigerHandler := handler.NewTigerHandler()
	router.POST("/api/v1/tigers", middleware.ServeV1Endpoint(middleware.AuthMiddleware, tigerHandler.CreateTiger))
	router.GET("/api/v1/tigers", middleware.ServeV1Endpoint(middleware.AuthMiddleware, tigerHandler.ListTigers))
}
