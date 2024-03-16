package repository

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"tigerhall_kittens/internal/db"
	"tigerhall_kittens/internal/logger"
	"tigerhall_kittens/internal/model"
)

type GetUserOpts struct {
	Username string
	Email    string
}

type UserRepo interface {
	CreateUser(ctx context.Context, user *model.User) error
	GetUser(ctx context.Context, opts GetUserOpts) (*model.User, error)
}

type userRepo struct {
	DB *gorm.DB
}

func NewUserRepo() UserRepo {
	return &userRepo{DB: db.Get()}
}

func (t *userRepo) CreateUser(ctx context.Context, user *model.User) error {
	err := t.DB.Create(user).Error
	if err != nil {
		logger.E(ctx, err, "Error while creating user")
		return err
	}

	return nil
}

func (t *userRepo) GetUser(ctx context.Context, opts GetUserOpts) (*model.User, error) {
	var user model.User

	var condition string
	if opts.Username != "" {
		condition = fmt.Sprintf("username = '%s'", opts.Username)
	}

	if opts.Email != "" {
		condition = condition + fmt.Sprintf(" or email = '%s'", opts.Email)
	}

	err := t.DB.Where(condition).First(&user).Error
	if err != nil {
		logger.E(ctx, err, "Error while fetching user")
		return nil, err
	}

	return &user, nil
}
