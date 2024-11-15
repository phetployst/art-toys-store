package usecase

import (
	"github.com/phetployst/art-toys-store/config"
	"github.com/phetployst/art-toys-store/modules/user/entities"
)

type UserUsecase interface {
	CreateNewUser(user *entities.User) (*entities.UserAccount, error)
	Login(loginRequest *entities.Login, config *config.Config) (*entities.UserCredential, error)
	Logout(logoutRequest *entities.Logout) error
}

type userService struct {
	repo  UserRepository
	utils UserUtilsService
}

func NewUserService(repo UserRepository, utils UserUtilsService) UserUsecase {
	return &userService{repo, utils}
}
