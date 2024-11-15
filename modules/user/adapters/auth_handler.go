package adapters

import (
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/phetployst/art-toys-store/modules/user/entities"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

type CustomValidator struct {
	validator *validator.Validate
}

func (c *CustomValidator) Validate(i interface{}) error {
	if err := c.validator.Struct(i); err != nil {
		return err
	}
	return nil
}

func (h *httpUserHandler) Register(c echo.Context) error {
	user := new(entities.User)

	validator := validator.New()
	c.Echo().Validator = &CustomValidator{validator: validator}

	if err := c.Bind(&user); err != nil {
		log.Printf("failed to bind input %v", err)
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid request data"})
	}

	if err := c.Validate(user); err != nil {
		log.Printf("failed to validate input %v", err)
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: "request data validation failed"})
	}

	userAccount, err := h.usecase.CreateNewUser(user)
	if err != nil {
		if err.Error() == "email or username already exists" {
			return c.JSON(http.StatusConflict, err.Error())
		}

		log.Printf("failed to create new user: %v", err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, userAccount)
}

func (h *httpUserHandler) Login(c echo.Context) error {
	loginRequest := new(entities.Login)

	if err := c.Bind(&loginRequest); err != nil {
		log.Printf("failed to bind input %v", err)
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid request data"})
	}

	userCredential, err := h.usecase.Login(loginRequest, h.config)
	if err != nil {
		if err.Error() == "user not found" {
			return c.JSON(http.StatusNotFound, ErrorResponse{Message: "User not found"})
		}
		if err.Error() == "invalid password" {
			return c.JSON(http.StatusUnauthorized, ErrorResponse{Message: "Invalid password"})
		}
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "Internal server error"})
	}

	return c.JSON(http.StatusOK, userCredential)
}