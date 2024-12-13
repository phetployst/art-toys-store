package usecase

import "github.com/phetployst/art-toys-store/modules/product/entities"

type ProductRepository interface {
	InsertProduct(product *entities.Product) (*entities.Product, error)
	GetAllProduct() ([]entities.Product, error)
	GetProductById(id uint) (*entities.Product, error)
}
