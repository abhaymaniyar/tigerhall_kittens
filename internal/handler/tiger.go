package handler

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"tigerhall_kittens/internal/model"
	"tigerhall_kittens/internal/service"
)

type TigerHandler interface {
	CreateTiger() httprouter.Handle
	ListTigers() httprouter.Handle
}

type tigerHandler struct {
	tigerService service.TigerService
}

func NewTigerHandler() TigerHandler {
	return &tigerHandler{tigerService: service.NewTigerService()}
}

func (t *tigerHandler) CreateTiger() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		var tiger model.Tiger
		if err := json.NewDecoder(r.Body).Decode(&tiger); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err := t.tigerService.CreateTiger(&tiger)
		if err != nil {
			//	TODO: report error and make it more descriptive
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(tiger)
	}
}

func (t *tigerHandler) ListTigers() httprouter.Handle {
	return func(w http.ResponseWriter, request *http.Request, params httprouter.Params) {
		// TODO: add pagination support
		tigers, err := t.tigerService.ListTigers()
		if err != nil {
			//	TODO: report error and make it more descriptive
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(tigers)
	}
}
