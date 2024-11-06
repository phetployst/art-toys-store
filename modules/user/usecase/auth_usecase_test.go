package usecase

import (
	"errors"
	"reflect"
	"testing"

	"github.com/phetployst/art-toys-store/modules/user/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRegister_auth(t *testing.T) {

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
