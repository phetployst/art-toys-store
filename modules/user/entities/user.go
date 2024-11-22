package entities

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type (
	User struct {
		gorm.Model
		Username     string     `gorm:"unique;not null" json:"username" validate:"required"`
		Email        string     `gorm:"unique;not null" json:"email" validate:"required,email"`
		PasswordHash string     `json:"password" validate:"required,min=8"`
		Role         string     `gorm:"default:'user'" json:"role" validate:"required,oneof=user admin"`
		Credentials  Credential `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"credentials"`
	}

	UserAccount struct {
		UserID   uint   `json:"user_id"`
		Username string `json:"username"`
		Email    string `json:"email"`
	}

	Login struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	Logout struct {
		UserID uint `gorm:"not null;index" json:"user_id"`
	}

	Refresh struct {
		UserID uint `gorm:"not null;index" json:"user_id"`
	}

	Credential struct {
		gorm.Model
		UserID       uint   `gorm:"not null;index" json:"user_id"`
		RefreshToken string `gorm:"type:text;not null" json:"refresh_token"`
		ExpiresAt    time.Time
	}

	UserCredential struct {
		UserID      uint   `json:"user_id"`
		Username    string `json:"username"`
		Role        string `json:"role"`
		AccessToken string `gorm:"type:text;not null" json:"access_token"`
	}

	JwtCustomClaims struct {
		UserID   uint   `json:"user_id"`
		Username string `json:"username"`
		Role     string `json:"role"`
		Type     string `json:"type"`
		jwt.RegisteredClaims
	}

	UserProfile struct {
		gorm.Model
		UserID            uint    `gorm:"unique;not null" json:"user_id" validate:"required"`
		Username          string  `gorm:"type:varchar(50);unique;not null" json:"username" validate:"required,min=3,max=50"`
		FirstName         string  `gorm:"type:varchar(50);not null" json:"first_name" validate:"required"`
		LastName          string  `gorm:"type:varchar(50);not null" json:"last_name" validate:"required"`
		Email             string  `gorm:"type:varchar(100);unique;not null" json:"email" validate:"required,email"`
		Address           Address `gorm:"embedded" json:"address" validate:"required"`
		ProfilePictureURL string  `gorm:"type:text" json:"profile_picture_url,omitempty"`
	}

	Address struct {
		Street     string `gorm:"type:varchar(100);not null" json:"street" validate:"required"`
		City       string `gorm:"type:varchar(50);not null" json:"city" validate:"required"`
		State      string `gorm:"type:varchar(50);not null" json:"state" validate:"required"`
		PostalCode string `gorm:"type:varchar(20);not null" json:"postal_code" validate:"required"`
		Country    string `gorm:"type:varchar(50);not null" json:"country" validate:"required"`
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
