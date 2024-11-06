package adapters

import (
	"errors"

	"github.com/phetployst/art-toys-store/modules/user/entities"
	"github.com/phetployst/art-toys-store/modules/user/usecase"
	"gorm.io/gorm"
)

type gormUserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) usecase.UserRepository {
	return &gormUserRepository{db}
}

func (r *gormUserRepository) CreateUser(user *entities.User) (uint, error) {

	if result := r.db.Create(&user); result.Error != nil {
		return 0, result.Error
	}

	return user.ID, nil
}

func (r *gormUserRepository) IsUniqueUser(email, username string) bool {
	user := new(entities.User)

	if err := r.db.Where("email = ? OR username = ?", email, username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return true
		}
	}

	return false
}

func (r *gormUserRepository) GetUserAccountById(userID uint) (*entities.User, error) {
	user := new(entities.User)

	if err := r.db.First(&user, userID).Error; err != nil {
		return nil, err
	}

	return user, nil
}
