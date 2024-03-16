package service

import (
	"context"

	"tigerhall_kittens/internal/logger"
	"tigerhall_kittens/internal/model"
	"tigerhall_kittens/internal/repository"
)

type TigerService interface {
	ListTigers(ctx context.Context, opts repository.ListTigersOpts) ([]model.Tiger, error)
	CreateTiger(ctx context.Context, tiger *model.Tiger) error
}

type tigerService struct {
	tigerRepo repository.TigerRepo
}

type TigerServiceOption func(service *tigerService)

func NewTigerService(options ...TigerServiceOption) TigerService {
	service := &tigerService{tigerRepo: repository.NewTigerRepo()}

	for _, option := range options {
		option(service)
	}

	return service
}

func WithTigerRepo(repo repository.TigerRepo) TigerServiceOption {
	return func(s *tigerService) {
		s.tigerRepo = repo
	}
}

func (t *tigerService) ListTigers(ctx context.Context, opts repository.ListTigersOpts) ([]model.Tiger, error) {
	tigers, err := t.tigerRepo.GetTigers(ctx, opts)
	if err != nil {
		logger.E(ctx, err, "Error while fetching tigers", logger.Field("opts", opts))
		return nil, err
	}

	return tigers, nil
}

func (t *tigerService) CreateTiger(ctx context.Context, tiger *model.Tiger) error {
	err := t.tigerRepo.SaveTiger(ctx, tiger)
	if err != nil {
		logger.E(ctx, err, "Error while saving tiger", logger.Field("tiger_id", tiger.ID))
		return err
	}

	return nil
}
