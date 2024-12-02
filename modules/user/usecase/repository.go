package usecase

import (
	"github.com/phetployst/art-toys-store/modules/user/entities"
)

type UserRepository interface {
	CreateUser(user *entities.User) (uint, error)
	IsUniqueUser(email, username string) bool
	GetUserAccountById(userId uint) (*entities.User, error)
	GetUserByUsername(username string) (*entities.User, error)
	InsertUserCredential(credential *entities.Credential) error
	GetUserCredentialByUserId(userID uint) error
	DeleteUserCredential(userID uint) error
	GetRefreshTokenByUserID(userID uint) (string, error)
	GetUserProfileByID(userID uint) (*entities.UserProfile, error)
	UpdateUserProfile(userProfile *entities.UserProfile) (*entities.UserProfile, error)
	GetAllUserProfile() (int64, []entities.UserProfile, error)
	InsertUserProfile(profile *entities.UserProfile) error
}
