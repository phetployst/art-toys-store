package adapters

import (
	"log"
	"net/http"

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
	id := c.Param("id")

	products, err := h.usecase.GetProductById(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "internal server error"})
	}

	return c.JSON(http.StatusOK, products)
}

func (h *httpProductHandler) UpdateProduct(c echo.Context) error {
	product := new(entities.Product)
	id := c.Param("id")

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

	productUpdate, err := h.usecase.UpdateProduct(product, id)
	if err != nil {
		log.Printf("failed to create new user: %v", err)
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "internal server error"})
	}

	return c.JSON(http.StatusCreated, productUpdate)
}

func (h *httpProductHandler) DeductStock(c echo.Context) error {
	id := c.Param("id")
	count := new(entities.CountProduct)

	validator := validator.New()
	c.Echo().Validator = &CustomValidator{validator: validator}

	if err := c.Bind(&count); err != nil {
		log.Printf("failed to bind input %v", err)
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: "invalid request data"})
	}

	if err := c.Validate(count); err != nil {
		log.Printf("failed to validate input %v", err)
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: "request data validation failed"})
	}

	newStock, err := h.usecase.DeductStock(id, count)
	if err != nil {
		switch err.Error() {
		case "product not found":
			return c.JSON(http.StatusNotFound, ErrorResponse{Message: "product not found"})
		case "insufficient stock":
			return c.JSON(http.StatusBadRequest, ErrorResponse{Message: "insufficient stock"})
		default:
			return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "internal server error"})
		}
	}

	return c.JSON(http.StatusOK, newStock)
}

func (h *httpProductHandler) SearchProducts(c echo.Context) error {
	keyword := c.QueryParam("keyword")
	if keyword == "" {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: "keyword is required"})
	}

	products, err := h.usecase.SearchProducts(keyword)
	if err != nil {
		switch err.Error() {
		case "product not found":
			return c.JSON(http.StatusNotFound, ErrorResponse{Message: "product not found"})
		default:
			return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "internal server error"})
		}
	}

	return c.JSON(http.StatusOK, products)
}
