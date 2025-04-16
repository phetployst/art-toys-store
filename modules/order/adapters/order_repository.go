package adapters

import (
	"errors"

	"github.com/phetployst/art-toys-store/modules/order/entities"
	"github.com/phetployst/art-toys-store/modules/order/usecase"
	"gorm.io/gorm"
)

type gormOrderRepository struct {
	db *gorm.DB
}

func NewOrdertRepository(db *gorm.DB) usecase.OrderRepository {
	return &gormOrderRepository{db}
}

func (r *gormOrderRepository) InsertItemToCart(userID uint, productID uint, quantity int, price float64) error {
	var cart entities.Cart

	// ค้นหาว่าผู้ใช้มีตะกร้าที่ active หรือไม่
	if err := r.db.Where("user_id = ? AND status = ?", userID, "active").FirstOrCreate(&cart, entities.Cart{
		UserID: userID,
		Status: "active",
	}).Error; err != nil {
		return err
	}

	var cartItem entities.CartItem
	// ตรวจสอบว่าสินค้านี้อยู่ในตะกร้าอยู่แล้วหรือไม่
	if err := r.db.Where("cart_id = ? AND product_id = ?", cart.ID, productID).First(&cartItem).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// ถ้าไม่พบให้เพิ่มสินค้าใหม่ในตะกร้า
			cartItem = entities.CartItem{
				CartID:    cart.ID,
				ProductID: productID,
				Quantity:  quantity,
				Price:     price,
			}
			return r.db.Create(&cartItem).Error
		}
		return err
	}

	// ถ้าพบสินค้าอยู่แล้วให้เพิ่มจำนวน
	cartItem.Quantity += quantity
	return r.db.Save(&cartItem).Error
}
