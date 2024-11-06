package usecase

import (
	"github.com/phetployst/art-toys-store/modules/user/entities"
	"golang.org/x/crypto/bcrypt"
)

type UserUtilsService interface {
	HashedPassword(password string) ([]byte, error)
	GetUserAccountById(userID uint) (*entities.UserAccount, error)
}

type userUtils struct {
	repo UserRepository
}

func NewUserUtilsService(repo UserRepository) UserUtilsService {
	return &userUtils{repo}
}

func (h *userUtils) HashedPassword(password string) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return hashedPassword, nil
}

func (h *userUtils) GetUserAccountById(userID uint) (*entities.UserAccount, error) {
	result, err := h.repo.GetUserAccountById(userID)
	if err != nil {
		return nil, err
	}

	userAccount := &entities.UserAccount{
		UserID:   result.ID,
		Username: result.Username,
		Email:    result.Email,
	}

	return userAccount, nil
}
