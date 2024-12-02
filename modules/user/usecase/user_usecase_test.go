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

		mockRepo.On("GetUserProfileByID", uint(31)).Return(userProfile, nil)

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

		got, err := service.GetUserProfile(uint(31))

		assert.NoError(t, err)
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v but want %v", got, want)
		}
	})

	t.Run("returns error when user profile is not found", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		service := userService{repo: mockRepo}

		mockRepo.On("GetUserProfileByID", uint(223)).Return((*entities.UserProfile)(nil), gorm.ErrRecordNotFound)

		got, err := service.GetUserProfile(uint(223))

		assert.Nil(t, got)
		assert.Error(t, err)
		assert.Equal(t, "credential not found", err.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("returns error on internal server error", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		service := userService{repo: mockRepo}

		mockRepo.On("GetUserProfileByID", uint(16)).Return((*entities.UserProfile)(nil), errors.New("database connection failed"))

		got, err := service.GetUserProfile(uint(16))

		assert.Nil(t, got)
		assert.Error(t, err)
		assert.Equal(t, "internal server error", err.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestUpdateUserProfile_user(t *testing.T) {
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

		mockRepo.On("IsUniqueUser", userProfile.Email, userProfile.Username).Return(true)
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

		mockRepo.On("IsUniqueUser", mock.Anything, mock.Anything).Return(true)
		mockRepo.On("UpdateUserProfile", mock.AnythingOfType("*entities.UserProfile")).Return((*entities.UserProfile)(nil), errors.New("database error"))

		got, err := service.UpdateUserProfile(&entities.UserProfile{})

		assert.Nil(t, got)
		assert.Error(t, err)
		assert.Equal(t, "internal server error", err.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("update user profile given email or username already exists", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		userService := userService{repo: mockRepo}

		mockRepo.On("IsUniqueUser", mock.Anything, mock.Anything).Return(false)

		_, err := userService.UpdateUserProfile(&entities.UserProfile{})

		assert.Error(t, err)
		assert.EqualError(t, err, "email or username already exists")

	})
}

func TestGetAllUserProfile_user(t *testing.T) {
	t.Run("successfully retrieves all user profiles", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		service := userService{repo: mockRepo}

		mockProfiles := []entities.UserProfile{
			{UserID: 31, Username: "phetploy", FirstName: "Phet", LastName: "Ploy", Email: "phetploy@example.com",
				ProfilePictureURL: "https://example.com/profiles/14.jpg",
				Address:           entities.Address{Street: "123 Green Lane", City: "Bangkok", State: "Central", PostalCode: "10110", Country: "Thailand"}},
			{UserID: 32, Username: "tonytony chopper", FirstName: "Tony", LastName: "Chopper", Email: "tonychopper@example.com",
				ProfilePictureURL: "https://example.com/profiles/32.jpg",
				Address:           entities.Address{Street: "456 Blue Street", City: "Chiang Mai", State: "North", PostalCode: "50200", Country: "Thailand"}},
		}
		mockRepo.On("GetAllUserProfile").Return(int64(2), mockProfiles, nil)

		gotCount, gotProfiles, err := service.GetAllUserProfile()

		want := []entities.UserProfileResponse{
			{UserID: 31, Username: "phetploy", FirstName: "Phet", LastName: "Ploy", Email: "phetploy@example.com",
				ProfilePictureURL: "https://example.com/profiles/14.jpg",
				Address:           entities.Address{Street: "123 Green Lane", City: "Bangkok", State: "Central", PostalCode: "10110", Country: "Thailand"}},
			{UserID: 32, Username: "tonytony chopper", FirstName: "Tony", LastName: "Chopper", Email: "tonychopper@example.com",
				ProfilePictureURL: "https://example.com/profiles/32.jpg",
				Address:           entities.Address{Street: "456 Blue Street", City: "Chiang Mai", State: "North", PostalCode: "50200", Country: "Thailand"}},
		}

		assert.NoError(t, err)
		assert.Equal(t, int64(2), gotCount)
		assert.Equal(t, want, gotProfiles)
	})

	t.Run("error retrieving user profiles", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		service := userService{repo: mockRepo}

		mockRepo.On("GetAllUserProfile").Return(int64(0), ([]entities.UserProfile)(nil), errors.New("database error"))

		gotCount, gotProfiles, err := service.GetAllUserProfile()

		assert.Error(t, err)
		assert.EqualError(t, err, "internal server error")

		if gotCount != int64(0) {
			t.Errorf("got count %v but want %v", gotCount, int64(0))
		}
		if !reflect.DeepEqual(gotProfiles, []entities.UserProfileResponse(nil)) {
			t.Errorf("got profiles %v but want %v", gotProfiles, []entities.UserProfileResponse(nil))
		}
	})

	t.Run("no user profiles found", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		service := userService{repo: mockRepo}

		mockRepo.On("GetAllUserProfile").Return(int64(0), ([]entities.UserProfile)(nil), nil)

		gotCount, gotProfiles, err := service.GetAllUserProfile()

		assert.Error(t, err)
		assert.EqualError(t, err, "no user profiles found")

		if gotCount != int64(0) {
			t.Errorf("got count %v but want %v", gotCount, int64(0))
		}
		if !reflect.DeepEqual(gotProfiles, []entities.UserProfileResponse(nil)) {
			t.Errorf("got profiles %v but want %v", gotProfiles, []entities.UserProfileResponse(nil))
		}
	})

}
