package adapters

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/phetployst/art-toys-store/modules/product/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateNewProduct(t *testing.T) {
	t.Run("create new product given valid input should be successful", func(t *testing.T) {
		mockService := new(MockProductUsecase)
		handler := &httpProductHandler{usecase: mockService}

		e := echo.New()
		defer e.Close()

		mockService.On("CreateNewProduct", mock.AnythingOfType("*entities.Product")).Return(&entities.ProductResponse{ID: uint(30), Name: "Customizable Art Toy", Description: "A fully customizable art toy.", Price: 20.0, ImageURL: "https://example.com/images/dimoo-starry-night.jpg"}, nil)

		body := `{"name": "Dimoo Starry Night", "description": "Dimoo inspired by Van Gogh's 'Starry Night,' featuring a dreamy and artistic design.", "price": 49.99, "stock": 25, "image_url": "https://example.com/images/dimoo-starry-night.jpg", "active": true}`
		request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
		request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		response := httptest.NewRecorder()
		c := e.NewContext(request, response)

		err := handler.CreateNewProduct(c)

		expectedJSON := `{"id": 30, "name": "Customizable Art Toy", "description": "A fully customizable art toy.", "price": 20.0, "image_url": "https://example.com/images/dimoo-starry-night.jpg"}`
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, response.Code)
		assert.JSONEq(t, expectedJSON, response.Body.String())

	})

	t.Run("create new product given error during binding", func(t *testing.T) {
		mockService := new(MockProductUsecase)
		handler := &httpProductHandler{usecase: mockService}

		e := echo.New()
		defer e.Close()

		mockService.On("CreateNewProduct", mock.AnythingOfType("*entities.Product")).Return((*entities.ProductResponse)(nil), nil)

		body := `{hello!}`
		request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
		request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		response := httptest.NewRecorder()
		c := e.NewContext(request, response)

		err := handler.CreateNewProduct(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, response.Code)
	})

	t.Run("create new product given invalid input", func(t *testing.T) {
		mockService := new(MockProductUsecase)
		handler := &httpProductHandler{usecase: mockService}

		e := echo.New()
		defer e.Close()

		mockService.On("CreateNewProduct", mock.AnythingOfType("*entities.Product")).Return((*entities.ProductResponse)(nil), nil)

		body := `{"name": "Customizable Art Toy"}`
		request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
		request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		response := httptest.NewRecorder()
		c := e.NewContext(request, response)

		err := handler.CreateNewProduct(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, response.Code)
	})

	t.Run("create new product given internal server error", func(t *testing.T) {
		mockService := new(MockProductUsecase)
		handler := &httpProductHandler{usecase: mockService}

		e := echo.New()
		defer e.Close()

		mockService.On("CreateNewProduct", mock.AnythingOfType("*entities.Product")).Return((*entities.ProductResponse)(nil), errors.New("internal server error"))

		body := `{"name": "Dimoo Starry Night", "description": "Dimoo inspired by Van Gogh's 'Starry Night,' featuring a dreamy and artistic design.", "price": 49.99, "stock": 25, "image_url": "https://example.com/images/dimoo-starry-night.jpg", "active": true}`
		request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
		request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		response := httptest.NewRecorder()
		c := e.NewContext(request, response)

		err := handler.CreateNewProduct(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, response.Code)
	})
}

