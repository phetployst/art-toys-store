package usecase

import (
	"errors"
	"reflect"
	"testing"

	"github.com/phetployst/art-toys-store/modules/user/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestGetUserProfile_user(t *testing.T) {
	t.Run("successfully retrieves user profile", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		service := userService{repo: mockRepo}

		userProfile := &entities.UserProfile{
			UserID:            31,
			Username:          "phetploy",
			FirstName:         "Phet",
			LastName:          "Ploy",
			Email:             "phetploy@example.com",
			ProfilePictureURL: "https://example.com/profiles/14.jpg",
			Address: entities.Address{
				Street:     "123 Green Lane",
				City:       "Bangkok",
				State:      "Central",
				PostalCode: "10110",
				Country:    "Thailand",
			},
		}

		mockRepo.On("GetUserProfileByID", "31").Return(userProfile, nil)

		want := &entities.UserProfileResponse{
			UserID:            31,
			Username:          "phetploy",
			FirstName:         "Phet",
			LastName:          "Ploy",
			Email:             "phetploy@example.com",
			ProfilePictureURL: "https://example.com/profiles/14.jpg",
			Address: entities.Address{
				Street:     "123 Green Lane",
				City:       "Bangkok",
				State:      "Central",
				PostalCode: "10110",
				Country:    "Thailand",
			},
		}

		got, err := service.GetUserProfile("31")

		assert.NoError(t, err)
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v but want %v", got, want)
		}
	})

	t.Run("returns error when user profile is not found", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		service := userService{repo: mockRepo}

		mockRepo.On("GetUserProfileByID", "223").Return((*entities.UserProfile)(nil), gorm.ErrRecordNotFound)

		got, err := service.GetUserProfile("223")

		assert.Nil(t, got)
		assert.Error(t, err)
		assert.Equal(t, "credential not found", err.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("returns error on internal server error", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		service := userService{repo: mockRepo}

		mockRepo.On("GetUserProfileByID", "16").Return((*entities.UserProfile)(nil), errors.New("database connection failed"))

		got, err := service.GetUserProfile("16")

		assert.Nil(t, got)
		assert.Error(t, err)
		assert.Equal(t, "internal server error", err.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestUpdateUserProfile(t *testing.T) {
	t.Run("successfully update user profile", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		service := userService{repo: mockRepo}

		userProfile := &entities.UserProfile{
			UserID:            31,
			Username:          "phetploy",
			FirstName:         "Duangsamon",
			LastName:          "Jamfar",
			Email:             "phetploy@example.com",
			ProfilePictureURL: "https://example.com/profiles/14.jpg",
			Address: entities.Address{
				Street:     "123 Green Lane",
				City:       "Bangkok",
				State:      "Central",
				PostalCode: "10110",
				Country:    "Thailand",
			},
		}

		mockRepo.On("UpdateUserProfile", userProfile).Return(userProfile, nil)

		want := &entities.UserProfileResponse{
			UserID:            31,
			Username:          "phetploy",
			FirstName:         "Duangsamon",
			LastName:          "Jamfar",
			Email:             "phetploy@example.com",
			ProfilePictureURL: "https://example.com/profiles/14.jpg",
			Address: entities.Address{
				Street:     "123 Green Lane",
				City:       "Bangkok",
				State:      "Central",
				PostalCode: "10110",
				Country:    "Thailand",
			},
		}

		got, err := service.UpdateUserProfile(userProfile)

		assert.NoError(t, err)
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v but want %v", got, want)
		}
	})

	t.Run("returns error on internal server error", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		service := userService{repo: mockRepo}

		mockRepo.On("UpdateUserProfile", mock.AnythingOfType("*entities.UserProfile")).Return((*entities.UserProfile)(nil), errors.New("database error"))

		got, err := service.UpdateUserProfile(&entities.UserProfile{})

		assert.Nil(t, got)
		assert.Error(t, err)
		assert.Equal(t, "internal server error", err.Error())
		mockRepo.AssertExpectations(t)
	})
}
