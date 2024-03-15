package repository

import (
	"gorm.io/gorm"
	"tigerhall_kittens/internal/db"
	"tigerhall_kittens/internal/logger"
	"tigerhall_kittens/internal/model"
)

type ListSightingOpts struct {
	Username string
	Email    string
}

type GetSightingOpts struct {
	TigerID       uint
	Lat           float64
	Lon           float64
	RangeInMeters uint
}

type SightingRepo interface {
	ReportSighting(user *model.Sighting) error
	GetSightingsCountInRange(opts GetSightingOpts) (int64, error)
	ListSightings(opts GetUserOpts) (*[]model.Sighting, error)
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

func (t *sightingRepo) ListSightings(opts GetUserOpts) (*[]model.Sighting, error) {
	var user *[]model.Sighting
	err := t.DB.Where(&model.User{Email: opts.Email}).First(&user).Error
	if err != nil {
		logger.E(nil, err, "Error while fetching user")
		return nil, err
	}

	return nil, nil
}
