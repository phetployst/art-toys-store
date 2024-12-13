package usecase

import (
	"errors"
	"reflect"
	"testing"

	"github.com/phetployst/art-toys-store/modules/product/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateNewProduct(t *testing.T) {
	t.Run("create new product successfully", func(t *testing.T) {
		mockRepo := new(MockProductRepository)
		productService := ProductService{repo: mockRepo}

		newProduct := &entities.Product{
			Name:        "Molly Classic",
			Description: "The iconic Molly figure, loved by art toy collectors worldwide.",
			Price:       340.99,
			Stock:       30,
			ImageURL:    "https://example.com/images/molly-classic.jpg",
			Active:      true,
		}

		mockRepo.On("InsertProduct", mock.AnythingOfType("*entities.Product")).Return(newProduct, nil)

		got, err := productService.CreateNewProduct(newProduct)

		assert.NoError(t, err)

		if !reflect.DeepEqual(got, newProduct) {
			t.Errorf("got %v but want %v", got, newProduct)
		}
	})

	t.Run("create new product error during query", func(t *testing.T) {
		mockRepo := new(MockProductRepository)
		productService := ProductService{repo: mockRepo}

		newProduct := &entities.Product{
			Name:        "Skull Panda Rebel",
			Description: "A rebellious design from Skull Panda, combining gothic aesthetics with modern art.",
			Price:       59.99,
			Stock:       15,
			ImageURL:    "https://example.com/images/skull-panda-rebel.jpg",
			Active:      true,
		}

		mockRepo.On("InsertProduct", mock.AnythingOfType("*entities.Product")).Return((*entities.Product)(nil), errors.New("databse error"))

		_, err := productService.CreateNewProduct(newProduct)

		assert.Error(t, err)
		assert.EqualError(t, err, "insert product repo is error")
	})
}

func TestGetAllProducts(t *testing.T) {
	t.Run("get all product successfully", func(t *testing.T) {
		mockRepo := new(MockProductRepository)
		productService := ProductService{repo: mockRepo}

		products := []entities.Product{
			{Name: "Dimoo Starry Night", Description: "Dimoo inspired by Van Gogh's 'Starry Night,' featuring a dreamy and artistic design.", Price: 49.99, Stock: 25, ImageURL: "https://example.com/images/dimoo-starry-night.jpg", Active: true},
			{Name: "Pucky Forest Fairy", Description: "A magical art toy figure from Pucky, with a whimsical forest fairy design.", Price: 44.99, Stock: 40, ImageURL: "https://example.com/images/pucky-forest-fairy.jpg", Active: true},
		}

		mockRepo.On("GetAllProduct").Return(products, nil)

		got, err := productService.GetAllProducts()

		assert.NoError(t, err)

		if !reflect.DeepEqual(got, products) {
			t.Errorf("got %v but want %v", got, products)
		}

	})

	t.Run("get all product given database error", func(t *testing.T) {
		mockRepo := new(MockProductRepository)
		productService := ProductService{repo: mockRepo}

		mockRepo.On("GetAllProduct").Return(([]entities.Product)(nil), errors.New("database error"))

		_, err := productService.GetAllProducts()

		assert.Error(t, err)
		assert.EqualError(t, err, "database error")

	})
}

func TestGetProductById(t *testing.T) {
	t.Run("get product by id successfully", func(t *testing.T) {
		mockRepo := new(MockProductRepository)
		productService := ProductService{repo: mockRepo}

		product := &entities.Product{Name: "Pucky Forest Fairy", Description: "A magical art toy figure from Pucky, with a whimsical forest fairy design.",
			Price: 44.99, Stock: 40, ImageURL: "https://example.com/images/pucky-forest-fairy.jpg", Active: true}

		mockRepo.On("GetProductById", uint(12)).Return(product, nil)

		got, err := productService.GetProductById(uint(12))

		assert.NoError(t, err)
		assert.Equal(t, "Pucky Forest Fairy", got.Name)

	})

	t.Run("get product by id given error", func(t *testing.T) {
		mockRepo := new(MockProductRepository)
		productService := ProductService{repo: mockRepo}

		mockRepo.On("GetProductById", uint(13)).Return((*entities.Product)(nil), errors.New("database error"))

		_, err := productService.GetProductById(uint(13))

		assert.Error(t, err)
		assert.EqualError(t, err, "database error")
	})
}

type MockProductRepository struct {
	mock.Mock
}

func (m *MockProductRepository) InsertProduct(product *entities.Product) (*entities.Product, error) {
	args := m.Called(product)
	return args.Get(0).(*entities.Product), args.Error(1)
}

func (m *MockProductRepository) GetAllProduct() ([]entities.Product, error) {
	args := m.Called()
	return args.Get(0).([]entities.Product), args.Error(1)
}

func (m *MockProductRepository) GetProductById(id uint) (*entities.Product, error) {
	args := m.Called(id)
	return args.Get(0).(*entities.Product), args.Error(1)
}
