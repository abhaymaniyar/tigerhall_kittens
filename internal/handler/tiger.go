package handler

import (
	"encoding/json"
	"fmt"
	"strconv"
	"tigerhall_kittens/internal/model"
	"tigerhall_kittens/internal/repository"
	"tigerhall_kittens/internal/service"
	"tigerhall_kittens/internal/web"
	"tigerhall_kittens/utils"
)

type TigerHandler interface {
	CreateTiger(req *web.Request) (*web.JSONResponse, web.ErrorInterface)
	ListTigers(r *web.Request) (*web.JSONResponse, web.ErrorInterface)
}

type tigerHandler struct {
	tigerService service.TigerService
}

func NewTigerHandler() TigerHandler {
	return &tigerHandler{tigerService: service.NewTigerService()}
}

func (t *tigerHandler) CreateTiger(r *web.Request) (*web.JSONResponse, web.ErrorInterface) {
	var tiger model.Tiger
	if err := json.NewDecoder(r.Body).Decode(&tiger); err != nil {
		return nil, web.ErrBadRequest("Failed to decode request body")
	}

	err := t.tigerService.CreateTiger(r.Context(), &tiger)
	if err != nil {
		return nil, web.ErrInternalServerError(fmt.Sprintf("error while saving tiger : %s", err))
	}

	jsonResponse, err := utils.StructToMap(tiger)
	if err != nil {
		return nil, web.ErrInternalServerError(err.Error())
	}

	return (*web.JSONResponse)(&jsonResponse), nil
}

func (t *tigerHandler) ListTigers(r *web.Request) (*web.JSONResponse, web.ErrorInterface) {
	pageStr := r.URL.Query().Get("page")
	perPageStr := r.URL.Query().Get("per_page")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		return nil, web.ErrBadRequest("Invalid page number")
	}

	perPage, err := strconv.Atoi(perPageStr)
	if err != nil || perPage <= 0 {
		return nil, web.ErrBadRequest("Invalid per_page value")
	}

	offset := (page - 1) * perPage

	tigers, err := t.tigerService.ListTigers(r.Context(), repository.ListTigersOpts{Limit: perPage, Offset: offset})
	if err != nil {
		return nil, web.ErrInternalServerError(fmt.Sprintf("Error while fetching tigers : %s", err.Error()))
	}

	res := map[string]interface{}{
		"tigers":   tigers,
		"page":     page,
		"per_page": perPage,
	}

	return (*web.JSONResponse)(&res), nil
}
