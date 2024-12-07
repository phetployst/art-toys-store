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
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCreateNewProduct(t *testing.T) {
	t.Run("create new product given valid input should be successful", func(t *testing.T) {
		mockService := new(MockProductUsecase)
		handler := &httpProductHandler{usecase: mockService}

		e := echo.New()
		defer e.Close()

		mockService.On("CreateNewProduct", mock.AnythingOfType("*entities.Product")).Return(&entities.Product{Name: "Customizable Art Toy", Description: "A fully customizable art toy.", Price: 20.0, Category: "Customizable", Stock: 100}, nil)

		body := `{"name": "Customizable Art Toy","description": "A fully customizable art toy.","price": 20.0,"category": "Customizable","stock": 100}`
		request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
		request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		response := httptest.NewRecorder()
		c := e.NewContext(request, response)

		err := handler.CreateNewProduct(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, response.Code)
	})

	t.Run("create new product given error during binding", func(t *testing.T) {
		mockService := new(MockProductUsecase)
		handler := &httpProductHandler{usecase: mockService}

		e := echo.New()
		defer e.Close()

		mockService.On("CreateNewProduct", mock.AnythingOfType("*entities.Product")).Return(&entities.Product{Name: "Customizable Art Toy", Description: "A fully customizable art toy.", Price: 20.0, Category: "Customizable", Stock: 100}, nil)

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

		mockService.On("CreateNewProduct", mock.AnythingOfType("*entities.Product")).Return(&entities.Product{Name: "Customizable Art Toy", Description: "A fully customizable art toy.", Price: 20.0, Category: "Customizable", Stock: 100}, nil)

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

		mockService.On("CreateNewProduct", mock.AnythingOfType("*entities.Product")).Return((*entities.Product)(nil), errors.New("internal server error"))

		body := `{"name": "Customizable Art Toy","description": "A fully customizable art toy.","price": 20.0,"category": "Customizable","stock": 100}`
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

		products := []entities.Product{
			{ID: primitive.NewObjectID(), Name: "Customizable Art Toy", Description: "A fully customizable art toy.", Price: 20.0, Category: "Customizable", Stock: 100},
			{ID: primitive.NewObjectID(), Name: "Limited Edition Robot", Description: "A high-quality limited edition.", Price: 150.0, Category: "Collector's Item", Stock: 10},
		}

		mockService.On("GetAllProducts").Return(products, nil)

		request := httptest.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()
		c := e.NewContext(request, response)

		err := handler.GetAllProducts(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, response.Code)
	})

	t.Run("get all product given error", func(t *testing.T) {
		mockService := new(MockProductUsecase)
		handler := &httpProductHandler{usecase: mockService}

		e := echo.New()
		defer e.Close()

		mockService.On("GetAllProducts").Return(([]entities.Product)(nil), errors.New("database error"))

		request := httptest.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()
		c := e.NewContext(request, response)

		err := handler.GetAllProducts(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, response.Code)
	})
}

type MockProductUsecase struct {
	mock.Mock
}

func (m *MockProductUsecase) CreateNewProduct(product *entities.Product) (*entities.Product, error) {
	args := m.Called(product)
	return args.Get(0).(*entities.Product), args.Error(1)
}

func (m *MockProductUsecase) GetAllProducts() ([]entities.Product, error) {
	args := m.Called()
	return args.Get(0).([]entities.Product), args.Error(1)
}
