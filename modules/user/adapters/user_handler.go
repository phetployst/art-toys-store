package adapters

import (
	"github.com/phetployst/art-toys-store/modules/user/usecase"
)

type httpUserHandler struct {
	usecase usecase.UserUsecase
}

func NewUserHandler(usecase usecase.UserUsecase) *httpUserHandler {
	return &httpUserHandler{usecase}
}
