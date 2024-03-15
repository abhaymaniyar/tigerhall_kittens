package handler

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"tigerhall_kittens/internal/service"
)

type UserHandler interface {
	CreateUser() httprouter.Handle
}

type userHandler struct {
	userService service.UserService
}

func NewUserHandler() UserHandler {
	return &userHandler{userService: service.NewUserService()}
}

// CreateUser creates a new user
func (h *userHandler) CreateUser() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		var user service.CreateUserReq
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, "Failed to decode request body", http.StatusBadRequest)
			return
		}

		err := h.userService.CreateUser(user)
		if err != nil {
			// TODO: dont generalize the errors to be 400 here
			http.Error(w, fmt.Sprintf("Error while creating user: %s", err.Error()), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(user)
	}
}
