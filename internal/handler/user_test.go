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

func TestUserHandler_CreateUser(t *testing.T) {
	var req service.CreateUserReq

	createUserRequestBody := `
		{
			"username": "new_user21",
			"password": "new_user",
			"email": "new_user11@gmail.com"
		}
	`

	invalidCreateUserRequestBody := `
		{
			"username": "new_user21",
	`

	path := "/api/v1/users"

	t.Run("should return bad request when request decode fails", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		router := httprouter.New()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUserService := mock_service.NewMockUserService(ctrl)
		userHandler := MakeUserHandler(mockUserService)

		req, _ := http.NewRequestWithContext(context.Background(), http.MethodPost,
			path, bytes.NewBuffer([]byte(invalidCreateUserRequestBody)))

		var userCreateReq service.CreateUserReq
		err := json.Unmarshal([]byte(createUserRequestBody), &userCreateReq)
		assert.Nil(t, err)

		router.Handle(http.MethodPost, path, middleware.ServeV1Endpoint(middleware.EmptyMiddleware,
			userHandler.CreateUser))
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

		mockUserService := mock_service.NewMockUserService(ctrl)
		userHandler := MakeUserHandler(mockUserService)

		err := json.Unmarshal([]byte(createUserRequestBody), &req)
		if err != nil {
			panic(err)
		}

		req, _ := http.NewRequestWithContext(context.Background(), http.MethodPost,
			path, bytes.NewBuffer([]byte(createUserRequestBody)))

		var userCreateReq service.CreateUserReq
		err = json.Unmarshal([]byte(createUserRequestBody), &userCreateReq)
		assert.Nil(t, err)
		mockUserService.EXPECT().CreateUser(gomock.Any(), &userCreateReq).Return(errors.New("some error"))

		router.Handle(http.MethodPost, path, middleware.ServeV1Endpoint(middleware.EmptyMiddleware,
			userHandler.CreateUser))
		router.ServeHTTP(recorder, req)

		respCode := recorder.Code
		assert.Equal(t, respCode, http.StatusInternalServerError)
		respBody, _ := ioutil.ReadAll(recorder.Body)
		var resData map[string]interface{}
		_ = json.Unmarshal(respBody, &resData)
		assert.Equal(t, resData["error"].(map[string]interface{})["message"], "error while creating user : some error")
		assert.Equal(t, resData["success"], false)
	})

	t.Run("should return success response", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		router := httprouter.New()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUserService := mock_service.NewMockUserService(ctrl)
		userHandler := MakeUserHandler(mockUserService)

		err := json.Unmarshal([]byte(createUserRequestBody), &req)
		if err != nil {
			panic(err)
		}

		req, _ := http.NewRequestWithContext(context.Background(), http.MethodPost,
			path, bytes.NewBuffer([]byte(createUserRequestBody)))

		var userCreateReq service.CreateUserReq
		err = json.Unmarshal([]byte(createUserRequestBody), &userCreateReq)
		assert.Nil(t, err)
		mockUserService.EXPECT().CreateUser(gomock.Any(), &userCreateReq).Return(nil)

		router.Handle(http.MethodPost, path, middleware.ServeV1Endpoint(middleware.EmptyMiddleware,
			userHandler.CreateUser))
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
