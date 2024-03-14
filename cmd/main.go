package main

import (
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"tigerhall_kittens/internal/db"
	"tigerhall_kittens/internal/handler"
	"tigerhall_kittens/internal/handler/middleware"
)

func main() {
	db.ConnectAndMigrate()

	router := httprouter.New()

	tigerHandler := handler.NewTigerHandler()
	router.POST("/api/v1/tigers", middleware.LoggingMiddleware(tigerHandler.CreateTiger()))
	router.GET("/api/v1/tigers", middleware.LoggingMiddleware(tigerHandler.ListTigers()))

	log.Fatal(http.ListenAndServe(":8082", router))
}
