package adapters

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
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

func (h *httpUserHandler) GetProfile(c echo.Context) error {
	userID := c.Param("user_id")

	userProfile, err := h.usecase.GetUserProfile(userID)
	if err != nil {
		switch err.Error() {
		case "credential not found":
			return c.JSON(http.StatusNotFound, ErrorResponse{
				Message: "User credential not found",
			})
		default:
			log.Printf("unexpected error: %v", err)
			return c.JSON(http.StatusInternalServerError, ErrorResponse{
				Message: "Internal server error",
			})
		}
	}

	return c.JSON(http.StatusOK, userProfile)
}
