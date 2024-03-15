package main

import (
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"tigerhall_kittens/internal/db"
	"tigerhall_kittens/internal/handler"
	"tigerhall_kittens/internal/handler/middleware"
	"tigerhall_kittens/internal/logger"
)

func main() {
	// TODO: fix logger
	logger.SetupLogger("")

	// TODO: create a db init
	db.ConnectAndMigrate()

	// TODO: create a router file
	router := httprouter.New()

	tigerHandler := handler.NewTigerHandler()
	router.POST("/api/v1/tigers", middleware.AuthMiddleware(middleware.LoggingMiddleware(tigerHandler.CreateTiger())))
	router.GET("/api/v1/tigers", middleware.AuthMiddleware(middleware.LoggingMiddleware(tigerHandler.ListTigers())))

	userHandler := handler.NewUserHandler()
	router.POST("/api/v1/users", middleware.AuthMiddleware(middleware.LoggingMiddleware(userHandler.CreateUser())))

	authHandler := handler.NewAuthHandler()
	router.POST("/api/v1/login", middleware.AuthMiddleware(middleware.LoggingMiddleware(authHandler.Login())))

	sightingHandler := handler.NewSightingHandler()
	router.POST("/api/v1/sightings", middleware.AuthMiddleware(middleware.LoggingMiddleware(sightingHandler.ReportSighting())))
	router.GET("/api/v1/tigers/:tiger_id/sightings", middleware.AuthMiddleware(middleware.LoggingMiddleware(sightingHandler.GetSightings())))

	log.Fatal(http.ListenAndServe(":8082", router))
}
