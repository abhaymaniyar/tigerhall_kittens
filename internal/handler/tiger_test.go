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
	"tigerhall_kittens/internal/model"
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
