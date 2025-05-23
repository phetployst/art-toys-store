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

	userProfile := &entities.UserProfile{
		UserID:   userAccount.UserID,
		Username: userAccount.Username,
		Email:    userAccount.Email,
	}

	if err := s.repo.InsertUserProfile(userProfile); err != nil {
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
		UserID:       userAccount.ID,
		Username:     userAccount.Username,
		Role:         userAccount.Role,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *userService) Logout(userID uint) error {

	if err := s.repo.GetUserCredentialByUserId(userID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {

			return errors.New("credential not found")
		}
		return errors.New("internal server error")
	}

	if err := s.repo.DeleteUserCredential(userID); err != nil {
		return errors.New("internal server error")
	}

	return nil
}

func (s *userService) Refresh(request *entities.Refresh, config *config.Config) (*entities.UserCredential, error) {

	claims, err := s.utils.ParseAndValidateToken(request.RefreshToken, config.Jwt.RefreshTokenSecret, "refresh")
	if err != nil {
		return nil, errors.New("invalid token")
	}

	newAccessToken, err := s.utils.GenerateJWT(claims.UserID, claims.Username, claims.Role, config)
	if err != nil {
		return nil, errors.New("internal server error")
	}

	return &entities.UserCredential{
		UserID:       claims.UserID,
		Username:     claims.Username,
		Role:         claims.Role,
		AccessToken:  newAccessToken,
		RefreshToken: "",
	}, nil
}
