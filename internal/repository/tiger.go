package repository

import (
	"gorm.io/gorm"
	"tigerhall_kittens/internal/db"
	"tigerhall_kittens/internal/model"
)

type ListTigersOpts struct {
	Limit  int
	Offset int
}

type TigerRepo interface {
	SaveTiger(tiger *model.Tiger) error
	GetTigers(opts ListTigersOpts) ([]model.Tiger, error)
}

type tigerRepo struct {
	DB *gorm.DB
}

func NewTigerRepo() TigerRepo {
	return &tigerRepo{DB: db.Get()}
}

func (t *tigerRepo) SaveTiger(tiger *model.Tiger) error {
	return t.DB.Create(tiger).Error
}

func (t *tigerRepo) GetTigers(opts ListTigersOpts) ([]model.Tiger, error) {
	var tigers []model.Tiger

	queryErr := t.DB.Limit(opts.Limit).Offset(opts.Offset).Order("last_seen_timestamp desc").Find(&tigers).Error
	if queryErr != nil {
		return nil, queryErr
	}

	return tigers, nil
}
