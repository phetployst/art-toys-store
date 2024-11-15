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
			Email:       "phetploy@example.com",
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