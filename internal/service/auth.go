package service

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"tigerhall_kittens/internal/config"
	"tigerhall_kittens/internal/logger"
	"tigerhall_kittens/internal/repository"
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

var JWTSecretKey = []byte(config.Env.SecretKey)

type AuthService interface {
	LoginUser(ctx context.Context, req LoginUserReq) (*LoginUserResponse, error)
}

type authService struct {
	userRepo repository.UserRepo
}

type AuthServiceOption func(service *authService)

func NewAuthService(options ...AuthServiceOption) AuthService {
	service := &authService{userRepo: repository.NewUserRepo()}

	for _, option := range options {
		option(service)
	}

	return service
}

func WithUserRepoForAuthService(repo repository.UserRepo) AuthServiceOption {
	return func(s *authService) {
		s.userRepo = repo
	}
}

func (t *authService) LoginUser(ctx context.Context, req LoginUserReq) (*LoginUserResponse, error) {
	user, err := t.userRepo.GetUser(ctx, repository.GetUserOpts{Username: req.Username})

	if errors.Is(err, gorm.ErrRecordNotFound) {
		logger.W(ctx, "User does not exist", logger.Field("username", req.Username))
		return nil, err
	}

	if err != nil {
		logger.W(ctx, "Error while getting user details", logger.Field("username", req.Username))
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		logger.W(ctx, "Invalid username or password", logger.Field("username", req.Username))
		return nil, ErrInvalidUsernamePassword
	}

	token, err := generateJWTToken(user.ID)
	if err != nil {
		logger.E(ctx, err, "Failed to generate token", logger.Field("username", req.Username))
		return nil, ErrTokenGenerationFailed
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
