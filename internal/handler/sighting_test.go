package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"tigerhall_kittens/internal/handler/middleware"
	"tigerhall_kittens/internal/model"
	"tigerhall_kittens/internal/repository"
	"tigerhall_kittens/internal/service"
	mock_service "tigerhall_kittens/internal/service/mocks"
	"time"
)

func TestSightingHandler_GetSightings(t *testing.T) {
	tigerID := 1
	page := 1
	perPage := 10
	offset := (page - 1) * perPage

	t.Run("should return bad request if page is not passed in query param", func(t *testing.T) {
		path := "/api/v1/tigers/%v/sightings"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		recorder := httptest.NewRecorder()
		router := httprouter.New()

		path = fmt.Sprintf(path, tigerID)

		mockSightingService := mock_service.NewMockSightingService(ctrl)
		sightingHandler := MakeSightingHandler(mockSightingService)

		body := make(map[string]interface{})
		data, _ := json.Marshal(body)
		req, _ := http.NewRequest(http.MethodGet, path, bytes.NewBuffer(data))

		router.Handle(http.MethodGet, "/api/v1/tigers/:tiger_id/sightings", middleware.ServeV1Endpoint(middleware.EmptyMiddleware,
			sightingHandler.GetSightings))
		router.ServeHTTP(recorder, req)

		respCode := recorder.Code
		assert.Equal(t, http.StatusBadRequest, respCode)
		respBody, _ := ioutil.ReadAll(recorder.Body)
		var resData map[string]interface{}
		_ = json.Unmarshal(respBody, &resData)
		assert.Equal(t, "Invalid page number", resData["error"].(map[string]interface{})["message"])
		assert.Equal(t, false, resData["success"])
	})

	t.Run("should return bad request if per_page is not passed in query param", func(t *testing.T) {
		path := "/api/v1/tigers/%v/sightings?page=%v"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		recorder := httptest.NewRecorder()
		router := httprouter.New()

		path = fmt.Sprintf(path, tigerID, page)

		mockSightingService := mock_service.NewMockSightingService(ctrl)
		sightingHandler := MakeSightingHandler(mockSightingService)

		body := make(map[string]interface{})
		data, _ := json.Marshal(body)
		req, _ := http.NewRequest(http.MethodGet, path, bytes.NewBuffer(data))

		router.Handle(http.MethodGet, "/api/v1/tigers/:tiger_id/sightings", middleware.ServeV1Endpoint(middleware.EmptyMiddleware,
			sightingHandler.GetSightings))
		router.ServeHTTP(recorder, req)

		respCode := recorder.Code
		assert.Equal(t, http.StatusBadRequest, respCode)
		respBody, _ := ioutil.ReadAll(recorder.Body)
		var resData map[string]interface{}
		_ = json.Unmarshal(respBody, &resData)
		assert.Equal(t, "Invalid per_page value", resData["error"].(map[string]interface{})["message"])
		assert.Equal(t, false, resData["success"])
	})

	t.Run("should return ISE if service returns error", func(t *testing.T) {
		path := "/api/v1/tigers/%v/sightings?page=%v&per_page=%v"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		recorder := httptest.NewRecorder()
		router := httprouter.New()

		path = fmt.Sprintf(path, tigerID, page, perPage)

		mockSightingService := mock_service.NewMockSightingService(ctrl)
		mockSightingService.EXPECT().GetSightings(gomock.Any(), repository.GetSightingOpts{
			TigerID: uint(tigerID),
			Limit:   perPage,
			Offset:  offset,
		}).Return(nil, errors.New("some error"))

		sightingHandler := MakeSightingHandler(mockSightingService)

		body := make(map[string]interface{})
		data, _ := json.Marshal(body)
		req, _ := http.NewRequest(http.MethodGet, path, bytes.NewBuffer(data))

		router.Handle(http.MethodGet, "/api/v1/tigers/:tiger_id/sightings", middleware.ServeV1Endpoint(middleware.EmptyMiddleware,
			sightingHandler.GetSightings))
		router.ServeHTTP(recorder, req)

		respCode := recorder.Code
		assert.Equal(t, http.StatusInternalServerError, respCode)
		respBody, _ := ioutil.ReadAll(recorder.Body)
		var resData map[string]interface{}
		_ = json.Unmarshal(respBody, &resData)
		assert.Equal(t, "Error while fetching sightings : some error", resData["error"].(map[string]interface{})["message"])
		assert.Equal(t, false, resData["success"])
	})

	t.Run("should return success response", func(t *testing.T) {
		path := "/api/v1/tigers/%v/sightings?page=%v&per_page=%v"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		recorder := httptest.NewRecorder()
		router := httprouter.New()

		path = fmt.Sprintf(path, tigerID, page, perPage)

		mockSightings := []model.Sighting{
			{
				TigerID:   uint(tigerID),
				Lat:       1.2,
				Lon:       2.2,
				Timestamp: time.Now(),
			},
			{
				TigerID:   uint(tigerID),
				Lat:       2.2,
				Lon:       4.2,
				Timestamp: time.Now().Add(4 * time.Hour),
			},
		}

		mockSightingService := mock_service.NewMockSightingService(ctrl)
		mockSightingService.EXPECT().GetSightings(gomock.Any(), repository.GetSightingOpts{
			TigerID: uint(tigerID),
			Limit:   perPage,
			Offset:  offset,
		}).Return(mockSightings, nil)

		sightingHandler := MakeSightingHandler(mockSightingService)

		body := make(map[string]interface{})
		data, _ := json.Marshal(body)
		req, _ := http.NewRequest(http.MethodGet, path, bytes.NewBuffer(data))

		router.Handle(http.MethodGet, "/api/v1/tigers/:tiger_id/sightings", middleware.ServeV1Endpoint(middleware.EmptyMiddleware,
			sightingHandler.GetSightings))
		router.ServeHTTP(recorder, req)

		respCode := recorder.Code
		assert.Equal(t, http.StatusOK, respCode)
		respBody, _ := ioutil.ReadAll(recorder.Body)
		var resData map[string]interface{}
		_ = json.Unmarshal(respBody, &resData)
		assert.Equal(t, true, resData["success"])
	})
}

