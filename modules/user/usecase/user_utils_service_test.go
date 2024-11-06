package usecase

import (
	"errors"
	"reflect"
	"strings"
	"testing"

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
