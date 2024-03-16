package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	mock_notification_worker "tigerhall_kittens/cmd/notification_worker/mocks"
	"tigerhall_kittens/internal/model"
	"tigerhall_kittens/internal/repository"
	mock_repository "tigerhall_kittens/internal/repository/mocks"
)

func TestSightingService_GetSightings(t *testing.T) {
	var tigerOneID, tigerTwoID uint
	tigerOneID = 1
	tigerTwoID = 2

	t.Run("should return error when repo returns error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()
		mockErr := errors.New("error while fetching sightings")

		getSightingOpts := repository.GetSightingOpts{TigerID: tigerOneID}

		mockSightingRepo := mock_repository.NewMockSightingRepo(ctrl)
		mockSightingRepo.EXPECT().GetSightings(ctx, getSightingOpts).Return(nil, mockErr)

		sightingService := NewSightingService(
			WithSightingRepo(mockSightingRepo),
		)

		sightings, actualErr := sightingService.GetSightings(ctx, getSightingOpts)
		assert.Nil(t, sightings)
		assert.Equal(t, mockErr, actualErr)
	})

	t.Run("should return sightings", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()
		mockSightings := []model.Sighting{
			{
				ID:        uuid.New(),
				TigerID:   tigerOneID,
				SightedAt: time.Now(),
			},
			{
				ID:        uuid.New(),
				TigerID:   tigerTwoID,
				SightedAt: time.Now().Add(2 * time.Hour),
			},
		}

		getSightingOpts := repository.GetSightingOpts{TigerID: tigerOneID}

		mockSightingRepo := mock_repository.NewMockSightingRepo(ctrl)
		mockSightingRepo.EXPECT().GetSightings(ctx, getSightingOpts).Return(mockSightings, nil)

		sightingService := NewSightingService(
			WithSightingRepo(mockSightingRepo),
		)

		sightings, actualErr := sightingService.GetSightings(ctx, getSightingOpts)
		assert.Equal(t, mockSightings, sightings)
		assert.Nil(t, actualErr)
	})
}

