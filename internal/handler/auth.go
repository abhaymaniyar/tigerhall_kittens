package handler

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"tigerhall_kittens/internal/service"
)

type AuthHandler interface {
	Login() httprouter.Handle
}

type authHandler struct {
	authService service.AuthService
}

func NewAuthHandler() AuthHandler {
	return &authHandler{authService: service.NewAuthService()}
}

// Login creates a new user
func (h *authHandler) Login() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		var user service.LoginUserReq
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, "Failed to decode request body", http.StatusBadRequest)
			return
		}

		loginResp, err := h.authService.LoginUser(user)
		if err != nil {
			// TODO: dont generalize the errors to be 400 here
			http.Error(w, fmt.Sprintf("Error while creating user: %s", err.Error()), http.StatusUnauthorized)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(loginResp)
	}
}
