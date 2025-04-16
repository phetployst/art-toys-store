package entities

import (
	"gorm.io/gorm"
)

type (
	Cart struct {
		gorm.Model
		UserID   uint   `gorm:"not null" json:"user_id"`
		Status   string `gorm:"type:varchar(20);not null" json:"status"` // e.g., active, completed
		CartItem []CartItem
	}

	CartItem struct {
		CartID    uint    `gorm:"not null" json:"cart_id"`
		ProductID uint    `gorm:"not null" json:"product_id"`
		Quantity  int     `gorm:"not null" json:"quantity" validate:"gte=1"`
		Price     float64 `gorm:"not null" json:"price"` // Snapshot of product price at the time of adding to cart
	}

	// Order struct {
	// 	gorm.Model
	// 	UserID uint `json:"user_id" gorm:"not null"`
	// 	// OrderItems      []OrderItem `json:"order_items" gorm:"foreignkey:OrderID"`
	// 	TotalAmount     float64 `json:"total_amount"`
	// 	Status          string  `json:"status" gorm:"default:'pending'"`
	// 	ShippingAddress string  `json:"shipping_address" gorm:"type:text;not null"`
	// }

	// OrderItem struct {
	// 	gorm.Model
	// 	OrderID   uint `json:"order_id" gorm:"not null"`
	// 	ProductID uint `json:"product_id" gorm:"not null"`
	// 	// ProductName string  `json:"product_name"`
	// 	// Price       float64 `json:"price"`
	// 	Quantity   int     `json:"quantity"`
	// 	TotalPrice float64 `json:"total_price"`
	// }
)