func TestSightingService_ReportSighting(t *testing.T) {
	var tigerOneID uint
	tigerOneID = 1

	lat := 1.2
	lon := 2.2

	userID := uuid.New()

	reportSightingReq := ReportSightingReq{
		TigerID:   tigerOneID,
		Lat:       lat,
		Lon:       lon,
		Timestamp: time.Now().String(),
		ImageURL:  "imageurl.com",
	}

	t.Run("should return error when error in fetching existing sightings", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		expectedErr := ErrFetchingExistingSightings
		getSightingOpts := repository.GetSightingOpts{
			TigerID:       tigerOneID,
			RangeInMeters: DEFAULT_SIGHTING_RANGE_IN_METERS,
			Lat:           lat,
			Lon:           lon,
		}

		mockSightingRepo := mock_repository.NewMockSightingRepo(ctrl)
		mockSightingRepo.EXPECT().GetSightings(ctx, getSightingOpts).Return(nil, expectedErr)

		sightingService := NewSightingService(
			WithSightingRepo(mockSightingRepo),
		)

		reportSightingReq := ReportSightingReq{
			TigerID:   tigerOneID,
			Lat:       lat,
			Lon:       lon,
			Timestamp: time.Now().String(),
			ImageURL:  "imageurl.com",
		}

		actualErr := sightingService.ReportSighting(ctx, reportSightingReq)
		assert.Equal(t, expectedErr, actualErr)
	})

	t.Run("should return error when sighting already exists in default range", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		expectedErr := ErrSightingAlreadyReported

		// same tiger on same location one hour ago
		existingSightingsForSameTigerInDefaultRange := []model.Sighting{
			{
				TigerID:   tigerOneID,
				SightedAt: time.Now().Add(-time.Hour),
				Lat:       lat,
				Lon:       lon,
			},
		}

		getSightingOpts := repository.GetSightingOpts{
			TigerID:       tigerOneID,
			RangeInMeters: DEFAULT_SIGHTING_RANGE_IN_METERS,
			Lat:           lat,
			Lon:           lon,
		}
		mockSightingRepo := mock_repository.NewMockSightingRepo(ctrl)
		mockSightingRepo.EXPECT().GetSightings(ctx, getSightingOpts).Return(existingSightingsForSameTigerInDefaultRange, nil)

		sightingService := NewSightingService(
			WithSightingRepo(mockSightingRepo),
		)

		actualErr := sightingService.ReportSighting(ctx, reportSightingReq)
		assert.Equal(t, expectedErr, actualErr)
	})

	t.Run("should return error when report sighting fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()
		ctx = context.WithValue(ctx, "userID", userID.String())

		expectedErr := errors.New("something went wrong while reporting")

		// no existing sightings for the tiger
		var existingSightingsForSameTigerInDefaultRange []model.Sighting

		getSightingOpts := repository.GetSightingOpts{
			TigerID:       tigerOneID,
			RangeInMeters: DEFAULT_SIGHTING_RANGE_IN_METERS,
			Lat:           lat,
			Lon:           lon,
		}
		mockSightingRepo := mock_repository.NewMockSightingRepo(ctrl)
		mockSightingRepo.EXPECT().GetSightings(ctx, getSightingOpts).Return(existingSightingsForSameTigerInDefaultRange, nil)

		mockSightingRepo.EXPECT().ReportSighting(ctx, gomock.Any()).Return(expectedErr)

		sightingService := NewSightingService(
			WithSightingRepo(mockSightingRepo),
		)

		actualErr := sightingService.ReportSighting(ctx, reportSightingReq)
		assert.Equal(t, expectedErr, actualErr)
	})

	t.Run("should return error when reporting sighting through email to other users fail", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()
		ctx = context.WithValue(ctx, "userID", userID.String())

		expectedErr := ErrSendingEmailNotification

		// no existing sightings for the tiger
		var existingSightingsForSameTigerInDefaultRange []model.Sighting

		getSightingOpts := repository.GetSightingOpts{
			TigerID:       tigerOneID,
			RangeInMeters: DEFAULT_SIGHTING_RANGE_IN_METERS,
			Lat:           lat,
			Lon:           lon,
		}
		mockSightingRepo := mock_repository.NewMockSightingRepo(ctrl)
		mockSightingRepo.EXPECT().GetSightings(ctx, getSightingOpts).Return(existingSightingsForSameTigerInDefaultRange, nil)

		mockSightingRepo.EXPECT().ReportSighting(ctx, gomock.Any()).Return(nil)

		mockEmailNotifer := mock_notification_worker.NewMockSightingEmailNotifer(ctrl)
		mockEmailNotifer.EXPECT().
			ReportSightingToAllUsers(ctx, reportSightingReq.TigerID, existingSightingsForSameTigerInDefaultRange).
			Return(expectedErr)

		sightingService := NewSightingService(
			WithSightingRepo(mockSightingRepo),
			WithsightingEmailNotifer(mockEmailNotifer),
		)

		actualErr := sightingService.ReportSighting(ctx, reportSightingReq)
		assert.Equal(t, expectedErr, actualErr)
	})

	t.Run("should return nil error when reporting is successful", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()
		ctx = context.WithValue(ctx, "userID", userID.String())

		// no existing sightings for the tiger
		var existingSightingsForSameTigerInDefaultRange []model.Sighting

		getSightingOpts := repository.GetSightingOpts{
			TigerID:       tigerOneID,
			RangeInMeters: DEFAULT_SIGHTING_RANGE_IN_METERS,
			Lat:           lat,
			Lon:           lon,
		}
		mockSightingRepo := mock_repository.NewMockSightingRepo(ctrl)
		mockSightingRepo.EXPECT().GetSightings(ctx, getSightingOpts).Return(existingSightingsForSameTigerInDefaultRange, nil)

		mockSightingRepo.EXPECT().ReportSighting(ctx, gomock.Any()).Return(nil)

		mockEmailNotifer := mock_notification_worker.NewMockSightingEmailNotifer(ctrl)
		mockEmailNotifer.EXPECT().
			ReportSightingToAllUsers(ctx, reportSightingReq.TigerID, existingSightingsForSameTigerInDefaultRange).
			Return(nil)

		sightingService := NewSightingService(
			WithSightingRepo(mockSightingRepo),
			WithsightingEmailNotifer(mockEmailNotifer),
		)

		actualErr := sightingService.ReportSighting(ctx, reportSightingReq)
		assert.Nil(t, actualErr)
	})
}
