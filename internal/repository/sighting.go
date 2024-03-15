package repository

import (
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
	ReportSighting(user *model.Sighting) error
	GetSightingsCountInRange(opts GetSightingOpts) (int64, error)
	GetSightings(opts GetSightingOpts) (*[]model.Sighting, error)
}

type sightingRepo struct {
	DB *gorm.DB
}

func NewSightingRepo() SightingRepo {
	return &sightingRepo{DB: db.Get()}
}

func (t *sightingRepo) ReportSighting(user *model.Sighting) error {
	return t.DB.Create(user).Error
}

func (t *sightingRepo) GetSightingsCountInRange(opts GetSightingOpts) (int64, error) {
	var count int64
	err := t.DB.Model(&model.Sighting{}).
		Where("tiger_id = ? AND st_distancesphere(st_makepoint(lat, lon), st_makepoint(?, ?)) < 5000", opts.TigerID, opts.Lat, opts.Lon).Count(&count).
		Error

	if err != nil {
		logger.E(nil, err, "Error while fetching sightings count", logger.Field("tiger_id", opts.TigerID), logger.Field("range_in_meters", opts.RangeInMeters))
		return 0, err
	}

	return count, nil
}

func (t *sightingRepo) GetSightings(opts GetSightingOpts) (*[]model.Sighting, error) {
	var sightings *[]model.Sighting
	err := t.DB.Order("timestamp desc").Where(&model.Sighting{TigerID: opts.TigerID}).Find(&sightings).Offset(opts.Offset).Limit(opts.Limit).Error
	if err != nil {
		logger.E(nil, err, "Error while fetching sightings")
		return nil, err
	}

	return sightings, nil
}
