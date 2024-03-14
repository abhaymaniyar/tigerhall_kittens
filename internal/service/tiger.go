package service

import (
	"tigerhall_kittens/internal/model"
	"tigerhall_kittens/internal/repository"
)

type TigerService interface {
	CreateTiger(tiger *model.Tiger) error
	ListTigers() ([]*model.Tiger, error)
}

type tigerService struct {
	tigerRepo repository.TigerRepo
}

func NewTigerService() TigerService {
	return &tigerService{tigerRepo: repository.NewTigerRepo()}
}

func (t *tigerService) CreateTiger(tiger *model.Tiger) error {
	err := t.tigerRepo.SaveTiger(tiger)
	if err != nil {
		// TODO: add error reporting and logging
		return err
	}

	return nil
}

func (t *tigerService) ListTigers() ([]*model.Tiger, error) {
	tigers, err := t.tigerRepo.GetTigers()
	if err != nil {
		// TODO: add error reporting and logging
		return nil, err
	}

	return tigers, nil
}
