package entities

import (
	"gorm.io/gorm"
)

type (
	Product struct {
		gorm.Model
		Name        string  `gorm:"type:varchar(100);not null" json:"name" validate:"required,min=3,max=100"`
		Description string  `gorm:"type:text" json:"description" validate:"max=500"`
		Price       float64 `gorm:"type:decimal(10,2);not null" json:"price" validate:"required,gt=0"`
		Stock       int     `gorm:"type:int;not null;default:0" json:"stock" validate:"gte=0"`
		ImageURL    string  `gorm:"type:text" json:"image_url" validate:"required,url"`
		Active      bool    `gorm:"type:boolean;default:true" json:"active"`
	}

	ProductResponse struct {
		ID          uint    `json:"id"`
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
		ImageURL    string  `json:"image_url"`
	}
)
