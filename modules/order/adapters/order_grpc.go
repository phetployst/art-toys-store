package adapters

import "github.com/phetployst/art-toys-store/modules/order/usecase"

type grpcOrderHandler struct {
	usecase usecase.OrderUsecase
}

func NewOrderGrpcHandler(usecase usecase.OrderUsecase) *grpcOrderHandler {
	return &grpcOrderHandler{usecase}
}
