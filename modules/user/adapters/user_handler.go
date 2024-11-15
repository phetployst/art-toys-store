package adapters

import (
	"github.com/phetployst/art-toys-store/config"
	"github.com/phetployst/art-toys-store/modules/user/usecase"
)

type httpUserHandler struct {
	usecase usecase.UserUsecase
	config  *config.Config
}

func NewUserHandler(usecase usecase.UserUsecase, config *config.Config) *httpUserHandler {
	return &httpUserHandler{usecase, config}
}
