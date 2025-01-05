package entities

import "github.com/golang-jwt/jwt/v5"

type (
	UserAccount struct {
		UserID   uint   `json:"user_id"`
		Username string `json:"username"`
		Email    string `json:"email"`
	}

	Login struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	Refresh struct {
		RefreshToken string `json:"refresh_token" validate:"required"`
	}

	UserCredential struct {
		UserID       uint   `json:"user_id"`
		Username     string `json:"username"`
		Role         string `json:"role"`
		AccessToken  string `gorm:"type:text;not null" json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}

	JwtCustomClaims struct {
		UserID   uint   `json:"user_id"`
		Username string `json:"username"`
		Role     string `json:"role"`
		Type     string `json:"type"`
		jwt.RegisteredClaims
	}

	UserProfileResponse struct {
		UserID            uint    `gorm:"unique;not null" json:"user_id" validate:"required"`
		Username          string  `gorm:"type:varchar(50);unique;not null" json:"username" validate:"required,min=3,max=50"`
		FirstName         string  `gorm:"type:varchar(50);not null" json:"first_name" validate:"required"`
		LastName          string  `gorm:"type:varchar(50);not null" json:"last_name" validate:"required"`
		Email             string  `gorm:"type:varchar(100);unique;not null" json:"email" validate:"required,email"`
		Address           Address `gorm:"embedded" json:"address" validate:"required"`
		ProfilePictureURL string  `gorm:"type:text" json:"profile_picture_url,omitempty"`
	}
)
