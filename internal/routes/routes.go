package routes

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func Init(router *httprouter.Router) {
	router.GET("/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.Header().Add("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, "{ \"message\":\"Hello world!. I am tigerhall_kittens.\",\"success\":true,\"api_version\": 1 }")
	})

	RegisterAuthRoutes(router)
	RegisterUserRoutes(router)
	RegisterTigerRoutes(router)
	RegisterSightingRoutes(router)
}
