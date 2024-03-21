package repository

import (
	"context"
	"gorm.io/gorm"

	"tigerhall_kittens/internal/db"
	"tigerhall_kittens/internal/logger"
	"tigerhall_kittens/internal/model"
)

type ListTigersOpts struct {
	Limit  int
	Offset int
}

type GetTigerOpts struct {
	TigerID uint
}

type TigerRepo interface {
	SaveTiger(ctx context.Context, tiger *model.Tiger) error
	GetTiger(ctx context.Context, opts GetTigerOpts) (*model.Tiger, error)
	GetTigers(ctx context.Context, opts ListTigersOpts) ([]model.Tiger, error)
	UpdateTiger(ctx context.Context, tiger *model.Tiger) error
}

type tigerRepo struct {
	DB *gorm.DB
}

func NewTigerRepo() TigerRepo {
	return &tigerRepo{DB: db.Get()}
}

func (t *tigerRepo) SaveTiger(ctx context.Context, tiger *model.Tiger) error {
	err := t.DB.Create(tiger).Error
	if err != nil {
		logger.E(ctx, err, "Error while saving tiger")
		return err
	}

	return nil
}

func (t *tigerRepo) GetTiger(ctx context.Context, opts GetTigerOpts) (*model.Tiger, error) {
	var tiger model.Tiger

	queryErr := t.DB.Where(model.Tiger{ID: opts.TigerID}).Find(&tiger).Error
	if queryErr != nil {
		logger.E(ctx, queryErr, "Error while fetching tigers")
		return nil, queryErr
	}

	return &tiger, nil
}

func (t *tigerRepo) GetTigers(ctx context.Context, opts ListTigersOpts) ([]model.Tiger, error) {
	var tigers []model.Tiger

	queryErr := t.DB.Limit(opts.Limit).Offset(opts.Offset).Order("last_seen_timestamp desc").Find(&tigers).Error
	if queryErr != nil {
		logger.E(ctx, queryErr, "Error while fetching tigers")
		return nil, queryErr
	}

	return tigers, nil
}

func (t *tigerRepo) UpdateTiger(ctx context.Context, tiger *model.Tiger) error {
	err := t.DB.Updates(tiger).Error
	if err != nil {
		logger.E(ctx, err, "Error while updating tiger")
		return err
	}

	return nil
}
