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

func TestAuthHandler_Login(t *testing.T) {
	var req service.LoginUserReq

	loginRequestBody := `
			{
				"username": "user_eighty212",
				"password": "user_eighty212"
			}
	`

	successfulLoginResponse := &service.LoginUserResponse{
		AccessToken: "random_access_token",
		Error:       nil,
	}

	path := "/api/v1/login"

	t.Run("should return ISE when login user fails", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		router := httprouter.New()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAuthService := mock_service.NewMockAuthService(ctrl)
		authHandler := MakeAuthHandler(mockAuthService)

		err := json.Unmarshal([]byte(loginRequestBody), &req)
		if err != nil {
			panic(err)
		}

		req, _ := http.NewRequestWithContext(context.Background(), http.MethodPost,
			path, bytes.NewBuffer([]byte(loginRequestBody)))

		var loginReq service.LoginUserReq
		err = json.Unmarshal([]byte(loginRequestBody), &loginReq)
		assert.Nil(t, err)
		mockAuthService.EXPECT().LoginUser(gomock.Any(), loginReq).Return(nil, errors.New("some error"))

		router.Handle(http.MethodPost, path, middleware.ServeV1Endpoint(middleware.EmptyMiddleware,
			authHandler.Login))
		router.ServeHTTP(recorder, req)

		respCode := recorder.Code
		assert.Equal(t, respCode, http.StatusInternalServerError)
		respBody, _ := ioutil.ReadAll(recorder.Body)
		var resData map[string]interface{}
		_ = json.Unmarshal(respBody, &resData)
		assert.Equal(t, resData["error"].(map[string]interface{})["message"], "Error while logging in req : some error")
		assert.Equal(t, resData["success"], false)
	})

	t.Run("should return access token when login is successful", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		router := httprouter.New()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAuthService := mock_service.NewMockAuthService(ctrl)
		authHandler := MakeAuthHandler(mockAuthService)

		err := json.Unmarshal([]byte(loginRequestBody), &req)
		if err != nil {
			panic(err)
		}

		req, _ := http.NewRequestWithContext(context.Background(), http.MethodPost,
			path, bytes.NewBuffer([]byte(loginRequestBody)))

		var loginReq service.LoginUserReq
		err = json.Unmarshal([]byte(loginRequestBody), &loginReq)
		assert.Nil(t, err)
		mockAuthService.EXPECT().LoginUser(gomock.Any(), loginReq).Return(successfulLoginResponse, nil)

		router.Handle(http.MethodPost, path, middleware.ServeV1Endpoint(middleware.EmptyMiddleware,
			authHandler.Login))
		router.ServeHTTP(recorder, req)

		respCode := recorder.Code
		assert.Equal(t, respCode, http.StatusOK)
		respBody, _ := ioutil.ReadAll(recorder.Body)
		var resData map[string]interface{}
		_ = json.Unmarshal(respBody, &resData)
		assert.Equal(t, resData["data"].(map[string]interface{})["access_token"], "random_access_token")
		assert.Equal(t, resData["success"], true)
	})
}
