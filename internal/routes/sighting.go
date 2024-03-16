package routes

import (
	"github.com/julienschmidt/httprouter"

	"tigerhall_kittens/internal/handler"
	"tigerhall_kittens/internal/handler/middleware"
)

func RegisterSightingRoutes(router *httprouter.Router) {
	sightingHandler := handler.NewSightingHandler()
	router.POST("/api/v1/sightings", middleware.ServeV1Endpoint(middleware.AuthMiddleware, sightingHandler.ReportSighting))
	router.GET("/api/v1/tigers/:tiger_id/sightings", middleware.ServeV1Endpoint(middleware.AuthMiddleware, sightingHandler.GetSightings))
}
