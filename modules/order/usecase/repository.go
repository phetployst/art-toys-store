package usecase

type OrderRepository interface {
	InsertItemToCart(userID uint, productID uint, quantity int, price float64) error
}
