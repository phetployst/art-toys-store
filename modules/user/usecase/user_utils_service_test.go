package usecase

import (
	"errors"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/phetployst/art-toys-store/config"
	"github.com/phetployst/art-toys-store/modules/user/entities"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func TestHashedPassword_utils(t *testing.T) {

	t.Run("hashes the password successfully", func(t *testing.T) {
		userService := &userUtils{}

		password := "password1234"

		result, err := userService.HashedPassword(password)

		assert.NoError(t, err)
		assert.NotEmpty(t, result)

		err = bcrypt.CompareHashAndPassword(result, []byte(password))
		assert.NoError(t, err)
	})

	t.Run("hashes the password given invalid cost", func(t *testing.T) {
		userService := &userUtils{}

		password := strings.Repeat("a", 100)

		result, err := userService.HashedPassword(password)

		assert.Error(t, err)
		assert.Empty(t, result)
	})
}

func TestGetUserAccountById_utils(t *testing.T) {

	t.Run("get user by id successfully", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		userService := NewUserUtilsService(mockRepo)

		userID := uint(1)
		user := &entities.User{Model: gorm.Model{ID: userID}, Username: "phetploy", Email: "phetploy@example.com", Role: "user"}

		mockRepo.On("GetUserAccountById", userID).Return(user, nil)

		want := &entities.UserAccount{
			UserID:   1,
			Username: "phetploy",
			Email:    "phetploy@example.com",
		}

		got, err := userService.GetUserAccountById(userID)

		assert.NoError(t, err)

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v but want %v", got, want)
		}
	})

	t.Run("get user by given fails to retrieve user account", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		userService := NewUserUtilsService(mockRepo)

		userID := uint(2)

		mockRepo.On("GetUserAccountById", userID).Return((*entities.User)(nil), errors.New("database error"))

		result, err := userService.GetUserAccountById(userID)

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestCheckPassword_utils(t *testing.T) {
	t.Run("check password successfully", func(t *testing.T) {
		userUtils := &userUtils{}

		password := "tonytonypassword"
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		assert.NoError(t, err)

		err = userUtils.CheckPassword(string(hashedPassword), password)

		assert.NoError(t, err)
	})

	t.Run("check password given error", func(t *testing.T) {
		userUtils := &userUtils{}

		password := "tonytonypassword"
		incorrectPassword := "tonytony"
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		assert.NoError(t, err)

		err = userUtils.CheckPassword(string(hashedPassword), incorrectPassword)

		assert.Error(t, err)
		assert.Equal(t, bcrypt.ErrMismatchedHashAndPassword, err)

	})
}

func TestGenerateJWT_utils(t *testing.T) {
	t.Run("generate JWT successfully", func(t *testing.T) {
		config := &config.Config{
			Jwt: config.Jwt{
				AccessTokenSecret: "secret",
			},
		}

		userUtils := &userUtils{}

		token, err := userUtils.GenerateJWT(1, "phetploy", "user", config)

		assert.NoError(t, err)
		assert.NotEmpty(t, token)

	})
}

func TestGenerateRefreshToken_utils(t *testing.T) {
	t.Run("generate refresh token successfully", func(t *testing.T) {

		config := &config.Config{
			Jwt: config.Jwt{
				RefreshTokenSecret: "refreshsecret",
			},
		}
		userUtils := &userUtils{}

		token, expiresAt, err := userUtils.GenerateRefreshToken(1, "phetploy", "user", config)

		assert.NoError(t, err)
		assert.NotEmpty(t, token)
		assert.WithinDuration(t, time.Now().Add(24*time.Hour), expiresAt, time.Minute)
	})

}

func TestSaveUserCredentials_utils(t *testing.T) {
	t.Run("save user credentials successfully", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		userUtils := &userUtils{repo: mockRepo}

		userID := uint(1)
		refreshToken := "sample_refresh_token"
		expiresAt := time.Now().Add(24 * time.Hour)

		mockRepo.On("InsertUserCredential", &entities.Credential{
			UserID:       userID,
			RefreshToken: refreshToken,
			ExpiresAt:    expiresAt,
		}).Return(nil)

		err := userUtils.SaveUserCredentials(userID, refreshToken, expiresAt)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("save user credentials given error", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		userUtils := &userUtils{repo: mockRepo}

		userID := uint(1)
		refreshToken := "sample_refresh_token"
		expiresAt := time.Now().Add(24 * time.Hour)

		mockRepo.On("InsertUserCredential", &entities.Credential{
			UserID:       userID,
			RefreshToken: refreshToken,
			ExpiresAt:    expiresAt,
		}).Return(errors.New("insert error"))

		err := userUtils.SaveUserCredentials(userID, refreshToken, expiresAt)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestParseAndValidateToken_utils(t *testing.T) {

	t.Run("valid token", func(t *testing.T) {
		utils := &userUtils{}

		secret := "testSecret"
		expectedType := "refresh"

		claims := &JwtCustomClaims{
			UserID:   uint(1),
			Username: "phetploy",
			Role:     "user",
			Type:     expectedType,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(5 * time.Minute)),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, _ := token.SignedString([]byte(secret))

		got, err := utils.ParseAndValidateToken(tokenString, secret, expectedType)

		want := &JwtCustomClaims{
			UserID:   uint(1),
			Username: "phetploy",
			Role:     "user",
			Type:     "refresh",
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: got.RegisteredClaims.ExpiresAt,
			},
		}

		assert.NoError(t, err)

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v but want %v", got, want)
		}
	})

	t.Run("invalid token format", func(t *testing.T) {
		utils := &userUtils{}

		secret := "testSecret"
		expectedType := "refresh"

		tokenString := "invalid.token.string"

		parsedClaims, err := utils.ParseAndValidateToken(tokenString, secret, expectedType)

		assert.Error(t, err)
		assert.Nil(t, parsedClaims)
	})

	t.Run("incorrect secret", func(t *testing.T) {
		utils := &userUtils{}

		secret := "testSecret"
		expectedType := "refresh"

		claims := &JwtCustomClaims{
			UserID:   uint(12),
			Username: "phetploy",
			Role:     "user",
			Type:     expectedType,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(5 * time.Minute)),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, _ := token.SignedString([]byte("wrongSecret"))

		parsedClaims, err := utils.ParseAndValidateToken(tokenString, secret, expectedType)

		assert.Error(t, err)
		assert.Nil(t, parsedClaims)
	})

	t.Run("unexpected token type", func(t *testing.T) {
		utils := &userUtils{}

		secret := "testSecret"
		expectedType := "refresh"

		claims := &JwtCustomClaims{
			UserID:   uint(13),
			Username: "phetploy",
			Role:     "user",
			Type:     "refresh",
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(5 * time.Minute)),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, _ := token.SignedString(secret)

		parsedClaims, err := utils.ParseAndValidateToken(tokenString, secret, expectedType)

		assert.Error(t, err)
		assert.Nil(t, parsedClaims)
	})

	t.Run("expired token", func(t *testing.T) {
		utils := &userUtils{}

		secret := "testSecret"
		expectedType := "refresh"

		claims := &JwtCustomClaims{
			UserID:   uint(13),
			Username: "phetploy",
			Role:     "user",
			Type:     expectedType,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(-5 * time.Minute)), // Expired
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, _ := token.SignedString([]byte(secret))

		parsedClaims, err := utils.ParseAndValidateToken(tokenString, secret, expectedType)

		assert.Error(t, err)
		assert.Nil(t, parsedClaims)
	})
}
