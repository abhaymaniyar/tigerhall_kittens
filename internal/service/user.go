package service

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"tigerhall_kittens/internal/logger"
	"tigerhall_kittens/internal/model"
	"tigerhall_kittens/internal/repository"
	"time"
)

type CreateUserReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type UserService interface {
	CreateUser(ctx context.Context, user CreateUserReq) error
}

type userService struct {
	userRepo repository.UserRepo
}

func NewUserService() UserService {
	return &userService{userRepo: repository.NewUserRepo()}
}

func (t *userService) CreateUser(ctx context.Context, createUserReq CreateUserReq) error {
	user, err := t.userRepo.GetUser(ctx, repository.GetUserOpts{Email: createUserReq.Email})
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.W(ctx, "Error while checking existing users", logger.Field("email", createUserReq.Email))
		return err
	}

	if user != nil {
		logger.D(ctx, "User already exists", logger.Field("email", createUserReq.Email))
		return errors.New("user already exists")
	}

	user = &model.User{
		ID:        uuid.New(),
		Email:     createUserReq.Email,
		Username:  createUserReq.Username,
		CreatedAt: time.Now(),
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(createUserReq.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.E(ctx, err, "Failed to hash password")
		return err
	}

	user.Password = string(hashedPassword)

	if err := t.userRepo.CreateUser(ctx, user); err != nil {
		logger.E(ctx, err, "Failed to create user")
		return err
	}

	return nil
}
