package repository

import (
	"gorm.io/gorm"
	"tigerhall_kittens/internal/db"
	"tigerhall_kittens/internal/model"
)

type TigerRepo interface {
	SaveTiger(tiger *model.Tiger) error
	GetTigers(limit, offset int) ([]model.Tiger, error)
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

func (t *tigerRepo) GetTigers(limit, offset int) ([]model.Tiger, error) {
	var tigers []model.Tiger

	// TODO: check why pagination is not working
	queryErr := t.DB.Order("last_seen_timestamp desc").Find(&tigers).Offset(offset).Limit(limit).Error
	if queryErr != nil {
		return nil, queryErr
	}

	return tigers, nil
}
