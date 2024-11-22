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
	Logout(logoutRequest *entities.Logout) error
	Refresh(userID *entities.Refresh, config *config.Config) (*entities.UserCredential, error)
	GetUserProfile(userID string) (*entities.UserProfileResponse, error)
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
