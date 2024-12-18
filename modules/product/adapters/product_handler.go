package adapters

import (
	"log"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/phetployst/art-toys-store/modules/product/entities"
	"github.com/phetployst/art-toys-store/modules/product/usecase"
)

type httpProductHandler struct {
	usecase usecase.ProductUsecase
}

func NewProductHandler(usecase usecase.ProductUsecase) *httpProductHandler {
	return &httpProductHandler{usecase}
}

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

func (h *httpProductHandler) CreateNewProduct(c echo.Context) error {
	product := new(entities.Product)

	validator := validator.New()
	c.Echo().Validator = &CustomValidator{validator: validator}

	if err := c.Bind(&product); err != nil {
		log.Printf("failed to bind input %v", err)
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid request data"})
	}

	if err := c.Validate(product); err != nil {
		log.Printf("failed to validate input %v", err)
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: "request data validation failed"})
	}

	newProduct, err := h.usecase.CreateNewProduct(product)
	if err != nil {
		log.Printf("failed to create new user: %v", err)
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "internal server error"})
	}

	return c.JSON(http.StatusCreated, newProduct)

}

func (h *httpProductHandler) GetAllProducts(c echo.Context) error {
	products, err := h.usecase.GetAllProducts()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "internal server error"})
	}

	return c.JSON(http.StatusOK, products)
}

func (h *httpProductHandler) GetProductById(c echo.Context) error {
	productId := c.Param("id")
	id64, err := strconv.ParseUint(productId, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "invalid product ID"})
	}

	id := uint(id64)
	products, err := h.usecase.GetProductById(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "internal server error"})
	}

	return c.JSON(http.StatusOK, products)
}
