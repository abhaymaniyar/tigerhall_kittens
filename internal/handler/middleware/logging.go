package middleware

import (
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"time"
)

func LoggingMiddleware(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		start := time.Now()
		log.Printf("Started %s %s", r.Method, r.URL.Path)

		next(w, r, ps)

		log.Printf("Completed in %v", time.Since(start))
	}
}
