package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"tigerhall_kittens/internal/model"
	"tigerhall_kittens/internal/repository"
	mock_repository "tigerhall_kittens/internal/repository/mocks"
)

func TestTigerService_ListTigers(t *testing.T) {
	opts := repository.ListTigersOpts{
		Limit:  10,
		Offset: 0,
	}

	t.Run("should return error when repo returns error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		expectedErr := errors.New("error while fetching tigers")

		mockTigerRepo := mock_repository.NewMockTigerRepo(ctrl)
		mockTigerRepo.EXPECT().GetTigers(ctx, opts).Return(nil, expectedErr)

		tigerService := NewTigerService(
			WithTigerRepo(mockTigerRepo),
		)

		tigers, actualErr := tigerService.ListTigers(ctx, opts)
		assert.Nil(t, tigers)
		assert.Equal(t, expectedErr, actualErr)
	})

	t.Run("should return tigers when repo returns tigers", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		mockTigers := []model.Tiger{
			{
				ID:          1,
				Name:        "Bengal Tiger",
				DateOfBirth: time.Now(),
			},
		}

		mockTigerRepo := mock_repository.NewMockTigerRepo(ctrl)
		mockTigerRepo.EXPECT().GetTigers(ctx, opts).Return(mockTigers, nil)

		tigerService := NewTigerService(
			WithTigerRepo(mockTigerRepo),
		)

		tigers, actualErr := tigerService.ListTigers(ctx, opts)
		assert.Equal(t, mockTigers, tigers)
		assert.Nil(t, actualErr)
	})
}

func TestTigerService_CreateTiger(t *testing.T) {
	tiger := &model.Tiger{
		ID:          1,
		Name:        "Bengal Tiger",
		DateOfBirth: time.Now(),
	}

	t.Run("should return error when repo returns error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		expectedErr := errors.New("error while fetching tigers")

		mockTigerRepo := mock_repository.NewMockTigerRepo(ctrl)
		mockTigerRepo.EXPECT().SaveTiger(ctx, tiger).Return(expectedErr)

		tigerService := NewTigerService(
			WithTigerRepo(mockTigerRepo),
		)

		actualErr := tigerService.CreateTiger(ctx, tiger)
		assert.Equal(t, expectedErr, actualErr)
	})

	t.Run("should return tigers when repo returns tigers", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		mockTigerRepo := mock_repository.NewMockTigerRepo(ctrl)
		mockTigerRepo.EXPECT().SaveTiger(ctx, tiger).Return(nil)

		tigerService := NewTigerService(
			WithTigerRepo(mockTigerRepo),
		)

		actualErr := tigerService.CreateTiger(ctx, tiger)
		assert.Nil(t, actualErr)
	})
}
