package usecase

import (
	"errors"

	"github.com/phetployst/art-toys-store/config"
	"github.com/phetployst/art-toys-store/modules/user/entities"
	"gorm.io/gorm"
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

func (s *userService) Login(loginRequest *entities.Login, config *config.Config) (*entities.UserCredential, error) {

	userAccount, err := s.repo.GetUserByUsername(loginRequest.Username)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if err := s.utils.CheckPassword(userAccount.PasswordHash, loginRequest.Password); err != nil {
		return nil, errors.New("invalid password")
	}

	accessToken, err := s.utils.GenerateJWT(userAccount.ID, userAccount.Username, userAccount.Role, config)
	if err != nil {
		return nil, errors.New("internal server error")
	}

	refreshToken, refreshTokenExpiry, err := s.utils.GenerateRefreshToken(userAccount.ID, userAccount.Username, userAccount.Role, config)
	if err != nil {
		return nil, errors.New("internal server error")
	}

	if err := s.utils.SaveUserCredentials(userAccount.ID, refreshToken, refreshTokenExpiry); err != nil {
		return nil, errors.New("internal server error")
	}

	return &entities.UserCredential{
		UserID:      userAccount.ID,
		Username:    userAccount.Username,
		Email:       userAccount.Email,
		AccessToken: accessToken,
	}, nil
}

func (s *userService) Logout(logoutRequest *entities.Logout) error {

	if err := s.repo.GetUserCredentialByUserId(logoutRequest.UserID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {

			return errors.New("credential not found")
		}
		return errors.New("internal server error")
	}

	if err := s.repo.DeleteUserCredential(logoutRequest.UserID); err != nil {
		return errors.New("internal server error")
	}

	return nil
}
