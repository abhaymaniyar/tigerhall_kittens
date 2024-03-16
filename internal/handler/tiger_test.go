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
	mock_service "tigerhall_kittens/internal/service/mocks"
)

func TestTigerHandler_CreateTiger(t *testing.T) {
	var req model.Tiger

	createTigerRequestBody := `
		{
			"name": "Royal Bengal Tiger",
			"date_of_birth": "2018-05-15T08:00:00Z",
			"last_seen_timestamp": "2028-03-15T08:00:00Z",
			"last_seen_lat": 23.4567,
			"last_seen_lon": 45.6789
		}
	`

	invalidCreateTigerRequestBody := `
		{
			"name": "Royal Bengal Tiger",
			"date_of_birth": "2018-05-15T08:00:00Z",
	`

	path := "/api/v1/tigers"

	t.Run("should return bad request when request decode fails", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		router := httprouter.New()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockTigerService := mock_service.NewMockTigerService(ctrl)
		tigerHandler := MakeTigerHandler(mockTigerService)

		req, _ := http.NewRequestWithContext(context.Background(), http.MethodPost,
			path, bytes.NewBuffer([]byte(invalidCreateTigerRequestBody)))

		var tiger model.Tiger
		err := json.Unmarshal([]byte(createTigerRequestBody), &tiger)
		assert.Nil(t, err)

		router.Handle(http.MethodPost, path, middleware.ServeV1Endpoint(middleware.EmptyMiddleware,
			tigerHandler.CreateTiger))
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

		mockTigerService := mock_service.NewMockTigerService(ctrl)
		tigerHandler := MakeTigerHandler(mockTigerService)

		err := json.Unmarshal([]byte(createTigerRequestBody), &req)
		if err != nil {
			panic(err)
		}

		req, _ := http.NewRequestWithContext(context.Background(), http.MethodPost,
			path, bytes.NewBuffer([]byte(createTigerRequestBody)))

		var tiger model.Tiger
		err = json.Unmarshal([]byte(createTigerRequestBody), &tiger)
		assert.Nil(t, err)
		mockTigerService.EXPECT().CreateTiger(gomock.Any(), &tiger).Return(errors.New("some error"))

		router.Handle(http.MethodPost, path, middleware.ServeV1Endpoint(middleware.EmptyMiddleware,
			tigerHandler.CreateTiger))
		router.ServeHTTP(recorder, req)

		respCode := recorder.Code
		assert.Equal(t, respCode, http.StatusInternalServerError)
		respBody, _ := ioutil.ReadAll(recorder.Body)
		var resData map[string]interface{}
		_ = json.Unmarshal(respBody, &resData)
		assert.Equal(t, resData["error"].(map[string]interface{})["message"], "error while saving tiger : some error")
		assert.Equal(t, resData["success"], false)
	})

	t.Run("should return success response", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		router := httprouter.New()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockTigerService := mock_service.NewMockTigerService(ctrl)
		tigerHandler := MakeTigerHandler(mockTigerService)

		err := json.Unmarshal([]byte(createTigerRequestBody), &req)
		if err != nil {
			panic(err)
		}

		req, _ := http.NewRequestWithContext(context.Background(), http.MethodPost,
			path, bytes.NewBuffer([]byte(createTigerRequestBody)))

		var tiger model.Tiger
		err = json.Unmarshal([]byte(createTigerRequestBody), &tiger)
		assert.Nil(t, err)
		mockTigerService.EXPECT().CreateTiger(gomock.Any(), &tiger).Return(nil)

		router.Handle(http.MethodPost, path, middleware.ServeV1Endpoint(middleware.EmptyMiddleware,
			tigerHandler.CreateTiger))
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

func TestTigerHandler_ListTigers(t *testing.T) {
	tigerID := 1
	page := 1
	perPage := 10
	offset := (page - 1) * perPage

	t.Run("should return bad request if page is not passed in query param", func(t *testing.T) {
		path := "/api/v1/tigers"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		recorder := httptest.NewRecorder()
		router := httprouter.New()

		mockTigerService := mock_service.NewMockTigerService(ctrl)
		tigerHandler := MakeTigerHandler(mockTigerService)

		body := make(map[string]interface{})
		data, _ := json.Marshal(body)
		req, _ := http.NewRequest(http.MethodGet, path, bytes.NewBuffer(data))

		router.Handle(http.MethodGet, "/api/v1/tigers", middleware.ServeV1Endpoint(middleware.EmptyMiddleware,
			tigerHandler.ListTigers))
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
		path := "/api/v1/tigers?page=%v"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		recorder := httptest.NewRecorder()
		router := httprouter.New()

		path = fmt.Sprintf(path, page)

		mockTigerService := mock_service.NewMockTigerService(ctrl)
		tigerHandler := MakeTigerHandler(mockTigerService)

		body := make(map[string]interface{})
		data, _ := json.Marshal(body)
		req, _ := http.NewRequest(http.MethodGet, path, bytes.NewBuffer(data))

		router.Handle(http.MethodGet, "/api/v1/tigers", middleware.ServeV1Endpoint(middleware.EmptyMiddleware,
			tigerHandler.ListTigers))
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
		path := "/api/v1/tigers?page=%v&per_page=%v"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		recorder := httptest.NewRecorder()
		router := httprouter.New()

		path = fmt.Sprintf(path, page, perPage)

		mockTigerService := mock_service.NewMockTigerService(ctrl)
		mockTigerService.EXPECT().ListTigers(gomock.Any(), repository.ListTigersOpts{
			Limit:  perPage,
			Offset: offset,
		}).Return(nil, errors.New("some error"))

		tigerHandler := MakeTigerHandler(mockTigerService)

		body := make(map[string]interface{})
		data, _ := json.Marshal(body)
		req, _ := http.NewRequest(http.MethodGet, path, bytes.NewBuffer(data))

		router.Handle(http.MethodGet, "/api/v1/tigers", middleware.ServeV1Endpoint(middleware.EmptyMiddleware,
			tigerHandler.ListTigers))
		router.ServeHTTP(recorder, req)

		respCode := recorder.Code
		assert.Equal(t, http.StatusInternalServerError, respCode)
		respBody, _ := ioutil.ReadAll(recorder.Body)
		var resData map[string]interface{}
		_ = json.Unmarshal(respBody, &resData)
		assert.Equal(t, "Error while fetching tigers : some error", resData["error"].(map[string]interface{})["message"])
		assert.Equal(t, false, resData["success"])
	})

	t.Run("should return success response", func(t *testing.T) {
		path := "/api/v1/tigers?page=%v&per_page=%v"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		recorder := httptest.NewRecorder()
		router := httprouter.New()

		path = fmt.Sprintf(path, page, perPage)

		mockTigers := []model.Tiger{
			{
				ID:   uint(tigerID),
				Name: "tiger 1",
			},
			{
				ID:   uint(tigerID),
				Name: "tiger 2",
			},
		}

		mockTigerService := mock_service.NewMockTigerService(ctrl)
		mockTigerService.EXPECT().ListTigers(gomock.Any(), repository.ListTigersOpts{
			Limit:  perPage,
			Offset: offset,
		}).Return(mockTigers, nil)

		tigerHandler := MakeTigerHandler(mockTigerService)

		body := make(map[string]interface{})
		data, _ := json.Marshal(body)
		req, _ := http.NewRequest(http.MethodGet, path, bytes.NewBuffer(data))

		router.Handle(http.MethodGet, "/api/v1/tigers", middleware.ServeV1Endpoint(middleware.EmptyMiddleware,
			tigerHandler.ListTigers))
		router.ServeHTTP(recorder, req)

		respCode := recorder.Code
		assert.Equal(t, http.StatusOK, respCode)
		respBody, _ := ioutil.ReadAll(recorder.Body)
		var resData map[string]interface{}
		_ = json.Unmarshal(respBody, &resData)
		assert.Equal(t, true, resData["success"])
	})
}