func TestSightingHandler_ReportSighting(t *testing.T) {
	var req service.ReportSightingReq

	reportSightingRequestBody := `
			{
				"username": "user_eighty212",
				"password": "user_eighty212"
			}
	`

	invalidReportSightingRequestBody := `
			{
				"username": "user_eighty212",
				"password": "user_eighty212"
	`

	path := "/api/v1/sightings"

	t.Run("should return bad request when request decode fails", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		router := httprouter.New()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockSightingService := mock_service.NewMockSightingService(ctrl)
		authHandler := MakeSightingHandler(mockSightingService)

		req, _ := http.NewRequestWithContext(context.Background(), http.MethodPost,
			path, bytes.NewBuffer([]byte(invalidReportSightingRequestBody)))

		router.Handle(http.MethodPost, path, middleware.ServeV1Endpoint(middleware.EmptyMiddleware,
			authHandler.ReportSighting))
		router.ServeHTTP(recorder, req)

		respCode := recorder.Code
		assert.Equal(t, http.StatusBadRequest, respCode)
		respBody, _ := ioutil.ReadAll(recorder.Body)
		var resData map[string]interface{}
		_ = json.Unmarshal(respBody, &resData)
		assert.Equal(t, resData["error"].(map[string]interface{})["message"], "Failed to decode request body")
		assert.Equal(t, false, resData["success"])
	})

	t.Run("should return ISE when reporting fails", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		router := httprouter.New()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockSightingService := mock_service.NewMockSightingService(ctrl)
		authHandler := MakeSightingHandler(mockSightingService)

		err := json.Unmarshal([]byte(reportSightingRequestBody), &req)
		if err != nil {
			panic(err)
		}

		req, _ := http.NewRequestWithContext(context.Background(), http.MethodPost,
			path, bytes.NewBuffer([]byte(reportSightingRequestBody)))

		var reportSightingReq service.ReportSightingReq
		err = json.Unmarshal([]byte(reportSightingRequestBody), &reportSightingReq)
		assert.Nil(t, err)
		mockSightingService.EXPECT().ReportSighting(gomock.Any(), reportSightingReq).Return(errors.New("some error"))

		router.Handle(http.MethodPost, path, middleware.ServeV1Endpoint(middleware.EmptyMiddleware,
			authHandler.ReportSighting))
		router.ServeHTTP(recorder, req)

		respCode := recorder.Code
		assert.Equal(t, respCode, http.StatusInternalServerError)
		respBody, _ := ioutil.ReadAll(recorder.Body)
		var resData map[string]interface{}
		_ = json.Unmarshal(respBody, &resData)
		assert.Equal(t, resData["error"].(map[string]interface{})["message"], "Error while reporting sighting : some error")
		assert.Equal(t, false, resData["success"])
	})

	t.Run("should return success response", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		router := httprouter.New()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockSightingService := mock_service.NewMockSightingService(ctrl)
		sightingHandler := MakeSightingHandler(mockSightingService)

		err := json.Unmarshal([]byte(reportSightingRequestBody), &req)
		if err != nil {
			panic(err)
		}

		req, _ := http.NewRequestWithContext(context.Background(), http.MethodPost,
			path, bytes.NewBuffer([]byte(reportSightingRequestBody)))

		var reportSightingReq service.ReportSightingReq
		err = json.Unmarshal([]byte(reportSightingRequestBody), &reportSightingReq)
		assert.Nil(t, err)
		mockSightingService.EXPECT().ReportSighting(gomock.Any(), reportSightingReq).Return(nil)

		router.Handle(http.MethodPost, path, middleware.ServeV1Endpoint(middleware.EmptyMiddleware,
			sightingHandler.ReportSighting))
		router.ServeHTTP(recorder, req)

		respCode := recorder.Code
		assert.Equal(t, respCode, http.StatusOK)
		respBody, _ := ioutil.ReadAll(recorder.Body)
		var resData map[string]interface{}
		_ = json.Unmarshal(respBody, &resData)
		assert.Empty(t, resData["data"])
		assert.Equal(t, resData["success"], true)
	})
}
