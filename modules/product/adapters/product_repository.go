package adapters

import (
	"github.com/phetployst/art-toys-store/modules/product/entities"
	"github.com/phetployst/art-toys-store/modules/product/usecase"
	"gorm.io/gorm"
)

type gormProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) usecase.ProductRepository {
	return &gormProductRepository{db}
}

func (r *gormProductRepository) InsertProduct(product *entities.Product) (*entities.Product, error) {
	if result := r.db.Create(&product); result.Error != nil {
		return nil, result.Error
	}

	return product, nil
}

func (r *gormProductRepository) GetAllProduct() ([]entities.Product, error) {
	var products []entities.Product

	result := r.db.Find(&products)
	if result.Error != nil {
		return nil, result.Error
	}

	return products, nil
}

func (r *gormProductRepository) GetProductById(id string) (*entities.Product, error) {
	product := new(entities.Product)

	if err := r.db.First(&product, id).Error; err != nil {
		return nil, err
	}

	return product, nil
}
