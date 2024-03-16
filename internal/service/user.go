package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"tigerhall_kittens/internal/logger"
	"tigerhall_kittens/internal/model"
	"tigerhall_kittens/internal/repository"
)

type CreateUserReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type UserService interface {
	CreateUser(ctx context.Context, user *CreateUserReq) error
}

type userService struct {
	userRepo repository.UserRepo
}

type UserServiceOption func(service *userService)

func NewUserService(options ...UserServiceOption) UserService {
	service := &userService{userRepo: repository.NewUserRepo()}

	for _, option := range options {
		option(service)
	}

	return service
}

func WithUserRepo(repo repository.UserRepo) UserServiceOption {
	return func(s *userService) {
		s.userRepo = repo
	}
}

func (t *userService) CreateUser(ctx context.Context, createUserReq *CreateUserReq) error {
	// TODO: refactor to use FirstOrCreate
	user, err := t.userRepo.GetUser(ctx, repository.GetUserOpts{Email: createUserReq.Email, Username: createUserReq.Username})
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.W(ctx, "Error while checking existing users", logger.Field("email", createUserReq.Email))
		return err
	}

	if user != nil {
		logger.D(ctx, "User already exists",
			logger.Field("email", createUserReq.Email),
			logger.Field("username", createUserReq.Username))
		return errors.New("user already exists with same email/username")
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
		return errors.New("error while creating user")
	}

	return nil
}
