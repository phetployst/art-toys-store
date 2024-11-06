package usecase

import (
	"github.com/phetployst/art-toys-store/modules/user/entities"
)

type UserRepository interface {
	CreateUser(user *entities.User) (uint, error)
	IsUniqueUser(email, username string) bool
	GetUserAccountById(userId uint) (*entities.User, error)
}
