package usecase

import (
	"errors"

	"github.com/phetployst/art-toys-store/modules/product/entities"
)

type ProductUsecase interface {
	CreateNewProduct(product *entities.Product) (*entities.Product, error)
	GetAllProducts() ([]entities.Product, error)
	GetProductById(id uint) (*entities.Product, error)
}

type ProductService struct {
	repo ProductRepository
}

func NewProductService(repo ProductRepository) ProductUsecase {
	return &ProductService{repo}
}

func (s *ProductService) CreateNewProduct(product *entities.Product) (*entities.Product, error) {
	newProduct, err := s.repo.InsertProduct(product)
	if err != nil {
		return nil, errors.New("insert product repo is error")
	}

	return newProduct, nil

}

func (s *ProductService) GetAllProducts() ([]entities.Product, error) {
	products, err := s.repo.GetAllProduct()
	if err != nil {
		return nil, errors.New("database error")
	}

	return products, nil
}

func (s *ProductService) GetProductById(productId uint) (*entities.Product, error) {
	product, err := s.repo.GetProductById(productId)
	if err != nil {
		return nil, errors.New("database error")
	}

	return product, nil
}
