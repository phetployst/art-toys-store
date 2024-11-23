package usecase

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/phetployst/art-toys-store/config"
	"github.com/phetployst/art-toys-store/modules/user/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestRegisterUsecase_auth(t *testing.T) {

	t.Run("register user given successfuly", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		mockUtil := new(MockUserUtilsService)
		userService := userService{repo: mockRepo, utils: mockUtil}

		user := &entities.User{Email: "phetploy@example.com", Username: "phetploy", PasswordHash: "password1234", Role: "user"}

		mockRepo.On("IsUniqueUser", user.Email, user.Username).Return(true)
		mockUtil.On("HashedPassword", user.PasswordHash).Return([]byte("hashedpassword"), nil)
		mockRepo.On("CreateUser", mock.Anything).Return(uint(1), nil)
		mockUtil.On("GetUserAccountById", uint(1)).Return(&entities.UserAccount{UserID: uint(1), Username: user.Username, Email: user.Email}, nil)

		want := &entities.UserAccount{
			UserID:   uint(1),
			Username: "phetploy",
			Email:    "phetploy@example.com",
		}

		got, err := userService.CreateNewUser(user)

		assert.NoError(t, err)

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v but want %v", got, want)
		}
	})

	t.Run("register user given email or username already exists", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		userService := userService{repo: mockRepo}

		user := &entities.User{Email: "phetploy@example.com", Username: "phetploy", PasswordHash: "password1234", Role: "user"}

		mockRepo.On("IsUniqueUser", user.Email, user.Username).Return(false)

		_, err := userService.CreateNewUser(user)

		assert.Error(t, err)
		assert.EqualError(t, err, "email or username already exists")

	})

	t.Run("register user given password hashing error", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		mockUtil := new(MockUserUtilsService)
		userService := userService{repo: mockRepo, utils: mockUtil}

		user := &entities.User{Email: "phetploy@example.com", Username: "phetploy", PasswordHash: "password1234", Role: "user"}

		mockRepo.On("IsUniqueUser", user.Email, user.Username).Return(true)
		mockUtil.On("HashedPassword", user.PasswordHash).Return(nil, errors.New("hashed password fail"))

		_, err := userService.CreateNewUser(user)

		assert.Error(t, err)
		assert.EqualError(t, err, "could not register user")
	})

	t.Run("register user given create user due to repository error", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		mockUtil := new(MockUserUtilsService)
		userService := userService{repo: mockRepo, utils: mockUtil}

		user := &entities.User{Email: "phetploy@example.com", Username: "phetploy", PasswordHash: "password1234", Role: "user"}

		mockRepo.On("IsUniqueUser", user.Email, user.Username).Return(true)
		mockUtil.On("HashedPassword", user.PasswordHash).Return([]byte("hashedpassword"), nil)
		mockRepo.On("CreateUser", mock.Anything).Return(uint(0), errors.New("database error"))

		_, err := userService.CreateNewUser(user)

		assert.Error(t, err)
		assert.EqualError(t, err, "could not register user")
	})

	t.Run("register user given error to get user account by id ", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		mockUtil := new(MockUserUtilsService)
		userService := userService{repo: mockRepo, utils: mockUtil}

		user := &entities.User{Email: "phetploy@example.com", Username: "phetploy", PasswordHash: "password1234", Role: "user"}

		mockRepo.On("IsUniqueUser", user.Email, user.Username).Return(true)
		mockUtil.On("HashedPassword", user.PasswordHash).Return([]byte("hashedpassword"), nil)
		mockRepo.On("CreateUser", mock.Anything).Return(uint(1), nil)
		mockUtil.On("GetUserAccountById", uint(1)).Return((*entities.UserAccount)(nil), errors.New("database error"))

		_, err := userService.CreateNewUser(user)

		assert.Error(t, err)
		assert.EqualError(t, err, "internal server error")

	})

}

