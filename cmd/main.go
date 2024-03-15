package main

import (
	"context"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"runtime/debug"
	"tigerhall_kittens/cmd/notification_worker"
	"tigerhall_kittens/internal/config"
	"tigerhall_kittens/internal/db"
	"tigerhall_kittens/internal/logger"
	"tigerhall_kittens/internal/routes"
	"tigerhall_kittens/internal/web"
)

func main() {
	defer logger.Sync()
	defer db.Close()

	// TODO: improve this
	ctx, cancel := context.WithCancel(context.Background())
	defer handlePanic(ctx, cancel)

	if err := config.LoadEnv(); err != nil {
		panic(err)
	}

	config.SetupLogger(config.Env.Environment)

	config.SetupDBConnection(context.Background())

	router := httprouter.New()
	routes.Init(router)

	go notification_worker.SetupNotificationWorker()

	log.Fatal(http.ListenAndServe(config.Env.Port, router))
}

func handlePanic(ctx context.Context, cancel context.CancelFunc) {
	if recvr := recover(); recvr != nil {
		errorMessage := fmt.Sprintf("%v", recvr)
		err := web.ErrInternalServerError(errorMessage)
		logger.E(ctx, err, "panic",
			logger.Field("status", err.HTTPStatusCode()),
			logger.Field("stack", string(debug.Stack())),
		)
	}

	cancel()
}
