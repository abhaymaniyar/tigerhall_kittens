package handler

import (
	"encoding/json"
	"fmt"
	"strconv"
	"tigerhall_kittens/internal/repository"
	"tigerhall_kittens/internal/service"
	"tigerhall_kittens/internal/web"
)

type SightingHandler interface {
	ReportSighting(req *web.Request) (*web.JSONResponse, web.ErrorInterface)
	GetSightings(req *web.Request) (*web.JSONResponse, web.ErrorInterface)
}

type sightingHandler struct {
	sightingService service.SightingService
}

func NewSightingHandler() SightingHandler {
	return &sightingHandler{sightingService: service.NewSightingService()}
}

func MakeSightingHandler(sightingService service.SightingService) SightingHandler {
	return &sightingHandler{sightingService: sightingService}
}

// ReportSighting creates a new user
func (h *sightingHandler) ReportSighting(r *web.Request) (*web.JSONResponse, web.ErrorInterface) {
	var req service.ReportSightingReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, web.ErrBadRequest("Failed to decode request body")
	}

	err := h.sightingService.ReportSighting(r.Context(), req)
	if err != nil {
		//TODO: fix this to return 400 in case of already existing sighting
		return nil, web.ErrInternalServerError(fmt.Sprintf("Error while reporting sighting : %s", err.Error()))
	}

	return &web.JSONResponse{}, nil
}

func (h *sightingHandler) GetSightings(r *web.Request) (*web.JSONResponse, web.ErrorInterface) {
	tigerIDStr := r.GetPathParam("tiger_id")

	// TODO: validation
	tigerID, _ := strconv.ParseUint(tigerIDStr, 10, 0)

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

	sightings, err := h.sightingService.GetSightings(r.Context(), repository.GetSightingOpts{
		TigerID: uint(tigerID),
		Limit:   perPage,
		Offset:  offset,
	})

	if err != nil {
		return nil, web.ErrInternalServerError(fmt.Sprintf("Error while fetching sightings : %s", err.Error()))
	}

	res := map[string]interface{}{
		"tigers":   sightings,
		"page":     page,
		"per_page": perPage,
	}

	return (*web.JSONResponse)(&res), nil
}
