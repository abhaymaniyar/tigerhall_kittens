package repository

import (
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
	CreateUser(user *model.User) error
	GetUser(opts GetUserOpts) (*model.User, error)
}

type userRepo struct {
	DB *gorm.DB
}

func NewUserRepo() UserRepo {
	return &userRepo{DB: db.Get()}
}

func (t *userRepo) CreateUser(user *model.User) error {
	return t.DB.Create(user).Error
}

func (t *userRepo) GetUser(opts GetUserOpts) (*model.User, error) {
	var user model.User
	err := t.DB.Where(&model.User{Email: opts.Email}).First(&user).Error
	if err != nil {
		logger.E(nil, err, "Error while fetching user")
		return nil, err
	}

	return &user, nil
}
