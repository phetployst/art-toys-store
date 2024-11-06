package entities

import (
	"gorm.io/gorm"
)

type (
	User struct {
		gorm.Model
		Username     string `gorm:"unique;not null" json:"username" validate:"required"`
		Email        string `gorm:"unique;not null" json:"email" validate:"required,email"`
		PasswordHash string `json:"password" validate:"required,min=8"`
		Role         string `gorm:"default:'user'" json:"role" validate:"required,oneof=user admin"`
	}

	UserAccount struct {
		UserID   uint   `json:"user_id"`
		Username string `json:"username"`
		Email    string `json:"email"`
	}
)
