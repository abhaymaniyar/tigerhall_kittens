package middleware

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func AuthMiddleware(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		// TODO: uncomment this later
		//token := r.Header.Get("Authorization")
		//if token != "Bearer example-token" {
		//	http.Error(w, "Unauthorized", http.StatusUnauthorized)
		//	return
		//}

		next(w, r, ps)
	}
}
