package cmd

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"tigerhall_kittens/internal/handler"
)

func main() {
	router := httprouter.New()

	tigerHandler := handler.NewTigerHandler()
	router.POST("/api/v1/tigers", tigerHandler.CreateTiger())
	router.GET("/api/v1/tigers", tigerHandler.ListTigers())

	http.ListenAndServe(":8080", router)
}
