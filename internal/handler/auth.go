package handler

import (
	"encoding/json"
	"fmt"
	"tigerhall_kittens/internal/service"
	"tigerhall_kittens/internal/web"
	"tigerhall_kittens/utils"
)

type AuthHandler interface {
	Login(r *web.Request) (*web.JSONResponse, web.ErrorInterface)
}

type authHandler struct {
	authService service.AuthService
}

func NewAuthHandler() AuthHandler {
	return &authHandler{authService: service.NewAuthService()}
}

// Login creates a new user
func (h *authHandler) Login(r *web.Request) (*web.JSONResponse, web.ErrorInterface) {
	var user service.LoginUserReq
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		return nil, web.ErrBadRequest("Failed to decode request body")
	}

	loginResp, err := h.authService.LoginUser(r.Context(), user)
	if err != nil {
		// TODO: dont generalize the errors to be 400 here
		return nil, web.ErrInternalServerError(fmt.Sprintf("Error while logging in user : %s", err.Error()))
	}

	jsonResponse, err := utils.StructToMap(loginResp)
	if err != nil {
		return nil, web.ErrInternalServerError(err.Error())
	}

	return (*web.JSONResponse)(&jsonResponse), nil
}
