package adapters

import (
	"github.com/labstack/echo/v4"
	"github.com/phetployst/art-toys-store/modules/order/usecase"
)

type httpOrderHandler struct {
	usecase usecase.OrderUsecase
}

func NewOrderHandler(usecase usecase.OrderUsecase) *httpOrderHandler {
	return &httpOrderHandler{usecase}
}

func (h *httpOrderHandler) AddItemToCart(c echo.Context) {
	// 	รับ item

	//

}
