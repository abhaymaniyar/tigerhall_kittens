package service

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"tigerhall_kittens/cmd/notification_worker"
	"tigerhall_kittens/internal/logger"
	"tigerhall_kittens/internal/model"
	"tigerhall_kittens/internal/repository"
	"time"
)

const DEFAULT_SIGHTING_RANGE_IN_METERS = 5000

type ReportSightingReq struct {
	TigerID   uint    `json:"tiger_id"`
	Lat       float64 `json:"lat"`
	Lon       float64 `json:"lon"`
	Timestamp string  `json:"timestamp"`
	ImageURL  string  `json:"image_url,omitempty"`
}

type SightingService interface {
	ReportSighting(ctx context.Context, user ReportSightingReq) error
	GetSightings(ctx context.Context, opts repository.GetSightingOpts) (*[]model.Sighting, error)
}

type sightingService struct {
	sightingRepo         repository.SightingRepo
	sightingEmailNotifer notification_worker.SightingEmailNotifer
}

func NewSightingService() SightingService {
	return &sightingService{
		sightingRepo:         repository.NewSightingRepo(),
		sightingEmailNotifer: notification_worker.NewSightingEmailNotifer(),
	}
}

func (t *sightingService) ReportSighting(ctx context.Context, reportSightingReq ReportSightingReq) error {
	count, err := t.sightingRepo.GetSightingsCountInRange(ctx, repository.GetSightingOpts{
		TigerID:       reportSightingReq.TigerID,
		Lat:           reportSightingReq.Lat,
		Lon:           reportSightingReq.Lon,
		RangeInMeters: DEFAULT_SIGHTING_RANGE_IN_METERS,
	})

	if err != nil {
		logger.E(ctx, err, "Error while checking existing sightings", logger.Field("tiger_id", reportSightingReq.TigerID))
		return errors.New("error while checking existing sightings")
	}

	if count > 0 {
		logger.I(ctx, "A sighting of the same tiger within 5 kilometers already exists", logger.Field("tiger_id", reportSightingReq.TigerID))
		return errors.New("sighting already exists in range")
	}

	sightingTs, err := time.Parse(time.RFC3339, reportSightingReq.Timestamp)
	userID, parseErr := uuid.Parse(ctx.Value("userID").(string))
	if parseErr != nil {
		logger.E(ctx, err, "Error while fetching user id from context")
		return errors.New("error while fetching user id from context")
	}

	sighting := &model.Sighting{
		TigerID:          reportSightingReq.TigerID,
		ReportedByUserID: userID,
		Lat:              reportSightingReq.Lat,
		Lon:              reportSightingReq.Lon,
		Timestamp:        sightingTs,
		ImageURL:         reportSightingReq.ImageURL,
	}

	if err := t.sightingRepo.ReportSighting(ctx, sighting); err != nil {
		logger.E(ctx, err, "Failed to create reportSightingReq")
		return err
	}

	err = t.sightingEmailNotifer.ReportSightingToAllUsers(ctx, reportSightingReq.TigerID)
	if err != nil {
		logger.E(ctx, err, "Error while sending email notification for sightings", logger.Field("tiger_id", reportSightingReq.TigerID), logger.Field("user_id", userID))
		return errors.New("error while checking existing sightings")
	}

	return nil
}

func (t *sightingService) GetSightings(ctx context.Context, opts repository.GetSightingOpts) (*[]model.Sighting, error) {
	sightings, err := t.sightingRepo.GetSightings(ctx, opts)

	if err != nil {
		logger.E(ctx, err, "Error while fetching sightings", logger.Field("opts", opts))
		return nil, err
	}

	return sightings, nil
}
