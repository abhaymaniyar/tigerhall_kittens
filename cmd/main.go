package main

import (
	"context"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"tigerhall_kittens/internal/config"
	"tigerhall_kittens/internal/db"
	"tigerhall_kittens/internal/logger"
	"tigerhall_kittens/internal/routes"
)

func main() {
	defer logger.Sync()
	defer db.Close()

	if err := config.LoadEnv(); err != nil {
		panic(err)
	}

	// TODO: fix logger
	config.SetupLogger(config.Env.Environment)

	// TODO: create a db init
	config.SetupDBConnection(context.Background())

	// TODO: create a router file
	router := httprouter.New()
	routes.Init(router)

	log.Fatal(http.ListenAndServe(":8082", router))
}
