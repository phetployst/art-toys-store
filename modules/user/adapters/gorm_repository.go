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

func (r *gormUserRepository) GetUserByUsername(username string) (*entities.User, error) {
	user := new(entities.User)

	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}

	return user, nil
}

func (r *gormUserRepository) InsertUserCredential(credential *entities.Credential) error {
	if result := r.db.Create(&credential); result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *gormUserRepository) GetUserCredentialByUserId(userID uint) error {
	credential := new(entities.Credential)

	if err := r.db.Where("user_id = ? AND deleted_at IS NULL", userID).First(credential).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		return err
	}

	return nil
}

func (r *gormUserRepository) DeleteUserCredential(userID uint) error {
	credential := new(entities.Credential)

	if result := r.db.Unscoped().Where("user_id = ?", userID).Delete(credential); result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *gormUserRepository) GetRefreshTokenByUserID(userID uint) (string, error) {
	credential := new(entities.Credential)

	if err := r.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		First(&credential).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", err
		}

		return "", err
	}

	return credential.RefreshToken, nil
}

func (r *gormUserRepository) GetUserProfileByID(userID string) (*entities.UserProfile, error) {
	userProfile := new(entities.UserProfile)

	if err := r.db.Where("user_id = ? AND deleted_at IS NULL", userID).First(userProfile).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}

	return userProfile, nil
}
