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
	GetSightings(ctx context.Context, opts GetSightingOpts) ([]model.Sighting, error)
	ReportSighting(ctx context.Context, sighting *model.Sighting) error
}

type sightingRepo struct {
	DB *gorm.DB
}

func NewSightingRepo() SightingRepo {
	return &sightingRepo{DB: db.Get()}
}

func (t *sightingRepo) GetSightings(ctx context.Context, opts GetSightingOpts) ([]model.Sighting, error) {
	var sightings []model.Sighting
	query := t.DB.Order("timestamp desc")

	if opts.Limit != 0 {
		query = query.Limit(opts.Limit).Offset(opts.Offset)
	}

	if opts.RangeInMeters != 0 {
		query = query.Where("tiger_id = ? AND st_distancesphere(st_makepoint(lat, lon), st_makepoint(?, ?)) < ?", opts.TigerID, opts.Lat, opts.Lon, opts.RangeInMeters)
	}

	err := query.Find(&sightings).Error

	if err != nil {
		logger.E(ctx, err, "Error while fetching sightings")
		return nil, err
	}

	return sightings, nil
}

func (t *sightingRepo) ReportSighting(ctx context.Context, sighting *model.Sighting) error {
	err := t.DB.Create(sighting).Error
	if err != nil {
		logger.E(ctx, err, "Error while saving sighting")
		return err
	}

	return nil
}
