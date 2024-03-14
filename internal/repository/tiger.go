package repository

import (
	"gorm.io/gorm"
	"tigerhall_kittens/internal/db"
	"tigerhall_kittens/internal/model"
)

type TigerRepo interface {
	SaveTiger(tiger *model.Tiger) error
	GetTigers() ([]*model.Tiger, error)
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

func (t *tigerRepo) GetTigers() ([]*model.Tiger, error) {
	var tigers []*model.Tiger

	queryErr := t.DB.Find(tigers).Order("lastSeenTimestamp desc").Error
	if queryErr != nil {
		return nil, queryErr
	}

	return tigers, nil
}
