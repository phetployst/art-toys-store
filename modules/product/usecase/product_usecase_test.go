package usecase

import (
	"errors"
	"reflect"
	"testing"

	"github.com/phetployst/art-toys-store/modules/product/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCreateNewProduct(t *testing.T) {
	t.Run("create new product successfully", func(t *testing.T) {
		mockRepo := new(MockProductRepository)
		productService := ProductService{repo: mockRepo}

		newProduct := &entities.Product{
			Name:        "Customizable Art Toy",
			Description: "A fully customizable art toy allowing users to pick colors and designs to suit their style.",
			Price:       20.0,
			Category:    "Customizable",
			Stock:       100,
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
			Name:        "Eco-Friendly Art Toy",
			Description: "An environmentally friendly art toy made from recycled materials with sustainable packaging.",
			Price:       18.0,
			Category:    "Eco-Friendly",
			Stock:       80,
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
			{ID: primitive.NewObjectID(), Name: "Customizable Art Toy", Description: "A fully customizable art toy.", Price: 20.0, Category: "Customizable", Stock: 100},
			{ID: primitive.NewObjectID(), Name: "Limited Edition Robot", Description: "A high-quality limited edition.", Price: 150.0, Category: "Collector's Item", Stock: 10},
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
