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
		switch err.Error() {
		case "email or username already exists":
			return c.JSON(http.StatusConflict, ErrorResponse{
				Message: "email or username already exists",
			})
		default:
			log.Printf("unexpected error: %v", err)
			return c.JSON(http.StatusInternalServerError, ErrorResponse{
				Message: "Internal server error",
			})
		}
	}

	return c.JSON(http.StatusCreated, userProfileUpdate)
}

func (h *httpUserHandler) GetAllUserProfile(c echo.Context) error {
	count, userProfiles, err := h.usecase.GetAllUserProfile()
	if err != nil {
		switch err.Error() {
		case "no user profiles found":
			return c.JSON(http.StatusNotFound, ErrorResponse{
				Message: "no user profiles found",
			})
		default:
			log.Printf("unexpected error: %v", err)
			return c.JSON(http.StatusInternalServerError, ErrorResponse{
				Message: "internal server error",
			})
		}
	}

	response := map[string]interface{}{
		"count":        count,
		"userProfiles": userProfiles,
	}
	return c.JSON(http.StatusOK, response)
}
