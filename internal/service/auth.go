package service

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"tigerhall_kittens/internal/logger"
	"tigerhall_kittens/internal/repository"
	"time"
)

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

type AuthService interface {
	CreateUser(user CreateUserReq) error
	LoginUser(req LoginUserReq) (*LoginUserResponse, error)
}

type authService struct {
	userRepo repository.UserRepo
}

func NewAuthService() AuthService {
	return &userService{userRepo: repository.NewUserRepo()}
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
		Claims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		UserID: userID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(JWTSecretKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
