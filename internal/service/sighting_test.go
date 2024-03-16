package service

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
	mock_notification_worker "tigerhall_kittens/cmd/notification_worker/mocks"
	"tigerhall_kittens/internal/model"
	"tigerhall_kittens/internal/repository"
	mock_repository "tigerhall_kittens/internal/repository/mocks"
	"time"
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
				ID:        1,
				TigerID:   tigerOneID,
				Timestamp: time.Now(),
			},
			{
				ID:        1,
				TigerID:   tigerTwoID,
				Timestamp: time.Now().Add(2 * time.Hour),
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

		expectedErr := errors.New("error while checking existing sightings")
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

		expectedErr := errors.New("sighting already exists in range")

		// same tiger on same location one hour ago
		existingSightingsForSameTigerInDefaultRange := []model.Sighting{
			{
				TigerID:   tigerOneID,
				Timestamp: time.Now().Add(-time.Hour),
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
		ctx = context.WithValue(ctx, "userID", userID)

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

		sightingTs, _ := time.Parse(time.RFC3339, reportSightingReq.Timestamp)
		mockSightingRepo.EXPECT().ReportSighting(ctx, &model.Sighting{
			TigerID:          reportSightingReq.TigerID,
			ReportedByUserID: userID,
			Lat:              reportSightingReq.Lat,
			Lon:              reportSightingReq.Lon,
			Timestamp:        sightingTs,
			ImageURL:         reportSightingReq.ImageURL,
		}).Return(expectedErr)

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
		ctx = context.WithValue(ctx, "userID", userID)

		expectedErr := errors.New("error while sending email notifications about sighting")

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

		sightingTs, _ := time.Parse(time.RFC3339, reportSightingReq.Timestamp)
		sighting := &model.Sighting{
			TigerID:          reportSightingReq.TigerID,
			ReportedByUserID: userID,
			Lat:              reportSightingReq.Lat,
			Lon:              reportSightingReq.Lon,
			Timestamp:        sightingTs,
			ImageURL:         reportSightingReq.ImageURL,
		}
		mockSightingRepo.EXPECT().ReportSighting(ctx, sighting).Return(nil)

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
		ctx = context.WithValue(ctx, "userID", userID)

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

		sightingTs, _ := time.Parse(time.RFC3339, reportSightingReq.Timestamp)
		sighting := &model.Sighting{
			TigerID:          reportSightingReq.TigerID,
			ReportedByUserID: userID,
			Lat:              reportSightingReq.Lat,
			Lon:              reportSightingReq.Lon,
			Timestamp:        sightingTs,
			ImageURL:         reportSightingReq.ImageURL,
		}
		mockSightingRepo.EXPECT().ReportSighting(ctx, sighting).Return(nil)

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
