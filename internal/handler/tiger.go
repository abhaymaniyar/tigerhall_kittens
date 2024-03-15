package handler

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
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

		err := t.tigerService.CreateTiger(r.Context(), &tiger)
		if err != nil {
			//	TODO: report error and make it more descriptive
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(tiger)
	}
}

func (t *tigerHandler) ListTigers() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		// TODO: add pagination support
		pageStr := r.URL.Query().Get("page")
		perPageStr := r.URL.Query().Get("per_page")

		page, err := strconv.Atoi(pageStr)
		if err != nil || page <= 0 {
			http.Error(w, "Invalid page number", http.StatusBadRequest)
			return
		}

		perPage, err := strconv.Atoi(perPageStr)
		if err != nil || perPage <= 0 {
			http.Error(w, "Invalid per_page value", http.StatusBadRequest)
			return
		}

		offset := (page - 1) * perPage

		tigers, err := t.tigerService.ListTigers(r.Context(), service.ListTigersOpts{Limit: perPage, Offset: offset})
		if err != nil {
			//	TODO: report error and make it more descriptive
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(tigers)
	}
}