func TestLoginUsecase_auth(t *testing.T) {
	t.Run("login successfully", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		mockUtil := new(MockUserUtilsService)
		userService := userService{repo: mockRepo, utils: mockUtil}

		config := &config.Config{Jwt: config.Jwt{AccessTokenSecret: "accessSecret", RefreshTokenSecret: "refreshSecret"}}
		user := &entities.User{Model: gorm.Model{ID: 13}, Username: "phetploy", PasswordHash: "hashedPassword", Role: "user", Email: "phetploy@example.com"}
		loginRequest := &entities.Login{Username: "phetploy", Password: "password"}

		accessToken := "access_token"
		refreshToken := "refresh_token"
		expiry := time.Now().Add(24 * time.Hour)

		mockRepo.On("GetUserByUsername", loginRequest.Username).Return(user, nil)
		mockUtil.On("CheckPassword", user.PasswordHash, loginRequest.Password).Return(nil)
		mockUtil.On("GenerateJWT", user.ID, user.Username, user.Role, config).Return(accessToken, nil)
		mockUtil.On("GenerateRefreshToken", user.ID, user.Username, user.Role, config).Return(refreshToken, expiry, nil)
		mockUtil.On("SaveUserCredentials", user.ID, refreshToken, expiry).Return(nil)

		want := &entities.UserCredential{
			UserID:      uint(13),
			Username:    "phetploy",
			Role:        "user",
			AccessToken: "access_token",
		}

		got, err := userService.Login(loginRequest, config)

		assert.NoError(t, err)
		assert.NotNil(t, got)

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v but want %v", got, want)
		}
	})

	t.Run("login with given user not found", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		mockUtil := new(MockUserUtilsService)
		userService := userService{repo: mockRepo, utils: mockUtil}

		config := &config.Config{Jwt: config.Jwt{AccessTokenSecret: "accessSecret", RefreshTokenSecret: "refreshSecret"}}
		loginRequest := &entities.Login{Username: "nonexistentuser", Password: "password"}

		mockRepo.On("GetUserByUsername", loginRequest.Username).Return((*entities.User)(nil), errors.New("user not found"))

		result, err := userService.Login(loginRequest, config)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "user not found", err.Error())
	})

	t.Run("login with given invalid password", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		mockUtil := new(MockUserUtilsService)
		userService := userService{repo: mockRepo, utils: mockUtil}

		config := &config.Config{Jwt: config.Jwt{AccessTokenSecret: "accessSecret", RefreshTokenSecret: "refreshSecret"}}
		user := &entities.User{Model: gorm.Model{ID: 13}, Username: "phetploy", PasswordHash: "hashedPassword", Role: "user", Email: "phetploy@example.com"}
		loginRequest := &entities.Login{Username: "phetploy", Password: "wrongpassword"}

		mockRepo.On("GetUserByUsername", loginRequest.Username).Return(user, nil)
		mockUtil.On("CheckPassword", user.PasswordHash, loginRequest.Password).Return(errors.New("invalid password"))

		result, err := userService.Login(loginRequest, config)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "invalid password", err.Error())
	})

	t.Run("login with error on JWT generation", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		mockUtil := new(MockUserUtilsService)
		userService := userService{repo: mockRepo, utils: mockUtil}

		config := &config.Config{Jwt: config.Jwt{AccessTokenSecret: "accessSecret", RefreshTokenSecret: "refreshSecret"}}
		user := &entities.User{Model: gorm.Model{ID: 13}, Username: "phetploy", PasswordHash: "hashedPassword", Role: "user", Email: "phetploy@example.com"}
		loginRequest := &entities.Login{Username: "phetploy", Password: "password"}

		mockRepo.On("GetUserByUsername", loginRequest.Username).Return(user, nil)
		mockUtil.On("CheckPassword", user.PasswordHash, loginRequest.Password).Return(nil)
		mockUtil.On("GenerateJWT", user.ID, user.Username, user.Role, config).Return("", errors.New("internal server error"))

		result, err := userService.Login(loginRequest, config)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "internal server error", err.Error())
	})

	t.Run("login with error on refresh token generation", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		mockUtil := new(MockUserUtilsService)
		userService := userService{repo: mockRepo, utils: mockUtil}

		config := &config.Config{Jwt: config.Jwt{AccessTokenSecret: "accessSecret", RefreshTokenSecret: "refreshSecret"}}
		user := &entities.User{Model: gorm.Model{ID: 13}, Username: "phetploy", PasswordHash: "hashedPassword", Role: "user", Email: "phetploy@example.com"}
		loginRequest := &entities.Login{Username: "phetploy", Password: "password"}

		accessToken := "access_token"
		mockRepo.On("GetUserByUsername", loginRequest.Username).Return(user, nil)
		mockUtil.On("CheckPassword", user.PasswordHash, loginRequest.Password).Return(nil)
		mockUtil.On("GenerateJWT", user.ID, user.Username, user.Role, config).Return(accessToken, nil)
		mockUtil.On("GenerateRefreshToken", user.ID, user.Username, user.Role, config).Return("", time.Time{}, errors.New("internal server error"))

		result, err := userService.Login(loginRequest, config)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "internal server error", err.Error())
	})

	t.Run("login with error on save user credentials", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		mockUtil := new(MockUserUtilsService)
		userService := userService{repo: mockRepo, utils: mockUtil}

		config := &config.Config{Jwt: config.Jwt{AccessTokenSecret: "accessSecret", RefreshTokenSecret: "refreshSecret"}}
		user := &entities.User{Model: gorm.Model{ID: 13}, Username: "phetploy", PasswordHash: "hashedPassword", Role: "user", Email: "phetploy@example.com"}
		loginRequest := &entities.Login{Username: "phetploy", Password: "password"}

		accessToken := "access_token"
		refreshToken := "refresh_token"
		expiry := time.Now().Add(24 * time.Hour)

		mockRepo.On("GetUserByUsername", loginRequest.Username).Return(user, nil)
		mockUtil.On("CheckPassword", user.PasswordHash, loginRequest.Password).Return(nil)
		mockUtil.On("GenerateJWT", user.ID, user.Username, user.Role, config).Return(accessToken, nil)
		mockUtil.On("GenerateRefreshToken", user.ID, user.Username, user.Role, config).Return(refreshToken, expiry, nil)
		mockUtil.On("SaveUserCredentials", user.ID, refreshToken, expiry).Return(errors.New("internal server error"))

		result, err := userService.Login(loginRequest, config)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "internal server error", err.Error())
	})
}

