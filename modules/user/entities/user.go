package entities

import (
	"time"

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

	Credential struct {
		gorm.Model
		UserID       uint   `gorm:"not null;index" json:"user_id"`           // Foreign key เชื่อมโยงกับ User
		RefreshToken string `gorm:"type:text;not null" json:"refresh_token"` // Refresh token
		ExpiresAt    time.Time
	}

	UserCredential struct {
		UserID      uint   `json:"user_id"`
		Username    string `json:"username"`
		Email       string `json:"email"`
		AccessToken string `gorm:"type:text;not null" json:"access_token"`
	}
)
