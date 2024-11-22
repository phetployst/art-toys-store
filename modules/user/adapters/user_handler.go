package adapters

import (
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/phetployst/art-toys-store/config"
	"github.com/phetployst/art-toys-store/modules/user/entities"
	"github.com/phetployst/art-toys-store/modules/user/usecase"
)

type httpUserHandler struct {
	usecase usecase.UserUsecase
	config  *config.Config
}

func NewUserHandler(usecase usecase.UserUsecase, config *config.Config) *httpUserHandler {
	return &httpUserHandler{usecase, config}
}

func (h *httpUserHandler) GetUserProfile(c echo.Context) error {
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

func (h *httpUserHandler) UpdateUserProfile(c echo.Context) error {
	userProfile := new(entities.UserProfile)

	validator := validator.New()
	c.Echo().Validator = &CustomValidator{validator: validator}

	if err := c.Bind(&userProfile); err != nil {
		log.Printf("failed to bind input %v", err)
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid request data"})
	}

	if err := c.Validate(userProfile); err != nil {
		log.Printf("failed to validate input %v", err)
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: "request data validation failed"})
	}

	userProfileUpdate, err := h.usecase.UpdateUserProfile(userProfile)
	if err != nil {
		log.Printf("failed to create new user: %v", err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, userProfileUpdate)
}