func TestLogoutUsecase_auth(t *testing.T) {
	t.Run("successfully logs out when credential is found and deleted", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		service := userService{repo: mockRepo}

		logoutRequest := &entities.Logout{UserID: 1}

		mockRepo.On("GetUserCredentialByUserId", uint(1)).Return(nil)
		mockRepo.On("DeleteUserCredential", uint(1)).Return(nil)

		err := service.Logout(logoutRequest)

		assert.NoError(t, err)
	})

	t.Run("returns 'credential not found' when credential is not found", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		service := userService{repo: mockRepo}

		logoutRequest := &entities.Logout{UserID: 1}

		mockRepo.On("GetUserCredentialByUserId", uint(1)).Return(gorm.ErrRecordNotFound)

		err := service.Logout(logoutRequest)

		assert.EqualError(t, err, "credential not found")
	})

	t.Run("returns 'credential not found' when credential is not found", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		service := userService{repo: mockRepo}

		logoutRequest := &entities.Logout{UserID: 1}

		mockRepo.On("GetUserCredentialByUserId", uint(1)).Return(errors.New("error"))

		err := service.Logout(logoutRequest)

		assert.EqualError(t, err, "internal server error")
	})

	t.Run("returns 'internal server error' when an error occurs while deleting credential", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		service := userService{repo: mockRepo}

		logoutRequest := &entities.Logout{UserID: 1}

		mockRepo.On("GetUserCredentialByUserId", uint(1)).Return(nil)
		mockRepo.On("DeleteUserCredential", uint(1)).Return(errors.New("error"))

		err := service.Logout(logoutRequest)

		assert.EqualError(t, err, "internal server error")
	})
}

