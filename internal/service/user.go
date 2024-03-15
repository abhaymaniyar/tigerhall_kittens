package service

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
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

type LoginUserReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginUserResponse struct {
	AccessToken string `json:"access_token,omitempty"`
	Error       error  `json:"error,omitempty"`
}

// Claims represents the JWT claims
type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	jwt.Claims
}

var JWTSecretKey = []byte("your_secret_key")

type UserService interface {
	CreateUser(user CreateUserReq) error
	LoginUser(req LoginUserReq) (*LoginUserResponse, error)
}

type userService struct {
	userRepo repository.UserRepo
}

func NewUserService() UserService {
	return &userService{userRepo: repository.NewUserRepo()}
}

func (t *userService) CreateUser(createUserReq CreateUserReq) error {
	user, err := t.userRepo.GetUser(repository.GetUserOpts{Email: createUserReq.Email})
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.W(nil, "Error while checking existing users", logger.Field("email", createUserReq.Email))
		return err
	}

	if user != nil {
		logger.D(nil, "User already exists", logger.Field("email", createUserReq.Email))
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

func (t *userService) LoginUser(req LoginUserReq) (*LoginUserResponse, error) {
	user, err := t.userRepo.GetUser(repository.GetUserOpts{Username: req.Username})

	if errors.Is(err, gorm.ErrRecordNotFound) {
		logger.W(nil, "User does not exist", logger.Field("username", req.Username))
		return nil, err
	}

	if err != nil {
		logger.W(nil, "Error while getting user details", logger.Field("username", req.Username))
		return nil, err
	}

	//fmt.Println(user)
	fmt.Println(user.Password)
	//fmt.Println(req.Password)

	// TODO: fix password comparision
	//if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
	//	logger.W(nil, "Invalid username or password", logger.Field("username", req.Username))
	//	return nil, err
	//}

	token, err := generateJWTToken(user.ID)
	if err != nil {
		logger.E(nil, err, "Failed to generate token", logger.Field("username", req.Username))
		return nil, err
	}

	return &LoginUserResponse{
		AccessToken: token,
	}, nil
}

func generateJWTToken(userID uuid.UUID) (string, error) {
	claims := &Claims{
		UserID: userID,
		Claims: jwt.MapClaims{
			"expires_at": time.Now().Add(time.Hour * 24).Unix(),
			"issued_at":  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(JWTSecretKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
