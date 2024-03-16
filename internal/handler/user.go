package handler

import (
	"encoding/json"
	"tigerhall_kittens/internal/service"
	"tigerhall_kittens/internal/web"
)

type UserHandler interface {
	CreateUser(r *web.Request) (*web.JSONResponse, web.ErrorInterface)
}

type userHandler struct {
	userService service.UserService
}

func NewUserHandler() UserHandler {
	return &userHandler{userService: service.NewUserService()}
}

func MakeUserHandler(userService service.UserService) UserHandler {
	return &userHandler{userService: userService}
}

// CreateUser creates a new user
func (h *userHandler) CreateUser(r *web.Request) (*web.JSONResponse, web.ErrorInterface) {
	var user service.CreateUserReq
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		return nil, web.ErrBadRequest("Failed to decode request body")
	}

	err := h.userService.CreateUser(r.Context(), &user)
	if err != nil {
		return nil, errorResponse(err)
	}

	return &web.JSONResponse{}, nil
}
