package usecase

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/phetployst/art-toys-store/config"
	"github.com/phetployst/art-toys-store/modules/user/entities"
	"golang.org/x/crypto/bcrypt"
)

type UserUtilsService interface {
	HashedPassword(password string) ([]byte, error)
	GetUserAccountById(userID uint) (*entities.UserAccount, error)
	CheckPassword(hashedPassword, inputPassword string) error
	GenerateJWT(userID uint, username, role string, config *config.Config) (string, error)
	GenerateRefreshToken(userID uint, username, role string, config *config.Config) (string, time.Time, error)
	SaveUserCredentials(userID uint, refreshToken string, expiresAt time.Time) error
	ParseAndValidateToken(tokenString, secret, expectedType string) (*JwtCustomClaims, error)
}

type userUtils struct {
	repo UserRepository
}

func NewUserUtilsService(repo UserRepository) UserUtilsService {
	return &userUtils{repo}
}

func (h *userUtils) HashedPassword(password string) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return hashedPassword, nil
}

func (h *userUtils) GetUserAccountById(userID uint) (*entities.UserAccount, error) {
	result, err := h.repo.GetUserAccountById(userID)
	if err != nil {
		return nil, err
	}

	userAccount := &entities.UserAccount{
		UserID:   result.ID,
		Username: result.Username,
		Email:    result.Email,
	}

	return userAccount, nil
}

func (h *userUtils) CheckPassword(hashedPassword, inputPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(inputPassword))
	if err != nil {
		return err
	}
	return nil
}

type JwtCustomClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	Type     string `json:"type"`
	jwt.RegisteredClaims
}

func (h *userUtils) GenerateJWT(userID uint, username, role string, config *config.Config) (string, error) {

	claims := &JwtCustomClaims{
		UserID:   userID,
		Username: username,
		Role:     role,
		Type:     "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(5 * time.Minute)),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	accessTokenString, err := accessToken.SignedString([]byte(config.Jwt.AccessTokenSecret))
	if err != nil {
		return "", err
	}

	return accessTokenString, nil
}

func (h *userUtils) GenerateRefreshToken(userID uint, username, role string, config *config.Config) (string, time.Time, error) {

	refreshTokenClaims := &JwtCustomClaims{
		UserID:   userID,
		Username: username,
		Role:     role,
		Type:     "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(config.Jwt.RefreshTokenSecret))
	if err != nil {
		return "", time.Time{}, err
	}

	return refreshTokenString, refreshTokenClaims.ExpiresAt.Time, nil
}

func (h *userUtils) SaveUserCredentials(userID uint, refreshToken string, expiresAt time.Time) error {
	credential := &entities.Credential{
		UserID:       userID,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
	}

	if err := h.repo.InsertUserCredential(credential); err != nil {
		return err
	}

	return nil
}

func (h *userUtils) ParseAndValidateToken(tokenString, secret, expectedType string) (*JwtCustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JwtCustomClaims)
	if !ok || !token.Valid || claims.Type != expectedType {
		return nil, err
	}

	return claims, nil
}
