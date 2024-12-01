package usecase

import (
	"errors"

	"github.com/phetployst/art-toys-store/config"
	"github.com/phetployst/art-toys-store/modules/user/entities"
	"gorm.io/gorm"
)

type UserUsecase interface {
	CreateNewUser(user *entities.User) (*entities.UserAccount, error)
	Login(loginRequest *entities.Login, config *config.Config) (*entities.UserCredential, error)
	Logout(userID uint) error
	Refresh(request *entities.Refresh, config *config.Config) (*entities.UserCredential, error)
	GetUserProfile(userID string) (*entities.UserProfileResponse, error)
	UpdateUserProfile(userProfile *entities.UserProfile) (*entities.UserProfileResponse, error)
	GetAllUserProfile() (int64, []entities.UserProfileResponse, error)
}

type userService struct {
	repo  UserRepository
	utils UserUtilsService
}

func NewUserService(repo UserRepository, utils UserUtilsService) UserUsecase {
	return &userService{repo, utils}
}

func (s *userService) GetUserProfile(userID string) (*entities.UserProfileResponse, error) {
	userProfile, err := s.repo.GetUserProfileByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {

			return nil, errors.New("credential not found")
		}
		return nil, errors.New("internal server error")
	}

	return &entities.UserProfileResponse{
		UserID:            userProfile.UserID,
		Username:          userProfile.Username,
		FirstName:         userProfile.FirstName,
		LastName:          userProfile.LastName,
		Email:             userProfile.Email,
		Address:           userProfile.Address,
		ProfilePictureURL: userProfile.ProfilePictureURL,
	}, nil
}

func (s *userService) UpdateUserProfile(userProfile *entities.UserProfile) (*entities.UserProfileResponse, error) {

	if !s.repo.IsUniqueUser(userProfile.Email, userProfile.Username) {
		return nil, errors.New("email or username already exists")
	}

	userProfileUpdate, err := s.repo.UpdateUserProfile(userProfile)
	if err != nil {
		return nil, errors.New("internal server error")
	}

	return &entities.UserProfileResponse{
		UserID:            userProfileUpdate.UserID,
		Username:          userProfileUpdate.Username,
		FirstName:         userProfileUpdate.FirstName,
		LastName:          userProfileUpdate.LastName,
		Email:             userProfileUpdate.Email,
		Address:           userProfileUpdate.Address,
		ProfilePictureURL: userProfileUpdate.ProfilePictureURL,
	}, nil
}

func (s *userService) GetAllUserProfile() (int64, []entities.UserProfileResponse, error) {
	count, profiles, err := s.repo.GetAllUserProfile()
	if err != nil {
		return 0, nil, errors.New("internal server error")
	}

	if count == 0 {
		return 0, nil, errors.New("no user profiles found")
	}

	var userProfileResponses []entities.UserProfileResponse
	for _, profile := range profiles {
		userProfileResponses = append(userProfileResponses, entities.UserProfileResponse{
			UserID:            profile.UserID,
			Username:          profile.Username,
			FirstName:         profile.FirstName,
			LastName:          profile.LastName,
			Email:             profile.Email,
			Address:           profile.Address,
			ProfilePictureURL: profile.ProfilePictureURL,
		})
	}

	return count, userProfileResponses, nil
}
