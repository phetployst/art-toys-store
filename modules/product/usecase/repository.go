package usecase

import "github.com/phetployst/art-toys-store/modules/product/entities"

type ProductRepository interface {
	InsertProduct(product *entities.Product) (*entities.Product, error)
	GetAllProduct() ([]entities.Product, error)
	GetProductById(id string) (*entities.Product, error)
	UpdateProduct(product *entities.Product, id string) (*entities.Product, error)
	UpdateStock(id string, count int) (int, error)
	SearchProducts(keyword string) ([]entities.Product, error)
}
