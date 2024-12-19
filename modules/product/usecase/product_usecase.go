package usecase

import (
	"errors"

	"github.com/phetployst/art-toys-store/modules/product/entities"
)

type ProductUsecase interface {
	CreateNewProduct(product *entities.Product) (*entities.ProductResponse, error)
	GetAllProducts() ([]entities.ProductResponse, error)
	GetProductById(id string) (*entities.ProductResponse, error)
	UpdateProduct(product *entities.Product, id string) (*entities.ProductResponse, error)
}

type ProductService struct {
	repo ProductRepository
}

func NewProductService(repo ProductRepository) ProductUsecase {
	return &ProductService{repo}
}

func (s *ProductService) CreateNewProduct(product *entities.Product) (*entities.ProductResponse, error) {
	newProduct, err := s.repo.InsertProduct(product)
	if err != nil {
		return nil, errors.New("database error")
	}

	return &entities.ProductResponse{
		ID:          newProduct.ID,
		Name:        newProduct.Name,
		Description: newProduct.Description,
		Price:       newProduct.Price,
		ImageURL:    newProduct.ImageURL,
	}, nil

}

func (s *ProductService) GetAllProducts() ([]entities.ProductResponse, error) {
	products, err := s.repo.GetAllProduct()
	if err != nil {
		return nil, errors.New("database error")
	}

	var productList []entities.ProductResponse
	for _, product := range products {
		productList = append(productList, entities.ProductResponse{
			ID:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			ImageURL:    product.ImageURL,
		})
	}

	return productList, nil
}

func (s *ProductService) GetProductById(productId string) (*entities.ProductResponse, error) {
	product, err := s.repo.GetProductById(productId)
	if err != nil {
		return nil, errors.New("database error")
	}

	return &entities.ProductResponse{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		ImageURL:    product.ImageURL,
	}, nil
}

func (s *ProductService) UpdateProduct(product *entities.Product, id string) (*entities.ProductResponse, error) {
	productUpdated, err := s.repo.UpdateProduct(product, id)
	if err != nil {
		return nil, errors.New("database error")
	}

	return &entities.ProductResponse{
		ID:          productUpdated.ID,
		Name:        productUpdated.Name,
		Description: productUpdated.Description,
		Price:       productUpdated.Price,
		ImageURL:    productUpdated.ImageURL,
	}, nil
}
