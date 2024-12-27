package adapters

import "github.com/phetployst/art-toys-store/modules/product/usecase"

type grpcProductHandler struct {
	usecase usecase.ProductUsecase
}

func NewProductGrpcHandler(usecase usecase.ProductUsecase) *grpcProductHandler {
	return &grpcProductHandler{usecase}
}
