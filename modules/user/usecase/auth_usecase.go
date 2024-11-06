package usecase

import (
	"errors"

	"github.com/phetployst/art-toys-store/modules/user/entities"
)

func (s *userService) CreateNewUser(user *entities.User) (*entities.UserAccount, error) {

	if !s.repo.IsUniqueUser(user.Email, user.Username) {
		return nil, errors.New("email or username already exists")
	}

	hashedPassword, err := s.utils.HashedPassword(user.PasswordHash)
	if err != nil {
		return nil, errors.New("could not register user")
	}

	newUser := entities.User{
		Username:     user.Username,
		Email:        user.Email,
		PasswordHash: string(hashedPassword),
		Role:         user.Role,
	}

	userID, err := s.repo.CreateUser(&newUser)
	if err != nil {
		return nil, errors.New("could not register user")
	}

	userAccount, err := s.utils.GetUserAccountById(userID)
	if err != nil {
		return nil, errors.New("internal server error")
	}

	return userAccount, nil

}
