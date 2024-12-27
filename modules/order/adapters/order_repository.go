package adapters

import (
	"github.com/phetployst/art-toys-store/modules/order/usecase"
	"gorm.io/gorm"
)

type gormProductRepository struct {
	db *gorm.DB
}

func NewOrdertRepository(db *gorm.DB) usecase.OrderRepository {
	return &gormProductRepository{db}
}
