package service

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"tigerhall_kittens/internal/model"
	"tigerhall_kittens/internal/repository"
	mock_repository "tigerhall_kittens/internal/repository/mocks"
)

func TestUserService_CreateUser(t *testing.T) {
	createUserReq := CreateUserReq{
		Username: "username",
		Password: "incorrect",
		Email:    "usermail@gmail.com",
	}

	getUserOpts := repository.GetUserOpts{
		Email:    createUserReq.Email,
		Username: "username",
	}

	existingUserEntry := &model.User{
		ID:       uuid.New(),
		Username: "username",
		Password: "random_string",
		Email:    "usermail@gmail.com",
	}

	t.Run("should return error when repo returns error other than ErrRecordNotFound", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		expectedErr := errors.New("error while fetching tigers")

		mockUserRepo := mock_repository.NewMockUserRepo(ctrl)
		mockUserRepo.EXPECT().GetUser(ctx, getUserOpts).Return(nil, expectedErr)

		tigerService := NewUserService(
			WithUserRepo(mockUserRepo),
		)

		actualErr := tigerService.CreateUser(ctx, &createUserReq)
		assert.Equal(t, expectedErr, actualErr)
	})

	t.Run("should return error when user already exists", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		expectedErr := errors.New("user already exists with same email/username")

		mockUserRepo := mock_repository.NewMockUserRepo(ctrl)
		mockUserRepo.EXPECT().GetUser(ctx, getUserOpts).Return(existingUserEntry, nil)

		tigerService := NewUserService(
			WithUserRepo(mockUserRepo),
		)

		actualErr := tigerService.CreateUser(ctx, &createUserReq)
		assert.Equal(t, expectedErr, actualErr)
	})

	t.Run("should return error when user creation fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		mockUserRepo := mock_repository.NewMockUserRepo(ctrl)
		mockUserRepo.EXPECT().GetUser(ctx, getUserOpts).Return(nil, nil)

		mockDbErr := errors.New("some db error")
		mockUserRepo.EXPECT().CreateUser(ctx, gomock.Any()).Return(mockDbErr)

		expectedErr := errors.New("error while creating user")

		tigerService := NewUserService(
			WithUserRepo(mockUserRepo),
		)

		actualErr := tigerService.CreateUser(ctx, &createUserReq)
		assert.Equal(t, expectedErr, actualErr)
	})

	t.Run("should return nil error when user creation is successful", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		mockUserRepo := mock_repository.NewMockUserRepo(ctrl)
		mockUserRepo.EXPECT().GetUser(ctx, getUserOpts).Return(nil, nil)

		mockUserRepo.EXPECT().CreateUser(ctx, gomock.Any()).Return(nil)

		tigerService := NewUserService(
			WithUserRepo(mockUserRepo),
		)

		actualErr := tigerService.CreateUser(ctx, &createUserReq)
		assert.Equal(t, nil, actualErr)
	})
}
