package usecase

import (
	"errors"

	"github.com/phetployst/art-toys-store/modules/product/entities"
)

type ProductUsecase interface {
	CreateNewProduct(product *entities.Product) (*entities.Product, error)
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