func TestGetAllProducts(t *testing.T) {
	t.Run("get all product successfully", func(t *testing.T) {
		mockService := new(MockProductUsecase)
		handler := &httpProductHandler{usecase: mockService}

		e := echo.New()
		defer e.Close()

		products := []entities.ProductResponse{
			{ID: uint(13), Name: "Dimoo Starry Night", Description: "Dimoo inspired by Van Gogh's 'Starry Night,' featuring a dreamy and artistic design.", Price: 49.99, ImageURL: "https://example.com/images/dimoo-starry-night.jpg"},
			{ID: uint(14), Name: "Pucky Forest Fairy", Description: "A magical art toy figure from Pucky, with a whimsical forest fairy design.", Price: 44.99, ImageURL: "https://example.com/images/pucky-forest-fairy.jpg"},
		}

		mockService.On("GetAllProducts").Return(products, nil)

		request := httptest.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()
		c := e.NewContext(request, response)

		err := handler.GetAllProducts(c)

		expectedJSON := `[{"id":13,"description":"Dimoo inspired by Van Gogh's 'Starry Night,' featuring a dreamy and artistic design.","image_url":"https://example.com/images/dimoo-starry-night.jpg","name":"Dimoo Starry Night","price":49.99},
  			{"id":14,"description":"A magical art toy figure from Pucky, with a whimsical forest fairy design.","image_url":"https://example.com/images/pucky-forest-fairy.jpg","name":"Pucky Forest Fairy","price":44.99}]`

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, response.Code)
		assert.JSONEq(t, expectedJSON, response.Body.String())

	})

	t.Run("get all product given error", func(t *testing.T) {
		mockService := new(MockProductUsecase)
		handler := &httpProductHandler{usecase: mockService}

		e := echo.New()
		defer e.Close()

		mockService.On("GetAllProducts").Return(([]entities.ProductResponse)(nil), errors.New("database error"))

		request := httptest.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()
		c := e.NewContext(request, response)

		err := handler.GetAllProducts(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, response.Code)
	})
}

func TestGetProductById(t *testing.T) {
	t.Run("get product by id successfully", func(t *testing.T) {
		mockService := new(MockProductUsecase)
		handler := &httpProductHandler{usecase: mockService}

		e := echo.New()
		defer e.Close()

		product := &entities.ProductResponse{ID: uint(12), Name: "Pucky Forest Fairy", Description: "A magical art toy figure from Pucky, with a whimsical forest fairy design.", Price: 44.99, ImageURL: "https://example.com/images/pucky-forest-fairy.jpg"}
		mockService.On("GetProductById", "12").Return(product, nil)

		request := httptest.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()
		c := e.NewContext(request, response)
		c.SetParamNames("id")
		c.SetParamValues("12")

		err := handler.GetProductById(c)

		expectedJSON := `{"id":12,"description":"A magical art toy figure from Pucky, with a whimsical forest fairy design.","image_url":"https://example.com/images/pucky-forest-fairy.jpg","name":"Pucky Forest Fairy","price":44.99}`
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, response.Code)
		assert.JSONEq(t, expectedJSON, response.Body.String())
	})

	t.Run("get product by id with error", func(t *testing.T) {
		mockService := new(MockProductUsecase)
		handler := &httpProductHandler{usecase: mockService}

		e := echo.New()
		defer e.Close()

		mockService.On("GetProductById", "12").Return((*entities.ProductResponse)(nil), errors.New("database error"))

		request := httptest.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()
		c := e.NewContext(request, response)
		c.SetParamNames("id")
		c.SetParamValues("12")

		err := handler.GetProductById(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, response.Code)
	})
}

type MockProductUsecase struct {
	mock.Mock
}

func (m *MockProductUsecase) CreateNewProduct(product *entities.Product) (*entities.ProductResponse, error) {
	args := m.Called(product)
	return args.Get(0).(*entities.ProductResponse), args.Error(1)
}

func (m *MockProductUsecase) GetAllProducts() ([]entities.ProductResponse, error) {
	args := m.Called()
	return args.Get(0).([]entities.ProductResponse), args.Error(1)
}

func (m *MockProductUsecase) GetProductById(id string) (*entities.ProductResponse, error) {
	args := m.Called(id)
	return args.Get(0).(*entities.ProductResponse), args.Error(1)
}

func (m *MockProductUsecase) UpdateProduct(product *entities.Product, id string) (*entities.ProductResponse, error) {
	args := m.Called(product, id)
	return args.Get(0).(*entities.ProductResponse), args.Error(1)
}
