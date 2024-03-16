package service

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"tigerhall_kittens/internal/model"
	"tigerhall_kittens/internal/repository"
	mock_repository "tigerhall_kittens/internal/repository/mocks"
)

func TestAuthService_LoginUser(t *testing.T) {
	loginReq := LoginUserReq{
		Username: "username",
		Password: "password",
	}

	t.Run("should return error when repo returns ErrRecordNotFound error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		mockUserRepo := mock_repository.NewMockUserRepo(ctrl)
		mockUserRepo.EXPECT().GetUser(ctx, repository.GetUserOpts{Username: loginReq.Username}).Return(nil, gorm.ErrRecordNotFound)

		authService := NewAuthService(
			WithUserRepoForAuthService(mockUserRepo),
		)

		resp, actualErr := authService.LoginUser(ctx, loginReq)
		assert.Equal(t, gorm.ErrRecordNotFound, actualErr)
		assert.Nil(t, resp)
	})

	t.Run("should return error when repo returns any other error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		mockUserRepo := mock_repository.NewMockUserRepo(ctrl)
		mockUserRepo.EXPECT().GetUser(ctx, repository.GetUserOpts{Username: loginReq.Username}).Return(nil, errors.New("some db error"))

		authService := NewAuthService(
			WithUserRepoForAuthService(mockUserRepo),
		)

		resp, actualErr := authService.LoginUser(ctx, loginReq)
		assert.Equal(t, errors.New("some db error"), actualErr)
		assert.Nil(t, resp)
	})

	t.Run("should return invalid username/password error if password hash does not match", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		mockUser := &model.User{
			Username: loginReq.Username,
			Password: "random_hash",
		}

		mockUserRepo := mock_repository.NewMockUserRepo(ctrl)
		mockUserRepo.EXPECT().GetUser(ctx, repository.GetUserOpts{Username: loginReq.Username}).Return(mockUser, nil)

		authService := NewAuthService(
			WithUserRepoForAuthService(mockUserRepo),
		)

		resp, actualErr := authService.LoginUser(ctx, loginReq)
		assert.Equal(t, errors.New("invalid username or password"), actualErr)
		assert.Nil(t, resp)
	})

	t.Run("should return token", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		mockUser := &model.User{
			Username: loginReq.Username,
			Password: "random_hash",
		}

		mockUserRepo := mock_repository.NewMockUserRepo(ctrl)
		mockUserRepo.EXPECT().GetUser(ctx, repository.GetUserOpts{Username: loginReq.Username}).Return(mockUser, nil)

		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(loginReq.Password), bcrypt.DefaultCost)
		mockUser.Password = string(hashedPassword)

		authService := NewAuthService(
			WithUserRepoForAuthService(mockUserRepo),
		)

		resp, actualErr := authService.LoginUser(ctx, loginReq)
		assert.Nil(t, actualErr)
		assert.Nil(t, resp.Error)
		assert.NotNil(t, resp.AccessToken)
	})
}
