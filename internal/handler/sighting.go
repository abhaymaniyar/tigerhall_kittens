package handler

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
	"tigerhall_kittens/internal/repository"
	"tigerhall_kittens/internal/service"
)

type SightingHandler interface {
	ReportSighting() httprouter.Handle
	GetSightings() httprouter.Handle
}

type sightingHandler struct {
	sightingService service.SightingService
}

func NewSightingHandler() SightingHandler {
	return &sightingHandler{sightingService: service.NewSightingService()}
}

// ReportSighting creates a new user
func (h *sightingHandler) ReportSighting() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		var req service.ReportSightingReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Failed to decode request body", http.StatusBadRequest)
			return
		}

		err := h.sightingService.ReportSighting(r.Context(), req)
		if err != nil {
			// TODO: dont generalize the errors to be 400 here
			http.Error(w, fmt.Sprintf("Error while creating req: %s", err.Error()), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(req)
	}
}

func (h *sightingHandler) GetSightings() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		tigerIDStr := params.ByName("tiger_id")

		// TODO: validation
		tigerID, _ := strconv.ParseUint(tigerIDStr, 10, 0)

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

		sightings, err := h.sightingService.GetSightings(repository.GetSightingOpts{
			TigerID: uint(tigerID),
			Limit:   perPage,
			Offset:  offset,
		})

		if err != nil {
			// TODO: dont generalize the errors to be 400 here
			http.Error(w, fmt.Sprintf("Error while creating req: %s", err.Error()), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(sightings)
	}
}
