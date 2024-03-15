package service

import (
	"errors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"tigerhall_kittens/internal/logger"
	"tigerhall_kittens/internal/model"
	"tigerhall_kittens/internal/repository"
	"time"
)

// User represents a user entity
// TODO:
type CreateUserReq struct {
	Username string `json:"username"`
	Password string `json:"-"`
	Email    string `json:"email"`
}

type UserService interface {
	CreateUser(user CreateUserReq) error
}

type userService struct {
	userRepo repository.UserRepo
}

func NewUserService() UserService {
	return &userService{userRepo: repository.NewUserRepo()}
}

func (t *userService) CreateUser(createUserReq CreateUserReq) error {
	user, err := t.userRepo.ListSightings(repository.GetUserOpts{Email: createUserReq.Email})
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.W(nil, "Error while checking existing users", logger.Field("email", createUserReq.Email))
		return err
	}

	if user != nil {
		logger.D(nil, "User already exists", logger.Field("email", createUserReq.Email))
		return errors.New("user already exists")
	}

	user = &model.User{
		UserID:    uuid.New(),
		Email:     createUserReq.Email,
		Username:  createUserReq.Username,
		CreatedAt: time.Now(),
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(createUserReq.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.E(nil, err, "Failed to hash password")
		return err
	}

	user.Password = string(hashedPassword)

	if err := t.userRepo.CreateUser(user); err != nil {
		// TODO: add error reporting and logging
		logger.E(nil, err, "Failed to create user")
		return err
	}

	return nil
}
