package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"tigerhall_kittens/internal/handler/middleware"
	"tigerhall_kittens/internal/service"
	mock_service "tigerhall_kittens/internal/service/mocks"
)

//func TestSightingHandler_GetSightings(t *testing.T) {
//	path := "/api/v1/tigers/%v/sightings"
//	tigerID := 1
//
//	t.Run("should return bad request if page is not passed in query param", func(t *testing.T) {
//		ctrl := gomock.NewController(t)
//		defer ctrl.Finish()
//
//		recorder := httptest.NewRecorder()
//		router := httprouter.New()
//
//		path = fmt.Sprintf(path, tigerID)
//		fmt.Println(path)
//
//		mockSightingService := mock_service.NewMockSightingService(ctrl)
//		sightingHandler := MakeSightingHandler(mockSightingService)
//
//		body := make(map[string]interface{})
//		data, _ := json.Marshal(body)
//		req, _ := http.NewRequest(http.MethodGet, path, bytes.NewBuffer(data))
//
//		router.Handle(http.MethodGet, path, middleware.ServeV1Endpoint(middleware.EmptyMiddleware,
//			sightingHandler.GetSightings))
//		router.ServeHTTP(recorder, req)
//
//		respCode := recorder.Code
//		assert.Equal(t, respCode, http.StatusBadRequest)
//		respBody, _ := ioutil.ReadAll(recorder.Body)
//		var resData map[string]interface{}
//		_ = json.Unmarshal(respBody, &resData)
//		assert.Equal(t, resData["error"].(map[string]interface{})["message"], "Invalid page number")
//		assert.Equal(t, resData["success"], false)
//	})
//
//	t.Run("should return bad request if per_page is not passed in query param", func(t *testing.T) {
//		ctrl := gomock.NewController(t)
//		defer ctrl.Finish()
//
//		recorder := httptest.NewRecorder()
//		router := httprouter.New()
//
//		path = fmt.Sprintf(path, tigerID)
//
//		mockSightingService := mock_service.NewMockSightingService(ctrl)
//		sightingHandler := MakeSightingHandler(mockSightingService)
//
//		body := make(map[string]interface{})
//		data, _ := json.Marshal(body)
//		req, _ := http.NewRequest(http.MethodGet, path, bytes.NewBuffer(data))
//
//		router.Handle(http.MethodGet, path, middleware.ServeV1Endpoint(middleware.EmptyMiddleware,
//			sightingHandler.GetSightings))
//		router.ServeHTTP(recorder, req)
//
//		respCode := recorder.Code
//		assert.Equal(t, respCode, http.StatusBadRequest)
//		respBody, _ := ioutil.ReadAll(recorder.Body)
//		var resData map[string]interface{}
//		_ = json.Unmarshal(respBody, &resData)
//		assert.Equal(t, resData["error"].(map[string]interface{})["message"], "Invalid page number")
//		assert.Equal(t, resData["success"], false)
//	})
//
//	t.Run("should return ISE if service returns error", func(t *testing.T) {
//
//	})
//
//	t.Run("should return success response", func(t *testing.T) {
//
//	})
//}

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

		var reportSightingReq service.ReportSightingReq
		err := json.Unmarshal([]byte(reportSightingRequestBody), &reportSightingReq)
		assert.Nil(t, err)

		router.Handle(http.MethodPost, path, middleware.ServeV1Endpoint(middleware.EmptyMiddleware,
			authHandler.ReportSighting))
		router.ServeHTTP(recorder, req)

		respCode := recorder.Code
		assert.Equal(t, respCode, http.StatusBadRequest)
		respBody, _ := ioutil.ReadAll(recorder.Body)
		var resData map[string]interface{}
		_ = json.Unmarshal(respBody, &resData)
		assert.Equal(t, resData["error"].(map[string]interface{})["message"], "Failed to decode request body")
		assert.Equal(t, resData["success"], false)
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
		assert.Equal(t, resData["success"], false)
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
