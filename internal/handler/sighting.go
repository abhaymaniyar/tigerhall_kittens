package handler

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"tigerhall_kittens/internal/service"
)

type SightingHandler interface {
	ReportSighting() httprouter.Handle
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

		err := h.sightingService.ReportSighting(req)
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
