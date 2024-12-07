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

type MockProductRepository struct {
	mock.Mock
}

func (m *MockProductRepository) InsertProduct(product *entities.Product) (*entities.Product, error) {
	args := m.Called(product)
	return args.Get(0).(*entities.Product), args.Error(1)
}
