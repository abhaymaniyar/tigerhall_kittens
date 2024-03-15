package repository

import (
	"context"
	"gorm.io/gorm"
	"tigerhall_kittens/internal/db"
	"tigerhall_kittens/internal/logger"
	"tigerhall_kittens/internal/model"
)

type GetSightingOpts struct {
	TigerID       uint
	Lat           float64
	Lon           float64
	RangeInMeters uint
	Limit         int
	Offset        int
}

type SightingRepo interface {
	ReportSighting(ctx context.Context, sighting *model.Sighting) error
	GetSightingsCountInRange(ctx context.Context, opts GetSightingOpts) (int64, error)
	GetSightings(ctx context.Context, opts GetSightingOpts) (*[]model.Sighting, error)
}

type sightingRepo struct {
	DB *gorm.DB
}

func NewSightingRepo() SightingRepo {
	return &sightingRepo{DB: db.Get()}
}

func (t *sightingRepo) ReportSighting(ctx context.Context, sighting *model.Sighting) error {
	err := t.DB.Create(sighting).Error
	if err != nil {
		logger.E(ctx, err, "Error while saving sighting")
		return err
	}

	return nil
}

func (t *sightingRepo) GetSightingsCountInRange(ctx context.Context, opts GetSightingOpts) (int64, error) {
	var count int64
	err := t.DB.Model(&model.Sighting{}).
		Where("tiger_id = ? AND st_distancesphere(st_makepoint(lat, lon), st_makepoint(?, ?)) < 5000", opts.TigerID, opts.Lat, opts.Lon).Count(&count).
		Error

	if err != nil {
		logger.E(ctx, err, "Error while fetching sightings count",
			logger.Field("tiger_id", opts.TigerID),
			logger.Field("range_in_meters", opts.RangeInMeters))
		return 0, err
	}

	return count, nil
}

func (t *sightingRepo) GetSightings(ctx context.Context, opts GetSightingOpts) (*[]model.Sighting, error) {
	var sightings *[]model.Sighting
	err := t.DB.
		Limit(opts.Limit).
		Offset(opts.Offset).
		Order("timestamp desc").
		Where(&model.Sighting{TigerID: opts.TigerID}).
		Find(&sightings).Error

	if err != nil {
		logger.E(ctx, err, "Error while fetching sightings")
		return nil, err
	}

	return sightings, nil
}
