package service

import (
	"errors"
	"github.com/google/uuid"
	"tigerhall_kittens/internal/logger"
	"tigerhall_kittens/internal/model"
	"tigerhall_kittens/internal/repository"
	"time"
)

const DEFAULT_SIGHTING_RANGE_IN_METERS = 5000

// TODO:
type ReportSightingReq struct {
	TigerID   uint    `json:"tiger_id"`
	Lat       float64 `json:"lat"`
	Lon       float64 `json:"lon"`
	Timestamp string  `json:"timestamp"`
	ImageURL  string  `json:"image_url,omitempty"`
}

type SightingService interface {
	ReportSighting(user ReportSightingReq) error
	GetSightings(opts repository.GetSightingOpts) (*[]model.Sighting, error)
}

type sightingService struct {
	sightingRepo repository.SightingRepo
}

func NewSightingService() SightingService {
	return &sightingService{sightingRepo: repository.NewSightingRepo()}
}

func (t *sightingService) ReportSighting(reportSightingReq ReportSightingReq) error {
	count, err := t.sightingRepo.GetSightingsCountInRange(repository.GetSightingOpts{
		TigerID:       reportSightingReq.TigerID,
		Lat:           reportSightingReq.Lat,
		Lon:           reportSightingReq.Lon,
		RangeInMeters: DEFAULT_SIGHTING_RANGE_IN_METERS,
	})

	if err != nil {
		logger.E(nil, err, "Error while checking existing sightings", logger.Field("tiger_id", reportSightingReq.TigerID))
		return errors.New("error while checking existing sightings")
	}

	if count > 0 {
		logger.I(nil, "A sighting of the same tiger within 5 kilometers already exists", logger.Field("tiger_id", reportSightingReq.TigerID))
		return errors.New("sighting already exists in range")
	}

	sightingTs, err := time.Parse(time.RFC3339, reportSightingReq.Timestamp)

	sighting := &model.Sighting{
		TigerID: reportSightingReq.TigerID,
		// TODO: get user id either from access token or from request
		ReportedByUserID: uuid.MustParse("29c1fac8-a1e6-4859-95b9-3d7c0425b70c"),
		Lat:              reportSightingReq.Lat,
		Lon:              reportSightingReq.Lon,
		Timestamp:        sightingTs,
		ImageURL:         reportSightingReq.ImageURL,
	}

	if err := t.sightingRepo.ReportSighting(sighting); err != nil {
		// TODO: add error reporting and logging
		logger.E(nil, err, "Failed to create reportSightingReq")
		return err
	}

	return nil
}

func (t *sightingService) GetSightings(opts repository.GetSightingOpts) (*[]model.Sighting, error) {
	sightings, err := t.sightingRepo.GetSightings(opts)

	if err != nil {
		// TODO: add error reporting and logging
		return nil, err
	}

	return sightings, nil
}
