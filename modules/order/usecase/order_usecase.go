package usecase

type OrderUsecase interface {
}

type OrderService struct {
	repo OrderRepository
}

func NewOrderService(repo OrderRepository) OrderUsecase {
	return &OrderService{repo}
}

func (s *OrderService) AddItemToCart() {

	// ตรวจสอบว่าสินค้ามีอยู่จริงและมีสต็อกเพียงพอหรือไม่
	// get by id product service

	// เพิ่มสินค้าลงในตระกร้า
}
