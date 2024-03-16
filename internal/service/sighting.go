package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"tigerhall_kittens/cmd/notification_worker"
	"tigerhall_kittens/internal/logger"
	"tigerhall_kittens/internal/model"
	"tigerhall_kittens/internal/repository"
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
	GetSightings(ctx context.Context, opts repository.GetSightingOpts) ([]model.Sighting, error)
}

type sightingService struct {
	sightingRepo         repository.SightingRepo
	sightingEmailNotifer notification_worker.SightingEmailNotifer
}

type SightingServiceOption func(service *sightingService)

func NewSightingService(options ...SightingServiceOption) SightingService {
	service := &sightingService{
		sightingRepo:         repository.NewSightingRepo(),
		sightingEmailNotifer: notification_worker.NewSightingEmailNotifer(),
	}

	for _, option := range options {
		option(service)
	}

	return service
}

func WithSightingRepo(repo repository.SightingRepo) SightingServiceOption {
	return func(s *sightingService) {
		s.sightingRepo = repo
	}
}

func WithsightingEmailNotifer(emailNotifer notification_worker.SightingEmailNotifer) SightingServiceOption {
	return func(s *sightingService) {
		s.sightingEmailNotifer = emailNotifer
	}
}

func (t *sightingService) ReportSighting(ctx context.Context, reportSightingReq ReportSightingReq) error {
	sightings, err := t.sightingRepo.GetSightings(ctx, repository.GetSightingOpts{
		TigerID:       reportSightingReq.TigerID,
		Lat:           reportSightingReq.Lat,
		Lon:           reportSightingReq.Lon,
		RangeInMeters: DEFAULT_SIGHTING_RANGE_IN_METERS,
	})

	if err != nil {
		logger.E(ctx, err, "Error while checking existing sightings", logger.Field("tiger_id", reportSightingReq.TigerID))
		return errors.New("error while checking existing sightings")
	}

	if sightings != nil && len(sightings) > 0 {
		logger.I(ctx, "A sighting of the same tiger within range already exists",
			logger.Field("tiger_id", reportSightingReq.TigerID),
			logger.Field("range_in_meters", DEFAULT_SIGHTING_RANGE_IN_METERS))
		return errors.New("sighting already exists in range")
	}

	sightingTs, err := time.Parse(time.RFC3339, reportSightingReq.Timestamp)
	userID := ctx.Value("userID").(uuid.UUID)

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

	err = t.sightingEmailNotifer.ReportSightingToAllUsers(ctx, reportSightingReq.TigerID, sightings)
	if err != nil {
		logger.E(ctx, err, "Error while sending email notification for sightings", logger.Field("tiger_id", reportSightingReq.TigerID), logger.Field("user_id", userID))
		// TODO: should we ignore this error in case of failure in notification?
		return errors.New("error while sending email notifications about sighting")
	}

	return nil
}

func (t *sightingService) GetSightings(ctx context.Context, opts repository.GetSightingOpts) ([]model.Sighting, error) {
	sightings, err := t.sightingRepo.GetSightings(ctx, opts)

	if err != nil {
		logger.E(ctx, err, "Error while fetching sightings", logger.Field("opts", opts))
		return nil, err
	}

	return sightings, nil
}
