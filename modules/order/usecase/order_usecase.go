package usecase

type OrderUsecase interface {
}

type OrderService struct {
	repo OrderRepository
}

func NewOrderService(repo OrderRepository) OrderUsecase {
	return &OrderService{repo}
}