func TestRefresh_auth(t *testing.T) {

	t.Run("refresh successful", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		mockUtils := new(MockUserUtilsService)
		userService := userService{repo: mockRepo, utils: mockUtils}

		config := &config.Config{Jwt: config.Jwt{AccessTokenSecret: "accessSecret", RefreshTokenSecret: "refreshSecret"}}

		userID := &entities.Refresh{UserID: uint(13)}
		refreshToken := "validRefreshToken"
		claims := &entities.JwtCustomClaims{UserID: uint(13), Username: "tonytonychopper", Role: "user", Type: "refresh"}

		mockRepo.On("GetRefreshTokenByUserID", userID.UserID).Return(refreshToken, nil)
		mockUtils.On("ParseAndValidateToken", refreshToken, config.Jwt.RefreshTokenSecret, "refresh").Return(claims, nil)
		mockUtils.On("GenerateJWT", claims.UserID, claims.Username, claims.Role, config).Return("newAccessToken", nil)

		want := &entities.UserCredential{
			UserID:      uint(13),
			Username:    "tonytonychopper",
			Role:        "user",
			AccessToken: "newAccessToken",
		}

		got, err := userService.Refresh(userID, config)

		assert.NoError(t, err)

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v but want %v", got, want)
		}
	})

	t.Run("credential not found", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		mockUtils := new(MockUserUtilsService)
		userService := userService{repo: mockRepo, utils: mockUtils}

		config := &config.Config{Jwt: config.Jwt{AccessTokenSecret: "accessSecret", RefreshTokenSecret: "refreshSecret"}}

		userID := &entities.Refresh{UserID: uint(13)}
		mockRepo.On("GetRefreshTokenByUserID", userID.UserID).Return("", gorm.ErrRecordNotFound)

		result, err := userService.Refresh(userID, config)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.EqualError(t, err, "credential not found")
	})

	t.Run("internal server error on repo lookup", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		mockUtils := new(MockUserUtilsService)
		userService := userService{repo: mockRepo, utils: mockUtils}

		config := &config.Config{Jwt: config.Jwt{AccessTokenSecret: "accessSecret", RefreshTokenSecret: "refreshSecret"}}

		userID := &entities.Refresh{UserID: uint(12)}
		mockRepo.On("GetRefreshTokenByUserID", userID.UserID).Return("", errors.New("database error"))

		result, err := userService.Refresh(userID, config)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.EqualError(t, err, "internal server error")

	})

	t.Run("invalid refresh token", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		mockUtils := new(MockUserUtilsService)
		userService := userService{repo: mockRepo, utils: mockUtils}

		config := &config.Config{Jwt: config.Jwt{AccessTokenSecret: "accessSecret", RefreshTokenSecret: "refreshSecret"}}

		userID := &entities.Refresh{UserID: uint(13)}
		refreshToken := "invalidRefreshToken"
		mockRepo.On("GetRefreshTokenByUserID", userID.UserID).Return(refreshToken, nil)
		mockUtils.On("ParseAndValidateToken", refreshToken, config.Jwt.RefreshTokenSecret, "refresh").Return((*entities.JwtCustomClaims)(nil), errors.New("token invalid"))

		result, err := userService.Refresh(userID, config)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.EqualError(t, err, "invalid token")

		mockRepo.AssertExpectations(t)
		mockUtils.AssertExpectations(t)
	})

	t.Run("error generating new access token", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		mockUtils := new(MockUserUtilsService)
		userService := userService{repo: mockRepo, utils: mockUtils}

		config := &config.Config{Jwt: config.Jwt{AccessTokenSecret: "accessSecret", RefreshTokenSecret: "refreshSecret"}}

		userID := &entities.Refresh{UserID: uint(13)}
		refreshToken := "validRefreshToken"
		claims := &entities.JwtCustomClaims{UserID: uint(13), Username: "phetploy", Role: "user", Type: "refresh"}
		mockRepo.On("GetRefreshTokenByUserID", userID.UserID).Return(refreshToken, nil)
		mockUtils.On("ParseAndValidateToken", refreshToken, config.Jwt.RefreshTokenSecret, "refresh").Return(claims, nil)
		mockUtils.On("GenerateJWT", claims.UserID, claims.Username, claims.Role, config).Return("", errors.New("jwt error"))

		result, err := userService.Refresh(userID, config)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.EqualError(t, err, "internal server error")
	})
}

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) IsUniqueUser(email, username string) bool {
	args := m.Called(email, username)
	return args.Bool(0)
}

func (m *MockUserRepository) CreateUser(user *entities.User) (uint, error) {
	args := m.Called(user)
	return args.Get(0).(uint), args.Error(1)
}

func (m *MockUserRepository) GetUserAccountById(userID uint) (*entities.User, error) {
	args := m.Called(userID)
	return args.Get(0).(*entities.User), args.Error(1)
}

func (m *MockUserRepository) IsUserExists(username string) bool {
	args := m.Called(username)
	return args.Bool(0)
}

func (m *MockUserRepository) GetUserByUsername(username string) (*entities.User, error) {
	args := m.Called(username)
	return args.Get(0).(*entities.User), args.Error(1)
}

func (m *MockUserRepository) InsertUserCredential(credential *entities.Credential) error {
	args := m.Called(credential)
	return args.Error(0)
}

func (m *MockUserRepository) GetUserCredentialByUserId(userID uint) error {
	args := m.Called(userID)
	return args.Error(0)
}

func (m *MockUserRepository) DeleteUserCredential(userID uint) error {
	args := m.Called(userID)
	return args.Error(0)
}

func (m *MockUserRepository) GetRefreshTokenByUserID(userID uint) (string, error) {
	args := m.Called(userID)
	return args.String(0), args.Error(1)
}

func (m *MockUserRepository) GetUserProfileByID(userID string) (*entities.UserProfile, error) {
	args := m.Called(userID)
	return args.Get(0).(*entities.UserProfile), args.Error(1)
}

func (m *MockUserRepository) UpdateUserProfile(userProfile *entities.UserProfile) (*entities.UserProfile, error) {
	args := m.Called(userProfile)
	return args.Get(0).(*entities.UserProfile), args.Error(1)
}

func (m *MockUserRepository) GetAllUserProfile() (int64, []entities.UserProfile, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Get(1).([]entities.UserProfile), args.Error(2)
}

type MockUserUtilsService struct {
	mock.Mock
}

func (m *MockUserUtilsService) HashedPassword(password string) ([]byte, error) {
	args := m.Called(password)
	if args.Get(0) != nil {
		return args.Get(0).([]byte), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserUtilsService) GetUserAccountById(userID uint) (*entities.UserAccount, error) {
	args := m.Called(userID)
	return args.Get(0).(*entities.UserAccount), args.Error(1)
}

func (m *MockUserUtilsService) CheckPassword(hashedPassword, inputPassword string) error {
	args := m.Called(hashedPassword, inputPassword)
	return args.Error(0)
}

func (m *MockUserUtilsService) GenerateJWT(userID uint, username, role string, config *config.Config) (string, error) {
	args := m.Called(userID, username, role, config)
	return args.String(0), args.Error(1)
}

func (m *MockUserUtilsService) GenerateRefreshToken(userID uint, username, role string, config *config.Config) (string, time.Time, error) {
	args := m.Called(userID, username, role, config)
	return args.String(0), args.Get(1).(time.Time), args.Error(2)
}

func (m *MockUserUtilsService) SaveUserCredentials(userID uint, refreshToken string, expiresAt time.Time) error {
	args := m.Called(userID, refreshToken, expiresAt)
	return args.Error(0)
}

func (m *MockUserUtilsService) ParseAndValidateToken(tokenString, secret, expectedType string) (*entities.JwtCustomClaims, error) {
	args := m.Called(tokenString, secret, expectedType)
	return args.Get(0).(*entities.JwtCustomClaims), args.Error(1)
}
